package hihttp

import (
	"context"
	"io"
	"time"
)

// retry 请求失败后的重试方法
// 传入的retryFunc是在失败时需要重新执行的方法，如c.execute
// c.RetryError 是重试如果也失败了，需回调通知
func (c hiclient) retry(ctx context.Context, h HiHTTP, payload io.Reader) []byte {
	// 记录重试次数
	for count := 0; c.retryCount > count; count++ {
		time.Sleep(c.retryWaitTime)
		req, err := h.execute(ctx, payload)
		if err != nil {
			continue
		}
		return req
	}

	c.retryError(ctx, c)
	return nil
}
