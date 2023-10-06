package main

import (
	"fmt"
	"sync"
)

type Button struct {
	Clicked *sync.Cond
}

func main() {
	button := Button{
		Clicked: sync.NewCond(&sync.Mutex{}),
	}

	subscribe := func(c *sync.Cond, fn func()) {
		var goroutineRunning sync.WaitGroup
		goroutineRunning.Add(1)
		go func() {
			goroutineRunning.Done()
			c.L.Lock()
			defer c.L.Unlock()
			c.Wait()
			fn()
		}()
		goroutineRunning.Wait()
	}

	var clickedRegistered sync.WaitGroup
	clickedRegistered.Add(3)

	subscribe(button.Clicked, func() {
		fmt.Println("maximazing window.")
		clickedRegistered.Done()
	})

	subscribe(button.Clicked, func() {
		fmt.Println("Displaying annoying dialog box!")
		clickedRegistered.Done()
	})

	subscribe(button.Clicked, func() {
		fmt.Println("Mouse Clicked.")
		clickedRegistered.Done()
	})

	button.Clicked.Broadcast()
	clickedRegistered.Wait()
}
