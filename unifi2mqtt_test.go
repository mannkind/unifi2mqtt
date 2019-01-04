package main

import (
	"testing"

	mqttExtCfg "github.com/mannkind/paho.mqtt.golang.ext/cfg"
	mqttExtDI "github.com/mannkind/paho.mqtt.golang.ext/di"
)

func defaultUnifi2Mqtt() *Unifi2Mqtt {
	c := NewUnifi2Mqtt(NewConfig(mqttExtCfg.NewMQTTConfig()), mqttExtDI.NewMQTTFuncWrapper())
	return c
}

func TestMqttConnect(t *testing.T) {
	c := defaultUnifi2Mqtt()
	c.onConnect(c.client)
}
