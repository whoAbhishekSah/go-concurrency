1. channel is not a queue. There must be someone receiving from the channel before you can send to the channel, otherwise you will get deadlock error.

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

To fix this deadlock error, define a receiver of the channel.

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
	fmt.Println("received ", <-x)
	wg.Wait()
}
```

**Output:**
```
pushing 1
pushed 1
received  1
```
