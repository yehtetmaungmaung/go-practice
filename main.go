package main

import (
	"fmt"
	"math/rand"
	"time"
)

func notify(services ...string) {
	for _, service := range services {
		go func(s string) {
			fmt.Printf("Starting to notifing %s...\n", s)
			time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
			fmt.Printf("Finished notifying %s...\n", s)
		}(service)
	}
	fmt.Println("All services notified!")
}

func main() {
	notify("Service-1", "Service-2", "Service-3")
	// Running this outputs "All services notified!" but we
	// won't see any of the services outputting their finished messages!
}
