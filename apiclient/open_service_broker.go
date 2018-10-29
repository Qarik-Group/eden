package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/errwrap"
	"github.com/pivotal-cf/brokerapi"
)

// OpenServiceBroker is the client struct for connecting to remote Open Service Broker API
type OpenServiceBroker struct {
	url        string
	username   string
	password   string
	catalog    *brokerapi.CatalogResponse
	apiVersion string
}

// NewOpenServiceBroker constructs OpenServiceBroker
func NewOpenServiceBroker(url, client, clientSecret, apiVersion string) *OpenServiceBroker {
	return &OpenServiceBroker{
		url:        url,
		username:   client,
		password:   clientSecret,
		apiVersion: apiVersion,
	}
}

// Catalog fetches the available service catalog from remote broker
func (broker *OpenServiceBroker) Catalog() (catalogResp *brokerapi.CatalogResponse, err error) {
	if broker.catalog == nil {
		url := fmt.Sprintf("%s/v2/catalog", broker.url)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, errwrap.Wrapf("Cannot construct HTTP request: {{err}}", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Broker-Api-Version", broker.apiVersion)
		req.SetBasicAuth(broker.username, broker.password)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, errwrap.Wrapf("Failed doing HTTP request: {{err}}", err)
		}
		defer resp.Body.Close()

		resBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errwrap.Wrapf("Failed reading HTTP response body: {{err}}", err)
		}

		broker.catalog = &brokerapi.CatalogResponse{}
		err = json.Unmarshal(resBody, broker.catalog)
		if err != nil {
			return nil, errwrap.Wrapf("Failed unmarshalling catalog response: {{err}}", err)
		}
	}
	return broker.catalog, nil
}

// Provision attempts to provision a new service instance
func (broker *OpenServiceBroker) Provision(serviceID, planID, instanceID, organizationGUID, spaceGUID string) (provisioningResp *brokerapi.ProvisioningResponse, isAsync bool, err error) {
	url := fmt.Sprintf("%s/v2/service_instances/%s?accepts_incomplete=true", broker.url, instanceID)
	details := brokerapi.ProvisionDetails{
		ServiceID:        serviceID,
		PlanID:           planID,
		OrganizationGUID: organizationGUID,
		SpaceGUID:        spaceGUID,
		RawParameters:    nil,
	}

	buffer := &bytes.Buffer{}
	if err = json.NewEncoder(buffer).Encode(details); err != nil {
		return nil, false, errwrap.Wrapf("Cannot encode provisioning details: {{err}}", err)
	}
	req, err := http.NewRequest("PUT", url, buffer)
	if err != nil {
		return nil, false, errwrap.Wrapf("Cannot construct HTTP request: {{err}}", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Broker-Api-Version", broker.apiVersion)
	req.SetBasicAuth(broker.username, broker.password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, false, errwrap.Wrapf("Failed doing HTTP request: {{err}}", err)
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, false, errwrap.Wrapf("Failed reading HTTP response body: {{err}}", err)
	}
	provisioningResp = &brokerapi.ProvisioningResponse{}
	err = json.Unmarshal(resBody, provisioningResp)
	if err != nil {
		return nil, false, errwrap.Wrapf("Failed unmarshalling provisioning response: {{err}}", err)
	}
	if resp.StatusCode == http.StatusAccepted {
		isAsync = true
	} else if resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK {
		isAsync = false
	} else {
		errorResp := &brokerapi.ErrorResponse{}
		err = json.Unmarshal(resBody, errorResp)
		if err != nil {
			return nil, false, errwrap.Wrapf("Failed unmarshalling error response: {{err}}", err)
		}
		errMsg := fmt.Sprintf("%s: %s", errorResp.Error, errorResp.Description)
		if errorResp.Error == "" {
			errMsg = fmt.Sprintf("Unknown internal broker error (%s)", broker.url)
		}
		return nil, false, fmt.Errorf(errMsg)
	}
	return
}

// Bind requests new set of credentials to access service instance
func (broker *OpenServiceBroker) Bind(serviceID, planID, instanceID, bindingID string) (binding *brokerapi.Binding, err error) {
	url := fmt.Sprintf("%s/v2/service_instances/%s/service_bindings/%s", broker.url, instanceID, bindingID)
	details := brokerapi.BindDetails{
		ServiceID:     serviceID,
		PlanID:        planID,
		AppGUID:       "eden-unknown",
		RawParameters: nil,
	}

	buffer := &bytes.Buffer{}
	if err = json.NewEncoder(buffer).Encode(details); err != nil {
		return nil, errwrap.Wrapf("Cannot encode binding details: {{err}}", err)
	}
	req, err := http.NewRequest("PUT", url, buffer)
	if err != nil {
		return nil, errwrap.Wrapf("Cannot construct HTTP request: {{err}}", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Broker-Api-Version", broker.apiVersion)
	req.SetBasicAuth(broker.username, broker.password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errwrap.Wrapf("Failed doing HTTP request: {{err}}", err)
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errwrap.Wrapf("Failed reading HTTP response body: {{err}}", err)
	}
	binding = &brokerapi.Binding{}
	err = json.Unmarshal(resBody, binding)
	if err != nil {
		return nil, errwrap.Wrapf("Failed unmarshalling binding response: {{err}}", err)
	}
	return
}

// Unbind destroys a set of credentials to access the service instance
func (broker *OpenServiceBroker) Unbind(serviceID, planID, instanceID, bindingID string) (err error) {
	url := fmt.Sprintf("%s/v2/service_instances/%s/service_bindings/%s?service_id=%s&plan_id=%s",
		broker.url, instanceID, bindingID, serviceID, planID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return errwrap.Wrapf("Cannot construct HTTP request: {{err}}", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Broker-Api-Version", broker.apiVersion)
	req.SetBasicAuth(broker.username, broker.password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errwrap.Wrapf("Failed doing HTTP request: {{err}}", err)
	}
	defer resp.Body.Close()
	return
}

// Deprovision destroys the service instance
func (broker *OpenServiceBroker) Deprovision(serviceID, planID, instanceID string) (deprovisioningResp *brokerapi.DeprovisionResponse, isAsync bool, err error) {
	url := fmt.Sprintf("%s/v2/service_instances/%s?service_id=%s&plan_id=%s&accepts_incomplete=true",
		broker.url, instanceID, serviceID, planID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		isAsync = false
		err = errwrap.Wrapf("Cannot construct HTTP request: {{err}}", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Broker-Api-Version", broker.apiVersion)
	req.SetBasicAuth(broker.username, broker.password)

	client := &http.Client{}
	resp, err := client.Do(req)
	deprovisioningResp = &brokerapi.DeprovisionResponse{}

	resBody, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(resBody, deprovisioningResp)

	if err != nil {
		isAsync = false
		err = errwrap.Wrapf("Failed doing HTTP request: {{err}}", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusAccepted {
		isAsync = true
	}
	return
}

// LastOperation fetches the status of the last operation perform upon a service instance
func (broker *OpenServiceBroker) LastOperation(serviceID, planID, instanceID, operation string) (lastOpResp *brokerapi.LastOperationResponse, err error) {
	url := fmt.Sprintf("%s/v2/service_instances/%s/last_operation?operation=%s&service_id=%s&plan_id=%s", broker.url, instanceID, operation, serviceID, planID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errwrap.Wrapf("Cannot construct HTTP request: {{err}}", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Broker-Api-Version", broker.apiVersion)
	req.SetBasicAuth(broker.username, broker.password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errwrap.Wrapf("Failed doing HTTP request: {{err}}", err)
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errwrap.Wrapf("Failed reading HTTP response body: {{err}}", err)
	}

	lastOpResp = &brokerapi.LastOperationResponse{}
	err = json.Unmarshal(resBody, lastOpResp)
	if err != nil {
		lastOpResp.Description = fmt.Sprintf("Failed to unmarshal last operation response, assuming it has succeeded: %s\n", err)
		lastOpResp.State = brokerapi.Succeeded
		err = nil
	}

	return
}

// FindServiceByNameOrID looks thru all services in catalog for one that has
// a name or ID matching 'nameOrID'
func (broker *OpenServiceBroker) FindServiceByNameOrID(nameOrID string) (*brokerapi.Service, error) {
	catalog, err := broker.Catalog()
	if err != nil {
		return nil, errwrap.Wrapf("Could not fetch catalog: {{err}}", err)
	}
	for _, service := range catalog.Services {
		if service.ID == nameOrID || service.Name == nameOrID {
			return &service, nil
		}
	}
	return nil, fmt.Errorf("No service has name or ID '%s'", nameOrID)
}

// FindPlanByNameOrID looks thru all plans for a service for one that has
// a name or ID matching 'nameOrID'. Defaults to first plan if 'nameOrID' is empty.
func (broker *OpenServiceBroker) FindPlanByNameOrID(service *brokerapi.Service, nameOrID string) (*brokerapi.ServicePlan, error) {
	if nameOrID == "" {
		return &service.Plans[0], nil
	}
	for _, plan := range service.Plans {
		if plan.ID == nameOrID || plan.Name == nameOrID {
			return &plan, nil
		}
	}
	return nil, fmt.Errorf("No plan has name or ID '%s' within service '%s'", nameOrID, service.Name)
}
