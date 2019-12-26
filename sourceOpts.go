package main

import (
	"time"
)

type sourceOpts struct {
	globalOpts
	Host           string        `env:"UNIFI_HOST"                 envDefault:"unifi"`
	Port           string        `env:"UNIFI_PORT"                 envDefault:"8443"`
	Site           string        `env:"UNIFI_SITE"                 envDefault:"default"`
	Username       string        `env:"UNIFI_USERNAME"             envDefault:"unifi"`
	Password       string        `env:"UNIFI_PASSWORD"             envDefault:"unifi"`
	AwayTimeout    time.Duration `env:"UNIFI_AWAYTIMEOUT"          envDefault:"5m"`
	LookupInterval time.Duration `env:"UNIFI_LOOKUPINTERVAL"       envDefault:"10s"`
}
