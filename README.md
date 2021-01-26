# unifi2mqtt

[![Software
License](https://img.shields.io/badge/License-MIT-orange.svg?style=flat-square)](https://github.com/mannkind/unifi2mqtt/blob/main/LICENSE.md)
[![Build Status](https://github.com/mannkind/unifi2mqtt/workflows/Main%20Workflow/badge.svg)](https://github.com/mannkind/unifi2mqtt/actions)
[![Coverage Status](https://img.shields.io/codecov/c/github/mannkind/unifi2mqtt/main.svg)](http://codecov.io/github/mannkind/unifi2mqtt?branch=main)

An experiment to publish device statuses from the Unifi Controller to MQTT.

## Use

The application can be locally built using `dotnet build` or you can utilize the multi-architecture Docker image(s).

### Example

```bash
docker run \
-e UNIFI__HOST="https://unifi-controller.dns.name:8443" \
-e UNIFI__USERNAME="unifiUsername" \
-e UNIFI__PASSWORD="unifiPassword" \
-e UNIFI__AWAYTIMEOUT="0.00:05:01" \
-e UNIFI__RESOURCES__0__MACAddress="11:22:33:44:55:66" \
-e UNIFI__RESOURCES__0__Slug="identifierSlug" \
mannkind/unifi2mqtt:latest
```

OR

```bash
UNIFI__HOST="https://unifi-controller.dns.name:8443" \
UNIFI__USERNAME="unifiUsername" \
UNIFI__PASSWORD="unifiPassword" \
UNIFI__AWAYTIMEOUT="0.00:05:01" \
UNIFI__RESOURCES__0__MACAddress="11:22:33:44:55:66" \
UNIFI__RESOURCES__0__Slug="identifierSlug" \
./unifi2mqtt 
```


## Configuration

Configuration happens via environmental variables

```bash
UNIFI__HOST                               - The Unifi Controller Host URL
UNIFI__USERNAME                           - The Unifi Controller Username
UNIFI__PASSWORD                           - The Unifi Controller Password
UNIFI__AWAYTIMEOUT                        - [OPTIONAL] The delay between last seeing a device and marking it as away, defaults to "0.00:05:01"
UNIFI__POLLINGINTERVAL                    - [OPTIONAL] The delay between device lookups, defaults to "0.00:00:11"
UNIFI__DISABLESSLVALIDATION               - [OPTIONAL] The flag that disables SSL validation, defaults to true
UNIFI__RESOURCES__#__MACAddress           - The n-th iteration of a mac address for a specific device
UNIFI__RESOURCES__#__Slug                 - The n-th iteration of a slug to identify the specific mac address
UNIFI__MQTT__TOPICPREFIX                  - [OPTIONAL] The MQTT topic on which to publish the collection lookup results, defaults to "home/unifi"
UNIFI__MQTT__DISCOVERYENABLED             - [OPTIONAL] The MQTT discovery flag for Home Assistant, defaults to false
UNIFI__MQTT__DISCOVERYPREFIX              - [OPTIONAL] The MQTT discovery prefix for Home Assistant, defaults to "homeassistant"
UNIFI__MQTT__DISCOVERYNAME                - [OPTIONAL] The MQTT discovery name for Home Assistant, defaults to "unifi"
UNIFI__MQTT__BROKER                       - [OPTIONAL] The MQTT broker, defaults to "test.mosquitto.org"
UNIFI__MQTT__USERNAME                     - [OPTIONAL] The MQTT username, default to ""
UNIFI__MQTT__PASSWORD                     - [OPTIONAL] The MQTT password, default to ""
```

## Prior Implementations

### Golang
* Last Commit: [c39d32c5d0721d32f8ebf089b796461f514b4d71](https://github.com/mannkind/unifi2mqtt/commit/c39d32c5d0721d32f8ebf089b796461f514b4d71)
* Last Docker Image: [mannkind/unifi2mqtt:v0.8.20061.0158](https://hub.docker.com/layers/mannkind/unifi2mqtt/v0.8.20061.0158/images/sha256-7020736d44b64fe8b9cbc87887f20216b4539c32c9a5ae6145c10fe3c233b5bf?context=explore)