package server

type HttpServerConfig struct {
	Port            string
	Development     bool
	ShutdownTimeout int

	Resources    []string
	RateLimiting *RateLimitingConfig
	Name         string
}

type RateLimitingConfig struct {
	RateFormat string
}
