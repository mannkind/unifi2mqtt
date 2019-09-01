package main

import "github.com/mannkind/twomqtt"

type mqttClientConfig struct {
	globalClientConfig
	MQTTProxyConfig twomqtt.MQTTProxyConfig
}
