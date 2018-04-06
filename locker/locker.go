package locker

// Options interface
type Options interface {
	Parse(opts ...Option)
}

// Option closure
type Option = func(opt Options) error

// Locker interface to abstract bsm/redis-lock
type Locker interface {
	Run(key string, handler func(), options ...Option)
}
