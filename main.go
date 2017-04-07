package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/pborman/uuid"
	"github.com/pivotal-cf/brokerapi"
	"github.com/starkandwayne/eden-cli/apiclient"
)

// BrokerOpts describes subset of flags/options for selecting target service broker API
type BrokerOpts struct {
	URLOpt          string `long:"url"           description:"Open Service Broker URL"                env:"EDEN_BROKER_URL" required:"true"`
	ClientOpt       string `long:"client"        description:"Override username or UAA client"        env:"EDEN_BROKER_CLIENT" required:"true"`
	ClientSecretOpt string `long:"client-secret" description:"Override password or UAA client secret" env:"EDEN_BROKER_CLIENT_SECRET" required:"true"`
}

// EdenOpts describes the flags/options for the CLI
type EdenOpts struct {
	// Slice of bool will append 'true' each time the option
	// is encountered (can be set multiple times, like -vvv)
	Verbose []bool `short:"v" long:"verbose" description:"Show verbose debug information"                   env:"EDEN_TRACE"`

	// Example of a value name
	ServiceName string `short:"s" long:"service" description:"Service instance name"                        env:"EDEN_SERVICE"`

	Broker BrokerOpts `group:"Broker Options"`
}

func main() {
	rand.Seed(5000)

	var opts EdenOpts
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}
	fmt.Printf("%#v\n", opts)

	broker := apiclient.NewOpenServiceBroker(opts.Broker.URLOpt, opts.Broker.ClientOpt, opts.Broker.ClientSecretOpt)

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

	time.Sleep(1 * time.Second)
	// TODO - store allocated instanceID into local DB
	provisioningResp, isAsync, err := broker.Provision(serviceID, planID, instanceID)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	// TODO - update local DB with status

	if isAsync {
		fmt.Println("provision:   in-progress")
		// TODO: don't pollute brokerapi back into this level
		lastOpResp := &brokerapi.LastOperationResponse{State: brokerapi.InProgress}
		for lastOpResp.State == brokerapi.InProgress {
			time.Sleep(5 * time.Second)
			lastOpResp, err = broker.LastOperation(serviceID, planID, instanceID)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}
			fmt.Printf("  - %s: %s\n", lastOpResp.State, lastOpResp.Description)
		}
	}
	fmt.Printf("provision:   %v\n", provisioningResp)

	time.Sleep(1 * time.Second)
	// TODO - store allocated bindingID into local DB
	bindingResp, err := broker.Bind(serviceID, planID, instanceID, bindingID)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	// TODO - update local DB with status

	fmt.Printf("binding:     %v\n", bindingResp.Credentials)

	time.Sleep(1 * time.Second)
	err = broker.Unbind(serviceID, planID, instanceID, bindingID)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	fmt.Println("unbinding:   done")

	time.Sleep(1 * time.Second)
	isAsync, err = broker.Deprovision(serviceID, planID, instanceID)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if isAsync {
		fmt.Println("deprovision: in-progress")
		// TODO: don't pollute brokerapi back into this level
		lastOpResp := &brokerapi.LastOperationResponse{State: brokerapi.InProgress}
		for lastOpResp.State == brokerapi.InProgress {
			lastOpResp, err = broker.LastOperation(serviceID, planID, instanceID)
			time.Sleep(5 * time.Second)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}
			fmt.Printf("  - %s: %s\n", lastOpResp.State, lastOpResp.Description)
		}
	}
	fmt.Println("deprovision: done")
}
