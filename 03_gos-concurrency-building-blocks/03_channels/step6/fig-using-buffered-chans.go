package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	var stdoutBuff bytes.Buffer         // <1>
	defer stdoutBuff.WriteTo(os.Stdout) // <2>
	// test2
	stdoutBuff.Grow(256) // Preserve buffer Capacity

	intStream := make(chan int, 4) // <3>
	go func() {
		defer close(intStream)
		defer fmt.Fprintln(&stdoutBuff, "Producer Done.")
		for i := 0; i < 5; i++ {
			fmt.Fprintf(&stdoutBuff, "Sending: %d\n", i)
			intStream <- i
		}
	}()

	for integer := range intStream {
		fmt.Fprintf(&stdoutBuff, "Received %v.\n", integer)
	}
	// test1
	fmt.Println("stdoutBuff.Len():", stdoutBuff.Len())
}
