package cmd

import (
	"fmt"
	"os"

	"github.com/starkandwayne/eden-cli/apiclient"
)

// CatalogOpts represents the 'catalog' command
type CatalogOpts struct {
}

// Execute is callback from go-flags.Commander interface
func (c CatalogOpts) Execute(_ []string) (err error) {
	broker := apiclient.NewOpenServiceBroker(Opts.Broker.URLOpt, Opts.Broker.ClientOpt, Opts.Broker.ClientSecretOpt)

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
	return
}
