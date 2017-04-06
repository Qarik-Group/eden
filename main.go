package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/pborman/uuid"
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

	var serviceID string
	var planID string
	for _, service := range catalogResp.Services {
		if serviceID == "" {
			serviceID = service.ID
		}
		for _, plan := range service.Plans {
			if planID == "" {
				planID = plan.ID
			}
			fmt.Println(service.Name, "-", plan.Name, "-", plan.Description)
		}
	}

	instanceID := uuid.New()
	bindingID := uuid.New()

	provisioningResp, err := broker.Provision(serviceID, planID, instanceID)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Printf("%#v\n", provisioningResp)

	bindingResp, err := broker.Bind(serviceID, planID, instanceID, bindingID)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Printf("%#v\n", bindingResp)

	err = broker.Unbind(serviceID, planID, instanceID, bindingID)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	err = broker.Deprovision(serviceID, planID, instanceID)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

}
