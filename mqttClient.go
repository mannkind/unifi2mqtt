package main

import (
	"reflect"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mannkind/twomqtt"
	log "github.com/sirupsen/logrus"
)

type mqttClient struct {
	twomqtt.Observer
	*twomqtt.MQTTProxy
	mqttClientConfig
}

func newMQTTClient(mqttClientCfg mqttClientConfig, client *twomqtt.MQTTProxy) *mqttClient {
	c := mqttClient{
		MQTTProxy:        client,
		mqttClientConfig: mqttClientCfg,
	}

	c.Initialize(
		c.onConnect,
		c.onDisconnect,
	)

	c.LogSettings()

	return &c
}

func (c *mqttClient) run() {
	c.Run()
}

func (c *mqttClient) onConnect(client mqtt.Client) {
	log.Info("Finished connecting to MQTT")
	c.Publish(c.AvailabilityTopic(), "online")
	c.publishDiscovery()
}

func (c *mqttClient) onDisconnect(client mqtt.Client, err error) {
	log.WithFields(log.Fields{
		"error": err,
	}).Error("Disconnected from MQTT")
}

func (c *mqttClient) publishDiscovery() {
	if !c.Discovery {
		return
	}

	for _, deviceName := range c.Devices {
		sensor := strings.ToLower(deviceName)
		mqd := c.NewMQTTDiscovery("", sensor, "binary_sensor")
		mqd.DeviceClass = "presence"

		c.PublishDiscovery(mqd)
	}
}

func (c *mqttClient) ReceiveCommand(cmd twomqtt.Command, e twomqtt.Event) {}
func (c *mqttClient) ReceiveState(e twomqtt.Event) {
	if e.Type != reflect.TypeOf(unifiDevice{}) {
		msg := "Unexpected event type; skipping"
		log.WithFields(log.Fields{
			"type": e.Type,
		}).Error(msg)
		return
	}

	info := e.Payload.(unifiDevice)
	topic := c.StateTopic("", info.name)

	c.Publish(topic, info.state)
}
