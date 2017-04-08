package cmd

import (
  "fmt"
)

// ServicesOpts represents the 'services' command
type ServicesOpts struct {
}

// Execute is callback from go-flags.Commander interface
func (c ServicesOpts) Execute(_ []string) (err error) {
  instances := Opts.config().ServiceInstances()
  for _, inst := range instances {
    bindingName := "n/a"
    if len(inst.Bindings) > 0 {
      bindingName = inst.Bindings[0].Name
    }
    fmt.Printf("%s\t%s\t%s\t%s", inst.Name,
      inst.ServiceName, inst.PlanName,
      bindingName)
  }
  return
}
