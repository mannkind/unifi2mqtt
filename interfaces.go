package main

type event struct {
	version int64
	key     string
	data    string
}

type observer interface {
	receiveState(event)
	receiveCommand(int64, event)
}

type publisher interface {
	register(observer)
}
