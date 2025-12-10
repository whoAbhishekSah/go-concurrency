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
