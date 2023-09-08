package main

import (
	"context"
	"log"
	"net/http"
	"runtime"
	"time"
)

type result struct {
	url     string
	err     error
	latency time.Duration
}

func get(ctx context.Context, url string, ch chan<- result) {
	var r result

	start := time.Now()
	ticker := time.NewTicker(1 * time.Second).C
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if resp, err := http.DefaultClient.Do(req); err != nil {
		r = result{url, err, 0}
	} else {
		t := time.Since(start).Round(time.Millisecond)
		r = result{url, nil, t}
		resp.Body.Close()
	}

	for {
		select {
		case ch <- r:
			return
		case <-ticker:
			log.Println("tick", r)
		}
	}
}

func first(ctx context.Context, urls []string) (*result, error) {
	results := make(chan result, len(urls))
	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

	for _, url := range urls {
		go get(ctx, url, results)
	}

	select {
	case r := <-results:
		return &r, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func main() {
	list := []string{
		"https://www.google.com",
		"https://www.amazon.com",
		"https://nytimes.com",
		"https://wsj.com",
		"https://facebook.com", // This will block forever since Facebook is blocked in Myanmar.
	}

	r, _ := first(context.Background(), list)

	if r.err != nil {
		log.Printf("%-30s %s", r.url, r.err.Error())
	} else {
		log.Printf("%-30s %s", r.url, r.latency.String())
	}

	time.Sleep(9 * time.Second)
	log.Println("quit anyway...", runtime.NumGoroutine(), "still running")
}
