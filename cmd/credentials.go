package cmd

import (
  "encoding/json"
  "fmt"
  "strings"
  "os"

  "github.com/hashicorp/errwrap"
)

// CredentialsOpts represents the 'credentials' command
type CredentialsOpts struct {
  BindingID string `short:"b" long:"bind" description:"Binding to display"`
  Attribute string `short:"a" long:"attribute" description:"Only diplay a single attribute from credentials"`
}

// Execute is callback from go-flags.Commander interface
func (c CredentialsOpts) Execute(_ []string) (err error) {
  instanceNameOrID := Opts.Instance.NameOrID
  if instanceNameOrID == "" {
    return fmt.Errorf("credentials command requires --instance [NAME|GUID], or $EDEN_INSTANCE")
  }
  inst := Opts.config().FindServiceInstance(instanceNameOrID)
  if inst.ServiceID == "" {
    return fmt.Errorf("credentials --instance [NAME|GUID] was not found")
  }
  if len(inst.Bindings) > 0 {
    binding := inst.Bindings[0]

    // convert binding.Credentials into nested map[string]map[string]interface{}
    if err := c.displayBinding(binding.CredentialsJSON(), c.Attribute); err != nil {
      return err
    }
  } else {
    fmt.Println("No bindings.")
  }
  return
}

func (c CredentialsOpts) displayBinding(credentials map[string]interface{}, attribute string) error {
  if attribute == "" {
    b, err := json.MarshalIndent(credentials, "", "  ")
    if err != nil {
      return errwrap.Wrapf("Could not marshal credentials: {{err}}", err)
    }
    os.Stdout.Write(b)
    return nil
  }
  if val, ok := credentials[attribute]; ok {
    fmt.Printf("%v", val)
    return nil
  }
  attributes := make([]string, 0, len(credentials))
  for key := range credentials {
    attributes = append(attributes, key)
  }

  return fmt.Errorf("credentials --attribute key was unknown; try: %s", strings.Join(attributes, ", "))
}
