package networkcheck

import (
	"NetworkWatcherCheck/pkg/interfaces"
	"NetworkWatcherCheck/pkg/models"
)

const (
	successResponse = "Connected"
)

type connectionTest struct {
	source             *models.VMDetails
	target             *models.VMDetails
	port               uint
	netWatcherProvider interfaces.NetworkWatcherCheckProvider
}

//NewConnectionTest create an instace of ConnectionTest
func NewConnectionTest(sourceVMId string, destinyVMId string, destinyPort uint, netWatcherProvider interfaces.NetworkWatcherCheckProvider) (*connectionTest, error) {
	connectionTest := &connectionTest{}
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
	connectionTest.netWatcherProvider = netWatcherProvider
	connectionTest.netWatcherProvider.Initialize(connectionTest.source, connectionTest.target, destinyPort)

	return connectionTest, nil

}

func (con *connectionTest) CheckVMConnection() (bool, error) {

	err := con.netWatcherProvider.QueueConnectionTestJob()
	if err != nil {
		return false, err
	}

	responseObj, err := con.netWatcherProvider.GetConnectionTestResult()
	if err != nil {
		return false, err
	}

	if responseObj.ConnectionStatus == successResponse {
		return true, nil
	}

	return false, nil
}
