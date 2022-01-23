package hihttp

import (
	"context"
	"time"
)

func (c hiclient) retry(ctx context.Context, retryFunc func(ctx context.Context) ([]byte, error)) []byte {
	// 记录重试次数
	for count := 0; c.RetryCount > count; count++ {
		time.Sleep(c.RetryWaitTime)
		req, err := retryFunc(ctx)
		if err != nil {
			continue
		}
		return req
	}

	c.RetryError(ctx, c)
	return nil
}
