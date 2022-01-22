package validation

type Options struct {
	Level level
}

type Option func(o *Options)

func WithLevel(l level) Option {
	return func(o *Options) {
		o.Level = l
	}
}
