package networkworkcheck

import (
	"NetworkWatcherCheck/pkg/azure"
	"NetworkWatcherCheck/pkg/models"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2017-12-01/compute"
)

const (
	networkWatcherRegion = "NetworkWatcher_%s"

	logInfoNetWatcherPrefix = "NetworkWatcher"

	submitJobEndpoint           = "https://management.azure.com/subscriptions/%s/resourceGroups/NetworkWatcherRG/providers/Microsoft.Network/networkWatchers/%s/connectivityCheck?api-version=2019-02-01"
	statusJobEndpoint           = "https://management.azure.com/subscriptions/%s/providers/Microsoft.Network/locations/%s/operationResults/%s?api-version=2019-02-01"
	errorNotAbleToCreateJob     = "Not able to submit a new job. Error: %s"
	errorNotAbleToCreateRequest = "Not able to create a new request. Error: %s"
	errorNotAbleToDoRequest     = "Not able to execute the request. Error: %s"
	errorGettingVMInfo          = "Not able to get VM info. Error: %s"
	errorFailGetResponse        = "Fail to get the response! %s"
	errorNotInitialized         = "Network Watcher provider not initialized!"
	errorJobStatusTimeout       = "Not able to retrieve the job status, after many retries!"

	authorizationHeader       = "Authorization"
	bearerToken               = "Bearer %s"
	contentTypeHeader         = "Content-type"
	contentAppJSON            = "application/json"
	netWatcherNumberOfRetries = 20
	netWatcherSleep           = 5
	jobIDLength               = 36
	retryableErrorType        = "RetryableError"
	retryableErrorMessage     = "Operation ConnectivityOperation (" //Used to retrieve the JodID from the message body

	queueJobRequestID = "X-Ms-Request-Id"
	missingHeader     = "Missing header %s"

	emptyResponse = "null"

	missingJobID = "Missing jobID"
)

type networkWatcherProvider struct {
	source      *models.VMDetails
	target      *models.VMDetails
	port        uint
	token       string
	jobID       string
	initialized bool
}

//NewNetworkWatcherProvider return an instace of networkWatcherProvider
func NewNetworkWatcherProvider() *networkWatcherProvider {
	return &networkWatcherProvider{
		initialized: false,
	}
}

func (net *networkWatcherProvider) Initialize(source, target *models.VMDetails, port uint) {
	net.source = source
	net.target = target
	net.port = port
	net.initialized = true
}

//QueueConnectionTestJob uses the Azure Network Watcher API to queue a connection test job
func (net *networkWatcherProvider) QueueConnectionTestJob() error {

	if !net.initialized {
		return errors.New(errorNotInitialized)
	}

	fmt.Println(fmt.Sprintf("%s - QueueConnectionTestJob \t Source: %s \t Target: %s \t Port: %d", logInfoNetWatcherPrefix, net.source.GetVMName(), net.target.GetVMName(), net.port))

	var err error = nil
	//Update VM Region - retrieve from azure
	region, err := net.getVMRegion()
	if err != nil {
		return err
	}
	net.source.SetRegion(region)

	conRequest := models.NewConnectivyRequestBody(*net.source, *net.target, net.port)
	requestURL := net.getAddressToSubmitJob()
	marshaleld, _ := json.Marshal(conRequest)

	auth := azure.NewAzureAuthorizing()
	net.token, err = auth.GetNewTokenForManagerEndpoint()
	if err != nil {
		return err
	}

	request, _ := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(marshaleld))
	request.Header.Set(contentTypeHeader, contentAppJSON)
	request.Header.Set(authorizationHeader, fmt.Sprintf(bearerToken, net.token))

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return fmt.Errorf(errorNotAbleToDoRequest, err.Error())
	}

	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		fmt.Println(fmt.Sprintf("%s - QueueConnectionTestJob Error received form the request: %d", logInfoNetWatcherPrefix, resp.StatusCode))
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		jobid := getJobIDFromRetryableError(string(bodyBytes))
		if len(jobid) > 0 {
			fmt.Println(fmt.Sprintf("%s - QueueConnectionTestJob The job was already queued!", logInfoNetWatcherPrefix))
			net.jobID = jobid
			return nil
		}

		fmt.Println(fmt.Sprintf("%s - QueueConnectionTestJob:  %s", logInfoNetWatcherPrefix, string(bodyBytes)))
		return fmt.Errorf(errorFailGetResponse, string(bodyBytes))
	}

	headerValue := resp.Header

	value, found := headerValue[queueJobRequestID]
	if !found {
		fmt.Println(fmt.Sprintf(missingHeader, queueJobRequestID))
		return fmt.Errorf(missingHeader, queueJobRequestID)
	}

	net.jobID = value[0]

	return nil
}

func isRetryableError(strError string) bool {

	if strings.Contains(strError, retryableErrorType) {
		return true
	}
	return false

}

func getJobIDFromRetryableError(strError string) string {
	retryableJob := ""
	if isRetryableError(strError) {
		idx := strings.Index(strError, retryableErrorMessage)
		if idx > -1 {
			init := idx + len(retryableErrorMessage)
			final := init + jobIDLength
			retryableJob = strError[init:final]
		}
	}
	return retryableJob
}

func (net *networkWatcherProvider) GetConnectionTestResult() (models.ConnectivityResponseBody, error) {

	if !net.initialized {
		return models.ConnectivityResponseBody{}, errors.New(errorNotInitialized)
	}

	if len(net.jobID) == 0 {
		return models.ConnectivityResponseBody{}, fmt.Errorf(errorNotAbleToDoRequest, missingJobID)
	}

	request, err := http.NewRequest(http.MethodGet, net.getAddressToJobStatus(), nil)
	if err != nil {
		return models.ConnectivityResponseBody{}, fmt.Errorf(errorNotAbleToCreateRequest, err.Error())
	}
	request.Header.Set(contentTypeHeader, contentAppJSON)
	request.Header.Set(authorizationHeader, fmt.Sprintf(bearerToken, net.token))

	var responseObj models.ConnectivityResponseBody
	var errObj error = nil

	retries := 0
	for retries < netWatcherNumberOfRetries {
		client := http.Client{}
		fmt.Println(fmt.Sprintf("%s \t  GetConnectionTestResult - Try: %d \t Source: %s \t Target: %s \t Port: %d", logInfoNetWatcherPrefix, retries, net.source.GetVMName(), net.target.GetVMName(), net.port))
		resp, _ := client.Do(request)
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err == nil {
				bodyString := string(bodyBytes)

				if bodyString != emptyResponse {
					json.Unmarshal(bodyBytes, &responseObj)
					break
				}
			}
		}
		time.Sleep(netWatcherSleep * time.Second)
		retries++
	}

	if len(responseObj.ConnectionStatus) == 0 {
		errObj = errors.New(errorJobStatusTimeout)
	}

	return responseObj, errObj
}

func (net *networkWatcherProvider) getVMRegion() (string, error) {
	auth := azure.NewAzureAuthorizing()
	vmClient := compute.NewVirtualMachinesClient(net.source.GetSubscriptionID())
	authorizer, err := auth.GetResourceManagementAuthorizer()
	vmClient.Authorizer = authorizer

	ctx := context.Background()
	vmName := net.source.GetVMName()
	rgName := net.source.GetResourceGroupName()
	result, err := vmClient.Get(ctx, rgName, vmName, compute.InstanceView)

	if err != nil {
		return "", fmt.Errorf(errorGettingVMInfo, err.Error())
	}

	return *result.Location, nil
}

func (net *networkWatcherProvider) getAddressToSubmitJob() string {
	return fmt.Sprintf(submitJobEndpoint, net.source.GetSubscriptionID(), fmt.Sprintf(networkWatcherRegion, net.source.GetRegion()))
}

func (net *networkWatcherProvider) getAddressToJobStatus() string {
	return fmt.Sprintf(statusJobEndpoint, net.source.GetSubscriptionID(), net.source.GetRegion(), net.jobID)
}
