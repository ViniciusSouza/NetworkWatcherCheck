package networkworkcheck

import (
	"NetworkWatcherCheck/pkg/mockup"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	sourceVM = "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/ResourceGroup/providers/Microsoft.Compute/virtualMachines/vm1"
	targetVM = "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/ResourceGroup/providers/Microsoft.Compute/virtualMachines/vm2"
)

var (
	networkProvider = mockup.NewNetworkWatcherProvider(callback)
)

func TestUN_CheckVMConnection(t *testing.T) {
	con, err := mockup.NewConnectionTestMockup(sourceVM, targetVM, port, networkProvider)
	assert.Nil(t, err)
	assert.NotNil(t, con)

	connected, err := con.CheckVMConnection()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Mockup")
	assert.False(t, connected)
}

func callback() {
	fmt.Println("netwatcher mock callback")
}
