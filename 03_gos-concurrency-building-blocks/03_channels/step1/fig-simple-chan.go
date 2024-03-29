package main

import (
	"fmt"
)

func main() {
	stringStream := make(chan string)
	go func() {
		// test2
		// if 0 != 1 {
		// 	return
		// }
		//traceBlockChanRecv (7) = 0x7
		//waitReasonChanReceive (14) = 0xe
		stringStream <- "Hello channels!" // <1>
	}()
	fmt.Println(<-stringStream) // <2>

	// test1
	// writeStream := make(chan<- interface{})
	// readStream := make(<-chan interface{})
	// <-writeStream            //invalid operation: cannot receive from send-only channel writeStream (variable of type chan<-
	// readStream <- struct{}{} //invalid operation: cannot send to receive-only channel readStream (variable of type <-chan
}
