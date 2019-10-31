package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	sourceVMResourceID      = "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/ResourceGroup/providers/Microsoft.Compute/virtualMachines/sourceVM"
	targetVMResourceID      = "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/ResourceGroup/providers/Microsoft.Compute/virtualMachines/targetVM"
	port               uint = 80
)

func TestUN_NewConnectivyRequestBody(t *testing.T) {

	sourceVM, err := NewVMDetails(sourceVMResourceID)
	assert.Nil(t, err)
	assert.NotNil(t, sourceVM)

	targetVM, err := NewVMDetails(targetVMResourceID)
	assert.Nil(t, err)
	assert.NotNil(t, sourceVM)

	requestBody := NewConnectivyRequestBody(*sourceVM, *targetVM, port)
	assert.NotEqual(t, connectivityRequestBody{}, requestBody)

	assert.Equal(t, sourceVM.GetResourceID(), requestBody.GetSourceResourceID())
	assert.Equal(t, targetVM.GetResourceID(), requestBody.GetDestinationResourceID())
	assert.Equal(t, port, requestBody.GetDestinationPort())

}
