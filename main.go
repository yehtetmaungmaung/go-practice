// Prime Sieve implementation with go channels
package main

import "fmt"

func generator(limit int, ch chan<- int) {
	for i := 2; i < limit; i++ {
		ch <- i
	}
	close(ch)
}

func filter(src <-chan int, dest chan<- int, prime int) {
	for i := range src {
		if i%prime != 0 {
			dest <- i
		}
	}
	close(dest)
}

func sieve(limit int) {
	ch := make(chan int)
	go generator(limit, ch)
	for {
		prime, ok := <-ch
		if !ok {
			break
		}

		ch1 := make(chan int)
		go filter(ch, ch1, prime)
		ch = ch1
		fmt.Print(prime, " ")
	}
	fmt.Println()
}

func main() {
	sieve(100)
}
