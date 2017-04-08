package cmd

import (
  "fmt"

  "github.com/kr/pretty"
)

// ServicesOpts represents the 'services' command
type ServicesOpts struct {
}

// Execute is callback from go-flags.Commander interface
func (c ServicesOpts) Execute(_ []string) (err error) {
  instances := Opts.config().ServiceInstances()
  fmt.Printf("%# v\n", pretty.Formatter(instances))
  return
}
