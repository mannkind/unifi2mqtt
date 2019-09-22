package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

// Name - Set during compilation when using included Makefile
var Name = "x2mqtt"

// Version - Set during compilation when using included Makefile
var Version = "X.X.X"

func main() {
	log.Infof("Version: %s", Version)

	x := initialize()
	x.run()

	select {}
}
