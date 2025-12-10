package main

import (
	"fmt"
	"time"
)

func await(fn func() any) any {
	out := make(chan any)
	go func() {
		out <- fn()
	}()
	return <-out
}

type result struct {
	idx int
	val any
}

// gather runs all passed functions concurrently
// and returns the results when they are ready.
func gather(funcs []func() any) []any {
	results := make(chan result, len(funcs))
	res := make([]any, len(funcs))
	for idx, fn := range funcs {
		go func() {
			results <- result{idx, fn()}
		}()
	}
	for range funcs {
		item := <-results
		res[item.idx] = item.val
	}
	return res
}

// squared returns a function that returns
// the square of the input number.
func squared(n int) func() any {
	return func() any {
		return n * n
	}
}

func main() {
	funcs := []func() any{squared(2), squared(3), squared(4)}

	start := time.Now()
	nums := gather(funcs)
	elapsed := float64(time.Since(start)) / 1_000_000

	fmt.Println(nums)
	fmt.Printf("Took %.0f ms\n", elapsed)
}
