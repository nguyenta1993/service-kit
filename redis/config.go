package redis

type Config struct {
	Addrs    []string
	Password string
	DB       int
	PoolSize int
}
