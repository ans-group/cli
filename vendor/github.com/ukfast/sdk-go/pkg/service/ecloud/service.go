package ecloud

import (
	"github.com/ukfast/sdk-go/pkg/connection"
)

// ECloudService is an interface for managing eCloud
type ECloudService interface {
	// Virtual Machine
	GetVirtualMachines(parameters connection.APIRequestParameters) ([]VirtualMachine, error)
	GetVirtualMachinesPaginated(parameters connection.APIRequestParameters) ([]VirtualMachine, error)
	GetVirtualMachine(vmID int) (VirtualMachine, error)
	CreateVirtualMachine(req CreateVirtualMachineRequest) (int, error)
	PatchVirtualMachine(vmID int, patch PatchVirtualMachineRequest) error
	CloneVirtualMachine(vmID int, req CloneVirtualMachineRequest) (int, error)
	DeleteVirtualMachine(vmID int) error
	PowerOnVirtualMachine(vmID int) error
	PowerOffVirtualMachine(vmID int) error
	PowerResetVirtualMachine(vmID int) error
	PowerShutdownVirtualMachine(vmID int) error
	PowerRestartVirtualMachine(vmID int) error
	CreateVirtualMachineTemplate(vmID int, req CreateVirtualMachineTemplateRequest) error
	GetVirtualMachineTags(vmID int, parameters connection.APIRequestParameters) ([]Tag, error)
	GetVirtualMachineTagsPaginated(vmID int, parameters connection.APIRequestParameters) ([]Tag, error)
	GetVirtualMachineTag(vmID int, tagKey string) (Tag, error)
	CreateVirtualMachineTag(vmID int, req CreateTagRequest) error
	PatchVirtualMachineTag(vmID int, tagKey string, patch PatchTagRequest) error
	DeleteVirtualMachineTag(vmID int, tagKey string) error

	// Solution
	GetSolutions(parameters connection.APIRequestParameters) ([]Solution, error)
	GetSolutionsPaginated(parameters connection.APIRequestParameters) ([]Solution, error)
	GetSolution(solutionID int) (Solution, error)
	PatchSolution(solutionID int, patch PatchSolutionRequest) (int, error)
	GetSolutionVirtualMachines(solutionID int, parameters connection.APIRequestParameters) ([]VirtualMachine, error)
	GetSolutionVirtualMachinesPaginated(solutionID int, parameters connection.APIRequestParameters) ([]VirtualMachine, error)
	GetSolutionSites(solutionID int, parameters connection.APIRequestParameters) ([]Site, error)
	GetSolutionSitesPaginated(solutionID int, parameters connection.APIRequestParameters) ([]Site, error)
	GetSolutionDatastores(solutionID int, parameters connection.APIRequestParameters) ([]Datastore, error)
	GetSolutionDatastoresPaginated(solutionID int, parameters connection.APIRequestParameters) ([]Datastore, error)
	GetSolutionHosts(solutionID int, parameters connection.APIRequestParameters) ([]Host, error)
	GetSolutionHostsPaginated(solutionID int, parameters connection.APIRequestParameters) ([]Host, error)
	GetSolutionNetworks(solutionID int, parameters connection.APIRequestParameters) ([]Network, error)
	GetSolutionNetworksPaginated(solutionID int, parameters connection.APIRequestParameters) ([]Network, error)
	GetSolutionFirewalls(solutionID int, parameters connection.APIRequestParameters) ([]Firewall, error)
	GetSolutionFirewallsPaginated(solutionID int, parameters connection.APIRequestParameters) ([]Firewall, error)
	GetSolutionTemplates(solutionID int, parameters connection.APIRequestParameters) ([]Template, error)
	GetSolutionTemplatesPaginated(solutionID int, parameters connection.APIRequestParameters) ([]Template, error)
	GetSolutionTemplate(solutionID int, templateName string) (Template, error)
	DeleteSolutionTemplate(solutionID int, templateName string) error
	RenameSolutionTemplate(solutionID int, templateName string, req RenameTemplateRequest) error
	GetSolutionTags(solutionID int, parameters connection.APIRequestParameters) ([]Tag, error)
	GetSolutionTagsPaginated(solutionID int, parameters connection.APIRequestParameters) ([]Tag, error)
	GetSolutionTag(solutionID int, tagKey string) (Tag, error)
	CreateSolutionTag(solutionID int, req CreateTagRequest) error
	PatchSolutionTag(solutionID int, tagKey string, patch PatchTagRequest) error
	DeleteSolutionTag(solutionID int, tagKey string) error

	// Site
	GetSites(parameters connection.APIRequestParameters) ([]Site, error)
	GetSitesPaginated(parameters connection.APIRequestParameters) ([]Site, error)
	GetSite(siteID int) (Site, error)

	// Host
	GetHosts(parameters connection.APIRequestParameters) ([]Host, error)
	GetHostsPaginated(parameters connection.APIRequestParameters) ([]Host, error)
	GetHost(hostID int) (Host, error)

	// Datastore
	GetDatastores(parameters connection.APIRequestParameters) ([]Datastore, error)
	GetDatastoresPaginated(parameters connection.APIRequestParameters) ([]Datastore, error)
	GetDatastore(datastoreID int) (Datastore, error)

	// Firewall
	GetFirewalls(parameters connection.APIRequestParameters) ([]Firewall, error)
	GetFirewallsPaginated(parameters connection.APIRequestParameters) ([]Firewall, error)
	GetFirewall(firewallID int) (Firewall, error)
	GetFirewallConfig(firewallID int) (FirewallConfig, error)

	// Pod
	GetPods(parameters connection.APIRequestParameters) ([]Pod, error)
	GetPodsPaginated(parameters connection.APIRequestParameters) ([]Pod, error)
	GetPod(podID int) (Pod, error)
	GetPodTemplates(podID int, parameters connection.APIRequestParameters) ([]Template, error)
	GetPodTemplatesPaginated(podID int, parameters connection.APIRequestParameters) ([]Template, error)
}

// Service implements ECloudService for managing
// eCloud via the UKFast API
type Service struct {
	connection connection.Connection
}

// NewService returns a new instance of eCloud Service
func NewService(connection connection.Connection) *Service {
	return &Service{
		connection: connection,
	}
}
