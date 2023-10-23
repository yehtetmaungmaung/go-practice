// As we covered in the section “Goroutines” on page 37, we know goroutines are
// cheap and easy to create; it’s one of the things that makes Go such a
// productive language. The runtime handles multiplexing the goroutines onto
// any number of operating system threads so that we don’t often have to worry
// about that level of abstraction. But they do cost resources, and goroutines
// are not garbage collected  by the runtime, so regardless of how small their
// memory footprint is, we don’t want to leave them lying about our process.
// So how do we go about ensuring they’re cleaned up? Let’s start from the
// beginning and think about this step by step:
//
// why would a goroutine exist? In Chapter 2, we established that goroutines
// rep resent units of work that may or may not run in parallel with each
// other. The go routine has a few paths to termination:
//
//   - When it has completed its work.
//   - When it cannot continue its work due to an unrecoverable error.
//   - When it’s told to stop working.
//
// We get the first two paths for free—these paths are your algorithm—but
// what about work cancellation? This turns out to be the most important bit
// because of the net work effect: if you’ve begun a goroutine, it’s most
// likely cooperating with several other goroutines in some sort of organized
// fashion. We could even represent this interconnectedness as a graph:
// whether or not a child goroutine should continue executing might be
// predicated on knowledge of the st ate of many other goroutines. The parent
// goroutine (often the main goroutine) with this full contextual knowledge
// should be able to tell its child goroutines to terminate. We’ll continue
// looking at large-scale goroutine interdependence in the next chapter, but
// for now let’s consider how to ensure a single child goroutine is guaranteed
// to be cleaned up. Let’s start with a simple example of a goroutine leak:
package main

import "fmt"

func main() {
	doWork := func(strings <-chan string) <-chan interface{} {
		completed := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited.")
			defer close(completed)
			for s := range strings {
				// Do something interesting
				fmt.Println(s)
			}
		}()
		return completed
	}

	// Here we see that the main goroutine passes a nil channel into doWork.
	// Therefore, the strings channel will never actually gets any strings
	// written onto it, and the goroutine containing doWork will remain in
	// memory for the lifetime of this process (we would even deadlock if we
	// joined the goroutine within doWork and the main goroutine)
	doWork(nil)
}
