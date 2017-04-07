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
	// fmt.Printf("%#v\n", opts)

	// instanceID := uuid.New()
	// bindingID := uuid.New()
	//
	// time.Sleep(1 * time.Second)
	// // TODO - store allocated instanceID into local DB
	// provisioningResp, isAsync, err := broker.Provision(serviceID, planID, instanceID)
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, err.Error())
	// 	os.Exit(1)
	// }
	// // TODO - update local DB with status
	//
	// if isAsync {
	// 	fmt.Println("provision:   in-progress")
	// 	// TODO: don't pollute brokerapi back into this level
	// 	lastOpResp := &brokerapi.LastOperationResponse{State: brokerapi.InProgress}
	// 	for lastOpResp.State == brokerapi.InProgress {
	// 		time.Sleep(5 * time.Second)
	// 		lastOpResp, err = broker.LastOperation(serviceID, planID, instanceID)
	// 		if err != nil {
	// 			fmt.Fprintln(os.Stderr, err.Error())
	// 			os.Exit(1)
	// 		}
	// 		fmt.Printf("  - %s: %s\n", lastOpResp.State, lastOpResp.Description)
	// 	}
	// }
	// fmt.Printf("provision:   %v\n", provisioningResp)
	//
	// time.Sleep(1 * time.Second)
	// // TODO - store allocated bindingID into local DB
	// bindingResp, err := broker.Bind(serviceID, planID, instanceID, bindingID)
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, err.Error())
	// 	os.Exit(1)
	// }
	// // TODO - update local DB with status
	//
	// fmt.Printf("binding:     %v\n", bindingResp.Credentials)
	//
	// time.Sleep(1 * time.Second)
	// err = broker.Unbind(serviceID, planID, instanceID, bindingID)
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, err.Error())
	// 	os.Exit(1)
	// }
	// fmt.Println("unbinding:   done")
	//
	// time.Sleep(1 * time.Second)
	// isAsync, err = broker.Deprovision(serviceID, planID, instanceID)
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, err.Error())
	// 	os.Exit(1)
	// }
	//
	// if isAsync {
	// 	fmt.Println("deprovision: in-progress")
	// 	// TODO: don't pollute brokerapi back into this level
	// 	lastOpResp := &brokerapi.LastOperationResponse{State: brokerapi.InProgress}
	// 	for lastOpResp.State == brokerapi.InProgress {
	// 		lastOpResp, err = broker.LastOperation(serviceID, planID, instanceID)
	// 		time.Sleep(5 * time.Second)
	// 		if err != nil {
	// 			fmt.Fprintln(os.Stderr, err.Error())
	// 			os.Exit(1)
	// 		}
	// 		fmt.Printf("  - %s: %s\n", lastOpResp.State, lastOpResp.Description)
	// 	}
	// }
	// fmt.Println("deprovision: done")
}
