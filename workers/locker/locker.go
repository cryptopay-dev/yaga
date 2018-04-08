package locker

type Locker interface {
	Run(key string, handler func())
}
