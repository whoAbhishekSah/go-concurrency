package main

import (
	"context"
	"fmt"
	"time"
)


func execute(ctx context.Context, fn func() int) (int, error){
	ch := make(chan int, 1)
	go func ()  {
		ch <- fn()
	}()
	select {
	case res := <-ch:
		return res, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

func work() int{
	time.Sleep(100 * time.Millisecond)
	fmt.Println("work done")
	return 42
}

func main_bak_8(){
	timeout := 42 * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	
	res, err := execute(ctx, work)
	fmt.Println(res, err)
}
