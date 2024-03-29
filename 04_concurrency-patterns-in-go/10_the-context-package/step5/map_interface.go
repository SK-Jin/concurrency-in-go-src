package main

import (
	"fmt"
)

type foo int
type bar int

func main() {
	m := make(map[interface{}]int)
	m[foo(1)] = 1
	m[bar(1)] = 2

	fmt.Printf("%v", m)
}
