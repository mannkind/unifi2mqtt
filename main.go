package main

import (
	"log"
)

// Version - Set during compilation when using included Makefile
var Version = "X.X.X"

func main() {
	log.Printf("Version: %s", Version)

	x := Initialize()
	if err := x.Run(); err != nil {
		log.Panicf("Run Error: %s", err)
	}
}
