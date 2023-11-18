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