package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/pivotal-cf/brokerapi"
	"github.com/starkandwayne/eden/apiclient"
)

// DeprovisionOpts represents the 'deprovision' command
type DeprovisionOpts struct {
}

// Execute is callback from go-flags.Commander interface
func (c DeprovisionOpts) Execute(_ []string) (err error) {
  instanceNameOrID := Opts.Instance.NameOrID
  if instanceNameOrID == "" {
    return fmt.Errorf("deprovision command requires --instance [NAME|GUID]")
  }
	instance := Opts.config().FindServiceInstance(instanceNameOrID)

	broker := apiclient.NewOpenServiceBroker(Opts.Broker.URLOpt, Opts.Broker.ClientOpt, Opts.Broker.ClientSecretOpt)
	isAsync, err := broker.Deprovision(instance.ServiceID, instance.PlanID, instance.ID)
	if err != nil {
		return errwrap.Wrapf("Failed to deprovision service instance {{err}}", err)
	}

	fmt.Printf("deprovision: %s/%s - guid: %s\n", instance.ServiceName, instance.PlanName, instance.ID)
	if isAsync {
		fmt.Println("deprovision: in-progress")
		// TODO: don't pollute brokerapi back into this level
		lastOpResp := &brokerapi.LastOperationResponse{State: brokerapi.InProgress}
		for lastOpResp.State == brokerapi.InProgress {
			time.Sleep(5 * time.Second)
			lastOpResp, err = broker.LastOperation(instance.ServiceID, instance.PlanID, instance.ID)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}
			fmt.Printf("deprovision: %s - %s\n", lastOpResp.State, lastOpResp.Description)
		}
	}
	Opts.config().DeprovisionServiceInstance(instance.ID)
	fmt.Println("deprovision: done")

	return
}
