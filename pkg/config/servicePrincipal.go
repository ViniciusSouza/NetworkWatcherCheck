package config

import (
	"fmt"

	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/kelseyhightower/envconfig"
)

const cloudName string = "AzurePublicCloud"

//ServicePrincipalConfig this struct provide azure authentication configuration properties using service principal
type ServicePrincipalConfig struct {
	ClientID       string `envconfig:"ARM_CLIENT_ID" required:"true"`
	ClientSecret   string `envconfig:"ARM_CLIENT_SECRET" required:"true"`
	SubscriptionID string `envconfig:"ARM_SUBSCRIPTION_ID" required:"true"`
	TenantID       string `envconfig:"ARM_TENANT_ID" required:"true"`
	UserAgent      string `default:"TerrainInsights"`
}

//NewServicePrincipalConfig return a new Azure configuration structure using the environemt variables
func NewServicePrincipalConfig() *ServicePrincipalConfig {

	var s ServicePrincipalConfig
	err := envconfig.Process("", &s)
	if err != nil {
		panic(err.Error())
	}

	return &s
}

//Environment return the azure environment configuration
func Environment() *azure.Environment {
	env, err := azure.EnvironmentFromName(cloudName)
	if err != nil {
		panic(fmt.Sprintf("invalid cloud name '%s' specified, cannot continue\n", cloudName))
	}
	return &env
}
