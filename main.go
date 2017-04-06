package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/pborman/uuid"
	"github.com/starkandwayne/eden-cli/apiclient"
	edenconfig "github.com/starkandwayne/eden-cli/config"
)

// EdenOpts describes the flags/options for the CLI
type EdenOpts struct {
	// Slice of bool will append 'true' each time the option
	// is encountered (can be set multiple times, like -vvv)
	Verbose []bool `short:"v" long:"verbose" description:"Show verbose debug information"                   env:"EDEN_TRACE"`

	// Example of a value name
	ServiceName string `short:"s" long:"service" description:"Service instance name"                        env:"EDEN_SERVICE"`

	BrokerURLOpt          string `long:"url"           description:"Open Service Broker URL"                env:"EDEN_BROKER_URL"`
	BrokerClientOpt       string `long:"client"        description:"Override username or UAA client"        env:"EDEN_BROKER_CLIENT"`
	BrokerClientSecretOpt string `long:"client-secret" description:"Override password or UAA client secret" env:"EDEN_BROKER_CLIENT_SECRET"`
}

func main() {
	rand.Seed(5000)

	// var opts EdenOpts
	// args, err := flags.Parse(&opts)
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, err.Error())
	// 	os.Exit(1)
	// }
	// fmt.Printf("%#v\n", args)

	// TODO: replace with fetching same data from "args" above
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

	// TODO - store allocated instanceID into local DB
	provisioningResp, isAsync, err := broker.Provision(serviceID, planID, instanceID)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	// TODO - update local DB with status

	fmt.Printf("provision: %#v\n", provisioningResp)
	fmt.Printf("provision is async = %v\n", isAsync)
	if isAsync {
		lastOpResp, err2 := broker.LastOperation(serviceID, planID, instanceID)
		if err2 != nil {
			fmt.Fprintln(os.Stderr, err2.Error())
			os.Exit(1)
		}
		fmt.Println(lastOpResp.State, lastOpResp.Description)
	}

	// TODO - store allocated bindingID into local DB
	bindingResp, err := broker.Bind(serviceID, planID, instanceID, bindingID)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	// TODO - update local DB with status

	fmt.Printf("binding: %#v\n", bindingResp)

	err = broker.Unbind(serviceID, planID, instanceID, bindingID)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	fmt.Println("unbinding: done")

	isAsync, err = broker.Deprovision(serviceID, planID, instanceID)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	fmt.Printf("deprovision is async = %v\n", isAsync)

	if isAsync {
		lastOpResp, err2 := broker.LastOperation(serviceID, planID, instanceID)
		if err2 != nil {
			fmt.Fprintln(os.Stderr, err2.Error())
			os.Exit(1)
		}
		fmt.Println(lastOpResp.State, lastOpResp.Description)
	}
	fmt.Println("deprovision: done")
}
