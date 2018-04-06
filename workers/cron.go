package workers

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/cryptopay-dev/yaga/logger"
	"github.com/cryptopay-dev/yaga/logger/log"
	"go.uber.org/atomic"
)

// Cron keeps track of any number of entries, invoking the associated func as
// specified by the schedule. It may be started, stopped, and the entries may
// be inspected while running.
type Cron struct {
	entries []*Entry
	done    chan struct{}
	add     chan *Entry
	state   *atomic.Int32
	logger  logger.Logger
	lockers map[TypeJob]Locker
}

// The Schedule describes a job's duty cycle.
type Schedule interface {
	// Return the next activation time, later than the given time.
	// Next is invoked initially, and then each time the job is run.
	Next(time.Time) time.Time
}

// Entry consists of a schedule and the func to execute on that schedule.
type Entry struct {
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
type byTime []*Entry

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

// New returns a new Cron job runner, in the Local time zone.
func New(lockers ...Locker) *Cron {
	c := &Cron{
		entries: nil,
		add:     make(chan *Entry),
		done:    make(chan struct{}),
		state:   atomic.NewInt32(0),
		logger:  log.Logger(),
		lockers: make(map[TypeJob]Locker),
	}
	for _, locker := range lockers {
		c.lockers[locker.TypeJob()] = locker
	}
	return c
}

func (c *Cron) schedule(entry *Entry) {
	if c.state.Load() == 0 {
		c.entries = append(c.entries, entry)
		return
	}

	select {
	case c.add <- entry:
	case <-c.done:
	}
}

// Start the cron scheduler in its own go-routine, or no-op if already started.
func (c *Cron) Start(ctx context.Context) {
	if c.state.CAS(0, 1) {
		go c.run(ctx)
	}
}

// Run the scheduler. this is private just due to the need to synchronize
// access to the 'running' state variable.
func (c *Cron) run(ctx context.Context) {
	// Figure out the next activation times for each entry.
	now := time.Now()
	wg := new(sync.WaitGroup)
	for _, entry := range c.entries {
		entry.Next = entry.Schedule.Next(now)
	}

	for {
		// Determine the next entry to run.
		sort.Sort(byTime(c.entries))

		var timer *time.Timer
		if len(c.entries) == 0 || c.entries[0].Next.IsZero() {
			// If there are no entries yet, just sleep - it still handles new entries
			// and stop requests.
			timer = time.NewTimer(100000 * time.Hour)
		} else {
			timer = time.NewTimer(c.entries[0].Next.Sub(now))
		}

		select {
		case now = <-timer.C:
			// Run every entry whose next time was less than now
			for _, e := range c.entries {
				if e.Next.After(now) || e.Next.IsZero() {
					break
				}
				wg.Add(1)
				go func(e2 *Entry) {
					defer func() {
						wg.Done()
						if r := recover(); r != nil {
							c.logger.Errorf("workers `%s` panic: %v", e2.Name, r)
						}
					}()
					e2.Job(ctx)
				}(e)
				e.Prev = e.Next
				e.Next = e.Schedule.Next(now)
			}

		case newEntry := <-c.add:
			timer.Stop()
			now = time.Now()
			newEntry.Next = newEntry.Schedule.Next(now)
			c.entries = append(c.entries, newEntry)

		case <-ctx.Done():
			c.state.Store(2)
			timer.Stop()
			wg.Wait()
			close(c.done)
			return
		}
	}
}

func (c *Cron) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-c.done:
	}

	return nil
}
