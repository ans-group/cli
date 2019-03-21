package client

import (
	"github.com/ukfast/sdk-go/pkg/connection"
	"github.com/ukfast/sdk-go/pkg/service/account"
	"github.com/ukfast/sdk-go/pkg/service/ddosx"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
	"github.com/ukfast/sdk-go/pkg/service/safedns"
	"github.com/ukfast/sdk-go/pkg/service/ssl"
)

type Client interface {
	SafeDNSService() safedns.SafeDNSService
	ECloudService() ecloud.ECloudService
	SSLService() ssl.SSLService
	DDoSXService() ddosx.DDoSXService
	AccountService() account.AccountService
}

type UKFastClient struct {
	connection connection.Connection
}

func NewClient(connection connection.Connection) *UKFastClient {
	return &UKFastClient{
		connection: connection,
	}
}

func (c *UKFastClient) SafeDNSService() safedns.SafeDNSService {
	return safedns.NewService(c.connection)
}

func (c *UKFastClient) ECloudService() ecloud.ECloudService {
	return ecloud.NewService(c.connection)
}

func (c *UKFastClient) SSLService() ssl.SSLService {
	return ssl.NewService(c.connection)
}

func (c *UKFastClient) DDoSXService() ddosx.DDoSXService {
	return ddosx.NewService(c.connection)
}

func (c *UKFastClient) AccountService() account.AccountService {
	return account.NewService(c.connection)
}
