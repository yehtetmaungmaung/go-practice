package main

import (
	"fmt"
	"log"
	"strconv"
)

func add(i, j int) int      { return i + j }
func sub(i, j int) int      { return i - j }
func multiply(i, j int) int { return i * j }
func div(i, j int) int      { return i / j }

var opsMap = map[string]func(int, int) int{
	"+": add,
	"-": sub,
	"*": multiply,
	"/": div,
}

func main() {
	expression := [][]string{
		{"2", "+", "3"},
		{"2", "-", "3"},
		{"2", "*", "3"},
		{"2", "/", "3"},
		{"2", "%", "3"},
		{"two", "+", "three"},
		{"5"},
	}

	for _, expr := range expression {
		if len(expr) < 3 {
			fmt.Println("invalid expression: ", expression)
			continue
		}
		p1, err := strconv.Atoi(expr[0])
		if err != nil {
			fmt.Println(err)
		}
		op := expr[1]
		opFunc, ok := opsMap[op]
		if !ok {
			fmt.Println("unspported operation: ", op)
		}

		p2, err := strconv.Atoi(expr[2])
		if err != nil {
			fmt.Println(err)
		}

		result := opFunc(p1, p2)
		log.Println(result)
	}
}
