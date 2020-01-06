package mqtt

import (
	"github.com/caarlos0/env/v6"
	"github.com/mannkind/twomqtt"
	"github.com/mannkind/unifi2mqtt/shared"
	log "github.com/sirupsen/logrus"
)

// Opts is for package related settings
type Opts struct {
	shared.Opts
	MQTTOpts twomqtt.MQTTOpts
}

// NewOpts creates a Opts based on environment variables
func NewOpts(opts shared.Opts) Opts {
	c := Opts{
		Opts: opts,
	}

	if err := env.Parse(&c); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Unable to unmarshal configuration")
	}

	// Defaults
	if c.MQTTOpts.DiscoveryName == "" {
		c.MQTTOpts.DiscoveryName = "unifi"
	}

	if c.MQTTOpts.TopicPrefix == "" {
		c.MQTTOpts.TopicPrefix = "home/unifi"
	}

	return c
}
