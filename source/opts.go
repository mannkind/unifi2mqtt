package source

import (
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/mannkind/unifi2mqtt/shared"
	log "github.com/sirupsen/logrus"
)

// Opts is for package related settings
type Opts struct {
	shared.Opts
	Host           string        `env:"UNIFI_HOST"                 envDefault:"unifi"`
	Port           string        `env:"UNIFI_PORT"                 envDefault:"8443"`
	Site           string        `env:"UNIFI_SITE"                 envDefault:"default"`
	Username       string        `env:"UNIFI_USERNAME"             envDefault:"unifi"`
	Password       string        `env:"UNIFI_PASSWORD"             envDefault:"unifi"`
	AwayTimeout    time.Duration `env:"UNIFI_AWAYTIMEOUT"          envDefault:"5m"`
	LookupInterval time.Duration `env:"UNIFI_LOOKUPINTERVAL"       envDefault:"10s"`
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

	return c
}
