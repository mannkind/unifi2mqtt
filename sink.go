package main

import (
	"strings"

	"github.com/mannkind/twomqtt"
)

type sink struct {
	*twomqtt.MQTT
	config   sinkOpts
	incoming <-chan sourceRep
}

func newSink(mqtt *twomqtt.MQTT, config sinkOpts, incoming <-chan sourceRep) *sink {
	c := sink{
		MQTT:     mqtt,
		config:   config,
		incoming: incoming,
	}

	c.MQTT.
		SetDiscoveryHandler(c.discovery).
		SetReadIncomingChannelHandler(c.read).
		Initialize()

	return &c
}

func (c *sink) run() {
	c.Run()
}

func (c *sink) discovery() []twomqtt.MQTTDiscovery {
	mqds := []twomqtt.MQTTDiscovery{}
	if !c.Discovery {
		return mqds
	}

	for _, deviceName := range c.config.Devices {
		sensorName := strings.ToLower(deviceName)
		sensorType := "binary_sensor"
		mqd := twomqtt.NewMQTTDiscovery(c.config.MQTTOpts, "", sensorName, sensorType)
		mqd.DeviceClass = "presence"
		mqd.Device.Name = Name
		mqd.Device.SWVersion = Version

		mqds = append(mqds, *mqd)
	}

	return mqds
}

func (c *sink) read() {
	for info := range c.incoming {
		c.publish(info)
	}
}

func (c *sink) publish(info sourceRep) twomqtt.MQTTMessage {
	topic := c.StateTopic("", info.name)
	payload := info.state

	return c.Publish(topic, payload)
}
