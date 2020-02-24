# unifi2mqtt

[![Software
License](https://img.shields.io/badge/License-MIT-orange.svg?style=flat-square)](https://github.com/mannkind/unifi2mqtt/blob/master/LICENSE.md)
[![Build Status](https://github.com/mannkind/unifi2mqtt/workflows/Main%20Workflow/badge.svg)](https://github.com/mannkind/unifi2mqtt/actions)
[![Coverage Status](https://img.shields.io/codecov/c/github/mannkind/unifi2mqtt/master.svg)](http://codecov.io/github/mannkind/unifi2mqtt?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/mannkind/unifi2mqtt)](https://goreportcard.com/report/github.com/mannkind/unifi2mqtt)

An experiment to publish device statuses from the Unifi Controller to MQTT.

## Use

The application can be locally built using `mage` or you can utilize the multi-architecture Docker image(s).

### Example

```bash
docker run \
-e UNIFI_HOST="unifi-controller.dns.name" \
-e UNIFI_PORT="8443" \
-e UNIFI_USERNAME="unifiUsername" \
-e UNIFI_PASSWORD="unifiPassword" \
-e UNIFI_DEVICEMAPPING="11:22:33:44:55:66;identifierSlug" \
-e MQTT_BROKER="tcp://localhost:1883" \
-e MQTT_DISCOVERY="true" \
mannkind/unifi2mqtt:latest
```

OR

```bash
UNIFI_HOST="unifi-controller.dns.name" \
UNIFI_PORT="8443" \
UNIFI_USERNAME="unifiUsername" \
UNIFI_PASSWORD="unifiPassword" \
UNIFI_DEVICEMAPPING="11:22:33:44:55:66;identifierSlug" \
MQTT_BROKER="tcp://localhost:1883" \
MQTT_DISCOVERY="true" \
./unifi2mqtt 
```

## Environment Variables

```bash
UNIFI_HOST              - The hostname of the controller, defaults to "unifi.local"
UNIFI_PORT              - The the port of the controller, defaults to "8443"
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
MQTT_CLIENTID           - [OPTIONAL] The clientId, defaults to ""
MQTT_BROKER             - [OPTIONAL] The MQTT broker, defaults to "tcp://mosquitto.org:1883"
MQTT_USERNAME           - [OPTIONAL] The MQTT username, default to ""
MQTT_PASSWORD           - [OPTIONAL] The MQTT password, default to ""
```
