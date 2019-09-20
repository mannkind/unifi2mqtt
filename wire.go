//+build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/mannkind/twomqtt"
)

func initialize() *bridge {
	wire.Build(
		newBridge,
		newStateChannel,
		newMQTTClient,
		newServiceClient,
		newConfig,
		wire.FieldsOf(new(config), "MQTTClientConfig"),
		wire.FieldsOf(new(config), "ServiceClientConfig"),
		wire.FieldsOf(new(mqttClientConfig), "MQTTProxyConfig"),
		twomqtt.NewMQTTProxy,
	)

	return &bridge{}
}
