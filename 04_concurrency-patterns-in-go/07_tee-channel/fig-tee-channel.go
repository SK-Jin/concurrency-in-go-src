package main

import (
	"fmt"
	"time"
)

func main() {
	orDone := func(done, c <-chan interface{}) <-chan interface{} {
		valStream := make(chan interface{})
		go func() {
			defer close(valStream)
			for {
				select {
				case <-done:
					return
				case v, ok := <-c:
					//if ok == false {
					if !ok {
						return
					}
					select {
					case valStream <- v:
					case <-done:
					}
				}
			}
		}()
		return valStream
	}
	tee := func(
		done <-chan interface{},
		in <-chan interface{},
	) (_, _ <-chan interface{}) {
		out1 := make(chan interface{})
		out2 := make(chan interface{})
		go func() {
			defer close(out1)
			defer close(out2)
			start := time.Now()
			for val := range orDone(done, in) {
				//startInner := time.Now()
				//fmt.Println("tee, for out1:", out1, ", out2:", out2)
				var out1, out2 = out1, out2
				for i := 0; i < 2; i++ {
					select {
					case <-done:
					case out1 <- val:
						//fmt.Println("tee, out1:", val)
						out1 = nil
					case out2 <- val:
						//fmt.Println("tee, out2:", val)
						out2 = nil
					}
				}
				//fmt.Printf("Search took: %v", time.Since(start))
				//fmt.Printf("tee, for val:%d, time:%v\n", val, time.Since(startInner))
			}
			fmt.Printf("tee, time:%v\n", time.Since(start))
		}()
		return out1, out2
	}
	repeat := func(
		done <-chan interface{},
		values ...interface{},
	) <-chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				for _, v := range values {
					select {
					case <-done:
						return
					case valueStream <- v:
					}
				}
			}
		}()
		return valueStream
	}
	take := func(
		done <-chan interface{},
		valueStream <-chan interface{},
		num int,
	) <-chan interface{} {
		takeStream := make(chan interface{})
		go func() {
			defer close(takeStream)
			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				case takeStream <- <-valueStream:
				}
			}
		}()
		return takeStream
	}
	done := make(chan interface{})
	defer close(done)

	out1, out2 := tee(done, take(done, repeat(done, 2, 3, 5, 7, 9), 10))
	for val1 := range out1 {
		fmt.Printf("===== out1: %v, out2: %v\n", val1, <-out2)
	}
}
