package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func main() {
	// Open and read the log file
	file, err := os.Open("logs.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Regular expression pattern to match "DeletedCount" values
	pattern := `DeletedCount":(\d+)`
	re := regexp.MustCompile(pattern)

	// Initialize a variable to hold the sum
	sum := 0

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Read and sum the "DeletedCount" values
	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if len(matches) == 2 {
			fmt.Println("match: ", matches)
			count, err := strconv.Atoi(matches[1])
			if err == nil {
				sum += count
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	fmt.Printf("Sum of DeletedCount values: %d\n", sum)
}
