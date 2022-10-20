package factory

import (
	"github.com/ans-group/sdk-go/pkg/client"
)

type ClientFactory interface {
	NewClient() (client.Client, error)
}

type ANSClientFactory struct {
	connectionFactory ConnectionFactory
}

func NewANSClientFactory(connectionFactory ConnectionFactory) *ANSClientFactory {
	return &ANSClientFactory{
		connectionFactory: connectionFactory,
	}
}

func (f *ANSClientFactory) NewClient() (client.Client, error) {
	conn, err := f.connectionFactory.NewConnection()
	if err != nil {
		return nil, err
	}
	return client.NewClient(conn), nil
}
