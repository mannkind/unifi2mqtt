package main

import (
	"time"

	"github.com/caarlos0/env"
	mqttExtCfg "github.com/mannkind/paho.mqtt.golang.ext/cfg"
	log "github.com/sirupsen/logrus"
)

type config struct {
	MQTT           *mqttExtCfg.MQTTConfig
	Host           string        `env:"UNIFI_HOST"                 envDefault:"unifi"`
	Port           string        `env:"UNIFI_PORT"                 envDefault:"8443"`
	Site           string        `env:"UNIFI_SITE"                 envDefault:"default"`
	Username       string        `env:"UNIFI_USERNAME"             envDefault:"unifi"`
	Password       string        `env:"UNIFI_PASSWORD"             envDefault:"unifi"`
	AwayTimeout    time.Duration `env:"UNIFI_AWAYTIMEOUT"          envDefault:"5m"`
	LookupInterval time.Duration `env:"UNIFI_LOOKUPINTERVAL"       envDefault:"10s"`
	DeviceMapping  []string      `env:"UNIFI_DEVICEMAPPING"        envDefault:"11:22:33:44:55:66;MyPhone,12:23:34:45:56:67;AnotherPhone"`
	DebugLogLevel  bool          `env:"UNIFI_DEBUG" envDefault:"false"`
}

func newConfig(mqttCfg *mqttExtCfg.MQTTConfig) *config {
	c := config{}
	c.MQTT = mqttCfg

	if c.MQTT.ClientID == "" {
		c.MQTT.ClientID = "DefaultUnifi2MqttClientID"
	}

	if c.MQTT.DiscoveryName == "" {
		c.MQTT.DiscoveryName = "unifi"
	}

	if c.MQTT.TopicPrefix == "" {
		c.MQTT.TopicPrefix = "home/unifi"
	}

	if err := env.Parse(&c); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Unable to unmarshal configuration")
	}

	redactedPassword := ""
	if len(c.MQTT.Password) > 0 || len(c.Password) > 0 {
		redactedPassword = "<REDACTED>"
	}

	log.WithFields(log.Fields{
		"MQTT.ClientID":        c.MQTT.ClientID,
		"MQTT.Broker":          c.MQTT.Broker,
		"MQTT.Username":        c.MQTT.Username,
		"MQTT.Password":        redactedPassword,
		"MQTT.Discovery":       c.MQTT.Discovery,
		"MQTT.DiscoveryPrefix": c.MQTT.DiscoveryPrefix,
		"MQTT.DiscoveryName":   c.MQTT.DiscoveryName,
		"MQTT.TopicPrefix":     c.MQTT.TopicPrefix,
		"Unifi.AwayTimeout":    c.AwayTimeout,
		"Unifi.LookupInterval": c.LookupInterval,
		"Unifi.Host":           c.Host,
		"Unifi.Port":           c.Port,
		"Unifi.Username":       c.Username,
		"Unifi.Password":       redactedPassword,
		"Unifi.DebugLogLevel":  c.DebugLogLevel,
	}).Info("Environmental Settings")

	if c.DebugLogLevel {
		log.SetLevel(log.DebugLevel)
		log.Debug("Enabling the debug log level")
	}

	return &c
}
