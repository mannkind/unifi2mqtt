//+build wireinject

package main

import (
	"github.com/google/wire"
	mqttExtCfg "github.com/mannkind/paho.mqtt.golang.ext/cfg"
	mqttExtDI "github.com/mannkind/paho.mqtt.golang.ext/di"
)

func initialize() *bridge {
	wire.Build(mqttExtCfg.NewMQTTConfig, mqttExtDI.NewMQTTFuncWrapper, newConfig, newBridge, newMQTTClient, newClient)

	return &bridge{}
}
