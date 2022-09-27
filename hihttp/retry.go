package hihttp

import (
	"context"
	"io"
	"time"
)

// retry 请求失败后的重试方法
// 传入的retryFunc是在失败时需要重新执行的方法，如c.execute
// c.RetryError 是重试如果也失败了，需回调通知
func (r *Request) retry(ctx context.Context, payload io.Reader) []byte {
	// 记录重试次数
	for count := 0; r.opt.retryCount > count; count++ {
		time.Sleep(r.opt.retryWait)
		req, err := r.execute(ctx, payload)
		if err != nil {
			continue
		}
		return req
	}

	r.opt.retryError(ctx, r)
	return nil
}
