package models

import (
	"errors"
	"fmt"
	"strings"
)

//VMDetails struct that holds the information of a VM within an Azure resource id
type VMDetails struct {
	resourceID     string
	resourceGroup  string
	subscriptionID string
	region         string
	name           string
}

const (
	defaultRegion         = "eastus"
	emptyResourceID       = "VM resource ID is empty!"
	unformattedResourceID = "The resourceID: %s is unformatted!"
)

//NewVMDetails the resourceid returning a vmDetails structure if success`
func NewVMDetails(vmResourceID string) (*VMDetails, error) {
	vm := &VMDetails{}

	if len(vmResourceID) == 0 {
		return nil, errors.New(emptyResourceID)
	}

	if strings.Index(vmResourceID, "/subscriptions/") < 0 ||
		strings.Index(vmResourceID, "/resourceGroups/") < 0 ||
		strings.Index(vmResourceID, "/providers/Microsoft.Compute/virtualMachines/") < 0 {
		return nil, fmt.Errorf(unformattedResourceID, vmResourceID)
	}

	vm.resourceID = vmResourceID
	vm.subscriptionID = getSubscriptionID(vmResourceID)
	vm.resourceGroup = getResourceGroupName(vmResourceID)
	vm.name = getVMName(vmResourceID)
	vm.region = defaultRegion

	if len(vm.subscriptionID) == 0 || len(vm.resourceGroup) == 0 || len(vm.name) == 0 {
		return nil, fmt.Errorf(unformattedResourceID, vmResourceID)
	}

	return vm, nil

}

func getSubscriptionID(virtualMachineID string) string {
	beginIndex := strings.Index(virtualMachineID, "/subscriptions/") + len("/subscriptions/")
	endIndex := strings.Index(virtualMachineID, "/resourceGroups/")
	splited := virtualMachineID[beginIndex:endIndex]
	return splited
}

func getVMName(virtualMachineID string) string {
	lastIndexOfSlash := strings.LastIndex(virtualMachineID, "/") + 1
	return virtualMachineID[lastIndexOfSlash:]
}

func getResourceGroupName(virtualMachineID string) string {
	beginIndex := strings.Index(virtualMachineID, "/resourceGroups/") + len("/resourceGroups/")
	endIndex := strings.Index(virtualMachineID, "/providers/")

	return virtualMachineID[beginIndex:endIndex]
}

//GetResourceID return the resourceID
func (vm *VMDetails) GetResourceID() string {
	return vm.resourceID
}

//GetResourceGroupName return the resource group
func (vm *VMDetails) GetResourceGroupName() string {
	return vm.resourceGroup
}

//GetVMName return the VM Name
func (vm *VMDetails) GetVMName() string {
	return vm.name
}

//GetSubscriptionID return the subscription from id
func (vm *VMDetails) GetSubscriptionID() string {
	return vm.subscriptionID
}

//GetRegion return the region
func (vm *VMDetails) GetRegion() string {
	return vm.region
}

//SetRegion updates the region with the information retrieved from azure
func (vm *VMDetails) SetRegion(region string) {
	vm.region = region
}
