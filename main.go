package main

import (
	"fmt"

	edenconfig "github.com/starkandwayne/eden-cli/config"
)

func main() {
	fmt.Println(edenconfig.BrokerEnv())
}
