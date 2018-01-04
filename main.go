package main

import (
	"./goddcci"
	"log"
	"fmt"
)

func main() {
	ddcci, err := goddcci.InitDDCci()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ddcci)
}