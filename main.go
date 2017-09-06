package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/jessevdk/go-flags"
	edencmd "github.com/starkandwayne/eden/cmd"
)

var Version = ""

func main() {
	rand.Seed(5000)

	if len(os.Args) > 1 {
		if os.Args[1] == "-v" || os.Args[1] == "--version" {
			if Version == "" {
				fmt.Printf("eden (development)\n")
			} else {
				fmt.Printf("eden v%s\n", Version)
			}
			os.Exit(0)
		}
	}

	parser := flags.NewParser(&edencmd.Opts, flags.Default)

	if len(os.Args) == 1 {
		_, err := parser.ParseArgs([]string{"--help"})
		if err != nil {
			os.Exit(1)
		}
	} else {
		_, err := parser.Parse()
		if err != nil {
			os.Exit(1)
		}
	}
}
