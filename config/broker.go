package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

// BrokerEnvDiscovery describes a target Open Service Broker API
type BrokerEnvDiscovery struct {
	URL      string `envconfig:"url"`
	Username string `envconfig:"client"`
	Password string `envconfig:"client_secret"`
}

var brokerEnv *BrokerEnvDiscovery

// BrokerEnv describes a target Open Service Broker API via environment variables
func BrokerEnv() *BrokerEnvDiscovery {
	if brokerEnv == nil {
		brokerEnv = &BrokerEnvDiscovery{}
		err := envconfig.Process("eden_broker", brokerEnv)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	return brokerEnv
}
