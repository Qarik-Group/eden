package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/pivotal-cf/brokerapi"
	"github.com/starkandwayne/eden/apiclient"
)

// UpdateOpts represents the 'update' command
type UpdateOpts struct {
	Parameters      string `short:"P" long:"parameters" description:"parameters in json format"`
}

// Execute is callback from go-flags.Commander interface
func (c UpdateOpts) Execute(_ []string) (err error) {
	instanceNameOrID := Opts.Instance.NameOrID
	if instanceNameOrID == "" {
		return fmt.Errorf("update command requires --instance [NAME|GUID], or $SB_INSTANCE")
	}
	instance := Opts.config().FindServiceInstance(instanceNameOrID)

	broker := apiclient.NewOpenServiceBroker(
		Opts.Broker.URLOpt,
		Opts.Broker.ClientOpt,
		Opts.Broker.ClientSecretOpt,
		Opts.Broker.APIVersion,
	)

	var parameters json.RawMessage
	if len(c.Parameters) > 0 {
		var input []byte
		if strings.HasPrefix(c.Parameters, "@") {
			input, err = ioutil.ReadFile(c.Parameters[1:])
			if err != nil {
				return errwrap.Wrapf("Could not read file: {{err}}", err)
			}
		} else {
			input = []byte(c.Parameters)
		}
		if err := json.Unmarshal(input, &parameters); err != nil {
			return errwrap.Wrapf("Could not unmarshal parameters: {{err}}", err)
		}
	}
	updateResp, isAsync, err := broker.Update(instance.ServiceID, instance.PlanID, instance.ID, parameters)
	if err != nil {
		return errwrap.Wrapf("Failed to update service instance {{err}}", err)
	}

	fmt.Printf("update:   name: %s\n", instanceNameOrID)
	if isAsync {
		fmt.Println("update:   in-progress")
		// TODO: don't pollute brokerapi back into this level
		lastOpResp := &brokerapi.LastOperationResponse{State: brokerapi.InProgress}
		for lastOpResp.State == brokerapi.InProgress {
			time.Sleep(5 * time.Second)
			lastOpResp, err = broker.LastOperation(instance.ServiceID, instance.PlanID, instance.ID, updateResp.OperationData)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}
			fmt.Printf("update:   %s - %s\n", lastOpResp.State, lastOpResp.Description)
		}
	}

	fmt.Println("update:   done")
	return
}
