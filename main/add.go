package main

import "fmt"

// go build -buildmode=plugin -o add.so add.go
func Add(a, b int) int {
	fmt.Println(a + b)
	return a + b
}
