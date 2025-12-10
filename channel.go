package main

import (
	"fmt"
	"strings"
)

func main_bak_1() {
	str := "one,two,,four"
	stream := make(chan string)
	go submit(str, stream)
	print(stream)
}

func submit(str string, stream chan<- string) {
	words := strings.Split(str, ",")
	for _, word := range words {
		stream <- word
	}
	close(stream)
}

func print(stream <-chan string) {
	for word := range stream {
		if word != "" {
			fmt.Printf("%s ", word)
		}
	}
	fmt.Println()
}
