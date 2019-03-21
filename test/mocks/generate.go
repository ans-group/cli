package mocks

//go:generate mockgen -package mocks -destination mock_accountservice.go github.com/ukfast/sdk-go/pkg/service/account AccountService
//go:generate mockgen -package mocks -destination mock_ddosxservice.go github.com/ukfast/sdk-go/pkg/service/ddosx DDoSXService
//go:generate mockgen -package mocks -destination mock_ecloudservice.go github.com/ukfast/sdk-go/pkg/service/ecloud ECloudService
//go:generate mockgen -package mocks -destination mock_safednsservice.go github.com/ukfast/sdk-go/pkg/service/safedns SafeDNSService
//go:generate mockgen -package mocks -destination mock_sslservice.go github.com/ukfast/sdk-go/pkg/service/ssl SSLService
//go:generate mockgen -package mocks -destination mock_resourcelocatorprovider.go github.com/ukfast/cli/internal/pkg/resource ResourceLocatorProvider
