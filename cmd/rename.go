package cmd

import (
	"fmt"
)

// RenameOpts represents the 'rename' command
type RenameOpts struct {
}

// Execute is callback from go-flags.Commander interface
func (c RenameOpts) Execute(args []string) (err error) {
	if len(args) != 1 {
		return fmt.Errorf("USAGE: eden rename -i [old-name] [new-name]")
	}
	newName := args[0]

	instanceNameOrID := Opts.Instance.NameOrID
	if instanceNameOrID == "" {
		return fmt.Errorf("rename command requires --instance [NAME|GUID], or $SB_INSTANCE")
	}
	inst := Opts.config().FindServiceInstance(instanceNameOrID)
	if inst.ServiceID == "" {
		return fmt.Errorf("rename --instance '%s' was not found", instanceNameOrID)
	}
	fmt.Printf("Renaming '%s' to '%s'\n", inst.Name, newName)
	Opts.config().RenameServiceInstance(instanceNameOrID, newName)
	return
}
