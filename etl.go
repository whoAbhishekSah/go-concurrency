package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	ErrMatchSubstr = "[ERROR]"
)

func aggregateLog(filepath string) int {
	pending := reader(filepath)       // return a channel where each line is written line by line
	collector := transformer(pending) // returns a channel where the only the error lines are put
	return sinker(collector)          // counts the error lines
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

func transformer(pending <-chan string) <-chan string {
	collector := make(chan string)
	go func() {
		for checkStr := range pending {
			if strings.Contains(checkStr, ErrMatchSubstr) {

				collector <- checkStr
			}

		}
		close(collector)
	}()
	return collector
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
