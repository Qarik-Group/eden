package cmd

import (
  "encoding/json"
  "fmt"
  "os"

  "github.com/hashicorp/errwrap"
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
  fmt.Printf("%s\t%s\t%s\n", inst.Name,
    inst.ServiceName, inst.PlanName)
  if len(inst.Bindings) > 0 {
    binding := inst.Bindings[0]

    b, err := json.MarshalIndent(binding.Credentials, "", "  ")
  	if err != nil {
  		return errwrap.Wrapf("Could not marshal credentials: {{err}}", err)
  	}
    os.Stdout.Write(b)
  } else {
    fmt.Println("No bindings.")
  }
  return
}
