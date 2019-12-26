module github.com/mannkind/unifi2mqtt

go 1.13

require (
	github.com/caarlos0/env/v6 v6.0.0
	github.com/dim13/unifi v0.0.0-20191217043323-6a9899784288
	github.com/google/wire v0.4.0
	github.com/magefile/mage v1.9.0
	github.com/mannkind/twomqtt v0.4.6
	github.com/robfig/cron/v3 v3.0.0
	github.com/sirupsen/logrus v1.4.2
)

// local development
// replace github.com/mannkind/twomqtt => ../twomqtt
