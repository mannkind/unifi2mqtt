package main

import (
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mannkind/twomqtt"
	log "github.com/sirupsen/logrus"
)

type mqttClient struct {
	mqttClientConfig
	*twomqtt.MQTTProxy
	stateUpdateChan stateChannel
}

func newMQTTClient(mqttClientCfg mqttClientConfig, client *twomqtt.MQTTProxy, stateUpdateChan stateChannel) *mqttClient {
	c := mqttClient{
		MQTTProxy:        client,
		mqttClientConfig: mqttClientCfg,
		stateUpdateChan:  stateUpdateChan,
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
	go c.receive()
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
		mqd.Device.Name = Name
		mqd.Device.SWVersion = Version

		c.PublishDiscovery(mqd)
	}
}

func (c *mqttClient) receive() {
	for info := range c.stateUpdateChan {
		c.receiveState(info)
	}
}

func (c *mqttClient) receiveState(info unifiDevice) {
	topic := c.StateTopic("", info.name)

	c.Publish(topic, info.state)
}
