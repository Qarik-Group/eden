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
	Strict bool `long:"strict" description:"Validate the catalog using the same heuristics as Cloud Foundry" env:"EDEN_STRICT"`
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

	if Opts.Catalog.Strict {
		errors := make([]error, 0)
		seen := make(map[string] bool)
		for _, service := range catalogResp.Services {
			for _, plan := range service.Plans {
				if _, exists := seen[plan.ID]; exists {
					errors = append(errors, fmt.Errorf("Service '%s' Plan '%s' ID '%s' is not unique", service.Name, plan.Name, plan.ID))
				}
				seen[plan.ID] = true
			}
		}

		if len(errors) != 0 {
			fmt.Fprintf(os.Stderr, "Catalog validation failed:\n")
			for _, err := range errors {
				fmt.Fprintf(os.Stderr, "  - %s\n", err)
			}
			os.Exit(1)
		}
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
