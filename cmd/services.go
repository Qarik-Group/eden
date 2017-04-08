package cmd

import (
  "fmt"
)

// ServicesOpts represents the 'services' command
type ServicesOpts struct {
}

// Execute is callback from go-flags.Commander interface
func (c ServicesOpts) Execute(_ []string) (err error) {
  instanceNameOrID := Opts.InstanceName
  if instanceNameOrID != "" {
    return c.showService(instanceNameOrID)
  }
  return c.showAllServices()
}

func (c ServicesOpts) showAllServices() (err error) {
  instances := Opts.config().ServiceInstances()
  for _, inst := range instances {
    bindingName := "n/a"
    if len(inst.Bindings) > 0 {
      bindingName = inst.Bindings[0].Name
    }
    fmt.Printf("%s\t%s\t%s\t%s\n", inst.Name,
      inst.ServiceName, inst.PlanName,
      bindingName)
  }
  return
}

func (c ServicesOpts) showService(instanceNameOrID string) (err error) {
  inst := Opts.config().FindServiceInstance(instanceNameOrID)
  if inst.ServiceID == "" {
    return fmt.Errorf("services --instance [NAME|GUID] was not found")
  }
  bindingName := "n/a"
  if len(inst.Bindings) > 0 {
    bindingName = inst.Bindings[0].Name
  }
  fmt.Printf("%s\t%s\t%s\t%s\n", inst.Name,
    inst.ServiceName, inst.PlanName,
    bindingName)
  return
}
