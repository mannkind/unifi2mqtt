package main

type stateChannel struct {
	input  <-chan sourceRep
	output chan<- sourceRep
}

func newStateChannel() stateChannel {
	c := make(chan sourceRep, 100)
	return stateChannel{
		c,
		c,
	}
}
