package main

import (
	"fmt"
	"math/rand"

	edenconfig "github.com/starkandwayne/eden-cli/config"
)

func main() {
	rand.Seed(5000)
	fmt.Println(edenconfig.BrokerEnv())
}
