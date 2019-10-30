package models

import (
	"strconv"
)

const (
	connectivityResourceID = "resourceId"
	connectivityPort       = "port"
)

//ConnectivityRequestBody request body structure for networkwatcher api
type connectivityRequestBody struct {
	Source      map[string]string `json:"source"`
	Destination map[string]string `json:"destination"`
}

//ConnectivityResponseBody response body structure for networkwatcher api
type ConnectivityResponseBody struct {
	ConnectionStatus string `json:"connectionStatus"`
}

//NewConnectivyRequestBody return a struct of models.ConnectivityRequestBody
func NewConnectivyRequestBody(source, target VMDetails, port uint) *connectivityRequestBody {
	sourceMap := make(map[string]string)
	sourceMap[connectivityResourceID] = source.GetResourceID()

	destMap := make(map[string]string)
	destMap[connectivityResourceID] = target.GetResourceID()
	destMap[connectivityPort] = strconv.FormatUint(uint64(port), 10)

	conRequest := &connectivityRequestBody{
		Source:      sourceMap,
		Destination: destMap,
	}

	return conRequest
}

//GetSourceResourceID return the Source vm's resource id fromn the structure
func (req *connectivityRequestBody) GetSourceResourceID() string {
	return req.Source[connectivityResourceID]
}

//GetDestinationResourceID return the Destination resource id from ConnectivityRequestBody structure
func (req *connectivityRequestBody) GetDestinationResourceID() string {
	return req.Destination[connectivityResourceID]
}

//GetDestinationPort return the Destination port from ConnectivityRequestBody structure
func (req *connectivityRequestBody) GetDestinationPort() uint {
	value, err := strconv.ParseUint(req.Destination[connectivityPort], 10, 32)
	if err != nil {
		value = 0
	}
	return uint(value)
}
