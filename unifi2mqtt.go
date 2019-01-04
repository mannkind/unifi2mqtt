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
)

// Unifi2Mqtt - Lookup device information from the unifi controller
type Unifi2Mqtt struct {
	discovery       bool
	discoveryPrefix string
	discoveryName   string
	topicPrefix     string
	host            string
	port            string
	site            string
	username        string
	password        string
	awayTimeout     time.Duration
	lookupInterval  time.Duration
	devices         map[string]string
	knownDevices    map[string]time.Time

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
	x.devices = make(map[string]string, 0)
	x.knownDevices = make(map[string]time.Time, 0)

	// Create a mapping between mac_addresses and names
	for _, m := range config.DeviceMapping {
		parts := strings.Split(m, ";")
		if len(parts) != 2 {
			continue
		}

		deviceMacAddress := parts[0]
		deviceName := parts[1]
		x.devices[deviceMacAddress] = deviceName
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
	u, err := unifi.Login(t.username, t.password, t.host, t.port, t.site, 5)
	if err != nil {
		log.Fatalf("Couldn't login to %s:%s with %s:<redacted>; %s", t.host, t.port, t.username, err)
		return
	}
	defer u.Logout()

	site, err := u.Site(t.site)
	if err != nil {
		log.Fatalf("Couldn't find site %s; %s", t.site, err)
		return
	}

	for {

		clients, err := u.Sta(site)
		if err != nil {
			log.Fatalf("Couldn't find clients; %s", err)
			continue
		}

		for _, client := range clients {
			slug, ok := t.devices[client.Mac]
			if !ok {
				// log.Printf("%s is not a device that is cared about", client.Mac)
				continue
			}

			ts := time.Unix(int64(client.LastSeen), 0)
			old := time.Since(ts) > t.awayTimeout
			topic := fmt.Sprintf(sensorTopicTemplate, t.topicPrefix, slug)
			if _, ok := t.knownDevices[client.Mac]; ok && old {
				log.Printf("%s (%s) is away", client.Mac, slug)
				t.publish(topic, "not_home")
				delete(t.knownDevices, client.Mac)
			} else if !ok && !old {
				log.Printf("%s (%s) is home", client.Mac, slug)
				t.publish(topic, "home")
				t.knownDevices[client.Mac] = ts
			}
		}

		time.Sleep(t.lookupInterval)
	}
}

func (t *Unifi2Mqtt) publishDiscovery() {
	for _, deviceName := range t.devices {
		sensor := strings.ToLower(deviceName)
		mqd := mqttExtHA.MQTTDiscovery{
			DiscoveryPrefix: t.discoveryPrefix,
			Component:       "binary_sensor",
			NodeID:          t.discoveryName,
			ObjectID:        sensor,

			Name:       fmt.Sprintf("%s %s", t.discoveryName, sensor),
			StateTopic: fmt.Sprintf(sensorTopicTemplate, t.topicPrefix, sensor),
			UniqueID:   fmt.Sprintf("%s.%s", t.discoveryName, sensor),
			Icon:       "mdi:phone",
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
