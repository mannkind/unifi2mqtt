package main

import (
	"os"
	"testing"

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
		Devices             string
		DiscoveryName       string
		TopicPrefix         string
		ExpectedName        string
		ExpectedStateTopic  string
		ExpectedUniqueID    string
		ExpectedDeviceClass string
	}{
		{
			knownMAC + ";" + knownMACName,
			defaultDiscoveryName,
			defaultTopicPrefix,
			defaultDiscoveryName + " " + knownMACName,
			defaultTopicPrefix + "/" + knownMACName + "/state",
			defaultDiscoveryName + "." + knownMACName,
			"presence",
		},
		{
			knownMAC + ";" + knownMACName,
			knownDiscoveryName,
			knownTopicPrefix,
			knownDiscoveryName + " " + knownMACName,
			knownTopicPrefix + "/" + knownMACName + "/state",
			knownDiscoveryName + "." + knownMACName,
			"presence",
		},
	}

	for _, v := range tests {
		setEnvs("true", v.DiscoveryName, v.TopicPrefix, v.Devices)

		c := initialize()
		mqds := c.sink.discovery()

		for _, mqd := range mqds {
			if mqd.Name != v.ExpectedName {
				t.Errorf("discovery Name does not match; %s vs %s", mqd.Name, v.ExpectedName)
			}
			if mqd.StateTopic != v.ExpectedStateTopic {
				t.Errorf("discovery StateTopic does not match; %s vs %s", mqd.StateTopic, v.ExpectedStateTopic)
			}
			if mqd.UniqueID != v.ExpectedUniqueID {
				t.Errorf("discovery UniqueID does not match; %s vs %s", mqd.UniqueID, v.ExpectedUniqueID)
			}
			if mqd.DeviceClass != v.ExpectedDeviceClass {
				t.Errorf("discovery DeviceClass does not match; %s vs %s", mqd.DeviceClass, v.ExpectedDeviceClass)
			}
		}
	}
}

func TestPublish(t *testing.T) {
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

		obj := sourceRep{
			name:  v.Name,
			state: "ON",
		}

		c := initialize()
		published := c.sink.publish(obj)

		if published.Payload != v.ExpectedPayload {
			t.Errorf("Actual:%s\nExpected:%s", published.Payload, v.ExpectedPayload)
		}
	}
}
