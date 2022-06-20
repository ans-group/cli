package mocks

//go:generate mockgen -package mocks -destination mock_accountservice.go -imports=github.com/ukfast/sdk-go/pkg/service/account github.com/ukfast/sdk-go/pkg/service/account AccountService
//go:generate mockgen -package mocks -destination mock_billingservice.go -imports=github.com/ukfast/sdk-go/pkg/service/billing github.com/ukfast/sdk-go/pkg/service/billing BillingService
//go:generate mockgen -package mocks -destination mock_ddosxservice.go -imports=github.com/ukfast/sdk-go/pkg/service/ddosx github.com/ukfast/sdk-go/pkg/service/ddosx DDoSXService
//go:generate mockgen -package mocks -destination mock_draasservice.go -imports=github.com/ukfast/sdk-go/pkg/service/draas github.com/ukfast/sdk-go/pkg/service/draas DRaaSService
//go:generate mockgen -package mocks -destination mock_ecloudservice.go -imports=github.com/ukfast/sdk-go/pkg/service/ecloud github.com/ukfast/sdk-go/pkg/service/ecloud ECloudService
//go:generate mockgen -package mocks -destination mock_ecloudflexservice.go -imports=github.com/ukfast/sdk-go/pkg/service/ecloudflex github.com/ukfast/sdk-go/pkg/service/ecloudflex ECloudFlexService
//go:generate mockgen -package mocks -destination mock_loadbalancerservice.go -imports=github.com/ukfast/sdk-go/pkg/service/loadbalancer github.com/ukfast/sdk-go/pkg/service/loadbalancer LoadBalancerService
//go:generate mockgen -package mocks -destination mock_cloudflare.go -imports=github.com/ukfast/sdk-go/pkg/service/cloudflare github.com/ukfast/sdk-go/pkg/service/cloudflare CloudflareService
//go:generate mockgen -package mocks -destination mock_pssservice.go -imports=github.com/ukfast/sdk-go/pkg/service/pss github.com/ukfast/sdk-go/pkg/service/pss PSSService
//go:generate mockgen -package mocks -destination mock_registrarservice.go -imports=github.com/ukfast/sdk-go/pkg/service/registrar github.com/ukfast/sdk-go/pkg/service/registrar RegistrarService
//go:generate mockgen -package mocks -destination mock_resourcelocatorprovider.go github.com/ukfast/cli/internal/pkg/resource ResourceLocatorProvider
//go:generate mockgen -package mocks -destination mock_safednsservice.go -imports=github.com/ukfast/sdk-go/pkg/service/safedns github.com/ukfast/sdk-go/pkg/service/safedns SafeDNSService
//go:generate mockgen -package mocks -destination mock_sharedexchangeservice.go -imports=github.com/ukfast/sdk-go/pkg/service/sharedexchange github.com/ukfast/sdk-go/pkg/service/sharedexchange SharedExchangeService
//go:generate mockgen -package mocks -destination mock_sslservice.go -imports=github.com/ukfast/sdk-go/pkg/service/ssl github.com/ukfast/sdk-go/pkg/service/ssl SSLService
//go:generate mockgen -package mocks -destination mock_storageservice.go -imports=github.com/ukfast/sdk-go/pkg/service/storage github.com/ukfast/sdk-go/pkg/service/storage StorageService
