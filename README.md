### 1. Unbuffered channel

`x := make (chan int)`

channel is not a queue. There must be someone receiving from the channel before you can send to the channel, otherwise you will get deadlock error.

In general, when we say channel, we mean unbuffered channel. Length 0.

**Example:**

```go
func a(wg *sync.WaitGroup, x chan int){
	fmt.Println("pushing 1")
	fmt.Println("pushed 1")
	defer wg.Done()
}

func main(){
	wg := &sync.WaitGroup{}
	x := make(chan int)
	wg.Add(1)
	go a(wg, x)
	wg.Wait()
}
```

**Output:**

```
pushing 1
pushed 1

```

```go
func a(wg *sync.WaitGroup, x chan int){
	fmt.Println("pushing 1")
	x<-1
	fmt.Println("pushed 1")
	defer wg.Done()
}

func main(){
	wg := &sync.WaitGroup{}
	x := make(chan int)
	wg.Add(1)
	go a(wg, x)
	wg.Wait()
}

```

Output:

```text
pushing 1
fatal error: all goroutines are asleep - deadlock!
```

To fix this deadlock error, define a receiver of the channel or created a buffered channel.

```go
func a(wg *sync.WaitGroup, x chan int){
	fmt.Println("pushing 1")
	x<-1                            // This will block until a receiver is ready
	fmt.Println("pushed 1")
	defer wg.Done()
}

func main(){
	wg := &sync.WaitGroup{}
	x := make(chan int)  // Creates an unbuffered channel
	wg.Add(1)
	go a(wg, x)
	fmt.Println("received ", <-x)       // This will block until a sender sends a value
	wg.Wait()
}
```

**Output:**

```
pushing 1
pushed 1
received  1
```

Sends and receives are synchronous: Every send operation on an unbuffered channel will block the sending goroutine until a corresponding receive operation is ready to take the value. Conversely, a receive operation will block until a value is sent to the channel.

Example buffered channel:

```go
func a(wg *sync.WaitGroup, x chan int){
	fmt.Println("pushing 1")
	x<-1                            // This will block until a receiver is ready
	fmt.Println("pushed 1")
	defer wg.Done()
}

func main(){
	wg := &sync.WaitGroup{}
	x := make(chan int, 1)  // Creates a buffered channel
	wg.Add(1)
	go a(wg, x)
	wg.Wait()
}
```

**Output:**

```
pushing 1
pushed 1
```
