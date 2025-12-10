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

// gather runs all passed functions concurrently
// and returns the results when they are ready.
func gather(funcs []func() any) []any {
	res := make([]any, len(funcs))
	done := make(chan struct{})
	for idx, fn := range funcs{
		go func(){
			res[idx] = fn()
			done <- struct{}{}
		}()
	}
	for range funcs {
		<-done
	}
	return res
}


func main() {
	slowpoke := func() any {
		fmt.Print("I'm so..& ")
		time.Sleep(100 * time.Millisecond)
		fmt.Println("slow")
		return "okay"
	}
	res := await(slowpoke)
	fmt.Println(res.(string))
}
