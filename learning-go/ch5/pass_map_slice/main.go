package main

import "fmt"

func modMap(m map[int]string) {
	m[2] = "hello"
	m[3] = "goodbye"
	delete(m, 1)
}

func modSlice(s []int) {
	for k, v := range s {
		s[k] = v * 2
	}
	// this value of s is never used
	// You can modify any element in the slice, but you can’t lengthen
	// the slice. This is true for maps and slices that are passed directly
	// into functions as well as map and slice fields in structs.
	s = append(s, 10)
}

func main() {
	m := map[int]string{
		1: "first",
		2: "second",
	}
	modMap(m)
	fmt.Println(m)

	s := []int{1, 2, 3}
	modSlice(s)
	fmt.Println(s)
}
