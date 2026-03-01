package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrFailed = errors.New("Failed")
	ErrManual = errors.New("Stopped Manually")
)

type Worker struct{
	fn func() error
	ctx context.Context
	cancel func(cause error)
	afterFuncs []func()
}

func (w *Worker) Start() {
	if w.ctx !=nil {
		return 
	}
	w.ctx, w.cancel = context.WithCancelCause(context.Background())
	for _, fn := range w.afterFuncs {
		context.AfterFunc(w.ctx, fn)
	}
	go w.work()
}

func (w *Worker) AfterStop(fn func()) {
	if w.ctx != nil {
		return
	}
	w.afterFuncs = append(w.afterFuncs, fn)
}

func NewWorker(fn func () error) *Worker{
	return &Worker{
		fn: fn,
	}
}

func (w * Worker) work(){
	for{
		select{
		case <- w.ctx.Done():
			return
		default:
			err := w.fn()
			if err!=nil {
				w.cancel(ErrFailed)
				return
			}
		}
	}
}

func (w *Worker) Stop() {
	if w.ctx == nil {
		return
	}
	w.cancel(ErrManual)
}

func (w *Worker) Err() error {
	return context.Cause(w.ctx)
}

func main() {
	count := 9
	fn:= func () error  {
		fmt.Print(count, " ")
		count--
		time.Sleep(10*time.Millisecond)
		return nil
	}

	worker := NewWorker(fn)
	worker.AfterStop(func ()  {
		fmt.Println("called after stop")
	})
	worker.Start()
	time.Sleep(2 * time.Millisecond)
	worker.Stop()
	time.Sleep(10 * time.Millisecond)
	fmt.Println(worker.Err())
}
