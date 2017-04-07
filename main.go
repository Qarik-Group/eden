package main

import (
	"math/rand"
	"os"

	"github.com/jessevdk/go-flags"
	edencmd "github.com/starkandwayne/eden-cli/cmd"
)

func main() {
	rand.Seed(5000)

	parser := flags.NewParser(&edencmd.Opts, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		os.Exit(1)
	}
}
