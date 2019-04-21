package main

import (
	"log"
)

// Version - Set during compilation when using included Makefile
var Version = "X.X.X"

func main() {
	log.Printf("Version: %s", Version)

	x := initialize()
	x.run()

	select {}
}
