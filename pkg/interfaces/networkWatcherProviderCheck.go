package interfaces

import "NetworkWatcherCheck/pkg/models"

type NetworkWatcherCheckProvider interface {
	QueueConnectionTestJob() error
	GetConnectionTestResult() (models.ConnectivityResponseBody, error)
	Initialize(source, target *models.VMDetails, port uint)
}
