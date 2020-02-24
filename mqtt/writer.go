package mqtt

import (
	"strings"

	"github.com/mannkind/twomqtt"
	"github.com/mannkind/unifi2mqtt/shared"
)

// Writer is for writing a shared representation to MQTT
type Writer struct {
	*twomqtt.MQTT
	opts     Opts
	incoming <-chan shared.Representation
}

// NewWriter creates a new Writer for writing a shared representation to MQTT
func NewWriter(mqtt *twomqtt.MQTT, opts Opts, incoming <-chan shared.Representation) *Writer {
	c := Writer{
		MQTT:     mqtt,
		opts:     opts,
		incoming: incoming,
	}

	c.MQTT.
		SetDiscoveryHandler(c.discovery).
		SetReadIncomingChannelHandler(c.read).
		Initialize()

	return &c
}

// discovery objects to publish when MQTT discovery is enabled
func (c *Writer) discovery() []twomqtt.MQTTDiscovery {
	mqds := []twomqtt.MQTTDiscovery{}
	if !c.Discovery {
		return mqds
	}

	for _, deviceName := range c.opts.Devices {
		sensorName := strings.ToLower(deviceName)
		sensorType := "binary_sensor"
		mqd := twomqtt.NewMQTTDiscovery(c.opts.MQTTOpts, "", sensorName, sensorType)
		mqd.DeviceClass = "presence"
		mqd.Device.Name = shared.Name
		mqd.Device.SWVersion = shared.Version

		mqds = append(mqds, *mqd)
	}

	return mqds
}

// read incoming shared representations and publish them to MQTT
func (c *Writer) read() {
	for info := range c.incoming {
		c.publish(info)
	}
}

// publish a shared representation to MQTT
func (c *Writer) publish(info shared.Representation) twomqtt.MQTTMessage {
	topic := c.StateTopic("", info.Name)
	payload := info.State

	return c.Publish(topic, payload)
}
