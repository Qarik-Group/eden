package cmd

import (
	"fmt"
	"os"

	boshtbl "github.com/cloudfoundry/bosh-cli/ui/table"
)

// ServicesOpts represents the 'services' command
type ServicesOpts struct {
}

// Execute is callback from go-flags.Commander interface
func (c ServicesOpts) Execute(_ []string) (err error) {
	instanceNameOrID := Opts.Instance.NameOrID
	if instanceNameOrID != "" {
		return c.showService(instanceNameOrID)
	}
	return c.showAllServices()
}

func (c ServicesOpts) showAllServices() (err error) {
	table := boshtbl.Table{
		Content: "services",

		Header: []boshtbl.Header{
			boshtbl.NewHeader("Name"),
			boshtbl.NewHeader("Service Name"),
			boshtbl.NewHeader("Plan Name"),
			boshtbl.NewHeader("Binding Name"),
			boshtbl.NewHeader("Broker URL"),
		},

		SortBy: []boshtbl.ColumnSort{
			{Column: 1, Asc: true},
		},
	}

	instances := Opts.config().ServiceInstances()
	for _, inst := range instances {
		bindingName := "n/a"
		if len(inst.Bindings) > 0 {
			bindingName = inst.Bindings[0].Name
		}
		table.Rows = append(table.Rows, []boshtbl.Value{
			boshtbl.NewValueString(inst.Name),
			boshtbl.NewValueString(inst.ServiceName),
			boshtbl.NewValueString(inst.PlanName),
			boshtbl.NewValueString(bindingName),
			boshtbl.NewValueString(inst.BrokerURL),
		})

	}
	table.Print(os.Stdout)
	return
}

func (c ServicesOpts) showService(instanceNameOrID string) (err error) {
	inst := Opts.config().FindServiceInstance(instanceNameOrID)
	if inst.ServiceID == "" {
		return fmt.Errorf("services --instance '%s' was not found", instanceNameOrID)
	}
	fmt.Printf("Instance Name: %s\n", inst.Name)
	fmt.Printf("Service/Plan:  %s/%s\n", inst.ServiceName, inst.PlanName)
	if len(inst.Bindings) > 0 {
		fmt.Println("Bindings:")
		for _, binding := range inst.Bindings {
			fmt.Printf("- %s\n", binding.Name)
		}
	} else {
		fmt.Println("No bindings.")
	}
	return
}
