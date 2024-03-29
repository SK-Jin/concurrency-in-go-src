package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// https://stackoverflow.com/questions/65446569/concurrency-and-replicated-requests-what-is-time-after-for
func main() {
	doWork := func(
		done <-chan interface{},
		id int,
		wg *sync.WaitGroup,
		result chan<- int,
	) {
		started := time.Now()
		defer wg.Done()

		// Simulate random load
		simulatedLoadTime := time.Duration(1*rand.Intn(5)) * time.Second

		/** use two separate select blocks because we want to send/receive two different values, the time.After (receive) and the id (send).
		  / if they were in the same select block, then we could only use one value at a time, the other will get lost. */
		select {
		// do not want to return on <-done because we still want to log the time it took
		case <-done:
		case <-time.After(simulatedLoadTime):
		}

		select {
		case <-done:
		case result <- id:
		}

		took := time.Since(started)
		// Display how long handlers would have taken
		if took < simulatedLoadTime {
			took = simulatedLoadTime
		}
		fmt.Printf("%v took %v\n", id, took)
	}

	start := time.Now()
	done := make(chan interface{})
	result := make(chan int)

	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go doWork(done, i, &wg, result)
	}

	firstReturned := <-result
	close(done)
	fmt.Printf("1st, Received an answer from #%v, %v\n", firstReturned, time.Since(start))

	wg.Wait()

	//fmt.Printf("wait, Received an answer from #%v\n", firstReturned)
}
