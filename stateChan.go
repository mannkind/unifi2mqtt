package main

type stateChannel = chan unifiDevice

func newStateChannel() stateChannel {
	return make(stateChannel, 100)
}
