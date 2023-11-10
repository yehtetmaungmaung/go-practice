/*
Confinement is the simple yet powerful idea of ensuring information is only ever
available from one concurrent process. When this is achieved, a concurrent program
is implicitly safe and no synchronization is needed. There are two kinds of confine‐
ment possible: ad hoc and lexical.

Ad hoc confinement is when you achieve confinement through a convention—
whether it be set by the languages community, the group you work within, or the
codebase you work within. In my opinion, sticking to convention is difficult to ach‐
ieve on projects of any size unless you have tools to perform static analysis on your
code every time someone commits some code. Here’s an example of ad hoc confine‐
ment that demonstrates why:
*/
package main

import "fmt"

func main() {
	data := make([]int, 4)

	loopData := func(handleData chan<- int) {
		defer close(handleData)
		for i := range data {
			handleData <- data[i]
		}
	}

	handleData := make(chan int)
	go loopData(handleData)

	for num := range handleData {
		fmt.Println(num)
	}
}

/*
We can see that the data slice of integers is available from both the loopData function
and the loop over the handleData channel; however, by convention we’re only access‐
ing it from the loopData function. But as the code is touched by many people, and
deadlines loom, mistakes might be made, and the confinement might break down
and cause issues. As I mentioned, a static-analysis tool might catch these kinds of
issues, but static analysis on a Go codebase suggests a level of maturity that not many
teams achieve. This is why I prefer lexical confinement: it wields the compiler to
enforce the confinement.
*/
