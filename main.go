package main

import (
	"fmt"
	"log"
	"net/http"
)

// Now: let's implement object oriented version of the previous version.

type nextCh chan int

// Every request will read from nextID channel, then the counter go
// routine will be able to send
func (ch nextCh) handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>You got %d!</h1>", <-ch)
}

// The function will run forever, but nextID will be blocked until
// handler reads from the nextID channel
func counter(ch chan<- int) {
	for i := 0; ; i++ {
		ch <- i
	}
}

func main() {
	var nextID nextCh = make(chan int)
	go counter(nextID)
	http.HandleFunc("/", nextID.handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
