# Error Handling

In concurrent programs, error handling can be difficult to get right. Sometimes, we
spend so much time thinking about how our various processes will be sharing infor‐
mation and coordinating, we forget to consider how they’ll gracefully handle errored
states. When Go eschewed the popular exception model of errors, it made a statement that error handling was important, and that as we develop our programs, we
should give our error paths the same attention we give our algorithms. In that spirit,
let’s take a look at how we do that when working with multiple concurrent processes.<br/><br/> 
The most fundamental question when thinking about error handling is, “Who should
be responsible for handling the error?” At some point, the program needs to stop ferrying the error up the stack and actually do something with it. What is responsible
for this?<br/><br/>
With concurrent processes, this question becomes a little more complex. Because a
concurrent process is operating independently of its parent or siblings, it can be difficult for it to reason about what the right thing to do with the error is. Take a look at
the following code for an example of this issue:

```go
package main

import (
	"fmt"
	"net/http"
)

func main() {
	done := make(chan interface{})
	defer close(done)

	urls := []string{"https://www.google.com", "https://badhost"}
	for response := range checkStatus(done, urls...) {
		fmt.Printf("Response: %v\n", response.Status)
	}
}

func checkStatus(done <-chan interface{}, urls ...string) <-chan *http.Response {
	responses := make(chan *http.Response)
	go func() {
		defer close(responses)
		for _, url := range urls {
			resp, err := http.Get(url)
			if err != nil {
				// Here we see the goroutine doing its best to signal that
				// there's an error. What else can it do? It can't pass it
				// back! How many errors is too many? Does it continue making
				// request?
				fmt.Println(err)
				continue
			}
			select {
			case <-done:
				return
			case responses <- resp:
			}
		}
	}()
	return responses
}
```
<br/>
Here we see that the goroutine has been given no choice in the matter. It can’t simply
swallow the error, and so it does the only sensible thing: it prints the error and hopes
something is paying attention. Don’t put your goroutines in this awkward position. I
suggest you separate your concerns: in general, your concurrent processes should
send their errors to another part of your program that has complete information
about the state of your program, and can make a more informed decision about what
to do.

```
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
```

The key thing to note here is how we’ve coupled the potential result with the potential
error. This represents the complete set of possible outcomes created from the goroutine checkStatus, and allows our main goroutine to make decisions about what to do
when errors occur. In broader terms, we’ve successfully separated the concerns of
error handling from our producer goroutine. This is desirable because the goroutine
that spawned the producer goroutine—in this case our main goroutine—has more
context about the running program, and can make more intelligent decisions about
what to do with errors.
<br/><br/>
In the previous example, we simply wrote errors out to stdio, but we could do something else. Let’s alter our program slightly so that it stops trying to check for status if
three or more errors occur:
```go
	done := make(chan interface{})
	defer close(done)

	errCount := 0
	urls := []string{"a", "https://www.google.com", "b", "c", "d"}

	for result := range checkStatus(done, urls...) {
		if result.Error != nil {
			fmt.Printf("Error: %v\n", result.Error)
			errCount++
			if errCount >= 3 {
				fmt.Println("Too many errors, breaking!")
				break
			}
			continue
		}
		fmt.Printf("Response: %v\n", result.Response.Status)
	}
```
You can see that because errors are returned from checkStatus and not handled
internally within the goroutine, error handling follows the familiar Go pattern. This
is a simple example, but it’s not hard to imagine situations where the main goroutine
is coordinating results from multiple goroutines and building up more complex rules
for continuing or canceling child goroutines. Again, the main takeaway here is that
errors should be considered first-class citizens when constructing values to return
from goroutines. If your goroutine can produce errors, those errors should be tightly
coupled with your result type, and passed along through the same lines of communication—just like regular synchronous functions.
