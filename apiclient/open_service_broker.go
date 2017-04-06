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
	}

	buffer := &bytes.Buffer{}
	if err = json.NewEncoder(buffer).Encode(details); err != nil {
		return nil, errwrap.Wrapf("Cannot construct HTTP request: {{err}}", err)
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
