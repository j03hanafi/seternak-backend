package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Starting program...")

	fmt.Println("Stopping program in...")
	for i := 3; i > 0; i-- {
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}
}
