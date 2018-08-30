package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jhunt/go-table"
	"github.com/starkandwayne/eden/apiclient"
)

// CatalogOpts represents the 'catalog' command
type CatalogOpts struct {
}

// Execute is callback from go-flags.Commander interface
func (c CatalogOpts) Execute(_ []string) (err error) {
	broker := apiclient.NewOpenServiceBroker(
		Opts.Broker.URLOpt,
		Opts.Broker.ClientOpt,
		Opts.Broker.ClientSecretOpt,
		Opts.Broker.APIVersion,
	)

	catalogResp, err := broker.Catalog()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if Opts.JSON {
		b, err := json.Marshal(catalogResp)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		fmt.Printf("%s\n", string(b))
		os.Exit(0)
	}

	table := table.NewTable("Service", "Plan", "Description")

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
			/* FIXME service descriptions are ignored */
			table.Row(nil, service.Name, plan.Name, plan.Description)
		}
	}

	table.Output(os.Stdout)
	return
}
