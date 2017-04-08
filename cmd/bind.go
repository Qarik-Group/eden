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
}

// Execute is callback from go-flags.Commander interface
func (c BindOpts) Execute(_ []string) (err error) {
  instanceNameOrID := Opts.InstanceName
  if instanceNameOrID == "" {
    return fmt.Errorf("bind command requires --instance [NAME|GUID]")
  }
	broker := apiclient.NewOpenServiceBroker(Opts.Broker.URLOpt, Opts.Broker.ClientOpt, Opts.Broker.ClientSecretOpt)

	instance := Opts.config().FindServiceInstance(instanceNameOrID)

  bindingID := uuid.New()
	bindingName := fmt.Sprintf("%s-%s", instance.ServiceName, bindingID)

	// TODO - store allocated bindingIDs into local DB
	bindingResp, err := broker.Bind(instance.ServiceID, instance.PlanID, instance.ID, bindingID)
	if err != nil {
		return errwrap.Wrapf("Failed to bind to service instance {{err}}", err)
	}
	Opts.config().BindServiceInstance(instance.ID, bindingID, bindingName, bindingResp.Credentials)

  fmt.Printf("%# v\n", pretty.Formatter(bindingResp))
	return
}
