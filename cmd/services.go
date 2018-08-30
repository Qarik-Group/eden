package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jhunt/go-table"
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
	if Opts.JSON {
		b, err := json.Marshal(Opts.config().ServiceInstances())
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		fmt.Printf("%s\n", string(b))
		os.Exit(0)
	}

	table := table.NewTable("Name", "Service", "Plan", "Binding", "Broker URL")

	instances := Opts.config().ServiceInstances()
	for _, inst := range instances {
		bindingName := "n/a"
		if len(inst.Bindings) > 0 {
			bindingName = inst.Bindings[0].Name
		}
		table.Row(nil, inst.Name, inst.ServiceName, inst.PlanName, bindingName, inst.BrokerURL)
	}
	table.Output(os.Stdout)
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
