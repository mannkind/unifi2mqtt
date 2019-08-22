module github.com/mannkind/unifi2mqtt

require (
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/dim13/unifi v0.0.0-20190324114524-e97746adf746
	github.com/eclipse/paho.mqtt.golang v1.2.0
	github.com/google/wire v0.3.0
	github.com/mannkind/paho.mqtt.golang.ext v0.3.0
	github.com/sirupsen/logrus v1.4.2
)

// local development
// replace github.com/mannkind/paho.mqtt.golang.ext => ../paho.mqtt.golang.ext
