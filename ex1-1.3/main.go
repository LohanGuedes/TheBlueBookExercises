package main

import (
	"fmt"
	"os"
	"time"
)

/* medium: 0.043us */
func main() {
	start := time.Now()
	fmt.Println(os.Args[0])

	for k, arg := range os.Args[1:] {
		fmt.Printf("[%v]:  %v\n", k, arg)
	}

	fmt.Printf("Time elapsed: %v\n", time.Since(start).Abs().Microseconds())
}
