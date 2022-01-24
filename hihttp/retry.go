package hihttp

import (
	"context"
	"time"
)

// retry 请求失败后的重试方法
// 传入的retryFunc是在失败时需要重新执行的方法，如c.execute
// c.RetryError 是重试如果也失败了，需回调通知
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
