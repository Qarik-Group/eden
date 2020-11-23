package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
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
	Parameters      string `short:"P" long:"parameters" description:"parameters in json format. To use a file as input, prepend the filename with '@' (-P=@data.json)"`
}

// Execute is callback from go-flags.Commander interface
func (c ProvisionOpts) Execute(_ []string) (err error) {
	broker := apiclient.NewOpenServiceBroker(
		Opts.Broker.URLOpt,
		Opts.Broker.ClientOpt,
		Opts.Broker.ClientSecretOpt,
		Opts.Broker.APIVersion,
	)

	service, err := broker.FindServiceByNameOrID(c.ServiceNameOrID)
	if err != nil {
		return errwrap.Wrapf("Could not find service in catalog: {{err}}", err)
	}
	plan, err := broker.FindPlanByNameOrID(service, c.PlanNameOrID)
	if err != nil {
		return errwrap.Wrapf("Could not find plan in service: {{err}}", err)
	}

	instanceName := Opts.Instance.NameOrID
	instanceID := uuid.New()
	if instanceName == "" {
		instanceName = fmt.Sprintf("%s-%s-%s", service.Name, plan.Name, instanceID)
	}
	prexisting := Opts.config().FindServiceInstance(instanceName)
	if prexisting.ServiceName != "" {
		return fmt.Errorf("Service instance '%s' already exists", instanceName)
	}

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
	provisioningResp, isAsync, err := broker.Provision(service.ID, plan.ID, instanceID, parameters)
	if err != nil {
		return errwrap.Wrapf("Failed to provision service instance: {{err}}", err)
	}
	Opts.config().ProvisionNewServiceInstance(instanceID, instanceName,
		service.ID, service.Name,
		plan.ID, plan.Name,
		Opts.Broker.URLOpt)

	fmt.Printf("provision:   %s/%s - name: %s\n", service.Name, plan.Name, instanceName)
	if isAsync {
		fmt.Println("provision:   in-progress")
		// TODO: don't pollute brokerapi back into this level
		lastOpResp := &brokerapi.LastOperationResponse{State: brokerapi.InProgress}
		for lastOpResp.State == brokerapi.InProgress {
			time.Sleep(5 * time.Second)
			lastOpResp, err = broker.LastOperation(service.ID, plan.ID, instanceID, provisioningResp.OperationData)
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
