package hihttp

import (
	"context"
	"time"
)

type Config struct {
	Timeout       time.Duration
	Prefix        string
	RetryCount    int
	RetryWaitTime time.Duration
	RetryError    func(ctx context.Context, c hiclient) error
}

// New method creates a new Resty client.
func New(cfg Config) {
	client.retryCount = cfg.RetryCount
	client.retryWaitTime = cfg.RetryWaitTime
	if cfg.RetryError == nil {
		cfg.RetryError = func(ctx context.Context, c hiclient) error {
			return nil
		}
	}
	client.retryError = cfg.RetryError
	if cfg.Timeout == 0 {
		cfg.Timeout = 10 * time.Second
	}
	client.Timeout = cfg.Timeout
}
