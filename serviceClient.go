package main

import (
	"time"

	"github.com/dim13/unifi"
	log "github.com/sirupsen/logrus"
)

const (
	home    = "ON"
	notHome = "OFF"
)

type serviceClient struct {
	serviceClientConfig
	stateUpdateChan stateChannel
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
	go c.loop()
}

func (c *serviceClient) loop() {
	for {
		log.Info("Connecting to Unifi")
		u, err := unifi.Login(c.Username, c.Password, c.Host, c.Port, c.Site, 5)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Fatal("Couldn't login")
			return
		}
		defer u.Logout()
		log.Info("Connected to Unifi")

		log.Info("Selecting Site")
		site, err := u.Site(c.Site)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Fatal("Couldn't find site")
			return
		}
		log.Info("Selected Site")

		log.Info("Forever Fetching Clients")
		for {
			// Default the status of all known macSlugMapping to not_home
			for _, slug := range c.Devices {
				c.deviceStatus[slug] = notHome
			}

			// Ask the unifi controller for all known clients
			log.Debug("Fetching Clients")
			clients, err := u.Sta(site)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Error("Couldn't find any clients")
				break
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
				old := time.Since(ts) > c.AwayTimeout
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
			time.Sleep(c.LookupInterval)
		}
	}
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
