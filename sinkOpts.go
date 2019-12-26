package main

import "github.com/mannkind/twomqtt"

type sinkOpts struct {
	globalOpts
	MQTTOpts twomqtt.MQTTOpts
}
