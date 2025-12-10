package main

import (
	"fmt"
	"strings"
)

func encode(str string) string {
	pending := submitter(str)
	encoded := encoder(pending)
	words := receiver(encoded)
	return strings.Join(words, " ")
}

func submitter(str string) <-chan string {
	pending := make(chan string)
	go func() {
		words := strings.Fields(str)
		for _, item := range words {
			pending <- item
		}
		close(pending)
	}()
	return pending
}

func encoder(pending <-chan string) chan string {
	encoded := make(chan string)
	go func() {
		for word := range pending {
			encoded <- fmt.Sprintf("encoded_%s", word)
		}
		close(encoded)
	}()
	return encoded
}

func receiver(in <-chan string) []string {
	res := make([]string, 0)
	for word := range in {
		res = append(res, word)
	}
	return res
}

func main() {
	src := "go is awesome"
	res := encode(src)
	fmt.Println(res)
}
