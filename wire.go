//+build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/mannkind/twomqtt"
)

func initialize() *app {
	wire.Build(
		newOpts,
		newApp,
		newStateChannel,
		newSink,
		newSource,
		wire.FieldsOf(new(stateChannel), "Input"),
		wire.FieldsOf(new(stateChannel), "Output"),
		wire.FieldsOf(new(opts), "Sink"),
		wire.FieldsOf(new(opts), "Source"),
		wire.FieldsOf(new(sinkOpts), "MQTTOpts"),
		twomqtt.NewMQTT,
	)

	return &app{}
}
