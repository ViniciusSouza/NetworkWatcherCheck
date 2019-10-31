package mockup

import "TerraformTestFramework/pkg/models"

type networkWatcherMockup struct {
	initialized bool
	jobID       string
	callback    func()
}

//NewNetworkWatcherProvider Creates a new mockup structure for the network watcher
func NewNetworkWatcherProvider(callback func()) *networkWatcherMockup {
	return &networkWatcherMockup{
		initialized: false,
		callback:    callback,
	}
}

func (net *networkWatcherMockup) QueueConnectionTestJob() error {
	net.jobID = "mockup"
	return nil
}

func (net *networkWatcherMockup) GetConnectionTestResult() (models.ConnectivityResponseBody, error) {
	net.callback()
	return models.ConnectivityResponseBody{}, nil
}

func (net *networkWatcherMockup) Initialize(source, target *models.VMDetails, port uint) {
	net.initialized = true
}
