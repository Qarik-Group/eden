package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/errwrap"
	"github.com/pivotal-cf/brokerapi"
	edenconfig "github.com/starkandwayne/eden-cli/config"
)

// OpenServiceBroker is the client struct for connecting to remote Open Service Broker API
type OpenServiceBroker struct {
	url      string
	username string
	password string
}

// NewOpenServiceBrokerFromBrokerEnv constructs OpenServiceBroker
func NewOpenServiceBrokerFromBrokerEnv(brokerEnv *edenconfig.BrokerEnvDiscovery) *OpenServiceBroker {
	return &OpenServiceBroker{
		url:      brokerEnv.URL,
		username: brokerEnv.Username,
		password: brokerEnv.Password,
	}
}

// Catalog fetches the available service catalog from remote broker
func (broker *OpenServiceBroker) Catalog() (catalogResp *brokerapi.CatalogResponse, err error) {
	url := fmt.Sprintf("%s/v2/catalog", broker.url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errwrap.Wrapf("Cannot construct HTTP request: {{err}}", err)
	}
	req.Header.Set("Content-Type", "application/json")
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

	catalogResp = &brokerapi.CatalogResponse{}
	err = json.Unmarshal(resBody, catalogResp)
	if err != nil {
		return nil, errwrap.Wrapf("Failed unmarshalling catalog response: {{err}}", err)
	}
	return
}

// Provision attempts to provision a new service instance
func (broker *OpenServiceBroker) Provision(serviceID, planID, instanceID string) (provisioningResp *brokerapi.ProvisioningResponse, err error) {
	url := fmt.Sprintf("%s/v2/service_instances/%s", broker.url, instanceID)
	details := brokerapi.ProvisionDetails{
		ServiceID:        serviceID,
		PlanID:           planID,
		OrganizationGUID: "eden-unknown-guid",
		SpaceGUID:        "eden-unknown-space",
		RawParameters:    nil,
	}

	buffer := &bytes.Buffer{}
	if err = json.NewEncoder(buffer).Encode(details); err != nil {
		return nil, errwrap.Wrapf("Cannot encode provisioning details: {{err}}", err)
	}
	req, err := http.NewRequest("PUT", url, buffer)
	if err != nil {
		return nil, errwrap.Wrapf("Cannot construct HTTP request: {{err}}", err)
	}
	req.Header.Set("Content-Type", "application/json")
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
	provisioningResp = &brokerapi.ProvisioningResponse{}
	err = json.Unmarshal(resBody, provisioningResp)
	if err != nil {
		return nil, errwrap.Wrapf("Failed unmarshalling provisioning response: {{err}}", err)
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
func (broker *OpenServiceBroker) Deprovision(serviceID, planID, instanceID string) (err error) {
	url := fmt.Sprintf("%s/v2/service_instances/%s?service_id=%s&plan_id=%s",
		broker.url, instanceID, serviceID, planID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return errwrap.Wrapf("Cannot construct HTTP request: {{err}}", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(broker.username, broker.password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errwrap.Wrapf("Failed doing HTTP request: {{err}}", err)
	}
	defer resp.Body.Close()
	return
}
