package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/pborman/uuid"
	"github.com/pivotal-cf/brokerapi"
	"github.com/starkandwayne/eden/apiclient"
)

// ProvisionOpts represents the 'provision' command
type ProvisionOpts struct {
	ServiceNameOrID string `short:"s" long:"service-name" description:"Service name/ID from catalog" required:"true"`
	PlanNameOrID    string `short:"p" long:"plan-name" description:"Plan name/ID from catalog (default: first)"`
}

// Execute is callback from go-flags.Commander interface
func (c ProvisionOpts) Execute(_ []string) (err error) {
	broker := apiclient.NewOpenServiceBroker(Opts.Broker.URLOpt, Opts.Broker.ClientOpt, Opts.Broker.ClientSecretOpt)

	service, err := broker.FindServiceByNameOrID(c.ServiceNameOrID)
  if err != nil {
		return errwrap.Wrapf("Could not find service in catalog: {{err}}", err)
	}
	plan, err := broker.FindPlanByNameOrID(service, c.PlanNameOrID)
	if err != nil {
		return errwrap.Wrapf("Could not find plan in service: {{err}}", err)
	}
	instanceID := uuid.New()
  name := fmt.Sprintf("%s-%s-%s", service.Name, plan.Name, instanceID)
  Opts.config().ProvisionNewServiceInstance(instanceID, name,
    service.ID, service.Name,
    plan.ID, plan.Name,
    Opts.Broker.URLOpt)

	// TODO - store allocated instanceID into local DB
	provisioningResp, isAsync, err := broker.Provision(service.ID, plan.ID, instanceID)
	if err != nil {
		return errwrap.Wrapf("Failed to provision service instance {{err}}", err)
	}
	// TODO - update local DB with status

	fmt.Printf("provision:   %s/%s - guid: %s\n", service.Name, plan.Name, instanceID)
	if isAsync {
		fmt.Println("provision:   in-progress")
		// TODO: don't pollute brokerapi back into this level
		lastOpResp := &brokerapi.LastOperationResponse{State: brokerapi.InProgress}
		for lastOpResp.State == brokerapi.InProgress {
			time.Sleep(5 * time.Second)
			lastOpResp, err = broker.LastOperation(service.ID, plan.ID, instanceID)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}
			fmt.Printf("provision:   %s - %s\n", lastOpResp.State, lastOpResp.Description)
		}
	}
	if provisioningResp.DashboardURL == "" {
		fmt.Println("provision:   done")
	} else {
		fmt.Printf("provision:   done - %s\n", provisioningResp.DashboardURL)
	}

	return
}
