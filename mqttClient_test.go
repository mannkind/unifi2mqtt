package main

import (
	"os"
	"reflect"
	"testing"

	"github.com/mannkind/twomqtt"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.PanicLevel)
}

func setEnvs(d, dn, tp, a string) {
	os.Setenv("MQTT_DISCOVERY", d)
	os.Setenv("MQTT_DISCOVERYNAME", dn)
	os.Setenv("MQTT_TOPICPREFIX", tp)
	os.Setenv("UNIFI_DEVICEMAPPING", a)
}

func clearEnvs() {
	setEnvs("false", "", "", "")
}

const defaultDiscoveryName = "unifi"
const defaultTopicPrefix = "home/unifi"
const knownMAC = "11:22:33:44:55:66"
const knownMACName = "myhouse"
const knownDiscoveryName = "unifiDiscoveryName"
const knownTopicPrefix = "home/unifiMQTTTopicPrefix"

func TestDiscovery(t *testing.T) {
	defer clearEnvs()

	var tests = []struct {
		Devices         string
		DiscoveryName   string
		TopicPrefix     string
		ExpectedTopic   string
		ExpectedPayload string
	}{
		{
			knownMAC + ";" + knownMACName,
			defaultDiscoveryName,
			defaultTopicPrefix,
			"homeassistant/binary_sensor/" + defaultDiscoveryName + "/" + knownMACName + "/config",
			"{\"availability_topic\":\"" + defaultTopicPrefix + "/status\",\"device_class\":\"presence\",\"name\":\"" + defaultDiscoveryName + " " + knownMACName + "\",\"state_topic\":\"" + defaultTopicPrefix + "/" + knownMACName + "/state\",\"unique_id\":\"" + defaultDiscoveryName + "." + knownMACName + "\"}",
		},
		{
			knownMAC + ";" + knownMACName,
			knownDiscoveryName,
			knownTopicPrefix,
			"homeassistant/binary_sensor/" + knownDiscoveryName + "/" + knownMACName + "/config",
			"{\"availability_topic\":\"" + knownTopicPrefix + "/status\",\"device_class\":\"presence\",\"name\":\"" + knownDiscoveryName + " " + knownMACName + "\",\"state_topic\":\"" + knownTopicPrefix + "/" + knownMACName + "/state\",\"unique_id\":\"" + knownDiscoveryName + "." + knownMACName + "\"}",
		},
	}

	for _, v := range tests {
		setEnvs("true", v.DiscoveryName, v.TopicPrefix, v.Devices)

		c := initialize()
		c.mqttClient.publishDiscovery()

		actualPayload := c.mqttClient.LastPublishedOnTopic(v.ExpectedTopic)
		if actualPayload != v.ExpectedPayload {
			t.Errorf("Actual:%s\nExpected:%s", actualPayload, v.ExpectedPayload)
		}
	}
}

func TestReceieveState(t *testing.T) {
	defer clearEnvs()

	var tests = []struct {
		Devices         string
		Name            string
		TopicPrefix     string
		ExpectedTopic   string
		ExpectedPayload string
	}{
		{
			knownMAC + ";" + knownMACName,
			knownMACName,
			defaultTopicPrefix,
			defaultTopicPrefix + "/" + knownMACName + "/state",
			"ON",
		},
	}

	for _, v := range tests {
		setEnvs("true", "", v.TopicPrefix, v.Devices)

		obj := unifiDevice{
			name:  v.Name,
			state: "ON",
		}
		event := twomqtt.Event{
			Type:    reflect.TypeOf(obj),
			Payload: obj,
		}

		c := initialize()
		c.mqttClient.ReceiveState(event)

		actualPayload := c.mqttClient.LastPublishedOnTopic(v.ExpectedTopic)
		if actualPayload != v.ExpectedPayload {
			t.Errorf("Actual:%s\nExpected:%s", actualPayload, v.ExpectedPayload)
		}
	}
}
