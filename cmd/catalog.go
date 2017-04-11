package cmd

import (
	"fmt"
	"os"

	boshtbl "github.com/cloudfoundry/bosh-cli/ui/table"
	"github.com/starkandwayne/eden/apiclient"
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

	table := boshtbl.Table{
		Content: "services",

		Header: []boshtbl.Header{
			boshtbl.NewHeader("Service Name"),
			boshtbl.NewHeader("Plan Name"),
			boshtbl.NewHeader("Description"),
		},

		SortBy: []boshtbl.ColumnSort{
			{Column: 1, Asc: true},
		},
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
			table.Rows = append(table.Rows, []boshtbl.Value{
				boshtbl.NewValueString(service.Name),
				boshtbl.NewValueString(plan.Name),
				boshtbl.NewValueString(plan.Description),
			})
		}
	}

	table.Print(os.Stdout)
	return
}
