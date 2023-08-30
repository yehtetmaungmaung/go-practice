package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func notify(services ...string) {

	// Step-1: Declare the sync.WaitGroup
	var wg sync.WaitGroup

	for _, service := range services {

		// Step-2: Add to the WaitGroup queue
		wg.Add(1)
		go func(s string) {
			fmt.Printf("Starting to notifying %s...\n", s)
			time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
			fmt.Printf("Finished notifying %s...\n", s)

			// Step-4: Inside each goroutine, mark items in the queue as done
			wg.Done()
		}(service)
	}

	// Step-3: Tell our code to wait on the WaitGroup queue to reach zero before proceeding
	wg.Wait()
	fmt.Println("All services notified!")
}

func main() {
	notify("Service-1", "Service-2", "Service-3", "Service-4")
}
