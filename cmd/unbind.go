package cmd

import (
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/starkandwayne/eden/apiclient"
)

// UnbindOpts represents the 'unbind' command
type UnbindOpts struct {
  BindingID       string `short:"b" long:"bind" description:"Binding ID to destroy" required:"true"`
}

// Execute is callback from go-flags.Commander interface
func (c UnbindOpts) Execute(_ []string) (err error) {
	instanceNameOrID := Opts.Instance.NameOrID
	if instanceNameOrID == "" {
		return fmt.Errorf("unbind command requires --instance [NAME|GUID], or $SB_INSTANCE")
	}
	instance := Opts.config().FindServiceInstance(instanceNameOrID)
	// TODO: convert c.BindingID into ID if its a name

	broker := apiclient.NewOpenServiceBroker(
		Opts.Broker.URLOpt,
		Opts.Broker.ClientOpt,
		Opts.Broker.ClientSecretOpt,
		Opts.Broker.APIVersion,
	)
	err = broker.Unbind(instance.ServiceID, instance.PlanID, instance.ID, c.BindingID)
	if err != nil {
		return errwrap.Wrapf("Failed to unbind to service instance {{err}}", err)
	}
	Opts.config().UnbindServiceInstance(instance.ID, c.BindingID)

	fmt.Println("Success")
	return
}
