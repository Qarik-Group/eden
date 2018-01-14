package cmd

import (
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/pborman/uuid"
	"github.com/starkandwayne/eden/apiclient"
)

// BindOpts represents the 'bind' command
type BindOpts struct {
}

// Execute is callback from go-flags.Commander interface
func (c BindOpts) Execute(_ []string) (err error) {
	instanceNameOrID := Opts.Instance.NameOrID
	if instanceNameOrID == "" {
		return fmt.Errorf("bind command requires --instance [NAME|GUID], or $SB_INSTANCE")
	}
	instance := Opts.config().FindServiceInstance(instanceNameOrID)

	broker := apiclient.NewOpenServiceBroker(
		Opts.Broker.URLOpt,
		Opts.Broker.ClientOpt,
		Opts.Broker.ClientSecretOpt,
		Opts.Broker.APIVersion,
	)

	bindingID := uuid.New()
	bindingName := fmt.Sprintf("%s-%s", instance.ServiceName, bindingID)

	// TODO - store allocated bindingIDs into local DB
	bindingResp, err := broker.Bind(instance.ServiceID, instance.PlanID, instance.ID, bindingID)
	if err != nil {
		return errwrap.Wrapf("Failed to bind to service instance {{err}}", err)
	}
	err = Opts.config().BindServiceInstance(instance.ID, bindingID, bindingName, bindingResp.Credentials)
	if err != nil {
		return errwrap.Wrapf("Failed to store binding {{err}}", err)
	}

	fmt.Println("Success")
	fmt.Println("")
	fmt.Printf("Run 'eden credentials -i %s -b %s' to see credentials\n", instance.Name, bindingName)
	return
}
