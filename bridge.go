package main

type bridge struct {
	mqttClient *mqttClient
	client     *client
}

func newBridge(config *config, mqttClient *mqttClient, client *client) *bridge {
	bridge := bridge{
		mqttClient: mqttClient,
		client:     client,
	}

	return &bridge
}

func (b *bridge) run() {
	b.client.register(b.mqttClient)

	b.mqttClient.run()
	b.client.run()
}
