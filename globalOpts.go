package main

type globalOpts struct {
	Devices sourceMapping `env:"UNIFI_DEVICEMAPPING" envDefault:"11:22:33:44:55:66;MyPhone,12:23:34:45:56:67;AnotherPhone"`
}
