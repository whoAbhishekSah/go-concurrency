package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func say(done chan<- struct{}, id int, phrase string) {
	for _, word := range strings.Fields(phrase) {
		fmt.Printf("Worker #%d says: %s..&\n", id, word)
		dur := time.Duration(rand.Intn(100)) * time.Millisecond
		time.Sleep(dur)
	}
	done <- struct{}{}
}

func main_bak_3() {
	phrases := []string{
		"go is awesome",
		"cats are cute",
		"rain is wet",
		"channels are hard",
		"floor is lava",
	}
	done := make(chan struct{})
	for idx, phrase := range phrases {
		go say(done, idx+1, phrase)
	}
	for range len(phrases) {
		<-done
	}
}
