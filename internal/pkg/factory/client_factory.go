package factory

import (
	"crypto/tls"
	"errors"
	"net/http"
	"time"

	"github.com/spf13/viper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/client"
	"github.com/ukfast/sdk-go/pkg/connection"
	"github.com/ukfast/sdk-go/pkg/logging"
)

type ClientFactory interface {
	NewClient() (client.Client, error)
}

type UKFastClientFactoryOption func(f *UKFastClientFactory)

type UKFastClientFactory struct {
	apiUserAgent string
}

func WithUserAgent(userAgent string) UKFastClientFactoryOption {
	return func(p *UKFastClientFactory) {
		p.apiUserAgent = userAgent
	}
}

func NewUKFastClientFactory(opts ...UKFastClientFactoryOption) *UKFastClientFactory {
	f := &UKFastClientFactory{}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

func (f *UKFastClientFactory) NewClient() (client.Client, error) {
	apiKey := viper.GetString("api_key")
	if len(apiKey) < 1 {
		return nil, errors.New("Missing api_key")
	}

	conn := connection.NewAPIConnection(&connection.APIKeyCredentials{APIKey: apiKey})
	conn.UserAgent = f.apiUserAgent
	apiURI := viper.GetString("api_uri")
	if apiURI != "" {
		conn.APIURI = apiURI
	}
	apiTimeoutSeconds := viper.GetInt("api_timeout_seconds")
	if apiTimeoutSeconds > 0 {
		conn.HTTPClient.Timeout = (time.Duration(apiTimeoutSeconds) * time.Second)
	}
	if viper.GetBool("api_insecure") {
		conn.HTTPClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}
	apiHeaders := viper.GetStringMapString("api_headers")
	if apiHeaders != nil {
		conn.Headers = http.Header{}
		for headerKey, headerValue := range apiHeaders {
			conn.Headers.Add(headerKey, headerValue)
		}
	}

	if viper.GetBool("api_debug") {
		logging.SetLogger(&output.DebugLogger{})
	}

	return client.NewClient(conn), nil
}
