package main

import (
	"fmt"
	"log"
	"net/http"
)

var nextID = make(chan int)

// Every request will read from nextID channel, then the counter go
// routine will be able to send
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>You got %d!</h1>", <-nextID)
}

// The function will run forever, but nextID will be blocked until
// handler reads from the nextID channel
func counter() {
	for i := 0; ; i++ {
		nextID <- i
	}
}

func main() {
	go counter()
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
