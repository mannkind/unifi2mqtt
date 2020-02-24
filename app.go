package main

import (
	"github.com/mannkind/unifi2mqtt/mqtt"
	"github.com/mannkind/unifi2mqtt/source"
)

type app struct {
	source *source.Reader
	sink   *mqtt.Writer
}

func newApp(src *source.Reader, snk *mqtt.Writer) *app {
	c := app{
		sink:   snk,
		source: src,
	}

	return &c
}

func (c *app) run() {
	c.sink.Run()
	c.source.Run()
}
