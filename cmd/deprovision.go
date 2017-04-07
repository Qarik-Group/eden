package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/pivotal-cf/brokerapi"
	"github.com/starkandwayne/eden-cli/apiclient"
)

// DeprovisionOpts represents the 'deprovision' command
type DeprovisionOpts struct {
	ServiceNameOrID string `short:"s" long:"service-name" description:"Service name/ID from catalog" required:"true"`
	PlanNameOrID    string `short:"p" long:"plan-name" description:"Plan name/ID from catalog (default: first)"`
}

// Execute is callback from go-flags.Commander interface
func (c DeprovisionOpts) Execute(_ []string) (err error) {
  instanceID := Opts.InstanceName
  if instanceID == "" {
    return fmt.Errorf("deprovision command requires --instance [NAME|GUID]")
  }

	broker := apiclient.NewOpenServiceBroker(Opts.Broker.URLOpt, Opts.Broker.ClientOpt, Opts.Broker.ClientSecretOpt)

	service, err := broker.FindServiceByNameOrID(c.ServiceNameOrID)
  if err != nil {
		return errwrap.Wrapf("Could not find service in catalog: {{err}}", err)
	}
	plan, err := broker.FindPlanByNameOrID(service, c.PlanNameOrID)
	if err != nil {
		return errwrap.Wrapf("Could not find plan in service: {{err}}", err)
	}

	isAsync, err := broker.Deprovision(service.ID, plan.ID, instanceID)
	if err != nil {
		return errwrap.Wrapf("Failed to deprovision service instance {{err}}", err)
	}
	// TODO - update local DB with status

	fmt.Printf("deprovision: %s/%s - guid: %s\n", service.Name, plan.Name, instanceID)
	if isAsync {
		fmt.Println("deprovision: in-progress")
		// TODO: don't pollute brokerapi back into this level
		lastOpResp := &brokerapi.LastOperationResponse{State: brokerapi.InProgress}
		for lastOpResp.State == brokerapi.InProgress {
			time.Sleep(5 * time.Second)
			lastOpResp, err = broker.LastOperation(service.ID, plan.ID, instanceID)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}
			fmt.Printf("deprovision: %s - %s\n", lastOpResp.State, lastOpResp.Description)
		}
	}
	fmt.Println("deprovision: done")

	return
}
