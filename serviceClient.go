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

type serviceClient struct {
	serviceClientConfig
	stateUpdateChan stateChannel
	unifiClient     *unifi.Unifi
	unifiSite       *unifi.Site
	deviceStatus    map[string]string
}

func newServiceClient(serviceClientCfg serviceClientConfig, stateUpdateChan stateChannel) *serviceClient {
	c := serviceClient{
		serviceClientConfig: serviceClientCfg,
		stateUpdateChan:     stateUpdateChan,
		deviceStatus:        map[string]string{},
	}

	redactedPassword := ""
	if len(c.Password) > 0 {
		redactedPassword = "<REDACTED>"
	}

	log.WithFields(log.Fields{
		"Unifi.AwayTimeout":    c.AwayTimeout,
		"Unifi.LookupInterval": c.LookupInterval,
		"Unifi.Host":           c.Host,
		"Unifi.Port":           c.Port,
		"Unifi.Username":       c.Username,
		"Unifi.Password":       redactedPassword,
		"Unifi.Devices":        c.Devices,
	}).Info("Service Environmental Settings")

	return &c
}

func (c *serviceClient) run() {
	// Run immediately
	go c.loop()

	// Schedule additional runs
	sched := cron.New()
	sched.AddFunc(fmt.Sprintf("@every %s", c.LookupInterval), c.loop)
	sched.Start()
}

func (c *serviceClient) loop() {
	log.Info("Looping")
	if c.unifiClient == nil {
		log.Info("Connecting to Unifi")
		u, err := unifi.Login(c.Username, c.Password, c.Host, c.Port, c.Site, 5)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Fatal("Couldn't login")
			return
		}
		log.Info("Connected to Unifi")

		log.Info("Selecting Site")
		s, err := u.Site(c.Site)
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
	for _, slug := range c.Devices {
		c.deviceStatus[slug] = notHome
	}

	// Ask the unifi controller for all known clients
	log.Debug("Fetching Clients")
	clients, err := c.unifiClient.Sta(c.unifiSite)
	if err != nil {
		c.unifiClient = nil
		c.unifiSite = nil

		log.WithFields(log.Fields{
			"error": err,
		}).Error("Couldn't find any clients")

		return
	}
	log.Debug("Fetched Clients")

	// Devices known to the controller will have their status set based on the Last Seen time
	// Devices missing will remain defaulted to not_home
	for _, client := range clients {
		slug, ok := c.Devices[client.Mac]
		if !ok {
			log.WithFields(log.Fields{
				"mac": client.Mac,
			}).Debug("The device is not known nor cared about")
			continue
		}

		ts := time.Unix(int64(client.LastSeen), 0)
		old := time.Since(ts) >= c.AwayTimeout
		payload := home
		if old {
			payload = notHome
		}

		c.deviceStatus[slug] = payload
	}

	// Publish known device statuses
	for name, state := range c.deviceStatus {
		obj, err := c.adapt(name, state)
		if err != nil {
			continue
		}

		c.stateUpdateChan <- obj
	}

	log.WithFields(log.Fields{
		"sleep": c.LookupInterval,
	}).Info("Finished looping; sleeping")
}

func (c *serviceClient) login() (*unifi.Unifi, error) {
	log.Info("Connecting to Unifi")
	u, err := unifi.Login(c.Username, c.Password, c.Host, c.Port, c.Site, 5)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Couldn't login")
		return nil, err
	}
	log.Info("Connected to Unifi")

	return u, err
}

func (c *serviceClient) adapt(name string, state string) (unifiDevice, error) {
	log.WithFields(log.Fields{
		"name":  name,
		"state": state,
	}).Debug("Adapting name/state information")

	obj := unifiDevice{
		name:  name,
		state: state,
	}

	log.Debug("Finished adapting name/state information")
	return obj, nil
}
