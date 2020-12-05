package gox

import (
	"time"
)

// Every executes a given function `f` each time `d` passes, until a boolean is received in the channel it returns
// Don't forget to manage yourself the concurrent access
func Every(d time.Duration, f func()) (stop RequestChan) {
	stop, shouldStop := NewRequestChan()
	Report(func() {
		ticker := time.NewTicker(d)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				f()
			case stopped := <-shouldStop:
				close(stopped)
				return
			}
		}
	})
	return stop
}

// NowAndEvery is like Every, but first runs the function right away in the
// caller goroutine.
func NowAndEvery(d time.Duration, f func()) (stop RequestChan) {
	f()
	return Every(d, f)
}
