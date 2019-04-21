package main

type event struct {
	version int64
	key     string
	data    string
}

type observer interface {
	receive(event)
}

type publisher interface {
	register(observer)
	publish(event)
}
