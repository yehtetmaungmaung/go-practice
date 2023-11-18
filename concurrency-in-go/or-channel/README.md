# The or-channel
At times you may find yourself wanting to combine one or more done channels into a
single done channel that closes if any of its component channels close. It is perfectly
acceptable, albeit verbose, to write a select statement that performs this coupling;
however, sometimes you can’t know the number of done channels you’re working
with at runtime. In this case, or if you just prefer a one-liner, you can combine these
channels together using the or-channel pattern.

This pattern creates a composite done channel through recursion and goroutines.

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	// or takes in a variadic slice of channels and returns a single
	// channel.
	var or func(channels ...<-chan interface{}) <-chan interface{}
	or = func(channels ...<-chan interface{}) <-chan interface{} {
		// Since this is a recursive function, we must set up
		// termination criteria. The first is that if the variadic
		// slice is empty, we simply return a nil channel. This is
		// consistent with the idea of passing in no channels; we
		// wouldn't expect a composite channel to do anything.
		switch len(channels) {
		case 0:
			return nil

		// Our second termination criteria states that if our variadic
		// slice only contains one element, we just return that element.
		case 1:
			return channels[0]
		}

		orDone := make(chan interface{})

		// Here is the main body of the function, and where the recursion
		// happens. We create a goroutine so that we can wait for messages
		// on our channels without blocking.
		go func() {
			defer close(orDone)

			switch len(channels) {
			// Because of how we're recursing, every recursive call to
			// `or` will at least have to channels. As an optimization to
			// keep the number of goroutines constrained, we place a
			// special case here for calls to `or` with only 2 channels.
			case 2:
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			// Here we recursively create an `or-channel` from all the
			// channels in our slice after the third index, and then select
			// from this. This recurrence relation will destructure the
			// rest of the slice into `or-channels` to form a tree from
			// which the first signal return. We also pass in the `orDone`
			// channel so that when goroutines up the tree exit, goroutines
			// down the tree also exit.
			default:
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
				case <-or(append(channels[3:], orDone)...):
				}
			}
		}()
		return orDone
	}

	// Here we keep track of roughly when the channel from the `or` function
	// begins to block.
	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	// And here we print the time it took for the read to occur.
	fmt.Printf("Done after %v", time.Since(start))
}

// sig simply creates a channel that will close when the time specified in
// the `after` elapses.
func sig(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}

```

Notice that despite placing several channels in our call to or that take various times to
close, our channel that closes after one second causes the entire channel created by
the call to or to close. This is because—despite its place in the tree the or function
builds—it will always close first and thus the channels that depend on its closure will
close as well.
We achieve this terseness at the cost of additional goroutines—f(x)=⌊x/2⌋ where x is
the number of goroutines—but remember that one of Go’s strengths is the ability to quickly create, schedule, and run goroutines, and the language actively encourages
using goroutines to model problems correctly. Worrying about the number of gorou‐
tines created here is probably a premature optimization. Further, if at compile time
you don’t know how many done channels you’re working with, there isn’t any other
way to combine done channels.
This pattern is useful to employ at the intersection of modules in your system. At
these intersections, you tend to have multiple conditions for canceling trees of gorou‐
tines through your call stack. Using the or function, you can simply combine these
together and pass it down the stack. We’ll take a look at another way of doing this in
“The context Package” on page 131 that is also very nice, and perhaps a bit more
descriptive.
