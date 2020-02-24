package source

import (
	"fmt"
	"time"

	"github.com/mannkind/unifi2mqtt/shared"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

const (
	home    = "ON"
	notHome = "OFF"
)

// Reader is for reading a shared representation out of a source system
type Reader struct {
	opts           Opts
	outgoing       chan<- shared.Representation
	service        *Service
	deviceAwayTime map[string]time.Time
}

// NewReader creates a new Reader for reading a shared representation out of a source system
func NewReader(opts Opts, outgoing chan<- shared.Representation, service *Service) *Reader {
	c := Reader{
		opts:           opts,
		outgoing:       outgoing,
		service:        service,
		deviceAwayTime: map[string]time.Time{},
	}

	service.SetOpts(opts)

	return &c
}

// Run starts the Reader
func (c *Reader) Run() {
	// Log service settings
	c.logSettings()

	// Run immediately
	c.poll()

	// Schedule additional runs
	sched := cron.New()
	sched.AddFunc(fmt.Sprintf("@every %s", c.opts.LookupInterval), c.poll)
	sched.Start()
}

// logSettings that are specific to reading the source system
func (c *Reader) logSettings() {
	log.WithFields(log.Fields{
		"Unifi.AwayTimeout":    c.opts.AwayTimeout,
		"Unifi.LookupInterval": c.opts.LookupInterval,
		"Unifi.Host":           c.opts.Host,
		"Unifi.Port":           c.opts.Port,
		"Unifi.Username":       c.opts.Username,
		"Unifi.Password":       "<REDACTED>",
		"Unifi.Devices":        c.opts.Devices,
	}).Info("Service Environmental Settings")
}

// poll the source system, adapt source system responses to the share representation, output data onto a channnel
func (c *Reader) poll() {
	log.Info("Polling")

	// Lookup "active" clients
	// @NOTE(mannkind) Clients remain "active" long after they disconnect
	clients, err := c.service.lookup()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Unable to find clients")
		return
	}

	// Set the Last Seen value for clients we care about
	// @NOTE(mannkind) It appears that we cannot trust client.LastSeen (anymore?)
	for _, client := range *clients {
		// Skip clients that we don't care about
		slug, ok := c.opts.Devices[client.Mac]
		if !ok {
			log.WithFields(log.Fields{
				"mac": client.Mac,
			}).Debug("The device is not known nor cared about")
			continue
		}

		// Make up our own Last Seen value
		c.deviceAwayTime[slug] = time.Now()
	}

	// Determine the status of clients we care about
	for _, slug := range c.opts.Devices {
		payload := notHome
		lastSeen, _ := c.deviceAwayTime[slug]
		if time.Now().Sub(lastSeen) < c.opts.AwayTimeout {
			payload = home
		}

		c.outgoing <- c.adapt(slug, payload)
	}

	log.WithFields(log.Fields{
		"sleep": c.opts.LookupInterval,
	}).Info("Finished polling; sleeping")
}

// adapt incoming value(s) to the shared representation
func (c *Reader) adapt(name string, state string) shared.Representation {
	return shared.Representation{
		Name:  name,
		State: state,
	}
}
