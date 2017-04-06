package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"

	"github.com/pivotal-cf/brokerapi"

	edenconfig "github.com/starkandwayne/eden-cli/config"
)

func main() {
	rand.Seed(5000)

	broker := edenconfig.BrokerEnv()

	client := &http.Client{}

	url := fmt.Sprintf("%s/v2/catalog", broker.URL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(broker.Username, broker.Password)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	catalogResp := &brokerapi.CatalogResponse{}

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(resBody, catalogResp)
	if err != nil {
		panic(err)
	}

	for _, service := range catalogResp.Services {
		for _, plan := range service.Plans {
			fmt.Println(service.Name, "-", plan.Name, "-", plan.Description)
		}
	}
}
