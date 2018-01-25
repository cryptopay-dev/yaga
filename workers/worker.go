package workers

type cronHandler func(*worker, func())

type worker struct {
	job     func()
	options Options
	pool    *pool
}

func newWorker(opts Options, p *pool, addToCron cronHandler) (*worker, error) {
	if opts.Schedule == nil || opts.Handler == nil {
		return nil, ErrWrongOptions
	}

	w, err := p.createWorker(opts)
	if err != nil {
		return nil, err
	}
	w.pool = p

	addToCron(w, w.job)

	return w, nil
}
