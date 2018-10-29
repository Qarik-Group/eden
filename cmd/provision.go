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

	SpaceGUID        string `long:"space-guid" description:"Explicitly provide a 'space_guid' provision field"`
	OrganizationGUID string `long:"organization-guid" description:"Explicitly provide a 'organization_guid' provision field"`
}

// Execute is callback from go-flags.Commander interface
func (c ProvisionOpts) Execute(_ []string) (err error) {
	broker := apiclient.NewOpenServiceBroker(
		Opts.Broker.URLOpt,
		Opts.Broker.ClientOpt,
		Opts.Broker.ClientSecretOpt,
		Opts.Broker.APIVersion,
	)

	if c.OrganizationGUID == "" {
		c.OrganizationGUID = "eden-unknown-org"
	}
	if c.SpaceGUID == "" {
		c.SpaceGUID = "eden-unknown-space"
	}

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

	provisioningResp, isAsync, err := broker.Provision(
		service.ID, plan.ID, instanceID,
		c.OrganizationGUID, c.SpaceGUID)
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
