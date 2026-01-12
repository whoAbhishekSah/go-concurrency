package main

import (
	"fmt"
	"time"
)

// schedule starts executing a function at the
// specified interval and provides a way to stop it.
// Once canceled, the function stops executing.
func schedule(dur time.Duration, fn func()) func() {
	done := make(chan struct{}, 1)
	timer := time.NewTicker(dur)
	cancel := func() {
		done <- struct{}{}
	}
	go func() {
		for {
			select {
			case <-timer.C:
				fn()
			case <-done:
				return
			}
		}
	}()

	return cancel
}

func main() {
	work := func() {
		at := time.Now()
		fmt.Printf("%s: work done\n", at.Format("15:04:05.000"))
	}
	cancel := schedule(50*time.Millisecond, work)
	defer cancel()
	// enough for 5 ticks
	time.Sleep(260 * time.Millisecond)
}
