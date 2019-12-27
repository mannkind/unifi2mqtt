package main

type comms struct {
	input  <-chan sourceRep
	output chan<- sourceRep
}

func newComms() comms {
	c := make(chan sourceRep, 100)
	return comms{
		c,
		c,
	}
}
