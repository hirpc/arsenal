package hsync

import "time"

type Options struct {
	timeout       time.Duration
	waitingPeriod time.Duration
	retry         bool
}

var options = Options{
	timeout:       time.Second * time.Duration(10),
	waitingPeriod: time.Millisecond * time.Duration(50),
	retry:         true,
}

type Option func(*Options)

func WithTimeout(t time.Duration) Option {
	return func(o *Options) {
		options.timeout = t
	}
}

func WithWaitingPeriod(t time.Duration) Option {
	return func(o *Options) {
		options.waitingPeriod = t
	}
}

func WithDisableRetry() Option {
	return func(o *Options) {
		options.retry = false
	}
}
