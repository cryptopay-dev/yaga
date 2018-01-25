package workers

type worker struct {
	job     func()
	options Options
	pool    *pool
}
