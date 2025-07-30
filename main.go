// Counting digits in words.
package main

import (
	"fmt"
	"time"
)

// After sending the message to the channel ➊, 
// goroutine B gets blocked. 
// Only when goroutine A receives the message ➌ does 
// goroutine B continue and print "message sent" ➋.
func main() {
    messages := make(chan string)

    go func() {
        fmt.Println("B: Sending message...")
        messages <- "ping"                    // (1)
        fmt.Println("B: Message sent!")       // (2)
    }()

    fmt.Println("A: Doing some work...")
    time.Sleep(500 * time.Millisecond)
    fmt.Println("A: Ready to receive a message...")

    <-messages                               //  (3)

    fmt.Println("A: Messege received!")
    time.Sleep(100 * time.Millisecond)
}
