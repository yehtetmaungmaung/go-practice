package main

import (
	"bytes"
	"fmt"
	"sync"
)

// So what’s the point? Why pursue confinement if we have synchronization
// available to us? The answer is improved performance and reduced cognitive
// load on developers. Synchronization comes with a cost, and if you can avoid
// it you won’t have any critical sections, and therefore you won’t have to
// pay the cost of synchronizing them. You also sidestep an entire class of
// issues possible with synchronization; developers simply don’t have to worry
// about these issues. Concurrent code that utilizes lexical confinement also
// has the benefit of usually being simpler to understand than concurrent
// code without lexically confined variables. This is because within the
// context of your lexical scope you can write synchronous code.
//
// Having said that, it can be difficult to establish confinement, and so
// sometimes we have to fall back to our wonderful Go concurrency primitives.
func main() {
	// Because printData doesn’t close around the data slice, it cannot access
	// it, and needs to take in a slice of byte to operate on. Because of the
	// lexical scope, we’ve made it impossible1 to do the wrong thing, and so
	// we don’t need to synchronize memory access or share data through
	// communication.(Ignoring the possibility of manually manipulating memory
	// via the unsafe package. It’s called unsafe for a reason!)
	printData := func(wg *sync.WaitGroup, data []byte) {
		defer wg.Done()

		var buff bytes.Buffer
		for _, b := range data {
			fmt.Fprintf(&buff, "%c", b)
		}
		fmt.Println(buff.String())
	}

	var wg sync.WaitGroup
	wg.Add(2)

	data := []byte("golang")
	go printData(&wg, data[:3])
	go printData(&wg, data[3:])

	wg.Wait()
}
