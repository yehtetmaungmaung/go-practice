/*
The way to successfully mitigate this is to establish a signal between the
parent goroutine and its children that allows the parent to signal
cancellation to its children. By convention, this signal is usually a
read-only channel named done. The parent goroutine passes this channel to
the child goroutine and then closes the channel when it wants to cancel
the child goroutine. Here’s an example:
*/

package main

import (
	"fmt"
	"time"
)

func main() {
	done := make(chan interface{})
	terminated := doWork(done, nil)

	// Here we create another goroutine that will cancel the goroutine
	// spawned in doWork if more than one second passes.
	go func() {
		// Cancel the operation after 1 second.
		time.Sleep(1 * time.Second)
		fmt.Println("Cancelling doWork gorouting...")
		close(done)
	}()

	// This is where we join the goroutine spawned from doWork with the main
	// goroutine.
	<-terminated
	fmt.Println("Done.")
}

// Here we pass the done channel to the doWork function. As a convention,
// this channel is the first parameter.
func doWork(done <-chan interface{}, strings <-chan string) <-chan interface{} {
	terminated := make(chan interface{})
	go func() {
		defer fmt.Println("doWork exited.")
		defer close(terminated)
		// On this line we see the ubiquitous for-select pattern in use.
		// One of our case statements is checking whether our done channel
		// has been signaled. If it has, we return from the goroutine.
		for {
			select {
			case s := <-strings:
				// Do something interesting
				fmt.Println(s)
			case <-done:
				return
			}
		}
	}()
	return terminated
}

/*
You can see that despite passing in nil for our strings channel, our goroutine
still exits successfully. Unlike the example before it, in this example we do
join the two goroutines, and yet do not receive a deadlock. This is because
before we join the two goroutines, we create a third goroutine to cancel the
goroutine within doWork after a second. We have successfully eliminated our
goroutine leak!
*/
