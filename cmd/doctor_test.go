package cmd

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/cli/test"
	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/sdk-go/pkg/client"
	"github.com/ans-group/sdk-go/pkg/config"
	"github.com/ans-group/sdk-go/pkg/connection"
	"github.com/ans-group/sdk-go/pkg/service/account"
	"github.com/ans-group/sdk-go/pkg/service/billing"
	"github.com/ans-group/sdk-go/pkg/service/cloudflare"
	"github.com/ans-group/sdk-go/pkg/service/ddosx"
	"github.com/ans-group/sdk-go/pkg/service/draas"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/ans-group/sdk-go/pkg/service/ecloudflex"
	"github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	"github.com/ans-group/sdk-go/pkg/service/pss"
	"github.com/ans-group/sdk-go/pkg/service/registrar"
	"github.com/ans-group/sdk-go/pkg/service/safedns"
	"github.com/ans-group/sdk-go/pkg/service/sharedexchange"
	"github.com/ans-group/sdk-go/pkg/service/ssl"
	"github.com/ans-group/sdk-go/pkg/service/storage"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// doctorTestClient is a hand-written implementation of client.Client returning per-service
// mocks. test/mocks contains no mock for the client interface itself, and the interface is
// small enough that hand-writing a fake is simpler than generating one.
//
// Note that no mutating (Create/Update/Patch/Delete) expectations are ever declared in
// this file. gomock fails on unexpected calls, so the absence of those expectations is
// what proves the doctor command is read-only
type doctorTestClient struct {
	account      *mocks.MockAccountService
	billing      *mocks.MockBillingService
	cloudflare   *mocks.MockCloudflareService
	ddosx        *mocks.MockDDoSXService
	draas        *mocks.MockDRaaSService
	ecloud       *mocks.MockECloudService
	loadbalancer *mocks.MockLoadBalancerService
	pss          *mocks.MockPSSService
	registrar    *mocks.MockRegistrarService
	safedns      *mocks.MockSafeDNSService
	ssl          *mocks.MockSSLService
	storage      *mocks.MockStorageService
}

func newDoctorTestClient(ctrl *gomock.Controller) *doctorTestClient {
	return &doctorTestClient{
		account:      mocks.NewMockAccountService(ctrl),
		billing:      mocks.NewMockBillingService(ctrl),
		cloudflare:   mocks.NewMockCloudflareService(ctrl),
		ddosx:        mocks.NewMockDDoSXService(ctrl),
		draas:        mocks.NewMockDRaaSService(ctrl),
		ecloud:       mocks.NewMockECloudService(ctrl),
		loadbalancer: mocks.NewMockLoadBalancerService(ctrl),
		pss:          mocks.NewMockPSSService(ctrl),
		registrar:    mocks.NewMockRegistrarService(ctrl),
		safedns:      mocks.NewMockSafeDNSService(ctrl),
		ssl:          mocks.NewMockSSLService(ctrl),
		storage:      mocks.NewMockStorageService(ctrl),
	}
}

func (c *doctorTestClient) AccountService() account.AccountService          { return c.account }
func (c *doctorTestClient) BillingService() billing.BillingService          { return c.billing }
func (c *doctorTestClient) CloudflareService() cloudflare.CloudflareService { return c.cloudflare }
func (c *doctorTestClient) DDoSXService() ddosx.DDoSXService                { return c.ddosx }
func (c *doctorTestClient) DRaaSService() draas.DRaaSService                { return c.draas }
func (c *doctorTestClient) ECloudService() ecloud.ECloudService             { return c.ecloud }
func (c *doctorTestClient) ECloudFlexService() ecloudflex.ECloudFlexService { return nil }
func (c *doctorTestClient) LoadBalancerService() loadbalancer.LoadBalancerService {
	return c.loadbalancer
}
func (c *doctorTestClient) PSSService() pss.PSSService                   { return c.pss }
func (c *doctorTestClient) RegistrarService() registrar.RegistrarService { return c.registrar }
func (c *doctorTestClient) SafeDNSService() safedns.SafeDNSService       { return c.safedns }
func (c *doctorTestClient) SharedExchangeService() sharedexchange.SharedExchangeService {
	return nil
}
func (c *doctorTestClient) SSLService() ssl.SSLService             { return c.ssl }
func (c *doctorTestClient) StorageService() storage.StorageService { return c.storage }

// doctorTestClientFactory implements factory.ClientFactory, returning a fixed client or,
// when err is set, a fixed error
type doctorTestClientFactory struct {
	client client.Client
	err    error
}

func (f *doctorTestClientFactory) NewClient() (client.Client, error) {
	if f.err != nil {
		return nil, f.err
	}

	return f.client, nil
}

// expectDoctorServiceProbes declares an expectation for each per-service probe, returning
// the error configured in errs for that service (nil, i.e. reachable, by default).
//
// Only the paginated variant of each list call is expected. The non-paginated variants
// walk every page, so gomock failing on an unexpected call is what proves each probe
// issues a single request
func expectDoctorServiceProbes(c *doctorTestClient, errs map[string]error, times int) {
	c.account.EXPECT().GetClientsPaginated(gomock.Any()).Return(doctorTestPage[account.Client](errs["account"])).Times(times)
	c.billing.EXPECT().GetCloudCostsPaginated(gomock.Any()).Return(doctorTestPage[billing.CloudCost](errs["billing"])).Times(times)
	c.cloudflare.EXPECT().GetAccountsPaginated(gomock.Any()).Return(doctorTestPage[cloudflare.Account](errs["cloudflare"])).Times(times)
	c.ddosx.EXPECT().GetDomainsPaginated(gomock.Any()).Return(doctorTestPage[ddosx.Domain](errs["ddosx"])).Times(times)
	c.draas.EXPECT().GetSolutionsPaginated(gomock.Any()).Return(doctorTestPage[draas.Solution](errs["draas"])).Times(times)
	c.ecloud.EXPECT().GetVPCsPaginated(gomock.Any()).Return(doctorTestPage[ecloud.VPC](errs["ecloud"])).Times(times)
	c.loadbalancer.EXPECT().GetClustersPaginated(gomock.Any()).Return(doctorTestPage[loadbalancer.Cluster](errs["loadbalancer"])).Times(times)
	c.pss.EXPECT().GetRequestsPaginated(gomock.Any()).Return(doctorTestPage[pss.Request](errs["pss"])).Times(times)
	c.registrar.EXPECT().GetDomainsPaginated(gomock.Any()).Return(doctorTestPage[registrar.Domain](errs["registrar"])).Times(times)
	c.safedns.EXPECT().GetZonesPaginated(gomock.Any()).Return(doctorTestPage[safedns.Zone](errs["safedns"])).Times(times)
	c.ssl.EXPECT().GetCertificatesPaginated(gomock.Any()).Return(doctorTestPage[ssl.Certificate](errs["ssl"])).Times(times)
}

// doctorTestPage returns an empty single page of results alongside the provided error, as
// returned by the paginated list call of each service. An empty page is a successful
// probe - it means the key was accepted and the account holds no resources of that type
func doctorTestPage[T any](err error) (*connection.Paginated[T], error) {
	return connection.NewPaginated(&connection.APIResponseBodyData[[]T]{}, connection.APIRequestParameters{}, nil), err
}

// initDoctorTestConfig initialises the config package from an in-memory config file
// containing the provided YAML, resetting it once the test completes.
//
// The config package wraps global viper state, so tests using it must not be run in
// parallel. ANS_API_KEY is cleared as viper reads it automatically, which would otherwise
// leak the developer's own key into the assertions below
func initDoctorTestConfig(t *testing.T, contents string) {
	t.Helper()

	t.Setenv("ANS_API_KEY", "")

	fs := afero.NewMemMapFs()
	config.SetFs(fs)

	err := afero.WriteFile(fs, "/tmp/doctorconfig.yml", []byte(contents), 0644)
	if err != nil {
		t.Fatalf("failed to write test config: %s", err)
	}

	err = config.Init("/tmp/doctorconfig.yml")
	if err != nil {
		t.Fatalf("failed to initialise test config: %s", err)
	}

	t.Cleanup(config.Reset)
}

// resetDoctorTestErrorLevel clears the output package's error level before and after the
// test, as it is a package global which would otherwise leak between tests asserting on
// the exit code
func resetDoctorTestErrorLevel(t *testing.T) {
	t.Helper()

	output.SetErrorLevel(0)
	t.Cleanup(func() { output.SetErrorLevel(0) })
}

// doctorResultByCheck returns the result with the given check name
func doctorResultByCheck(t *testing.T, results []DoctorResult, check string) DoctorResult {
	t.Helper()

	for _, result := range results {
		if result.Check == check {
			return result
		}
	}

	t.Fatalf("no result found for check '%s'", check)
	return DoctorResult{}
}

// apiResponseBodyError returns an SDK error carrying the given HTTP status code
func apiResponseBodyError(statusCode int) error {
	return &connection.APIResponseBodyError{
		Errors: []connection.APIResponseBodyErrorItem{
			{Status: statusCode, Title: "test error"},
		},
	}
}

// apiIPNotAllowedError returns an SDK error shaped as the API returns it when a request is
// rejected because the caller's IP address is not on the key's allow list - a plain message
// field rather than a structured errors array, wrapped as connection.APIResponse does
func apiIPNotAllowedError(ipAddress string) error {
	bodyErr := &connection.APIResponseBodyError{Message: fmt.Sprintf("IP address not allowed: %s", ipAddress)}
	return fmt.Errorf("unexpected status code (403): %w", bodyErr)
}

func Test_httpStatusFromError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected int
	}{
		{
			name:     "StructuredError_ReturnsItemStatus",
			err:      apiResponseBodyError(403),
			expected: 403,
		},
		{
			name:     "StructuredErrorWithoutItems_ReturnsZero",
			err:      &connection.APIResponseBodyError{Errors: nil},
			expected: 0,
		},
		{
			name:     "UnexpectedStatusCodeError_ReturnsParsedStatus",
			err:      errors.New("unexpected status code (401)"),
			expected: 401,
		},
		{
			name:     "TransportError_ReturnsZero",
			err:      errors.New("dial tcp: connection refused"),
			expected: 0,
		},
		{
			name:     "NilError_ReturnsZero",
			err:      nil,
			expected: 0,
		},
		{
			name:     "WrappedStructuredError_ReturnsItemStatus",
			err:      fmt.Errorf("failed to retrieve data: %w", apiResponseBodyError(404)),
			expected: 404,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, httpStatusFromError(tt.err))
		})
	}
}

func Test_doctorIPNotAllowed(t *testing.T) {
	tests := []struct {
		name            string
		err             error
		expectedBlocked bool
		expectedIP      string
	}{
		{
			name:            "IPNotAllowedError_ReturnsIP",
			err:             apiIPNotAllowedError("82.69.110.153"),
			expectedBlocked: true,
			expectedIP:      "82.69.110.153",
		},
		{
			name:            "RolePermissionError_NotBlocked",
			err:             apiResponseBodyError(403),
			expectedBlocked: false,
			expectedIP:      "",
		},
		{
			name:            "NilError_NotBlocked",
			err:             nil,
			expectedBlocked: false,
			expectedIP:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blocked, ipAddress := doctorIPNotAllowed(tt.err)

			assert.Equal(t, tt.expectedBlocked, blocked)
			assert.Equal(t, tt.expectedIP, ipAddress)
		})
	}
}

func Test_doctorProbeStatus(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus string
		expectedDetail string
	}{
		{
			name:           "NoError_OK",
			err:            nil,
			expectedStatus: doctorStatusOK,
			expectedDetail: "reachable",
		},
		{
			// A key scoped to only some services is correctly configured, so an
			// unauthorised probe is informational rather than a failure
			name:           "Unauthorised_Warns",
			err:            apiResponseBodyError(401),
			expectedStatus: doctorStatusWarn,
			expectedDetail: "401 - api key not permitted for this service",
		},
		{
			name:           "Forbidden_Warns",
			err:            apiResponseBodyError(403),
			expectedStatus: doctorStatusWarn,
			expectedDetail: "403 - api key lacks the required role for this service",
		},
		{
			// IP restriction is reported distinctly from a role/permission failure, as the
			// fix is different (allow-list the address rather than grant a role)
			name:           "IPNotAllowed_Warns",
			err:            apiIPNotAllowedError("82.69.110.153"),
			expectedStatus: doctorStatusWarn,
			expectedDetail: "403 - ip address 82.69.110.153 not allowed for this api key",
		},
		{
			name:           "NotFound_Warns",
			err:            apiResponseBodyError(404),
			expectedStatus: doctorStatusWarn,
			expectedDetail: "404 not found (service may not be provisioned)",
		},
		{
			name:           "ServerError_Fails",
			err:            apiResponseBodyError(500),
			expectedStatus: doctorStatusFail,
			expectedDetail: "500 server error",
		},
		{
			name:           "UnknownError_Warns",
			err:            errors.New("dial tcp: connection refused"),
			expectedStatus: doctorStatusWarn,
			expectedDetail: "dial tcp: connection refused",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, detail := doctorProbeStatus(tt.err)

			assert.Equal(t, tt.expectedStatus, status)
			assert.Equal(t, tt.expectedDetail, detail)
		})
	}
}

func Test_doctorConfigChecks(t *testing.T) {
	t.Run("MissingAPIKey_Fails", func(t *testing.T) {
		initDoctorTestConfig(t, "contexts:\n  testcontext:\n    api_uri: api.example.com\ncurrent_context: testcontext\n")

		results := doctorConfigChecks("/tmp/doctorconfig.yml")

		result := doctorResultByCheck(t, results, "api key")
		assert.Equal(t, doctorStatusFail, result.Status)
		assert.Equal(t, "not set", result.Detail)
	})

	t.Run("APIKeyNeverPrinted", func(t *testing.T) {
		initDoctorTestConfig(t, "contexts:\n  testcontext:\n    api_key: SUPERSECRETVALUE\ncurrent_context: testcontext\n")

		results := doctorConfigChecks("/tmp/doctorconfig.yml")

		for _, result := range results {
			assert.NotContains(t, result.Check, "SUPERSECRETVALUE")
			assert.NotContains(t, result.Status, "SUPERSECRETVALUE")
			assert.NotContains(t, result.Detail, "SUPERSECRETVALUE")
			assert.NotContains(t, result.Context, "SUPERSECRETVALUE")
		}

		result := doctorResultByCheck(t, results, "api key")
		assert.Equal(t, doctorStatusOK, result.Status)
		assert.Equal(t, "set (16 chars)", result.Detail)
	})

	t.Run("InvalidOutputDefault_Warns", func(t *testing.T) {
		initDoctorTestConfig(t, "output:\n  default: nonsense\napi_key: someapikey\n")

		results := doctorConfigChecks("/tmp/doctorconfig.yml")

		result := doctorResultByCheck(t, results, "config: output.default")
		assert.Equal(t, doctorStatusWarn, result.Status)
		assert.Contains(t, result.Detail, "nonsense")
	})

	t.Run("ValidOutputDefault_OK", func(t *testing.T) {
		initDoctorTestConfig(t, "output:\n  default: json\napi_key: someapikey\n")

		results := doctorConfigChecks("/tmp/doctorconfig.yml")

		result := doctorResultByCheck(t, results, "config: output.default")
		assert.Equal(t, doctorStatusOK, result.Status)
		assert.Equal(t, "json", result.Detail)
	})

	t.Run("ContextMissingAPIKey_Warns", func(t *testing.T) {
		initDoctorTestConfig(t, "contexts:\n  withkey:\n    api_key: someapikey\n  withoutkey:\n    api_uri: api.example.com\ncurrent_context: withkey\n")

		results := doctorConfigChecks("/tmp/doctorconfig.yml")

		result := doctorResultByCheck(t, results, "config: contexts")
		assert.Equal(t, doctorStatusWarn, result.Status)
		assert.Contains(t, result.Detail, "withoutkey")
		assert.NotContains(t, result.Detail, "withkey,")
	})

	t.Run("ContextsChecked_RestoresCurrentContext", func(t *testing.T) {
		initDoctorTestConfig(t, "contexts:\n  contexta:\n    api_key: someapikey\n  contextb:\n    api_key: someapikey\ncurrent_context: contextb\n")

		doctorConfigChecks("/tmp/doctorconfig.yml")

		assert.Equal(t, "contextb", config.GetCurrentContextName())
	})

	t.Run("MissingConfigFile_Warns", func(t *testing.T) {
		initDoctorTestConfig(t, "api_key: someapikey\n")

		results := doctorConfigChecks("/tmp/doesnotexist.yml")

		result := doctorResultByCheck(t, results, "config file")
		assert.Equal(t, doctorStatusWarn, result.Status)
	})

	t.Run("ExistingConfigFile_OK", func(t *testing.T) {
		initDoctorTestConfig(t, "api_key: someapikey\n")

		configPath := fmt.Sprintf("%s/doctorconfig.yml", t.TempDir())
		err := afero.WriteFile(afero.NewOsFs(), configPath, []byte("api_key: someapikey\n"), 0644)
		assert.Nil(t, err)

		results := doctorConfigChecks(configPath)

		result := doctorResultByCheck(t, results, "config file")
		assert.Equal(t, doctorStatusOK, result.Status)
		assert.Equal(t, configPath, result.Detail)
	})
}

func Test_doctorAuthChecks(t *testing.T) {
	t.Run("Success_AllOK", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)
		service.EXPECT().GetDetails().Return(account.Details{}, nil).Times(1)

		results := doctorAuthChecks(service)

		assert.Equal(t, doctorStatusOK, doctorResultByCheck(t, results, "api connectivity").Status)
		assert.Equal(t, doctorStatusOK, doctorResultByCheck(t, results, "api authentication").Status)
	})

	t.Run("Unauthorised_AuthenticationFails", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)
		service.EXPECT().GetDetails().Return(account.Details{}, apiResponseBodyError(401)).Times(1)

		results := doctorAuthChecks(service)

		// A status code implies the API was reached, so connectivity itself is fine
		assert.Equal(t, doctorStatusOK, doctorResultByCheck(t, results, "api connectivity").Status)

		authentication := doctorResultByCheck(t, results, "api authentication")
		assert.Equal(t, doctorStatusFail, authentication.Status)
		assert.Contains(t, authentication.Detail, "401")
	})

	t.Run("IPNotAllowed_AuthenticationFailsWithIPDetail", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)
		service.EXPECT().GetDetails().Return(account.Details{}, apiIPNotAllowedError("82.69.110.153")).Times(1)

		results := doctorAuthChecks(service)

		assert.Equal(t, doctorStatusOK, doctorResultByCheck(t, results, "api connectivity").Status)

		authentication := doctorResultByCheck(t, results, "api authentication")
		assert.Equal(t, doctorStatusFail, authentication.Status)
		assert.Contains(t, authentication.Detail, "82.69.110.153")
		assert.NotContains(t, authentication.Detail, "lacks permission")
	})

	t.Run("TransportError_ConnectivityFails", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)
		service.EXPECT().GetDetails().Return(account.Details{}, errors.New("dial tcp: connection refused")).Times(1)

		results := doctorAuthChecks(service)

		connectivity := doctorResultByCheck(t, results, "api connectivity")
		assert.Equal(t, doctorStatusFail, connectivity.Status)
		assert.Contains(t, connectivity.Detail, "connection refused")
		assert.Equal(t, doctorStatusWarn, doctorResultByCheck(t, results, "api authentication").Status)
	})
}

func Test_doctorServiceChecks(t *testing.T) {
	t.Run("AllServicesReachable_AllOK", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		c := newDoctorTestClient(mockCtrl)
		expectDoctorServiceProbes(c, nil, 1)

		results := doctorServiceChecks(c)

		assert.Len(t, results, 11)
		for _, result := range results {
			assert.Equal(t, doctorStatusOK, result.Status, result.Check)
			assert.Equal(t, "reachable", result.Detail)
		}
	})

	t.Run("EmptyListIsSuccess", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		c := newDoctorTestClient(mockCtrl)
		expectDoctorServiceProbes(c, nil, 1)

		results := doctorServiceChecks(c)

		// Every mock returns an empty slice, which means the key was accepted and the
		// account simply holds no resources of that type
		assert.Equal(t, doctorStatusOK, doctorResultByCheck(t, results, "service: ecloud").Status)
	})

	t.Run("ForbiddenService_Warns", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		c := newDoctorTestClient(mockCtrl)
		expectDoctorServiceProbes(c, map[string]error{"cloudflare": apiResponseBodyError(403)}, 1)

		results := doctorServiceChecks(c)

		cloudflareResult := doctorResultByCheck(t, results, "service: cloudflare")
		assert.Equal(t, doctorStatusWarn, cloudflareResult.Status)
		assert.Equal(t, "403 - api key lacks the required role for this service", cloudflareResult.Detail)

		for _, result := range results {
			if result.Check == "service: cloudflare" {
				continue
			}

			assert.Equal(t, doctorStatusOK, result.Status, result.Check)
		}
	})

	t.Run("ResultsAreDeterministic", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		c := newDoctorTestClient(mockCtrl)
		expectDoctorServiceProbes(c, map[string]error{"pss": apiResponseBodyError(404)}, 2)

		first := doctorServiceChecks(c)
		second := doctorServiceChecks(c)

		assert.Equal(t, first, second)
	})
}

func Test_doctorApplicationChecks(t *testing.T) {
	t.Run("Scopes_OK", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)
		service.EXPECT().GetApplicationServices("app-abcdef12").Return(account.ApplicationServiceMapping{
			Scopes: []account.ApplicationServiceScope{
				{Service: "ecloud", Roles: []string{"read", "write"}},
			},
		}, nil).Times(1)
		service.EXPECT().GetApplicationRestrictions("app-abcdef12").Return(account.ApplicationRestriction{}, nil).Times(1)

		results := doctorApplicationChecks(service, "app-abcdef12")

		services := doctorResultByCheck(t, results, "application services")
		assert.Equal(t, doctorStatusOK, services.Status)
		assert.Equal(t, "ecloud(read,write)", services.Detail)

		// The zero value represents an empty array from the API, meaning no restrictions
		restrictions := doctorResultByCheck(t, results, "application restrictions")
		assert.Equal(t, doctorStatusOK, restrictions.Status)
		assert.Equal(t, "no restrictions configured", restrictions.Detail)
	})

	t.Run("Restrictions_Warns", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)
		service.EXPECT().GetApplicationServices("app-abcdef12").Return(account.ApplicationServiceMapping{}, nil).Times(1)
		service.EXPECT().GetApplicationRestrictions("app-abcdef12").Return(account.ApplicationRestriction{
			IPRestrictionType: "allowlist",
			IPRanges:          []string{"1.2.3.4/32"},
		}, nil).Times(1)

		results := doctorApplicationChecks(service, "app-abcdef12")

		restrictions := doctorResultByCheck(t, results, "application restrictions")
		assert.Equal(t, doctorStatusWarn, restrictions.Status)
		assert.Contains(t, restrictions.Detail, "allowlist")
	})

	t.Run("Error_Fails", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)
		service.EXPECT().GetApplicationServices("app-abcdef12").Return(account.ApplicationServiceMapping{}, errors.New("test error")).Times(1)
		service.EXPECT().GetApplicationRestrictions("app-abcdef12").Return(account.ApplicationRestriction{}, errors.New("test error")).Times(1)

		results := doctorApplicationChecks(service, "app-abcdef12")

		assert.Equal(t, doctorStatusFail, doctorResultByCheck(t, results, "application services").Status)
		assert.Equal(t, doctorStatusFail, doctorResultByCheck(t, results, "application restrictions").Status)
	})
}

func Test_doctorVersionCheck(t *testing.T) {
	t.Run("UnknownVersion_Skips", func(t *testing.T) {
		result := doctorVersionCheck("UNKNOWN")

		assert.Equal(t, doctorStatusSkip, result.Status)
	})

	t.Run("EmptyVersion_Skips", func(t *testing.T) {
		result := doctorVersionCheck("")

		assert.Equal(t, doctorStatusSkip, result.Status)
	})

	t.Run("UnparseableVersion_Warns", func(t *testing.T) {
		result := doctorVersionCheck("not-a-version")

		assert.Equal(t, doctorStatusWarn, result.Status)
		assert.Contains(t, result.Detail, "not-a-version")
	})
}

func Test_doctorResultCollection_DefaultColumns(t *testing.T) {
	t.Run("WithoutContext_OmitsContextColumn", func(t *testing.T) {
		collection := DoctorResultCollection{{Check: "api key", Status: doctorStatusOK}}

		assert.Equal(t, []string{"check", "status", "detail"}, collection.DefaultColumns())
	})

	t.Run("WithContext_IncludesContextColumn", func(t *testing.T) {
		collection := DoctorResultCollection{{Check: "api key", Status: doctorStatusOK, Context: "somecontext"}}

		assert.Equal(t, []string{"context", "check", "status", "detail"}, collection.DefaultColumns())
	})
}

func Test_doctorTruncate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "ShortString_Unchanged",
			input:    "connection refused",
			expected: "connection refused",
		},
		{
			name:     "NewlinesReplacedAndTrimmed",
			input:    " failed to connect\nconnection refused ",
			expected: "failed to connect connection refused",
		},
		{
			name:     "LongString_Truncated",
			input:    strings.Repeat("a", 100),
			expected: strings.Repeat("a", doctorDetailMaxLength-3) + "...",
		},
		{
			// Multi-byte runes must not be cut in half, so the result is measured in
			// runes rather than bytes
			name:     "MultiByteString_TruncatedByRune",
			input:    strings.Repeat("é", 100),
			expected: strings.Repeat("é", doctorDetailMaxLength-3) + "...",
		},
		{
			name:     "MultiByteStringWithinLimit_Unchanged",
			input:    strings.Repeat("é", doctorDetailMaxLength),
			expected: strings.Repeat("é", doctorDetailMaxLength),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, doctorTruncate(tt.input))
		})
	}
}

// doctorTestCmd returns a doctor command with the given flags set, detached from the root
// command so that it can be executed directly
func doctorTestCmd(t *testing.T, f *doctorTestClientFactory, flags map[string]string) *cobra.Command {
	t.Helper()

	cmd := doctorCmd(f)
	for name, value := range flags {
		err := cmd.Flags().Set(name, value)
		if err != nil {
			t.Fatalf("failed to set flag '%s': %s", name, err)
		}
	}

	return cmd
}

func Test_doctor(t *testing.T) {
	t.Run("SkipServices_SkipsProbes", func(t *testing.T) {
		initDoctorTestConfig(t, "api_key: someapikey\n")

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		c := newDoctorTestClient(mockCtrl)
		c.account.EXPECT().GetDetails().Return(account.Details{}, nil).Times(1)

		// No service probe expectations are declared, so gomock will fail the test if any
		// probe is called
		f := &doctorTestClientFactory{client: c}
		cmd := doctorTestCmd(t, f, map[string]string{
			"skip-services":      "true",
			"skip-version-check": "true",
		})

		var err error
		stdOut := test.CatchStdOut(t, func() {
			err = doctor(f, cmd)
		})

		assert.Nil(t, err)
		assert.Contains(t, stdOut, "skipped (--skip-services)")
	})

	t.Run("AllChecks_RendersResults", func(t *testing.T) {
		initDoctorTestConfig(t, "contexts:\n  testcontext:\n    api_key: SUPERSECRETVALUE\ncurrent_context: testcontext\n")

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		c := newDoctorTestClient(mockCtrl)
		c.account.EXPECT().GetDetails().Return(account.Details{}, nil).Times(1)
		expectDoctorServiceProbes(c, nil, 1)

		f := &doctorTestClientFactory{client: c}
		cmd := doctorTestCmd(t, f, map[string]string{
			"skip-version-check": "true",
		})

		var err error
		stdOut := test.CatchStdOut(t, func() {
			err = doctor(f, cmd)
		})

		assert.Nil(t, err)
		assert.Contains(t, stdOut, "service: ecloud")
		assert.Contains(t, stdOut, "api authentication")
		// The API key must never be rendered, in any output format
		assert.NotContains(t, stdOut, "SUPERSECRETVALUE")
	})

	t.Run("ScopeRestrictedKey_ErrorLevelUnchanged", func(t *testing.T) {
		initDoctorTestConfig(t, "api_key: someapikey\n")
		resetDoctorTestErrorLevel(t)

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		// A key which is valid but scoped to only some services is correctly configured,
		// so neither the unreachable services nor the warnings affect the exit code
		c := newDoctorTestClient(mockCtrl)
		c.account.EXPECT().GetDetails().Return(account.Details{}, nil).Times(1)
		expectDoctorServiceProbes(c, map[string]error{
			"registrar":  apiResponseBodyError(404),
			"cloudflare": apiResponseBodyError(401),
			"ddosx":      apiResponseBodyError(403),
		}, 1)

		f := &doctorTestClientFactory{client: c}
		cmd := doctorTestCmd(t, f, map[string]string{
			"skip-version-check": "true",
		})

		exitCode := -1
		oldOutputExit := output.SetOutputExit(func(code int) { exitCode = code })
		defer func() { output.SetOutputExit(oldOutputExit) }()

		var stdErr string
		_, stdErr = test.CatchStdOutStdErr(t, func() {
			_ = doctor(f, cmd)
			output.ExitWithErrorLevel()
		})

		assert.Equal(t, 0, exitCode)
		assert.Empty(t, stdErr)
	})

	t.Run("FailSetsErrorLevel", func(t *testing.T) {
		initDoctorTestConfig(t, "api_key: someapikey\n")
		resetDoctorTestErrorLevel(t)

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		// An unauthorised account service call means the key itself was rejected
		c := newDoctorTestClient(mockCtrl)
		c.account.EXPECT().GetDetails().Return(account.Details{}, apiResponseBodyError(401)).Times(1)
		expectDoctorServiceProbes(c, map[string]error{"ecloud": apiResponseBodyError(401)}, 1)

		f := &doctorTestClientFactory{client: c}
		cmd := doctorTestCmd(t, f, map[string]string{
			"skip-version-check": "true",
		})

		exitCode := -1
		oldOutputExit := output.SetOutputExit(func(code int) { exitCode = code })
		defer func() { output.SetOutputExit(oldOutputExit) }()

		var stdOut, stdErr string
		stdOut, stdErr = test.CatchStdOutStdErr(t, func() {
			_ = doctor(f, cmd)
			output.ExitWithErrorLevel()
		})

		assert.Equal(t, 1, exitCode)
		assert.Contains(t, stdOut, "api authentication")
		// Failures are described in the rendered output, so nothing is written to stderr
		assert.Empty(t, stdErr)
	})

	// A client cannot be created without an api key, which is exactly the misconfiguration
	// this command exists to diagnose, so the report must still be rendered
	t.Run("ClientError_RendersReport", func(t *testing.T) {
		initDoctorTestConfig(t, "output:\n  default: nonsense\n")
		resetDoctorTestErrorLevel(t)

		f := &doctorTestClientFactory{err: errors.New("missing api_key")}
		cmd := doctorTestCmd(t, f, map[string]string{
			"skip-version-check": "true",
		})

		exitCode := -1
		oldOutputExit := output.SetOutputExit(func(code int) { exitCode = code })
		defer func() { output.SetOutputExit(oldOutputExit) }()

		var err error
		stdOut := test.CatchStdOut(t, func() {
			err = doctor(f, cmd)
			output.ExitWithErrorLevel()
		})

		assert.Nil(t, err)
		assert.Equal(t, 1, exitCode)
		// The config checks, which explain the cause, must survive the client failure
		assert.Contains(t, stdOut, "not set")
		assert.Contains(t, stdOut, "nonsense")
		assert.Contains(t, stdOut, "missing api_key")
	})

	t.Run("AllContexts_RestoresOriginalContext", func(t *testing.T) {
		initDoctorTestConfig(t, "contexts:\n  contexta:\n    api_key: someapikey\n  contextb:\n    api_key: someapikey\ncurrent_context: contextb\n")

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		c := newDoctorTestClient(mockCtrl)
		c.account.EXPECT().GetDetails().Return(account.Details{}, nil).Times(2)
		expectDoctorServiceProbes(c, nil, 2)

		f := &doctorTestClientFactory{client: c}
		cmd := doctorTestCmd(t, f, map[string]string{
			"all-contexts":       "true",
			"skip-version-check": "true",
		})

		var err error
		stdOut := test.CatchStdOut(t, func() {
			err = doctor(f, cmd)
		})

		assert.Nil(t, err)
		assert.Equal(t, "contextb", config.GetCurrentContextName())
		assert.Equal(t, 2, strings.Count(stdOut, "api authentication"))
		assert.Contains(t, stdOut, "contexta")
	})
}
