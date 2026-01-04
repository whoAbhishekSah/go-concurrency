package main

import (
	"bufio"
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

func aggregateLog(filepath string) int {
	pending := reader(filepath) // return a channel where each line is written line by line
	collector := transformDispatcher(pending)
	return sinker(collector) // counts the error lines
}

func transformDispatcher(pending <-chan string) <-chan string {
	collector := make(chan string)
	wg := &sync.WaitGroup{}
	for i := range MaxWorker {
		wg.Add(1)
		go transformer(i, pending, collector, wg)
	}
	go func() {
		wg.Wait()
		close(collector)
	}()
	return collector
}

func reader(filepath string) <-chan string {
	pending := make(chan string)
	go func() {
		file, err := os.Open(filepath)
		if err != nil {
			panic(fmt.Errorf("Failed to open file %s, err: %v", filepath, err))
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()

			pending <- line
		}
		if err := scanner.Err(); err != nil {
			panic(fmt.Errorf("scan error %v", err))
		}
		close(pending)
	}()
	return pending
}

func transformer(id int, pending <-chan string, collector chan<- string, wg *sync.WaitGroup) {
	fmt.Printf("staring worker %d\n", id)
	defer wg.Done()
	for checkStr := range pending {
		if strings.Contains(checkStr, ErrMatchSubstr) {
			// Simulate expensive processing (e.g., resizing an image)
			collector <- checkStr
			time.Sleep(500 * time.Millisecond)
			fmt.Printf("Worker %d finished str processing\n", id)
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

func main() {
	fmt.Println(aggregateLog("./app.log"))
}
