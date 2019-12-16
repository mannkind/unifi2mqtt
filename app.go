package main

type app struct {
	source *source
	sink   *sink
}

func newApp(src *source, snk *sink) *app {
	bridge := app{
		sink:   snk,
		source: src,
	}

	return &bridge
}

func (b *app) run() {
	b.sink.run()
	b.source.run()
}
