package main

import (
	"fmt"
	"strings"
)

const eof = "__EOF__"

func main(){
	str := "one,two,,four"
	in := make(chan string)
	go func(){
		words := strings.Split(str, ",")
		for _, word := range words {
			in <- word
		}
		in <- eof
	}()

	for {
		word := <- in
		if word == eof {
			break
		}
		if word != ""{
			fmt.Printf("%s ", word)
		}
	}
}
