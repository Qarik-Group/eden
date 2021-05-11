package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/pborman/uuid"
	"github.com/starkandwayne/eden/apiclient"
)

// BindOpts represents the 'bind' command
type BindOpts struct {
	Parameters      string `short:"P" long:"parameters" description:"parameters in json format. To use a file as input, prepend the filename with '@' (-P=@data.json)"`
}

// Execute is callback from go-flags.Commander interface
func (c BindOpts) Execute(_ []string) (err error) {
	instanceNameOrID := Opts.Instance.NameOrID
	if instanceNameOrID == "" {
		return fmt.Errorf("bind command requires --instance [NAME|GUID], or $SB_INSTANCE")
	}
	instance := Opts.config().FindServiceInstance(instanceNameOrID)
	bindingID := Opts.Binding.ID
	if bindingID == "" {
		bindingID = uuid.New()
	}

	broker := apiclient.NewOpenServiceBroker(
		Opts.Broker.URLOpt,
		Opts.Broker.ClientOpt,
		Opts.Broker.ClientSecretOpt,
		Opts.Broker.APIVersion,
	)

	bindingName := fmt.Sprintf("%s-%s", instance.ServiceName, bindingID)

	var parameters json.RawMessage
	if len(c.Parameters) > 0 {
		var input []byte
		if strings.HasPrefix(c.Parameters, "@") {
			input, err = ioutil.ReadFile(c.Parameters[1:])
			if err != nil {
				return errwrap.Wrapf("Could not read file: {{err}}", err)
			}
		} else {
			input = []byte(c.Parameters)
		}
		if err := json.Unmarshal(input, &parameters); err != nil {
			return errwrap.Wrapf("Could not unmarshal parameters: {{err}}", err)
		}
	}
	bindingResp, err := broker.Bind(instance.ServiceID, instance.PlanID, instance.ID, bindingID, parameters)
	if err != nil {
		return errwrap.Wrapf("Failed to bind to service instance {{err}}", err)
	}
	err = Opts.config().BindServiceInstance(instance.ID, bindingID, bindingName, bindingResp.Credentials)
	if err != nil {
		return errwrap.Wrapf("Failed to store binding {{err}}", err)
	}

	if Opts.JSON {
		var out struct {
			Instance    interface{} `json:"instance"`
			Binding     interface{} `json:"binding"`
			BindingID   string      `json:"binding_id"`
			BindingName string      `json:"binding_name"`
		}
		out.Instance = instance
		out.Binding = bindingResp
		out.BindingID = bindingID
		out.BindingName = bindingName
		b, err := json.Marshal(out)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", string(b))
		os.Exit(0)
	}
	fmt.Println("Success")
	fmt.Println("")
	fmt.Printf("Run 'eden credentials -i %s -b %s' to see credentials\n", instance.Name, bindingName)
	return
}
