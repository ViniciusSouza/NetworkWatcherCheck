package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	goodResourceID            = "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/ResourceGroup/providers/Microsoft.Compute/virtualMachines/agentvm"
	badResourceID             = "resourceGroups/ResourceGroup/providers/Microsoft.Compute/virtualMachines/agentvm"
	expectedSubscriptionID    = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	expectedResourceGroupName = "ResourceGroup"
	expectedVMName            = "agentvm"
)

func TestUN_ParseGoodResourceID(t *testing.T) {
	vmdetails, err := NewVMDetails(goodResourceID)
	assert.Nil(t, err)
	assert.NotEmpty(t, vmdetails.GetVMName())
}

func TestUN_FaildParseResourceID(t *testing.T) {
	_, err := NewVMDetails(badResourceID)
	assert.NotNil(t, err)
}

func TestUN_SubscriptionID(t *testing.T) {
	vmdetails, err := NewVMDetails(goodResourceID)
	assert.Nil(t, err)
	assert.Equal(t, vmdetails.GetSubscriptionID(), expectedSubscriptionID)
}

func TestUN_VMName(t *testing.T) {
	vmdetails, err := NewVMDetails(goodResourceID)
	assert.Nil(t, err)
	assert.Equal(t, vmdetails.GetVMName(), expectedVMName)
}
func TestUN_ResourceGroupName(t *testing.T) {
	vmdetails, err := NewVMDetails(goodResourceID)
	assert.Nil(t, err)
	assert.Equal(t, vmdetails.GetResourceGroupName(), expectedResourceGroupName)
}
