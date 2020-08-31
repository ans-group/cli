package mocks

//go:generate mockgen -package mocks -destination mock_accountservice.go github.com/ukfast/sdk-go/pkg/service/account AccountService
//go:generate mockgen -package mocks -destination mock_billingservice.go github.com/ukfast/sdk-go/pkg/service/billing BillingService
//go:generate mockgen -package mocks -destination mock_ddosxservice.go github.com/ukfast/sdk-go/pkg/service/ddosx DDoSXService
//go:generate mockgen -package mocks -destination mock_ecloudservice.go github.com/ukfast/sdk-go/pkg/service/ecloud ECloudService
//go:generate mockgen -package mocks -destination mock_ecloudflexservice.go github.com/ukfast/sdk-go/pkg/service/ecloudflex ECloudFlexService
//go:generate mockgen -package mocks -destination mock_safednsservice.go github.com/ukfast/sdk-go/pkg/service/safedns SafeDNSService
//go:generate mockgen -package mocks -destination mock_sslservice.go github.com/ukfast/sdk-go/pkg/service/ssl SSLService
//go:generate mockgen -package mocks -destination mock_registrarservice.go github.com/ukfast/sdk-go/pkg/service/registrar RegistrarService
//go:generate mockgen -package mocks -destination mock_resourcelocatorprovider.go github.com/ukfast/cli/internal/pkg/resource ResourceLocatorProvider
//go:generate mockgen -package mocks -destination mock_pssservice.go github.com/ukfast/sdk-go/pkg/service/pss PSSService
//go:generate mockgen -package mocks -destination mock_storageservice.go github.com/ukfast/sdk-go/pkg/service/storage StorageService
//go:generate mockgen -package mocks -destination mock_ltaasservice.go github.com/ukfast/sdk-go/pkg/service/ltaas LTaaSService
//go:generate mockgen -package mocks -destination mock_draasservice.go github.com/ukfast/sdk-go/pkg/service/draas DRaaSService
