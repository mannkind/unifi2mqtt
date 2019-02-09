package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dim13/unifi"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	mqttExtDI "github.com/mannkind/paho.mqtt.golang.ext/di"
	mqttExtHA "github.com/mannkind/paho.mqtt.golang.ext/ha"
)

const (
	sensorTopicTemplate = "%s/%s/state"
	home                = "ON"
	notHome             = "OFF"
)

// Unifi2Mqtt - Lookup device information from the unifi controller
type Unifi2Mqtt struct {
	discovery            bool
	discoveryPrefix      string
	discoveryName        string
	topicPrefix          string
	host                 string
	port                 string
	site                 string
	username             string
	password             string
	awayTimeout          time.Duration
	lookupInterval       time.Duration
	macSlugMapping       map[string]string
	deviceStatus         map[string]string
	devicePreviousStatus map[string]string

	client mqtt.Client
}

// NewUnifi2Mqtt - Returns a new reference to a fully configured object.
func NewUnifi2Mqtt(config *Config, mqttFuncWrapper *mqttExtDI.MQTTFuncWrapper) *Unifi2Mqtt {
	x := Unifi2Mqtt{
		discovery:       config.MQTT.Discovery,
		discoveryPrefix: config.MQTT.DiscoveryPrefix,
		discoveryName:   config.MQTT.DiscoveryName,
		topicPrefix:     config.MQTT.TopicPrefix,
		awayTimeout:     config.AwayTimeout,
		lookupInterval:  config.LookupInterval,
		host:            config.Host,
		port:            config.Port,
		site:            config.Site,
		username:        config.Username,
		password:        config.Password,
	}
	x.macSlugMapping = make(map[string]string, 0)
	x.deviceStatus = make(map[string]string, 0)
	x.devicePreviousStatus = make(map[string]string, 0)

	// Create a mapping between mac_addresses and names
	for _, m := range config.DeviceMapping {
		parts := strings.Split(m, ";")
		if len(parts) != 2 {
			continue
		}

		deviceMacAddress := parts[0]
		deviceName := parts[1]
		x.macSlugMapping[deviceMacAddress] = deviceName
		x.deviceStatus[deviceName] = notHome
		x.devicePreviousStatus[deviceName] = "init"
	}

	opts := mqttFuncWrapper.
		ClientOptsFunc().
		AddBroker(config.MQTT.Broker).
		SetClientID(config.MQTT.ClientID).
		SetOnConnectHandler(x.onConnect).
		SetConnectionLostHandler(x.onDisconnect).
		SetUsername(config.MQTT.Username).
		SetPassword(config.MQTT.Password)

	x.client = mqttFuncWrapper.ClientFunc(opts)

	return &x
}

// Run - Start the collection lookup process
func (t *Unifi2Mqtt) Run() error {
	log.Print("Connecting to MQTT")
	if token := t.client.Connect(); !token.Wait() || token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (t *Unifi2Mqtt) onConnect(client mqtt.Client) {
	log.Print("Connected to MQTT")

	if !client.IsConnected() {
		log.Print("Subscribe Error: Not Connected (Reloading Config?)")
		return
	}

	if t.discovery {
		t.publishDiscovery()
	}

	go t.loop()
}

func (t *Unifi2Mqtt) onDisconnect(client mqtt.Client, err error) {
	log.Printf("Disconnected from MQTT: %s.", err)
}

func (t *Unifi2Mqtt) loop() {
	log.Print("Connecting to Unifi")
	u, err := unifi.Login(t.username, t.password, t.host, t.port, t.site, 5)
	if err != nil {
		log.Fatalf("Couldn't login to %s:%s with %s:<redacted>; %s", t.host, t.port, t.username, err)
		return
	}
	defer u.Logout()
	log.Print("Connected to Unifi")

	log.Print("Selecting Site")
	site, err := u.Site(t.site)
	if err != nil {
		log.Fatalf("Couldn't find site %s; %s", t.site, err)
		return
	}
	log.Print("Selected Site")

	log.Print("Forever Fetching Clients")
	for {
		// Default the status of all known macSlugMapping to not_home
		for _, slug := range t.macSlugMapping {
			t.deviceStatus[slug] = notHome
		}

		// Ask the unifi controller for all known clients
		// log.Print("Fetching Clients")
		clients, err := u.Sta(site)
		if err != nil {
			log.Fatalf("Couldn't find clients; %s", err)
			continue
		}
		// log.Print("Fetched Clients")

		// Devices known to the controller will have their status set based on the Last Seen time
		// Devices missing will remain defaulted to not_home
		for _, client := range clients {
			slug, ok := t.macSlugMapping[client.Mac]
			if !ok {
				// log.Printf("%s is not a device that is cared about", client.Mac)
				continue
			}

			ts := time.Unix(int64(client.LastSeen), 0)
			old := time.Since(ts) > t.awayTimeout
			payload := home
			if old {
				payload = notHome
			}

			t.deviceStatus[slug] = payload
		}

		// Publish known device statuses
		// log.Print("Publishing Statuses")
		for slug, payload := range t.deviceStatus {

			// Don't publish the status if the device status hasn't changed
			previousPayload, _ := t.devicePreviousStatus[slug]
			if previousPayload == payload {
				continue
			}

			t.devicePreviousStatus[slug] = payload
			topic := fmt.Sprintf(sensorTopicTemplate, t.topicPrefix, slug)
			t.publish(topic, payload)
		}
		// log.Print("Published Statuses")

		time.Sleep(t.lookupInterval)
	}
}

func (t *Unifi2Mqtt) publishDiscovery() {
	for _, deviceName := range t.macSlugMapping {
		sensor := strings.ToLower(deviceName)
		mqd := mqttExtHA.MQTTDiscovery{
			DiscoveryPrefix: t.discoveryPrefix,
			Component:       "binary_sensor",
			NodeID:          t.discoveryName,
			ObjectID:        sensor,

			Name:        fmt.Sprintf("%s %s", t.discoveryName, sensor),
			StateTopic:  fmt.Sprintf(sensorTopicTemplate, t.topicPrefix, sensor),
			UniqueID:    fmt.Sprintf("%s.%s", t.discoveryName, sensor),
			Icon:        "mdi:phone",
			DeviceClass: "presence",
		}

		mqd.PublishDiscovery(t.client)
	}
}

func (t *Unifi2Mqtt) publish(topic string, payload string) {
	retain := true
	if token := t.client.Publish(topic, 0, retain, payload); token.Wait() && token.Error() != nil {
		log.Printf("Publish Error: %s", token.Error())
	}

	log.Print(fmt.Sprintf("Publishing - Topic: %s ; Payload: %s", topic, payload))
}
