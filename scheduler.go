package main

import (
	"fmt"
	"time"
)

// schedule starts executing a function at the
// specified interval and provides a way to stop it.
// Once canceled, the function stops executing.
func schedule(dur time.Duration, fn func()) func() {
	cancel := make(chan struct{}, 1)
	done := make(chan struct{}, 1)
	timer := time.NewTicker(dur)
	cancelFn := func() {
		cancel <- struct{}{}
	}
	go func() {
		done <- struct{}{}
		for {
			select {
			case <-timer.C:
				select {
				case <-done:
					go func() {
						fn()
						done <- struct{}{}
					}()
				case <-cancel:
					return
				}
			case <-cancel:
				return
			}
		}
	}()

	return cancelFn
}

func main_bak_7() {
	work := func() {
		at := time.Now()
		time.Sleep(60 * time.Millisecond)
		fmt.Printf("%s: work done\n", at.Format("15:04:05.000"))
	}
	cancel := schedule(50*time.Millisecond, work)
	defer cancel()
	// enough for 5 ticks
	time.Sleep(260 * time.Millisecond)
}
