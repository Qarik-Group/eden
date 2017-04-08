package cmd

import (
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/starkandwayne/eden-cli/apiclient"
)

// UnbindOpts represents the 'unbind' command
type UnbindOpts struct {
	// TODO: store these from ProvisionOpts
	ServiceNameOrID string `short:"s" long:"service-name" description:"Service name/ID from catalog" required:"true"`
	PlanNameOrID    string `short:"p" long:"plan-name" description:"Plan name/ID from catalog (default: first)"`
  BindingID       string `short:"b" long:"bind" description:"Binding ID to destroy" required:"true"`
}

// Execute is callback from go-flags.Commander interface
func (c UnbindOpts) Execute(_ []string) (err error) {
  instanceID := Opts.InstanceName
  if instanceID == "" {
    return fmt.Errorf("unbind command requires --instance [NAME|GUID]")
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

	err = broker.Unbind(service.ID, plan.ID, instanceID, c.BindingID)
	if err != nil {
		return errwrap.Wrapf("Failed to unbind to service instance {{err}}", err)
	}
	Opts.config().UnbindServiceInstance(instanceID, c.BindingID)

  fmt.Println("Success")
	return
}
