package main

import (
	"fmt"
	"math/rand"

	"github.com/starkandwayne/eden-cli/apiclient"
	edenconfig "github.com/starkandwayne/eden-cli/config"
)

func main() {
	rand.Seed(5000)

	broker := apiclient.NewOpenServiceBrokerFromBrokerEnv(edenconfig.BrokerEnv())

	catalogResp := broker.Catalog()

	for _, service := range catalogResp.Services {
		for _, plan := range service.Plans {
			fmt.Println(service.Name, "-", plan.Name, "-", plan.Description)
		}
	}
}
