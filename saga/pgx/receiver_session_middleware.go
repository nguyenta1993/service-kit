package pgx

import (
	"context"
	"fmt"

	"github.com/gogovan/ggx-kr-service-utils/logger"
	"github.com/gogovan/ggx-kr-service-utils/saga/msg"
	"go.uber.org/zap"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

func ReceiverSessionMiddleware(conn *pgxpool.Pool, logger logger.Logger) func(msg.MessageReceiver) msg.MessageReceiver {
	return func(next msg.MessageReceiver) msg.MessageReceiver {
		return msg.ReceiveMessageFunc(func(ctx context.Context, message msg.Message) (err error) {
			var tx pgx.Tx

			tx, err = conn.Begin(ctx)
			if err != nil {
				logger.Error("error while starting the request transaction", zap.Error(err))
				return fmt.Errorf("failed to start transaction: %s", err.Error())
			}

			txCtx := context.WithValue(ctx, pgxTxKey, tx)

			defer func() {
				p := recover()
				switch {
				case p != nil:
					txErr := tx.Rollback(ctx)
					if txErr != nil {
						logger.Error("error while rolling back the message receiver transaction during panic", zap.Error(txErr))
					}
					panic(p)
				case err != nil:
					txErr := tx.Rollback(ctx)
					if txErr != nil {
						logger.Error("error while rolling back the message receiver transaction", zap.Error(txErr))
					}
				default:
					txErr := tx.Commit(ctx)
					if txErr != nil {
						logger.Error("error while committing the message receiver transaction", zap.Error(txErr))
					}
				}
			}()

			err = next.ReceiveMessage(txCtx, message)

			return err
		})
	}
}
