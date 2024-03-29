package main

import (
	"fmt"
)

func main() {
	chanOwner := func() <-chan int {
		resultStream := make(chan int, 5) // <1>
		go func() {                       // <2>
			defer close(resultStream) // <3>
			for i := 0; i <= 5; i++ {
				resultStream <- i
			}
		}()
		return resultStream // <4>
	}

	// resultStream := chanOwner()
	// for result := range resultStream { // <5>
	// 	fmt.Printf("Received: %d\n", result)
	// }

	const chanSize int = 5
	//var resultStream [chanSize]<-chan int
	var resultStream <-chan int
	for i := 0; i < chanSize; i++ {
		resultStream = chanOwner()
		for result := range resultStream { // <5>
			fmt.Printf("Received%d: %d\n", i+1, result)
		}
	}

	fmt.Println("Done receiving!")
}
