# Go concurrency

Notes from :
- https://go.dev/tour/concurrency/

## Goroutines

> A goroutine is a lightweight thread managed by the Go runtime.

To start a goroutine:
```go
go f(x, y, z)
```

About goroutines:
- the evaluation of the parameters happens in the curren goroutine
- execution of the function happens in the new goroutine
- gorourtines run in the same address space
- access to shared memory must be sync
- `sync` package provides useful primitives (`sync.WaitGroup`)

## Channels

> Channels are a type of conduit through which you can send and receive values with the channel operator, <-.

```go
ch <- v         // Send v to channel ch.
v := <-ch       // Receive from ch, and assign value to v.
v, ok := <- ch  // Receive from ch and return true if there is a value, false if empty
```

- The data flows in the direction of the arrow.
- Channels must be created before use like maps and slices: 
```go
ch := make(chan int)
```

- Sends and receives block until the other side is ready
- Useful to sync goroutines without locks or condition variables

## Buffered channels

Channels can be buffered by adding the length as the second argument to make when initialising:
```go
ch := make(chan int)        // non-buffered channel
ch := make(chan int, 100)   // buffered channel
```

- Sends to a buffer channel block ONLY when the buffer is full
- Receives block when the buffer is empty

Example of code when buffer is overfill:
```go
package main

import "fmt"

func main() {
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	ch <- 2  // overfill the buffer of length 2 with a third value
	fmt.Println(<-ch)
	fmt.Println(<-ch)
}
```

Output:
```shell
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [chan send]:
main.main()
	/tmp/sandbox2595570494/prog.go:9 +0x5c
```

## Range & close

1. Sender
- A sender can `close` a channel => indicating no more values will be sent
- Should be used only to signal to the receiver there are no more values coming (g.g.: terminate a range loop)
```go
ch := make(chan int)
ch <- 2
close(ch)
```

- Only the sender should close the channel! Sending on a close channel ==> panic 😱
```shell
panic: send on closed channel

goroutine 6 [running]:
main.fibonacci(0xa, 0x0?)
	/tmp/sandbox1955149049/prog.go:15 +0x8b
created by main.main
	/tmp/sandbox1955149049/prog.go:20 +0x8a
```

2. Receiver
- A receiver can test if a channel has been closed by adding a parameter:
```go
v, ok := <- ch 
```

`ok = true`: channel is still opened with still values to be read  
`ok = false`: there are no mores values, channel is closed


- Range over a channel to receive values from the channel repeatedly until it is closed:
```go
for i := range c
```

Example of code:
```go
package main

import (
	"fmt"
)

func fibonacci(n int, c chan int) {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		c <- x
		x, y = y, x+y
	}
	close(c)
}

func main() {
	c := make(chan int, 10)
	go fibonacci(cap(c), c)
	for i := range c {
		fmt.Println(i)
	}
}
```

## Select

> The `select` statement lets a goroutine wait on multiple communication operators.

- `select` blocks until one of its cases can run
- then executes the case
- if multiple are ready, chooses one randomly

```go
package main

import "fmt"

func fibonacci(c chan int, quit <-chan struct{}) {
	x, y := 0, 1
	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}

func main() {
	c := make(chan int)
	quit := make(chan struct{})
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(<-c)
		}
		quit <- struct{}{}
	}()
	fibonacci(c, quit)
}
```

## Default select

- Default is run in a select if no other case is ready
- Use it to try a send or receive **without blocking**

```go
select {
case i := <-c:
    // use i
default:
    // receiving from c would block
}
```

Example:
```go
package main

import (
	"fmt"
	"time"
)

func main() {
	tick := time.Tick(100 * time.Millisecond)
	boom := time.After(500 * time.Millisecond)
	for {
		select {
		case <-tick:
			fmt.Println("tick.")
		case <-boom:
			fmt.Println("BOOM!")
			return
		default:
			fmt.Println("    .")
			time.Sleep(50 * time.Millisecond)
		}
	}
}
```

NB: `time.Tick` creates a channel cause Ticker holds a channel in this structure:
```go
// Tick is a convenience wrapper for NewTicker providing access to the ticking
// channel only. While Tick is useful for clients that have no need to shut down
// the Ticker, be aware that without a way to shut it down the underlying
// Ticker cannot be recovered by the garbage collector; it "leaks".
// Unlike NewTicker, Tick will return nil if d <= 0.
func Tick(d Duration) <-chan Time {
if d <= 0 {
return nil
}
return NewTicker(d).C
}

// A Ticker holds a channel that delivers “ticks” of a clock
// at intervals.
type Ticker struct {
	C <-chan Time // The channel on which the ticks are delivered.
	r runtimeTimer
}
```