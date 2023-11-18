# Chapter 12. The context

Servers need a way to handle metadata on individual requests. This metadata falls into two general categories: metadata that is required to correctly process the request, and metadata on when to stop processing the request. For example, an HTTP server might want to use a tracking ID to identify a chain of requests through a set of microservices. It also might want to set a timer that ends requests to other microservices if they take too long. Many languages use threadlocal variables to store this kind of information, associating data to a specific operating system thread of execution. This does’t work in Go because goroutines don’t have unique identities that can be used to look up values. More importantly, threadlocals feel like magic; values go in one place and pop up somewhere else.

Go solves the request metadata problem with a construct called the context. Let’s see how to use it correctly.

## What Is the context?

Rather than add a new feature to the language, a context is simply an instance that meets the context interface defined in the context package. As you know, idiomatic Go encourages explicit data passing via function parameters. The same is true for the context. It is just another parameter to your function. Just like Go has a convention that the last return value from a function is an error, there is another Go convention that the context is explicitly passed through your program as the first parameter of a function. The usual name for the context parameter is `ctx``:

```go
func logic(ctx context.Context, info string) (string, error) {
	// do sime interesting stuff here
	return "", nil
}
```

In addition to defining the Context interface, the context package also contains several factory functions for creating and wrapping contexts. When you don’t have an existing context, such as at the entry point to a command-line program, create an empty initial context with the function context.Background. This returns a variable of type context.Context. (Yes, this is an exception to the usual pattern of returning a concrete type from a function call.)

An empty context is a starting point; each time you add metadata to the context, you do so by wrapping the existing context using one of the factory functions in the context package:

```go
ctx := context.Background()
result, err := logic(ctx, "a string")
```
<hr>

***`Note`***:

There is another function, context.TODO, that also creates an empty context.Context. It is intended for temporary use during development. If you aren’t sure where the context is going to come from or how it’s going to be used, use context.TODO to put a placeholder in your code. Production code shouldn’t include context.TODO.
<hr>

When writing an HTTP server, you use a slightly different pattern for acquiring and passing the context through layers of middleware to the top-level http.Handler. Unfortunately, context was added to the Go APIs long after the net/http package was created. Due to the compatibility promise, there was no way to change the http.Handler interface to add a context.Context parameter.

The compatibility promise does allow new methods to be added to existing types, and that’s what the Go team did. There are two context-related methods on http.Request:

- Context returns the context.Context associated with the request.

- WithContext takes in a context.Context and returns a new http.Request with the old request’s state combined with the supplied context.Context.

Here's the general patern:

```go
func Middleware(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
        ctx := req.Context()
        // wrap the context with stuff -- we'll see how soon!
        req = req.WithContext(ctx)
        handler.ServeHTTP(rw, req)
    })
}
```
The first thing we do in our middleware is extract the existing context from the request using the Context method. After we put values into the context, we create a new request based on the old request and the now-populated context using the WithContext method. Finally, we call the handler and pass it our new request and the existing http.ResponseWriter.

When you get to the handler, you extract the context from the request using the Context method and call your business logic with the context as the first parameter, just like we saw previously:

```go
func handler(rw http.ResponseWriter, req *http.Request) {
    ctx := req.Context()
    err := req.ParseForm()
    if err != nil {
        rw.WriteHeader(http.StatusInternalServerError)
        rw.Write([]byte(err.Error()))
        return
    }
    data := req.FormValue("data")
    result, err := logic(ctx, data)
    if err != nil {
        rw.WriteHeader(http.StatusInternalServerError)
        rw.Write([]byte(err.Error()))
        return
    }
    rw.Write([]byte(result))
}
```

There’s one more situation where you use the WithContext method: when making an HTTP call from your application to another HTTP service. Just like we did when passing a context through middleware, you set the context on the outgoing request using WithContext:

```go
type ServiceCaller struct {
    client *http.Client
}

func (sc ServiceCaller) callAnotherService(ctx context.Context, data string)
                                          (string, error) {
    req, err := http.NewRequest(http.MethodGet,
                "http://example.com?data="+data, nil)
    if err != nil {
        return "", err
    }
    req = req.WithContext(ctx)
    resp, err := sc.client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("Unexpected status code %d",
                              resp.StatusCode)
    }
    // do the rest of the stuff to process the response
    id, err := processResponse(resp.Body)
    return id, err
}
```

Now that we know how to acquire and pass a context, let’s start making them useful. We’ll begin with cancellation.

## Cancellation

Imagine that you have a request that spawns several goroutines, each one calling a different HTTP service. If one service returns an error that prevents you from returning a valid result, there is no point in continuing to process the other goroutines. In Go, this is called cancellation and the context provides the mechanism for implementation.

To create a cancellable context, use the context.WithCancel function. It takes in a context.Context as a parameter and returns a context.Context and a context.CancelFunc. The returned context.Context is not the same context that was passed into the function. Instead, it is a child context that wraps the passed-in parent context.Context. A context.CancelFunc is a function that cancels the context, telling all of the code that’s listening for potential cancellation that it’s time to stop processing.

<hr />

***`Note`***

We’ll see this wrapping pattern several times. A context is treated as an immutable instance. Whenever we add information to a context, we do so by wrapping an existing parent context with a child context. This allows us to use contexts to pass information into deeper layers of the code. The context is never used to pass information out of deeper layers to higher layers.
<hr />

Let’s take a look at how it works. Because this code sets up a server, you can’t run it on The Go Playground, but you can [download](https://github.com/learning-go-book/context_cancel) it. First we’ll set up two servers in a file called servers.go:

```go
func slowServer() *httptest.Server {
    s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
         r *http.Request) {
        time.Sleep(2 * time.Second)
        w.Write([]byte("Slow response"))
    }))
    return s
}

func fastServer() *httptest.Server {
    s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
         r *http.Request) {
        if r.URL.Query().Get("error") == "true" {
            w.Write([]byte("error"))
            return
        }
        w.Write([]byte("ok"))
    }))
    return s
}
```

These functions launch servers when they are called. One server sleeps for two seconds and then returns the message Slow response. The other checks to see if there is a query parameter error set to true. If there is, it returns the message error. Otherwise, it returns the message ok.

***`NOTE`***

We are using the httptest.Server, which makes it easier to write unit tests for code that talks to remote servers. It’s useful here since both the client and the server are within the same program. We’ll learn more about httptest.Server in ["httptest"](https://learning.oreilly.com/library/view/learning-go/9781492077206/ch13.html#httptest).

Next, we’re going to write the client portion of the code in a file called client.go:

```go
var client = http.Client{}

func callBoth(ctx context.Context, errVal string, slowURL string,
              fastURL string) {
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()
    var wg sync.WaitGroup
    wg.Add(2)
    go func() {
        defer wg.Done()
        err := callServer(ctx, "slow", slowURL)
        if err != nil {
            cancel()
        }
    }()
    go func() {
        defer wg.Done()
        err := callServer(ctx, "fast", fastURL+"?error="+errVal)
        if err != nil {
            cancel()
        }
    }()
    wg.Wait()
    fmt.Println("done with both")
}

func callServer(ctx context.Context, label string, url string) error {
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        fmt.Println(label, "request err:", err)
        return err
    }
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println(label, "response err:", err)
        return err
    }
    data, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println(label, "read err:", err)
        return err
    }
    result := string(data)
    if result != "" {
        fmt.Println(label, "result:", result)
    }
    if result == "error" {
        fmt.Println("cancelling from", label)
        return errors.New("error happened")
    }
    return nil
}
```

All of the interesting stuff is in this file. First, our callBoth function creates a cancellable context and a cancellation function from the passed-in context. By convention, this function variable is named cancel. It is important to remember that any time you create a cancellable context, you must call the cancel function. It is fine to call it more than once; every invocation after the first is ignored. We use a defer to make sure that it is eventually called. Next, we set up two goroutines and pass the cancellable context, a label, and the URL to callServer, and wait for them both to complete. If either call to callServer returns an error, we call the cancel function.

The callServer function is a simple client. We create our requests with the cancellable context and make a call. If an error happens, or if we get the string error returned, we return the error.

Finally, we have the main function, which kicks off the program, in the file main.go:

```go
func main() {
    ss := slowServer()
    defer ss.Close()
    fs := fastServer()
    defer fs.Close()

    ctx := context.Background()
    callBoth(ctx, os.Args[1], ss.URL, fs.URL)
}
```
In main, we start the servers, create a context, and then call the clients with the context, the first argument to our program, and the URLs for our servers.

Here’s what happens if you run without an error:

```vim
$ make run-ok
go build
./context_cancel false
fast result: ok
slow result: Slow response
done with both
```

And here’s what happens if an error is triggered:

```vim
$ make run-cancel
go build
./context_cancel true
fast result: error
cancelling from fast
slow response err: Get "http://127.0.0.1:38804": context canceled
done with both
```
<hr />

***`NOTE`***

Any time you create a context that has an associated cancel function, you must call that cancel function when you are done processing, whether or not your processing ends in an error. If you do not, your program will leak resources (memory and goroutines) and eventually slow down or crash. There is no error if you call the cancel function more than once; any invocation after the first does nothing. The easiest way to make sure you call the cancel function is to use defer to invoke it right after the cancel function is returned.
<hr />
While manual cancellation is useful, it’s not your only option. In the next section, we’ll see how to automate cancellation with timeouts.
<hr />

### `Handling Server shutdown`