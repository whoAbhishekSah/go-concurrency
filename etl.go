package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	ErrMatchSubstr = "[ERROR]"
	MaxWorker      = 3
)

func aggregateLog(ctx context.Context, filepath string) int {
	pending := reader(ctx, filepath) // return a channel where each line is written line by line
	collector := transformDispatcher(ctx, pending)
	return sinker(collector) // counts the error lines
}

func transformDispatcher(ctx context.Context, pending <-chan string) <-chan string {
	collector := make(chan string)
	wg := &sync.WaitGroup{}
	for i := range MaxWorker {
		wg.Add(1)
		go transformer(ctx, i, pending, collector, wg)
	}
	go func() {
		wg.Wait()
		close(collector)
	}()
	return collector
}

// context cancellation must be checked each time we scan a new line
func reader(ctx context.Context, filepath string) <-chan string {
	pending := make(chan string)
	go func() {
		file, err := os.Open(filepath)
		if err != nil {
			fmt.Printf("Failed to open file %s, err: %v\n", filepath, err)
			return
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			select {
			case pending <- line:
			case <-ctx.Done():
				fmt.Println("Context error in reader")
				return
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Printf("scan error %v\n", err)
		}
		close(pending)
	}()
	return pending
}

// caveat: loss of order of logs
// The sinker could receive old line before new lines.
// context cancellation must be checked each time we read from channel
func transformer(ctx context.Context, id int, pending <-chan string, collector chan<- string, wg *sync.WaitGroup) {
	fmt.Printf("staring worker %d\n", id)
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Context error in transformer while reading from pending, worker id %d\n", id)
			return
		case checkStr, ok := <-pending:
			if !ok {
				return
			}
			if strings.Contains(checkStr, ErrMatchSubstr) {
				// Simulate expensive processing (e.g., resizing an image)
				time.Sleep(500 * time.Millisecond)
				fmt.Printf("Worker %d finished str processing\n", id)
				select {
				case <-ctx.Done():
					fmt.Printf("Context error in transformer while sending to collector, worker id %d\n", id)
					return
				case collector <- checkStr:
				}
			}
		}
	}
}

// The final stage (sinker) must act as the barrier.
// If we run it in a goroutine, main returns immediately and the program dies
// before processing happens.
func sinker(collector <-chan string) int {
	i := 0
	for range collector {
		i += 1
	}
	return i
}

func main_bak_5() {
	ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
	fmt.Println(aggregateLog(ctx, "./app.log"))
}
