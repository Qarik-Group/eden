package apiclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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
func (broker *OpenServiceBroker) Catalog() (catalogResp *brokerapi.CatalogResponse) {
	url := fmt.Sprintf("%s/v2/catalog", broker.url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(broker.username, broker.password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	catalogResp = &brokerapi.CatalogResponse{}

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(resBody, catalogResp)
	if err != nil {
		panic(err)
	}
	return
}
