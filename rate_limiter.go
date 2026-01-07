package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// We want 5 req/sec = 1 req every 200ms
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	requests := generator(ctx)
	// TODO: Implement this function below
	throttled := limiter(ctx, requests, 200*time.Millisecond)

	executor(throttled)
}

// 1. Generator: Pushes data as fast as possible
func generator(ctx context.Context) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 1; i <= 10; i++ {
			select {
			case out <- i:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

// 2. Limiter: The Gatekeeper
// YOUR TASK: Implement this function
func limiter(ctx context.Context, in <-chan int, interval time.Duration) <-chan int {
	out := make(chan int)
	ticker := time.NewTicker(interval)
	tokens := make(chan int, 3)
	tokens <- 0
	tokens <- 0
	tokens <- 0
	// add 1 token each interval
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				select {
				case <-ctx.Done():
					return
				case tokens <- 0:
				}
			}
		}
	}()
	go func() {
		defer close(out)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-tokens: // we have tokens left in our bucket
				select {
				case <-ctx.Done():
					return
				case item, ok := <-in:
					if !ok {
						return
					}
					select {
					case <-ctx.Done():
						return
					case out <- item:
					}
				}
			}
		}
	}()
	return out
}

// 3. Executor: Consumes the data
func executor(in <-chan int) {
	start := time.Now()
	for req := range in {
		fmt.Printf("Processed req %d at %v\n", req, time.Since(start))
	}
}
