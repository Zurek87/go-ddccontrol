package main

import (
	"./goddcci"
	"log"
	"fmt"
	"time"
)

func main() {
	ddcci, err := goddcci.InitDDCci()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Time to setup")

	fmt.Println("Time to setup")
	ddcci.SetBrightness(100)
	time.Sleep(5 * time.Second)
	ddcci.SetBrightness(0)
}