package hihttp

import "time"

// 几个公共参数
type Options struct {
	// 重试次数
	retryCount int
	// 重试等待时间
	retryWait  time.Duration
	retryError RetryErrorFunc
	// 超时时间
	timeout time.Duration
}
type Option func(*Options)

func WithRetryCount(retryCount int) Option {
	return func(o *Options) {
		o.retryCount = retryCount
	}
}

func WithRetryWait(retryWait time.Duration) Option {
	return func(o *Options) {
		o.retryWait = retryWait
	}
}
func WithRetryError(retryError RetryErrorFunc) Option {
	return func(o *Options) {
		o.retryError = retryError
	}
}
func WithTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.timeout = timeout
	}
}
