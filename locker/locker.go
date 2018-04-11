package locker

// Options interface
type Options interface {
	Parse(opts ...Option) error
}

// Option closure
type Option = func(opt Options) error

// Locker interface for drivers
type Locker interface {
	Run(key string, handler func(), options ...Option) error
}
