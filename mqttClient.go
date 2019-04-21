package main

import (
	"fmt"
	"log"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	mqttExtDI "github.com/mannkind/paho.mqtt.golang.ext/di"
	mqttExtHA "github.com/mannkind/paho.mqtt.golang.ext/ha"
)

const (
	sensorTopicTemplate = "%s/%s/state"
)

type mqttClient struct {
	discovery       bool
	discoveryPrefix string
	discoveryName   string
	topicPrefix     string

	macSlugMapping map[string]string
	previousState  map[string]string

	client mqtt.Client
}

func newMQTTClient(config *config, mqttFuncWrapper *mqttExtDI.MQTTFuncWrapper) *mqttClient {
	c := mqttClient{
		discovery:       config.MQTT.Discovery,
		discoveryPrefix: config.MQTT.DiscoveryPrefix,
		discoveryName:   config.MQTT.DiscoveryName,
		topicPrefix:     config.MQTT.TopicPrefix,
	}

	c.macSlugMapping = make(map[string]string, 0)
	c.previousState = make(map[string]string, 0)

	// Create a mapping between mac_addresses and names
	for _, m := range config.DeviceMapping {
		parts := strings.Split(m, ";")
		if len(parts) != 2 {
			continue
		}

		deviceMacAddress := parts[0]
		deviceName := parts[1]
		c.macSlugMapping[deviceMacAddress] = deviceName
		c.previousState[deviceName] = "init"
	}

	opts := mqttFuncWrapper.
		ClientOptsFunc().
		AddBroker(config.MQTT.Broker).
		SetClientID(config.MQTT.ClientID).
		SetOnConnectHandler(c.onConnect).
		SetConnectionLostHandler(c.onDisconnect).
		SetUsername(config.MQTT.Username).
		SetPassword(config.MQTT.Password).
		SetWill(c.availabilityTopic(), "offline", 0, true)

	c.client = mqttFuncWrapper.ClientFunc(opts)

	return &c
}

func (c *mqttClient) run() {
	log.Print("Connecting to MQTT")
	if token := c.client.Connect(); !token.Wait() || token.Error() != nil {
		log.Printf("Error connecting to MQTT: %s", token.Error())
		panic("Exiting...")
	}
}

func (c *mqttClient) onConnect(client mqtt.Client) {
	log.Print("Connected to MQTT")
	c.publish(c.availabilityTopic(), "online")
	c.publishDiscovery()
}

func (c *mqttClient) onDisconnect(client mqtt.Client, err error) {
	log.Printf("Disconnected from MQTT: %s.", err)
}

func (c *mqttClient) availabilityTopic() string {
	return fmt.Sprintf("%s/status", c.topicPrefix)
}

func (c *mqttClient) publishDiscovery() {
	if !c.discovery {
		return
	}

	for _, deviceName := range c.macSlugMapping {
		sensor := strings.ToLower(deviceName)
		mqd := mqttExtHA.MQTTDiscovery{
			DiscoveryPrefix: c.discoveryPrefix,
			Component:       "binary_sensor",
			NodeID:          c.discoveryName,
			ObjectID:        sensor,

			AvailabilityTopic: c.availabilityTopic(),
			Name:              fmt.Sprintf("%s %s", c.discoveryName, sensor),
			StateTopic:        fmt.Sprintf(sensorTopicTemplate, c.topicPrefix, sensor),
			UniqueID:          fmt.Sprintf("%s.%s", c.discoveryName, sensor),
			DeviceClass:       "presence",
		}

		mqd.PublishDiscovery(c.client)
	}
}

func (c *mqttClient) receive(e event) {
	slug := e.key
	payload := e.data
	topic := fmt.Sprintf(sensorTopicTemplate, c.topicPrefix, slug)

	// Don't publish the status if the device status hasn't changed
	if c.previousState[slug] == payload {
		return
	}

	c.previousState[slug] = payload
	c.publish(topic, payload)
}

func (c *mqttClient) publish(topic string, payload string) {
	retain := true
	if token := c.client.Publish(topic, 0, retain, payload); token.Wait() && token.Error() != nil {
		log.Printf("Publish Error: %s", token.Error())
	}

	log.Print(fmt.Sprintf("Publishing - Topic: %s ; Payload: %s", topic, payload))
}
