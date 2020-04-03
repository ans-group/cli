package factory

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/client"
	"github.com/ukfast/sdk-go/pkg/connection"
	"github.com/ukfast/sdk-go/pkg/logging"
)

type ClientFactory interface {
	NewClient() client.Client
}

type UKFastClientFactoryOption func(f *UKFastClientFactory)

type UKFastClientFactory struct {
	apiKey      string
	apiTimeout  int
	apiURI      string
	apiInsecure bool
	apiHeaders  map[string]string
	apiDebug    bool
}

func WithAPIKey(apiKey string) UKFastClientFactoryOption {
	return func(p *UKFastClientFactory) {
		p.apiKey = apiKey
	}
}

func WithTimeout(apiTimeout int) UKFastClientFactoryOption {
	return func(p *UKFastClientFactory) {
		p.apiTimeout = apiTimeout
	}
}

func WithURI(apiURI string) UKFastClientFactoryOption {
	return func(p *UKFastClientFactory) {
		p.apiURI = apiURI
	}
}

func WithInsecure(apiInsecure bool) UKFastClientFactoryOption {
	return func(p *UKFastClientFactory) {
		p.apiInsecure = apiInsecure
	}
}

func WithHeaders(apiHeaders map[string]string) UKFastClientFactoryOption {
	return func(p *UKFastClientFactory) {
		p.apiHeaders = apiHeaders
	}
}

func WithDebug(apiDebug bool) UKFastClientFactoryOption {
	return func(p *UKFastClientFactory) {
		p.apiDebug = apiDebug
	}
}

func NewUKFastClientFactory(opts ...UKFastClientFactoryOption) *UKFastClientFactory {
	f := &UKFastClientFactory{}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

func (f *UKFastClientFactory) NewClient() client.Client {
	conn := connection.NewAPIConnection(&connection.APIKeyCredentials{APIKey: f.apiKey})
	conn.UserAgent = "ukfast-cli"
	if f.apiURI != "" {
		conn.APIURI = f.apiURI
	}
	if f.apiTimeout > 0 {
		conn.HTTPClient.Timeout = (time.Duration(f.apiTimeout) * time.Second)
	}
	if f.apiInsecure {
		conn.HTTPClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}
	if f.apiHeaders != nil {
		conn.Headers = http.Header{}
		for headerKey, headerValue := range f.apiHeaders {
			conn.Headers.Add(headerKey, headerValue)
		}
	}

	if f.apiDebug {
		logging.SetLogger(&output.DebugLogger{})
	}

	return client.NewClient(conn)
}
