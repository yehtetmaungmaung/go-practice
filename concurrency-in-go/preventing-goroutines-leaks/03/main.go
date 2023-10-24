/*
If a goroutine is responsible for creating a goroutine, it is also
responsible for ensuring it can stop the goroutine.

The previous example handles the case for goroutines receiving on a
channel nicely, but what if we’re dealing with the reverse situation:
a goroutine blocked on attempt ing to write a value to a channel?
Here’s a quick example to demonstrate the issue:
*/
package main

import (
	"fmt"
	"math/rand"
	"time"
)

//func main() {
//	randStream := newRandStream()
//	fmt.Println("3 random ints:")
//	for i := 1; i <= 3; i++ {
//		fmt.Printf("%d: %d\n", i, <-randStream)
//	}
//}

// func newRandStream() <-chan int {
// 	randStream := make(chan int)
// 	go func() {
// 		// Here we print out a message when the goroutine successfully
// 		// terminates. Will this run?
// 		defer fmt.Println("newRandStream closure exited.")
// 		defer close(randStream)
// 		for {
// 			randStream <- rand.Int()
// 		}
// 	}()
// 	return randStream
// }

/*
You can see from the output that the deferred fmt.Println statement never
gets run. After the third iteration of our loop, our goroutine blocks
trying to send the next random integer to a channel that is no longer being
read from. We have no way of telling the producer it can stop. The solution,
just like for the receiving case, is to provide the producer goroutine
with a channel informing it to exit:
*/

func main() {
	done := make(chan interface{})
	randStream := newRandStream(done)
	fmt.Println("3 random ints:")
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}
	close(done)

	// Simulate ongoing work
	time.Sleep(1 * time.Second)
}

func newRandStream(done <-chan interface{}) <-chan int {
	randStream := make(chan int)
	go func() {
		defer fmt.Println("newRandStream closure exited.")
		defer close(randStream)
		for {
			select {
			case randStream <- rand.Int():
			case <-done:
				return
			}
		}
	}()
	return randStream
}
