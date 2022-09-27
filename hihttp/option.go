package hihttp

import "time"

// Options 几个公共参数
type Options struct {
	// retryCount 重试次数
	retryCount int
	// retryWait 重试等待时间
	retryWait  time.Duration
	retryError RetryErrorFunc
	// timeout 超时时间
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
