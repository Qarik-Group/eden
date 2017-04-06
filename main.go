package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/starkandwayne/eden-cli/apiclient"
	edenconfig "github.com/starkandwayne/eden-cli/config"
)

func main() {
	rand.Seed(5000)

	broker := apiclient.NewOpenServiceBrokerFromBrokerEnv(edenconfig.BrokerEnv())

	catalogResp, err := broker.Catalog()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	for _, service := range catalogResp.Services {
		for _, plan := range service.Plans {
			fmt.Println(service.Name, "-", plan.Name, "-", plan.Description)
		}
	}
}
