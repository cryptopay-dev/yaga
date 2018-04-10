package pool

import (
	"context"
	"sync"
)

func Run(ctx context.Context, jobCh <-chan func(context.Context)) {
	size := cap(jobCh)
	wg := new(sync.WaitGroup)
	wg.Add(size)
	for i := 0; i < size; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case job, ok := <-jobCh:
					if !ok {
						return
					}
					if job != nil {
						job(ctx)
					}
				case <-ctx.Done():
					return
				}
			}
		}()
	}
	wg.Wait()
}
