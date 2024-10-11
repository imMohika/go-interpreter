package main

import (
	"fmt"
	"go-interpreter/itop"
	"os"
)

func main() {
	fmt.Printf("Welcome to itop!\n")
	itop.Start(os.Stdin, os.Stdout)
}
