// Based on https://github.com/golang/go/blob/master/src/sync/waitgroup_test.go
//
// Changed tests:
// Now no skipping known-racy test under the race detector
// TestWaitGroupMisuse2 => TestWaitGroupNoMisuse2
// TestWaitGroupMisuse3 => TestWaitGroupNoMisuse3

package workers_test

import (
	"runtime"
	"sync/atomic"
	"testing"

	. "github.com/cryptopay-dev/yaga/workers"
)

func testWaitGroup(t *testing.T, wg1 WaitGroup, wg2 WaitGroup) {
	n := 16
	wg1.Add(n)
	wg2.Add(n)
	exited := make(chan bool, n)
	for i := 0; i != n; i++ {
		go func() {
			wg1.Done()
			wg2.Wait()
			exited <- true
		}()
	}
	wg1.Wait()
	for i := 0; i != n; i++ {
		select {
		case <-exited:
			t.Fatal("WaitGroup released group too soon")
		default:
		}
		wg2.Done()
	}
	for i := 0; i != n; i++ {
		<-exited // Will block if barrier fails to unlock someone.
	}
}

func TestWaitGroup(t *testing.T) {
	wg1 := NewWaitGroup()
	wg2 := NewWaitGroup()

	// Run the same test a few times to ensure barrier is in a proper state.
	for i := 0; i != 8; i++ {
		testWaitGroup(t, wg1, wg2)
	}
}

func TestWaitGroupMisuse(t *testing.T) {
	defer func() {
		err := recover()
		if err != "waitgroup: negative WaitGroup counter" {
			t.Fatalf("Unexpected panic: %#v", err)
		}
	}()
	wg := NewWaitGroup()
	wg.Add(1)
	wg.Done()
	wg.Done()
	t.Fatal("Should panic")
}

func TestWaitGroupNoMisuse2(t *testing.T) {
	if runtime.NumCPU() <= 4 {
		t.Skip("NumCPU<=4, skipping: this test requires parallelism")
	}
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(4))
	done := make(chan interface{}, 2)
	for i := 0; i < 1000; i++ {
		wg := NewWaitGroup()
		var here uint32
		wg.Add(1)
		go func() {
			defer func() {
				done <- recover()
			}()
			atomic.AddUint32(&here, 1)
			for atomic.LoadUint32(&here) != 3 {
				// spin
			}
			wg.Wait()
		}()
		go func() {
			defer func() {
				done <- recover()
			}()
			atomic.AddUint32(&here, 1)
			for atomic.LoadUint32(&here) != 3 {
				// spin
			}
			wg.Add(1) // This is the bad guy.
			wg.Done()
		}()
		atomic.AddUint32(&here, 1)
		for atomic.LoadUint32(&here) != 3 {
			// spin
		}
		wg.Done()
		for j := 0; j < 2; j++ {
			if err := <-done; err != nil {
				panic(err)
			}
		}
	}
}

func TestWaitGroupNoMisuse3(t *testing.T) {
	if runtime.NumCPU() <= 1 {
		t.Skip("NumCPU==1, skipping: this test requires parallelism")
	}
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(4))
	done := make(chan interface{}, 2)
	for i := 0; i < 1000; i++ {
		wg := NewWaitGroup()
		wg.Add(1)
		go func() {
			defer func() {
				done <- recover()
			}()
			wg.Done()
		}()
		go func() {
			defer func() {
				done <- recover()
			}()
			wg.Wait()
			// Start reusing the wg before waiting for the Wait below to return.
			wg.Add(1)
			go func() {
				wg.Done()
			}()
			wg.Wait()
		}()
		wg.Wait()
		for j := 0; j < 2; j++ {
			if err := <-done; err != nil {
				panic(err)
			}
		}
	}
}

func TestWaitGroupRace(t *testing.T) {
	// Run this test for about 1ms.
	for i := 0; i < 1000; i++ {
		wg := NewWaitGroup()
		n := new(int32)
		// spawn goroutine 1
		wg.Add(1)
		go func() {
			atomic.AddInt32(n, 1)
			wg.Done()
		}()
		// spawn goroutine 2
		wg.Add(1)
		go func() {
			atomic.AddInt32(n, 1)
			wg.Done()
		}()
		// Wait for goroutine 1 and 2
		wg.Wait()
		if atomic.LoadInt32(n) != 2 {
			t.Fatal("Spurious wakeup from Wait")
		}
	}
}

func TestWaitGroupAlign(t *testing.T) {
	type X struct {
		x  byte
		wg WaitGroup
	}
	var x X
	x.wg = NewWaitGroup()
	x.wg.Add(1)
	go func(x *X) {
		x.wg.Done()
	}(&x)
	x.wg.Wait()
}
