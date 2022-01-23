package hihttp

import (
	"context"
	"time"
)

type Config struct {
	Prefix        string
	RetryCount    int
	RetryWaitTime time.Duration
	RetryError    func(ctx context.Context, c hiclient) error
}

// New method creates a new Resty client.
func New(cfg Config) {
	client.prefix = cfg.Prefix
	client.RetryCount = cfg.RetryCount
	client.RetryWaitTime = cfg.RetryWaitTime
	client.RetryError = cfg.RetryError
}
