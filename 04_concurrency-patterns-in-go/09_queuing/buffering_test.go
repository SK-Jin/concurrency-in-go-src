package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func BenchmarkUnbufferedWrite(b *testing.B) {
	performWrite(b, tmpFileOrFatal())
}

func BenchmarkBufferedWrite(b *testing.B) {
	bufferredFile := bufio.NewWriter(tmpFileOrFatal())
	performWrite(b, bufio.NewWriter(bufferredFile))
	//performWrite(b, bufferredFile)
}

func tmpFileOrFatal() *os.File {
	file, err := ioutil.TempFile("", "tmp")
	if err != nil {
		//log.Fatal("error: %v", err)
		log.Fatalf("error: %v", err)
	}
	return file
}

func performWrite(b *testing.B, writer io.Writer) {
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

	// sleep := func(
	// 	done <-chan interface{},
	// 	sleepTime time.Duration,
	// 	valueStream <-chan interface{},
	// ) <-chan interface{} {
	// 	takeStream := make(chan interface{})
	// 	go func() {
	// 		defer close(takeStream)

	// 		select {
	// 		case <-done:
	// 		default:
	// 			time.Sleep(sleepTime)

	// 			select {
	// 			case <-done:
	// 			case takeStream <- <-valueStream:
	// 			}
	// 		}
	// 	}()
	// 	return takeStream
	// }

	done := make(chan interface{})
	defer close(done)

	b.ResetTimer()
	// zeros := take(done, repeat(done, 0), b.N)
	// short := sleep(done, 1*time.Second, zeros)
	// long := sleep(done, 4*time.Second, short)
	// pipeline := long

	for bt := range take(done, repeat(done, byte(0)), b.N) {
		//for bt := range pipeline {
		writer.Write([]byte{bt.(byte)})
	}
}
