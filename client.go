package main

import (
	"log"
	"strings"
	"time"

	"github.com/dim13/unifi"
)

const (
	home    = "ON"
	notHome = "OFF"
)

type client struct {
	observers map[observer]struct{}

	host           string
	port           string
	site           string
	username       string
	password       string
	awayTimeout    time.Duration
	lookupInterval time.Duration
	macSlugMapping map[string]string
	deviceStatus   map[string]string
}

func newClient(config *config) *client {
	c := client{
		observers: map[observer]struct{}{},

		awayTimeout:    config.AwayTimeout,
		lookupInterval: config.LookupInterval,
		host:           config.Host,
		port:           config.Port,
		site:           config.Site,
		username:       config.Username,
		password:       config.Password,
	}
	c.macSlugMapping = make(map[string]string, 0)
	c.deviceStatus = make(map[string]string, 0)

	// Create a mapping between mac_addresses and names
	for _, m := range config.DeviceMapping {
		parts := strings.Split(m, ";")
		if len(parts) != 2 {
			continue
		}

		deviceMacAddress := parts[0]
		deviceName := parts[1]
		c.macSlugMapping[deviceMacAddress] = deviceName
		c.deviceStatus[deviceName] = notHome
	}

	return &c
}

func (c *client) run() {
	go c.loop()
}

func (c *client) register(l observer) {
	c.observers[l] = struct{}{}
}

func (c *client) publish(e event) {
	for o := range c.observers {
		o.receive(e)
	}
}

func (c *client) loop() {
	log.Print("Connecting to Unifi")
	u, err := unifi.Login(c.username, c.password, c.host, c.port, c.site, 5)
	if err != nil {
		log.Fatalf("Couldn't login to %s:%s with %s:<redacted>; %s", c.host, c.port, c.username, err)
		return
	}
	defer u.Logout()
	log.Print("Connected to Unifi")

	log.Print("Selecting Site")
	site, err := u.Site(c.site)
	if err != nil {
		log.Fatalf("Couldn't find site %s; %s", c.site, err)
		return
	}
	log.Print("Selected Site")

	log.Print("Forever Fetching Clients")
	for {
		// Default the status of all known macSlugMapping to not_home
		for _, slug := range c.macSlugMapping {
			c.deviceStatus[slug] = notHome
		}

		// Ask the unifi controller for all known clients
		// log.Print("Fetching Clients")
		clients, err := u.Sta(site)
		if err != nil {
			log.Fatalf("Couldn't find clients; %s", err)
			continue
		}
		// log.Print("Fetched Clients")

		// Devices known to the controller will have their status set based on the Last Seen time
		// Devices missing will remain defaulted to not_home
		for _, client := range clients {
			slug, ok := c.macSlugMapping[client.Mac]
			if !ok {
				// log.Printf("%s is not a device that is cared about", client.Mac)
				continue
			}

			ts := time.Unix(int64(client.LastSeen), 0)
			old := time.Since(ts) > c.awayTimeout
			payload := home
			if old {
				payload = notHome
			}

			c.deviceStatus[slug] = payload
		}

		// Publish known device statuses
		for slug, payload := range c.deviceStatus {
			c.publish(event{
				version: 1,
				key:     slug,
				data:    payload,
			})
		}

		time.Sleep(c.lookupInterval)
	}
}
