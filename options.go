package tourniquet

// Option handles configurable options.
type Option interface {
	apply(option *funcOption)
	withCustomErrorOnCloseHandler() func(error)
}

type funcOption struct {
	f                         func(*funcOption)
	customErrorOnCloseHandler func(error)
}

func (f *funcOption) apply(option *funcOption) {
	f.f(option)
}

func (f *funcOption) withCustomErrorOnCloseHandler() func(error) {
	return f.customErrorOnCloseHandler
}

// WithCustomErrorOnCloseHandler allows to override the default behaviour when an
// error occurred while closing a connection when the TTL provided has expired.
// By default, getting a connection may return an error if the close failed.
func WithCustomErrorOnCloseHandler(handler func(error)) Option {
	return newFuncOption(func(options *funcOption) {
		options.customErrorOnCloseHandler = handler
	})
}

func newFuncOption(f func(*funcOption)) *funcOption {
	return &funcOption{
		f: f,
	}
}

func parseOptions(opts ...Option) Option {
	o := new(funcOption)
	for _, opt := range opts {
		opt.apply(o)
	}
	return o
}
