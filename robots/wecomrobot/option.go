package wecomrobot

type Options struct {
	defaultRegistration bool
}

type Option func(o *Options)

func WithDefaultRegistration() Option {
	return func(o *Options) {
		o.defaultRegistration = true
	}
}
