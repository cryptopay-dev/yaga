package workers

import (
	"context"
	"sort"
	"time"

	"github.com/cryptopay-dev/yaga/workers/pool"
)

// The Schedule describes a job's duty cycle.
type Schedule interface {
	// Return the next activation time, later than the given time.
	// Next is invoked initially, and then each time the job is run.
	Next(time.Time) time.Time
}

// Entry consists of a schedule and the func to execute on that schedule.
type entry struct {
	Name string

	// The schedule on which this job should be run.
	Schedule Schedule

	// The next time the job will run. This is the zero time if Cron has not been
	// started or this entry's schedule is unsatisfiable
	Next time.Time

	// The last time this job was run. This is the zero time if the job has never
	// been run.
	Prev time.Time

	// The Job to run.
	Job func(ctx context.Context)
}

// byTime is a wrapper for sorting the entry array by time
// (with zero time at the end).
type byTime []*entry

func (s byTime) Len() int      { return len(s) }
func (s byTime) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byTime) Less(i, j int) bool {
	// Two zero times should return false.
	// Otherwise, zero is "greater" than any other time.
	// (To sort it at the end of the list.)
	if s[i].Next.IsZero() {
		return false
	}
	if s[j].Next.IsZero() {
		return true
	}
	return s[i].Next.Before(s[j].Next)
}

func (w *Workers) schedule(entry *entry) {
	if w.state.Load() == 0 {
		w.entries = append(w.entries, entry)
		return
	}

	w.add <- entry
}

func (w *Workers) Start(ctx context.Context) {
	if w.state.CAS(0, 1) {
		jobCh := make(chan func(context.Context), w.size)
		go w.run(ctx, jobCh)
		pool.Run(ctx, jobCh)
	}
}

func (w *Workers) run(ctx context.Context, jobCh chan func(context.Context)) {
	// Figure out the next activation times for each entry.
	now := time.Now()
	for _, entry := range w.entries {
		entry.Next = entry.Schedule.Next(now)
	}

	for {
		// Determine the next entry to run.
		sort.Sort(byTime(w.entries))

		var timer *time.Timer
		if len(w.entries) == 0 || w.entries[0].Next.IsZero() {
			// If there are no entries yet, just sleep - it still handles new entries
			// and stop requests.
			timer = time.NewTimer(100000 * time.Hour)
		} else {
			timer = time.NewTimer(w.entries[0].Next.Sub(now))
		}

		select {
		case now = <-timer.C:
			// Run every entry whose next time was less than now
			for _, e := range w.entries {
				if e.Next.After(now) || e.Next.IsZero() {
					break
				}
				select {
				case jobCh <- e.Job:
				default:
					w.logger.Errorf("workers `%s` cannot run", e.Name)
				}
				e.Prev = e.Next
				e.Next = e.Schedule.Next(now)
			}

		case newEntry := <-w.add:
			timer.Stop()
			now = time.Now()
			newEntry.Next = newEntry.Schedule.Next(now)
			w.entries = append(w.entries, newEntry)

		case <-ctx.Done():
			w.state.Store(2)
			timer.Stop()
			return
		}
	}
}
