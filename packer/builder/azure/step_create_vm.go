// Copyright (c) Microsoft Open Technologies, Inc.
// All Rights Reserved.
// Licensed under the Apache License, Version 2.0.
// See License.txt in the project root for license information.
package azure

import (
	"fmt"

	"github.com/MSOpenTech/packer-azure/packer/builder/azure/constants"
	"github.com/MSOpenTech/packer-azure/packer/builder/azure/retry"

	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"

	"github.com/Azure/azure-sdk-for-go/management"
	vm "github.com/Azure/azure-sdk-for-go/management/virtualmachine"
)

type StepCreateVm struct{}

func (*StepCreateVm) Run(state multistep.StateBag) multistep.StepAction {
	client := state.Get(constants.RequestManager).(management.Client)
	ui := state.Get("ui").(packer.Ui)
	config := state.Get(constants.Config).(*Config)

	errorMsg := "Error Creating Temporary Azure VM: %s"

	ui.Say("Creating Temporary Azure VM...")

	role := state.Get("role").(*vm.Role)

	options := vm.CreateDeploymentOptions{}
	if config.VNet != "" && config.Subnet != "" {
		options.VirtualNetworkName = config.VNet
	}

	if err := retry.ExecuteAsyncOperation(client, func() (management.OperationID, error) {
		return vm.NewClient(client).CreateDeployment(*role, config.tmpServiceName, options)
	}); err != nil {
		err := fmt.Errorf(errorMsg, err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	state.Put(constants.VmExists, 1)
	state.Put(constants.DiskExists, 1)

	return multistep.ActionContinue
}

func (*StepCreateVm) Cleanup(multistep.StateBag) {}
