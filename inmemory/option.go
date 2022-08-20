package inmemory

import "time"

type Options struct {
	maxLife time.Duration
}

type Option func(o *Options)

func WithMaxLife(d time.Duration) Option {
	return func(o *Options) {
		o.maxLife = d
	}
}
