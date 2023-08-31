package main

import (
	"fmt"
	"net/http"
)

var nextID int

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>You got %d!</h1>", nextID)
	// This is bad!!!
	// The next commit shows the proper way by using channel
	nextID++
}

func main() {
	http.HandleFunc("/", handler)
}
