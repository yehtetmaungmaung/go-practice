package main

import (
	"fmt"
	"net/http"
)

// Result encompasses both the *http.Response and the error possible from
// an iteration of the loop within our goroutine
type Result struct {
	Error    error
	Response *http.Response
}

func main() {
	done := make(chan interface{})
	defer close(done)

	urls := []string{
		"https://www.google.com",
		"https://badhost",
	}

	for result := range checkStatus(done, urls...) {
		// Here, in our main goroutine, we are able to deal with errors
		// coming out of the goroutine started by checkStatus intelligently,
		// and with the full context of the larger program
		if result.Error != nil {
			fmt.Printf("error: %v", result.Error)
			continue
		}
		fmt.Printf("Response: %v\n", result.Response.Status)
	}
}

func checkStatus(done <-chan interface{}, urls ...string) <-chan Result {
	results := make(chan Result)
	go func() {
		defer close(results)

		for _, url := range urls {
			var result Result
			resp, err := http.Get(url)
			result = Result{
				Error:    err,
				Response: resp,
			}

			select {
			case <-done:
				return
			case results <- result:
			}
		}
	}()
	return results
}
