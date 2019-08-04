package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

// Version - Set during compilation when using included Makefile
var Version = "X.X.X"

func main() {
	log.Infof("Version: %s", Version)

	x := initialize()
	x.run()

	select {}
}
