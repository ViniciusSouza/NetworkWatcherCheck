package networkworkcheck

import (
	"NetworkWatcherCheck/pkg/interfaces"
	"NetworkWatcherCheck/pkg/models"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	provider interfaces.NetworkWatcherProviderCheck
)

const (
	goodResourceID      = "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/ResourceGroup/providers/Microsoft.Compute/virtualMachines/agentvm"
	port           uint = 8080
)

func TestUN_NewNetworkWatcherProvider(t *testing.T) {

	provider = NewNetworkWatcherProvider()
	assert.NotNil(t, provider)

	err := provider.QueueConnectionTestJob()
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), errorNotInitialized)

}

func TestUN_Initialize(t *testing.T) {
	TestUN_NewNetworkWatcherProvider(t)

	sourceVM, err := models.NewVMDetails(goodResourceID)
	assert.NotNil(t, sourceVM)
	assert.Nil(t, err)

	targetVM, err := models.NewVMDetails(goodResourceID)
	assert.NotNil(t, targetVM)
	assert.Nil(t, err)

	provider.Initialize(sourceVM, targetVM, port)

}

func TestUN_QueueConnectionTestJob(t *testing.T) {
	TestUN_Initialize(t)

	os.Setenv("ARM_CLIENT_ID", "")
	os.Setenv("ARM_CLIENT_SECRET", "")
	os.Setenv("ARM_SUBSCRIPTION_ID", "")
	os.Setenv("ARM_TENANT_ID", "")

	err := provider.QueueConnectionTestJob()
	assert.NotNil(t, err)
}

func TestUN_GetAddress(t *testing.T) {
	internalProvider := NewNetworkWatcherProvider()

	sourceVM, err := models.NewVMDetails(goodResourceID)
	assert.NotNil(t, sourceVM)
	assert.Nil(t, err)

	targetVM, err := models.NewVMDetails(goodResourceID)
	assert.NotNil(t, targetVM)
	assert.Nil(t, err)

	internalProvider.Initialize(sourceVM, targetVM, port)

	submitJobAddress := internalProvider.getAddressToSubmitJob()
	assert.NotEmpty(t, submitJobAddress)
	url, err := url.Parse(submitJobAddress)
	assert.NotNil(t, url)
	assert.Nil(t, err)

	statusJobAddress := internalProvider.getAddressToJobStatus()
	assert.NotEmpty(t, statusJobAddress)
	url, err = url.Parse(statusJobAddress)
	assert.NotNil(t, url)
	assert.Nil(t, err)

}
