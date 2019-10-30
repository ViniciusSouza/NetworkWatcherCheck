package mockup

import (
	"NetworkWatcherCheck/pkg/interfaces"
	"NetworkWatcherCheck/pkg/models"
	"errors"
)

type connectionTestMockup struct {
	source         *models.VMDetails
	target         *models.VMDetails
	port           uint
	networkWatcher interfaces.NetworkWatcherProviderCheck
}

//NewConnectionTestMockup create an instace of connectionTestMockup
func NewConnectionTestMockup(sourceVMId string, destinyVMId string, destinyPort uint, netWatcherProvider interfaces.NetworkWatcherProviderCheck) (*connectionTestMockup, error) {
	connectionTest := &connectionTestMockup{}
	var err error

	connectionTest.source, err = models.NewVMDetails(sourceVMId)
	if err != nil {
		return nil, err
	}

	connectionTest.target, err = models.NewVMDetails(destinyVMId)
	if err != nil {
		return nil, err
	}
	connectionTest.port = destinyPort
	connectionTest.networkWatcher = netWatcherProvider
	connectionTest.networkWatcher.Initialize(connectionTest.source, connectionTest.target, destinyPort)
	return connectionTest, nil

}

func (con *connectionTestMockup) CheckVMConnection() (bool, error) {
	con.networkWatcher.QueueConnectionTestJob()
	con.networkWatcher.GetConnectionTestResult()
	return false, errors.New("Mockup")
}
