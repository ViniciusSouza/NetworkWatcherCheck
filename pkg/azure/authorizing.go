package azure

import (
	"NetworkWatcherCheck/pkg/config"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
)

type authorizing struct {
	armAuthorizer autorest.Authorizer
	armToken      *adal.ServicePrincipalToken
	azureConfig   *config.ServicePrincipalConfig
	environment   *azure.Environment
}

//NewAzureAuthorizing Create a structure for authorizing
func NewAzureAuthorizing() *authorizing {
	auth := &authorizing{}
	auth.azureConfig = config.NewServicePrincipalConfig()
	auth.environment = config.Environment()
	return auth
}

// GetResourceManagementAuthorizer gets an OAuthTokenAuthorizer for Azure Resource Manager
func (auth *authorizing) GetResourceManagementAuthorizer() (autorest.Authorizer, error) {
	var authorizer autorest.Authorizer = nil
	var err error = nil

	if auth.armAuthorizer == nil {
		authorizer, err = auth.getAuthorizerForResource(auth.environment.ResourceManagerEndpoint)

		if err == nil {
			auth.armAuthorizer = authorizer
		}
	}

	return auth.armAuthorizer, err
}

//GetNewTokenForManagerEndpoint method that returns a new OAuth token string for the request
func (auth *authorizing) GetNewTokenForManagerEndpoint() (string, error) {

	var token *adal.ServicePrincipalToken = nil
	var err error = nil

	if auth.armToken == nil {
		token, err = auth.getServicePrincipalToken(true, auth.environment.ResourceManagerEndpoint)
		if err != nil {
			return "", err
		}
		auth.armToken = token
	}

	return auth.armToken.OAuthToken(), nil

}

func (auth *authorizing) getAuthorizerForResource(resource string) (autorest.Authorizer, error) {

	var authorizer autorest.Authorizer
	var err error

	if auth.armToken == nil {
		auth.armToken, err = auth.getServicePrincipalToken(true, resource)
		if err != nil {
			return nil, err
		}
	}

	authorizer = autorest.NewBearerAuthorizer(auth.armToken)

	return authorizer, err
}

func (auth *authorizing) getServicePrincipalToken(autoRefresh bool, resource string) (*adal.ServicePrincipalToken, error) {

	oauthConfig, err := adal.NewOAuthConfig(auth.environment.ActiveDirectoryEndpoint, auth.azureConfig.TenantID)
	if err != nil {
		return nil, err
	}

	serviceToken, err := adal.NewServicePrincipalToken(*oauthConfig, auth.azureConfig.ClientID, auth.azureConfig.ClientSecret, resource)
	if err != nil {
		return nil, err
	}
	serviceToken.SetAutoRefresh(autoRefresh)
	return serviceToken, nil
}
