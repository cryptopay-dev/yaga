package workers

import (
	"sync"

	"github.com/cryptopay-dev/yaga/logger/log"
	"github.com/robfig/cron"
)

type Schedule = cron.Schedule

type Cmd = func() error

type Entry = cron.Entry

type Workers struct {
	wg   sync.WaitGroup
	cron *cron.Cron
}

func New() *Workers {
	w := &Workers{
		wg:   sync.WaitGroup{},
		cron: cron.New(),
	}

	return w
}

func (w *Workers) AddFunc(spec string, cmd Cmd) error {
	return w.cron.AddFunc(spec, func() {
		w.wg.Add(1)
		if err := cmd(); err != nil {
			log.Error(err)
		}
		w.wg.Done()
	})
}

type runner struct{ handler func() }

func (r *runner) Run() { r.handler() }

func (w *Workers) AddJob(spec string, cmd Cmd) error {
	return w.cron.AddJob(spec, &runner{
		handler: func() {
			w.wg.Add(1)
			if err := cmd(); err != nil {
				log.Error(err)
			}
			w.wg.Done()
		},
	})
}

func (w *Workers) Schedule(schedule Schedule, cmd Cmd) {
	w.cron.Schedule(schedule, &runner{
		handler: func() {
			w.wg.Add(1)
			if err := cmd(); err != nil {
				log.Error(err)
			}
			w.wg.Done()
		},
	})
}

func (w *Workers) Start() { w.cron.Start() }

func (w *Workers) Stop() {
	w.cron.Stop()
	w.wg.Wait()
}

func (w *Workers) Entries() []*Entry { return w.cron.Entries() }
