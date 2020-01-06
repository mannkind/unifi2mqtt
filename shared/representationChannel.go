package shared

// RepresentationChannel is an inter-application communication channel
type RepresentationChannel = chan Representation

// NewRepresentationChannel creates the required object(s) for communication
func NewRepresentationChannel() RepresentationChannel {
	return make(chan Representation, 100)
}

// NewRepresentationChannelIncoming creates the required object(s) for incoming representations
func NewRepresentationChannelIncoming(c RepresentationChannel) <-chan Representation {
	return c
}

// NewRepresentationChannelOutgoing creates the required object(s) for outgoing representations
func NewRepresentationChannelOutgoing(c RepresentationChannel) chan<- Representation {
	return c
}
