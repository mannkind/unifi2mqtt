# unifi2mqtt

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
MQTT_TOPICPREFIX        - [OPTIONAL] The MQTT topic on which to publish the collection lookup results, defaults to "home/wsdot"
MQTT_DISCOVERY          - [OPTIONAL] The MQTT discovery flag for Home Assistant, defaults to false
MQTT_DISCOVERYPREFIX    - [OPTIONAL] The MQTT discovery prefix for Home Assistant, defaults to "homeassistant"
MQTT_DISCOVERYNAME      - [OPTIONAL] The MQTT discovery name for Home Assistant, defaults to "wsdot"
MQTT_CLIENTID           - [OPTIONAL] The clientId, defaults to "Defaultwsdot2MQTTClientID"
MQTT_BROKER             - [OPTIONAL] The MQTT broker, defaults to "tcp://mosquitto.org:1883"
MQTT_USERNAME           - [OPTIONAL] The MQTT username, default to ""
MQTT_PASSWORD           - [OPTIONAL] The MQTT password, default to ""
```
