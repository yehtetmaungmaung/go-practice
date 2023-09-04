package main

import (
	"log"
	"net/http"
	"time"
)

type result struct {
	url     string
	err     error
	latency time.Duration
}

func get(url string, ch chan<- result) {
	start := time.Now()

	if resp, err := http.Get(url); err != nil {
		ch <- result{url, err, 0}
	} else {
		t := time.Since(start).Round(time.Millisecond)
		ch <- result{url, nil, t}
		resp.Body.Close()
	}
}

func main() {
	results := make(chan result)
	list := []string{
		"https://www.google.com",
		"https://amazon.com",
		"https://nytimes.com",
		"https://wsj.com",
	}

	for _, url := range list {
		go get(url, results)
	}

	for range list {
		r := <-results

		if r.err != nil {
			log.Printf("%-30s %s\n", r.url, r.err)
		} else {
			log.Printf("%-30s %s\n", r.url, r.latency)
		}
	}
}
