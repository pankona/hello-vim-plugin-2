package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Print("Enter your message: ")
	var input string
	_, err := fmt.Scanf("%s\n", &input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Your message: %s\n", input)
}