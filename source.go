package main

import (
	"fmt"
	"time"

	"github.com/dim13/unifi"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

const (
	home    = "ON"
	notHome = "OFF"
)

type source struct {
	config       sourceOpts
	outgoing     chan<- sourceRep
	unifiClient  *unifi.Unifi
	unifiSite    *unifi.Site
	deviceStatus map[string]string
}

func newSource(config sourceOpts, outgoing chan<- sourceRep) *source {
	c := source{
		config:       config,
		outgoing:     outgoing,
		deviceStatus: map[string]string{},
	}

	return &c
}

func (c *source) run() {
	// Log service settings
	c.logSettings()

	// Run immediately
	c.poll()

	// Schedule additional runs
	sched := cron.New()
	sched.AddFunc(fmt.Sprintf("@every %s", c.config.LookupInterval), c.poll)
	sched.Start()
}

func (c *source) logSettings() {
	redactedPassword := ""
	if len(c.config.Password) > 0 {
		redactedPassword = "<REDACTED>"
	}

	log.WithFields(log.Fields{
		"Unifi.AwayTimeout":    c.config.AwayTimeout,
		"Unifi.LookupInterval": c.config.LookupInterval,
		"Unifi.Host":           c.config.Host,
		"Unifi.Port":           c.config.Port,
		"Unifi.Username":       c.config.Username,
		"Unifi.Password":       redactedPassword,
		"Unifi.Devices":        c.config.Devices,
	}).Info("Service Environmental Settings")
}

func (c *source) poll() {
	log.Info("Polling")
	if c.unifiClient == nil {
		log.Info("Connecting to Unifi")
		u, err := unifi.Login(c.config.Username, c.config.Password, c.config.Host, c.config.Port, c.config.Site, 5)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Fatal("Couldn't login")
			return
		}
		log.Info("Connected to Unifi")

		log.Info("Selecting Site")
		s, err := u.Site(c.config.Site)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Fatal("Couldn't find site")
			return
		}
		log.Info("Selected Site")

		c.unifiClient = u
		c.unifiSite = s
	}

	// Default the status of all known macSlugMapping to not_home
	for _, slug := range c.config.Devices {
		c.deviceStatus[slug] = notHome
	}

	// Ask the unifi controller for all known clients
	log.Debug("Fetching Clients")
	clients, err := c.unifiClient.Sta(c.unifiSite)
	if err != nil || len(clients) == 0 {
		c.unifiClient = nil
		c.unifiSite = nil

		log.WithFields(log.Fields{
			"error":      err,
			"numClients": len(clients),
		}).Error("Couldn't find any clients")

		return
	}
	log.Debug("Fetched Clients")

	// Devices known to the controller will have their status set based on the Last Seen time
	// Devices missing will remain defaulted to not_home
	for _, client := range clients {
		slug, ok := c.config.Devices[client.Mac]
		if !ok {
			log.WithFields(log.Fields{
				"mac": client.Mac,
			}).Debug("The device is not known nor cared about")
			continue
		}

		ts := time.Unix(int64(client.LastSeen), 0)
		old := time.Since(ts) >= c.config.AwayTimeout
		payload := home
		if old {
			payload = notHome
		}

		c.deviceStatus[slug] = payload
	}

	// Publish known device statuses
	for name, state := range c.deviceStatus {
		c.outgoing <- c.adapt(name, state)
	}

	log.WithFields(log.Fields{
		"sleep": c.config.LookupInterval,
	}).Info("Finished polling; sleeping")
}

func (c *source) adapt(name string, state string) sourceRep {
	return sourceRep{
		name:  name,
		state: state,
	}
}
