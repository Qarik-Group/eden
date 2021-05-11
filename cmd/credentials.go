package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

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
		return fmt.Errorf("credentials command requires --instance [NAME|GUID], or $SB_INSTANCE")
	}
	inst := Opts.config().FindServiceInstance(instanceNameOrID)
	if inst.ServiceID == "" {
		return fmt.Errorf("credentials --instance '%s' was not found", instanceNameOrID)
	}
	if len(inst.Bindings) > 0 {
		binding_idx := inst.FindServiceBinding(c.BindingID)
		if binding_idx == -1 {
			return fmt.Errorf("binding '%s' was not found for instance '%s'", c.BindingID, instanceNameOrID)
		}
		binding := inst.Bindings[binding_idx]

		// convert binding.Credentials into nested map[string]map[string]interface{}
		credentialsJSON, err := binding.CredentialsJSON()
		if err != nil {
			return err
		}
		if err := c.displayBinding(credentialsJSON, c.Attribute); err != nil {
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
		fmt.Printf("%s\n", string(b))
		return nil
	}
	if val, ok := credentials[attribute]; ok {
		fmt.Printf("%v\n", val)
		return nil
	}
	attributes := make([]string, 0, len(credentials))
	for key := range credentials {
		attributes = append(attributes, key)
	}

	return fmt.Errorf("credentials --attribute key was unknown; try: %s", strings.Join(attributes, ", "))
}
