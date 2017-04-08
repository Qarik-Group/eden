package cmd

import (
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/pborman/uuid"
  "github.com/kr/pretty"
	"github.com/starkandwayne/eden-cli/apiclient"
)

// BindOpts represents the 'bind' command
type BindOpts struct {
	// TODO: store these from ProvisionOpts
	ServiceNameOrID string `short:"s" long:"service-name" description:"Service name/ID from catalog" required:"true"`
	PlanNameOrID    string `short:"p" long:"plan-name" description:"Plan name/ID from catalog (default: first)"`
}

// Execute is callback from go-flags.Commander interface
func (c BindOpts) Execute(_ []string) (err error) {
  instanceID := Opts.InstanceName
  if instanceID == "" {
    return fmt.Errorf("bind command requires --instance [NAME|GUID]")
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

  bindingID := uuid.New()
	bindingName := fmt.Sprintf("%s-%s", service.Name, bindingID)

	// TODO - store allocated bindingIDs into local DB
	bindingResp, err := broker.Bind(service.ID, plan.ID, instanceID, bindingID)
	if err != nil {
		return errwrap.Wrapf("Failed to bind to service instance {{err}}", err)
	}
	Opts.config().BindServiceInstance(instanceID, bindingID, bindingName, bindingResp.Credentials)

  fmt.Printf("%# v\n", pretty.Formatter(bindingResp))
	return
}
