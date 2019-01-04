//+build wireinject

package main

import (
	"github.com/google/wire"
	mqttExtCfg "github.com/mannkind/paho.mqtt.golang.ext/cfg"
	mqttExtDI "github.com/mannkind/paho.mqtt.golang.ext/di"
)

// Initialize - Compile-time DI
func Initialize() *Unifi2Mqtt {
	wire.Build(mqttExtCfg.NewMQTTConfig, NewConfig, mqttExtDI.NewMQTTFuncWrapper, NewUnifi2Mqtt)

	return &Unifi2Mqtt{}
}
