package factory

import (
	"crypto/tls"
	"errors"
	"net/http"
	"time"

	"github.com/ans-group/cli/internal/pkg/config"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/connection"
	"github.com/ans-group/sdk-go/pkg/logging"
)

type ConnectionFactory interface {
	NewConnection() (connection.Connection, error)
}

type ANSConnectionFactoryOption func(f *ANSConnectionFactory)

type ANSConnectionFactory struct {
	apiUserAgent string
}

func WithUserAgent(userAgent string) ANSConnectionFactoryOption {
	return func(p *ANSConnectionFactory) {
		p.apiUserAgent = userAgent
	}
}

func NewANSConnectionFactory(opts ...ANSConnectionFactoryOption) *ANSConnectionFactory {
	f := &ANSConnectionFactory{}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

func (f *ANSConnectionFactory) NewConnection() (connection.Connection, error) {
	apiKey := config.GetString("api_key")
	if len(apiKey) < 1 {
		return nil, errors.New("Missing api_key")
	}

	conn := connection.NewAPIConnection(&connection.APIKeyCredentials{APIKey: apiKey})
	conn.UserAgent = f.apiUserAgent
	apiURI := config.GetString("api_uri")
	if apiURI != "" {
		conn.APIURI = apiURI
	}
	apiTimeoutSeconds := config.GetInt("api_timeout_seconds")
	if apiTimeoutSeconds > 0 {
		conn.HTTPClient.Timeout = (time.Duration(apiTimeoutSeconds) * time.Second)
	}
	if config.GetBool("api_insecure") {
		conn.HTTPClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}
	apiHeaders := config.GetStringMapString("api_headers")
	if apiHeaders != nil {
		conn.Headers = http.Header{}
		for headerKey, headerValue := range apiHeaders {
			conn.Headers.Add(headerKey, headerValue)
		}
	}

	if config.GetBool("api_debug") {
		logging.SetLogger(&output.DebugLogger{})
	}

	return conn, nil
}
