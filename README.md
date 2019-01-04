# unifi2mqtt

An experiment to publish device statuses from the Unifi Controller to MQTT.

[![Software
License](https://img.shields.io/badge/License-MIT-orange.svg?style=flat-square)](https://github.com/mannkind/unifi2mqtt/blob/master/LICENSE.md)
[![Travis CI](https://img.shields.io/travis/mannkind/unifi2mqtt/master.svg?style=flat-square)](https://travis-ci.org/mannkind/unifi2mqtt)
[![Coverage Status](https://img.shields.io/codecov/c/github/mannkind/unifi2mqtt/master.svg)](http://codecov.io/github/mannkind/unifi2mqtt?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/mannkind/unifi2mqtt)](https://goreportcard.com/report/github.com/mannkind/unifi2mqtt)

## Installation

### Via Docker

```bash
docker run -d --name="unifi2mqtt" -v /etc/localtime:/etc/localtime:ro mannkind/unifi2mqtt
```

### Via Make

```bash
git clone https://github.com/mannkind/unifi2mqtt
cd unifi2mqtt
make
./unifi2mqtt
```

## Configuration

Configuration happens via environmental variables

```bash
UNIFI_HOST              - The hostname of the controller, defaults to "unifi.local"
UNIFI_PORT              - The the port of the controller, defaults to "8843"
UNIFI_SITE              - The site of the controller, defaults to "default"
UNIFI_USERNAME          - The username used to access the controller, defaults to "unifi"
UNIFI_PASSWORD          - The password used to access the controller, defaults to "unifi"
UNIFI_AWAYTIMEOUT       - The timeout before marking a device as away, defaults to "5m"
UNIFI_LOOKUPINTERVAL    - The interval to lookup devices on the controller, defaults to "10s"
UNIFI_DEVICEMAPPING     - The map of mac addresses to names, defaults to "11:22:33:44:55:66;MyPhone,12:23:34:45:56:67;AnotherPhone"
MQTT_TOPICPREFIX        - [OPTIONAL] The MQTT topic on which to publish the collection lookup results, defaults to "home/unifi"
MQTT_DISCOVERY          - [OPTIONAL] The MQTT discovery flag for Home Assistant, defaults to false
MQTT_DISCOVERYPREFIX    - [OPTIONAL] The MQTT discovery prefix for Home Assistant, defaults to "homeassistant"
MQTT_DISCOVERYNAME      - [OPTIONAL] The MQTT discovery name for Home Assistant, defaults to "unifi"
MQTT_CLIENTID           - [OPTIONAL] The clientId, defaults to "DefaultUNIFI2MQTTClientID"
MQTT_BROKER             - [OPTIONAL] The MQTT broker, defaults to "tcp://mosquitto.org:1883"
MQTT_USERNAME           - [OPTIONAL] The MQTT username, default to ""
MQTT_PASSWORD           - [OPTIONAL] The MQTT password, default to ""
```
