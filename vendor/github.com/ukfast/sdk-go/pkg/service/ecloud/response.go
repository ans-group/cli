package ecloud

import "github.com/ukfast/sdk-go/pkg/connection"

// GetVirtualMachinesResponseBody represents the API response body from the GetVirtualMachines resource
type GetVirtualMachinesResponseBody struct {
	connection.APIResponseBody

	Data []VirtualMachine `json:"data"`
}

// GetVirtualMachineResponseBody represents the API response body from the GetVirtualMachine resource
type GetVirtualMachineResponseBody struct {
	connection.APIResponseBody

	Data VirtualMachine `json:"data"`
}

// GetTagsResponseBody represents the API response body from the GetTags resource
type GetTagsResponseBody struct {
	connection.APIResponseBody

	Data []Tag `json:"data"`
}

// GetTagResponseBody represents the API response body from the GetTag resource
type GetTagResponseBody struct {
	connection.APIResponseBody

	Data Tag `json:"data"`
}

// GetSolutionsResponseBody represents the API response body from the GetSolutions resource
type GetSolutionsResponseBody struct {
	connection.APIResponseBody

	Data []Solution `json:"data"`
}

// GetSolutionResponseBody represents the API response body from the GetSolution resource
type GetSolutionResponseBody struct {
	connection.APIResponseBody

	Data Solution `json:"data"`
}

// GetSitesResponseBody represents the API response body from the GetSites resource
type GetSitesResponseBody struct {
	connection.APIResponseBody

	Data []Site `json:"data"`
}

// GetSiteResponseBody represents the API response body from the GetSite resource
type GetSiteResponseBody struct {
	connection.APIResponseBody

	Data Site `json:"data"`
}

// GetHostsResponseBody represents the API response body from the GetHosts resource
type GetHostsResponseBody struct {
	connection.APIResponseBody

	Data []Host `json:"data"`
}

// GetHostResponseBody represents the API response body from the GetHost resource
type GetHostResponseBody struct {
	connection.APIResponseBody

	Data Host `json:"data"`
}

// GetDatastoresResponseBody represents the API response body from the GetDatastores resource
type GetDatastoresResponseBody struct {
	connection.APIResponseBody

	Data []Datastore `json:"data"`
}

// GetDatastoreResponseBody represents the API response body from the GetDatastore resource
type GetDatastoreResponseBody struct {
	connection.APIResponseBody

	Data Datastore `json:"data"`
}

// GetTemplatesResponseBody represents the API response body from the GetTemplates resource
type GetTemplatesResponseBody struct {
	connection.APIResponseBody

	Data []Template `json:"data"`
}

// GetTemplateResponseBody represents the API response body from the GetTemplate resource
type GetTemplateResponseBody struct {
	connection.APIResponseBody

	Data Template `json:"data"`
}

// GetNetworksResponseBody represents the API response body from the GetNetworks resource
type GetNetworksResponseBody struct {
	connection.APIResponseBody

	Data []Network `json:"data"`
}

// GetNetworkResponseBody represents the API response body from the GetNetwork resource
type GetNetworkResponseBody struct {
	connection.APIResponseBody

	Data Network `json:"data"`
}

// GetFirewallsResponseBody represents the API response body from the GetFirewalls resource
type GetFirewallsResponseBody struct {
	connection.APIResponseBody

	Data []Firewall `json:"data"`
}

// GetFirewallResponseBody represents the API response body from the GetFirewall resource
type GetFirewallResponseBody struct {
	connection.APIResponseBody

	Data Firewall `json:"data"`
}

// GetFirewallConfigResponseBody represents the API response body from the GetFirewallConfig resource
type GetFirewallConfigResponseBody struct {
	connection.APIResponseBody

	Data FirewallConfig `json:"data"`
}

// GetPodsResponseBody represents the API response body from the GetPods resource
type GetPodsResponseBody struct {
	connection.APIResponseBody

	Data []Pod `json:"data"`
}

// GetPodResponseBody represents the API response body from the GetPod resource
type GetPodResponseBody struct {
	connection.APIResponseBody

	Data Pod `json:"data"`
}

// GetAppliancesResponseBody represents the API response body from the GetAppliances resource
type GetAppliancesResponseBody struct {
	connection.APIResponseBody

	Data []Appliance `json:"data"`
}

// GetApplianceResponseBody represents the API response body from the GetAppliance resource
type GetApplianceResponseBody struct {
	connection.APIResponseBody

	Data Appliance `json:"data"`
}

// GetApplianceParametersResponseBody represents the API response body from the GetApplianceParameters resource
type GetApplianceParametersResponseBody struct {
	connection.APIResponseBody

	Data []ApplianceParameter `json:"data"`
}
