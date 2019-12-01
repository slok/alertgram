package main

import "fmt"

var (
	// Version will be populated in compilation time.
	Version = "dev"
)

func main() {
	fmt.Printf("Hello world %s", Version)
}
