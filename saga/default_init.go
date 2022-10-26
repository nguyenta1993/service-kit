package saga

import (
	"context"
	"sync"
	"time"

	"github.com/tikivn/s14e-backend-utils/logger"
	"github.com/tikivn/s14e-backend-utils/saga/msg"
	_ "github.com/tikivn/s14e-backend-utils/saga/msgpack"
	"github.com/tikivn/s14e-backend-utils/saga/saga"

	pgx "github.com/tikivn/s14e-backend-utils/saga/pgx"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"golang.org/x/sync/errgroup"
)

type SagaService struct {
	Logger            logger.Logger
	PgConn            pgx.Client
	SagaInstanceStore saga.InstanceStore
	Publisher         *msg.Publisher
	Subscriber        *msg.Subscriber
}

// Return the saga service and waiter function
func NewSagaStore(ctx context.Context, log logger.Logger, producer msg.Producer, consumer msg.Consumer, pgConnStr string, waitFuncs ...func(ctx context.Context) error) (*SagaService, func(context.Context) error, func()) {
	s := &SagaService{
		Logger: log,
	}

	var pgConn *pgxpool.Pool
	pgConn, err := pgxpool.Connect(ctx, pgConnStr)
	if err != nil {
		panic(err)
	}

	// 1. Outbox: Use session client which will fetch a transaction from the context
	s.PgConn = pgx.NewSessionClient()

	s.SagaInstanceStore = pgx.NewSagaInstanceStore(log, pgConn)
	s.Subscriber = msg.NewSubscriber(consumer, log)
	s.Subscriber.Use(
		MessageInstrumentation(),
		// 3. Outbox: Use a message receiver middleware to start a new transaction for each incoming message
		pgx.ReceiverSessionMiddleware(pgConn, s.Logger),
	)
	s.Publisher = msg.NewPublisher(producer, log)

	closeFunc := func() {
		if pgConn != nil {
			pgConn.Close()
		}
	}
	return s, s.waitForMessaging, closeFunc
}

func MessageInstrumentation() func(msg.MessageReceiver) msg.MessageReceiver {
	responseTime := promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "message_response_time",
		Help:    "Message response time in microseconds",
		Buckets: []float64{300, 600, 900, 1_500, 5_000, 10_000, 20_000},
	})

	return func(next msg.MessageReceiver) msg.MessageReceiver {
		return msg.ReceiveMessageFunc(func(ctx context.Context, message msg.Message) error {
			start := time.Now()
			err := next.ReceiveMessage(ctx, message)
			responseTime.Observe(float64(time.Since(start).Microseconds()))

			return err
		})
	}
}

func (s SagaService) waitForMessaging(ctx context.Context) error {

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		return s.Subscriber.Start(ctx)
	})

	group.Go(func() error {
		<-gCtx.Done()
		sCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			defer wg.Done()
			if s.Publisher != nil {
				if err := s.Publisher.Stop(sCtx); err != nil {
					s.Logger.Error("error while shutting down publisher")
				}
			}
		}()
		go func() {
			defer wg.Done()
			if s.Subscriber != nil {
				if err := s.Subscriber.Stop(sCtx); err != nil {
					s.Logger.Error("error while shutting down subscriber")

				}
			}
		}()
		done := make(chan struct{})
		go func() {
			defer close(done)
			wg.Wait()
		}()
		select {
		case <-done:
		case <-sCtx.Done():
		}
		return nil
	})

	return group.Wait()
}
