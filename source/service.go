package source

import (
	"fmt"
	"github.com/dim13/unifi"
	log "github.com/sirupsen/logrus"
)

// Service is for reading a directly from a source system
type Service struct {
	opts        Opts
	unifiClient *unifi.Unifi
	unifiSite   *unifi.Site
}

// NewService creates a new Service for reading a directly from a source system
func NewService() *Service {
	c := Service{
		opts: Opts{},
	}

	return &c
}

// SetOpts sets the required options to access the source system
func (c *Service) SetOpts(opts Opts) {
	c.opts = opts
}

// login to the unifi controller
func (c *Service) login() error {
	// No need to continually login
	if c.unifiClient != nil {
		return nil
	}

	log.Info("Connecting to Unifi")
	u, err := unifi.Login(c.opts.Username, c.opts.Password, c.opts.Host, c.opts.Port, c.opts.Site, 5)
	// @NOTE(mannkind) unifi.Login does not actually seem to return an error when the username/password is wrong :/
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Couldn't login to the controller")
		return err
	}
	log.Info("Connected to Unifi")

	// @NOTE(mannkind) If the username/password are wrong, this is where we'll error
	log.Info("Selecting site")
	s, err := u.Site(c.opts.Site)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Couldn't select the site specified")
		return err
	}
	log.Info("Selected site")

	c.unifiClient = u
	c.unifiSite = s

	return nil
}

// reset the store unifi client and site
func (c *Service) reset() {
	c.unifiClient = nil
	c.unifiSite = nil
}

// lookup data from the source system
func (c *Service) lookup() (*[]unifi.Sta, error) {
	err := c.login()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Unable to communicate with the Unifi Controller")
		return nil, err
	}

	// Ask the unifi controller for all "active" clients
	log.Debug("Fetching Clients")
	clients, err := c.unifiClient.Sta(c.unifiSite)
	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"numClients": len(clients),
		}).Error("Couldn't login or find clients")
		c.reset()
		return nil, fmt.Errorf("Couldn't find clients")
	}
	log.Debug("Fetched Clients")

	return &clients, nil
}
