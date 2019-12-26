package main

type app struct {
	source *source
	sink   *sink
}

func newApp(src *source, snk *sink) *app {
	c := app{
		sink:   snk,
		source: src,
	}

	return &c
}

func (c *app) run() {
	c.sink.run()
	c.source.run()
}
