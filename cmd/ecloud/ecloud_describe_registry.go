package ecloud

import (
	"slices"

	"github.com/ans-group/sdk-go/pkg/connection"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
)

// describeDescriptor describes how to fetch and present one eCloud resource type.
type describeDescriptor struct {
	// Prefix is the ID prefix, without the trailing hyphen (e.g. "i", "vpc").
	Prefix string
	// Name is the human-readable singular type name (e.g. "Instance").
	Name string
	// Fetch retrieves the resource itself.
	Fetch func(service ecloud.ECloudService, id string) (any, error)
	// Children retrieves related child collections.
	Children []describeChild
}

// describeChild describes one child collection of a resource.
type describeChild struct {
	// Title is the section heading (e.g. "NICs").
	Title string
	// Fetch retrieves a single page of the collection, returning the items as a
	// slice along with the total number of items available.
	Fetch func(service ecloud.ECloudService, id string, params connection.APIRequestParameters) (items any, total int, err error)
}

// describeFetch adapts a typed single-resource SDK getter to the untyped fetch
// signature used by describeDescriptor.
func describeFetch[T any](fn func(service ecloud.ECloudService, id string) (T, error)) func(ecloud.ECloudService, string) (any, error) {
	return func(service ecloud.ECloudService, id string) (any, error) {
		resource, err := fn(service, id)
		if err != nil {
			return nil, err
		}

		return resource, nil
	}
}

// describeChildren builds a describeChild from a typed paginated SDK getter. The
// collection type C is given explicitly so that items are returned wrapped in the
// relevant output collection type (e.g. NICCollection), allowing the renderer to
// reuse its DefaultColumns.
//
// The paginated getters are used deliberately: their non-paginated counterparts
// retrieve every page of a collection, so could not honour --child-limit.
func describeChildren[C ~[]T, T any](title string, fn func(service ecloud.ECloudService, id string, parameters connection.APIRequestParameters) (*connection.Paginated[T], error)) describeChild {
	return describeChild{
		Title: title,
		Fetch: func(service ecloud.ECloudService, id string, params connection.APIRequestParameters) (any, int, error) {
			paginated, err := fn(service, id, params)
			if err != nil {
				return nil, 0, err
			}

			return C(paginated.Items()), paginated.Total(), nil
		},
	}
}

// ecloudDescribeDescriptors is the registry of describable eCloud v2 resource types,
// keyed by the type-prefixed ID scheme used throughout the eCloud v2 API.
//
// eCloud v1 resources are deliberately absent, as they are identified by bare
// integer IDs which carry no type information. Discount plans (prefix "disc") are
// also absent, as the SDK provides no single-resource getter for them.
var ecloudDescribeDescriptors = []describeDescriptor{
	{
		Prefix: "ar",
		Name:   "Affinity rule",
		Fetch:  describeFetch(ecloud.ECloudService.GetAffinityRule),
		Children: []describeChild{
			describeChildren[AffinityRuleMemberCollection]("Members", ecloud.ECloudService.GetAffinityRuleMembersPaginated),
		},
	},
	{
		Prefix: "arm",
		Name:   "Affinity rule member",
		Fetch:  describeFetch(ecloud.ECloudService.GetAffinityRuleMember),
	},
	{
		Prefix: "az",
		Name:   "Availability zone",
		Fetch:  describeFetch(ecloud.ECloudService.GetAvailabilityZone),
		Children: []describeChild{
			describeChildren[IOPSTierCollection]("IOPS tiers", ecloud.ECloudService.GetAvailabilityZoneIOPSTiersPaginated),
		},
	},
	{
		Prefix: "bkupgs",
		Name:   "Backup gateway specification",
		Fetch:  describeFetch(ecloud.ECloudService.GetBackupGatewaySpecification),
	},
	{
		Prefix: "bkupgw",
		Name:   "Backup gateway",
		Fetch:  describeFetch(ecloud.ECloudService.GetBackupGateway),
	},
	{
		Prefix: "bm",
		Name:   "Billing metric",
		Fetch:  describeFetch(ecloud.ECloudService.GetBillingMetric),
	},
	{
		Prefix: "dhcp",
		Name:   "DHCP",
		Fetch:  describeFetch(ecloud.ECloudService.GetDHCP),
		Children: []describeChild{
			describeChildren[TaskCollection]("Tasks", ecloud.ECloudService.GetDHCPTasksPaginated),
		},
	},
	{
		Prefix: "fip",
		Name:   "Floating IP",
		Fetch:  describeFetch(ecloud.ECloudService.GetFloatingIP),
		Children: []describeChild{
			describeChildren[TaskCollection]("Tasks", ecloud.ECloudService.GetFloatingIPTasksPaginated),
		},
	},
	{
		Prefix: "fwp",
		Name:   "Firewall policy",
		Fetch:  describeFetch(ecloud.ECloudService.GetFirewallPolicy),
		Children: []describeChild{
			describeChildren[FirewallRuleCollection]("Firewall rules", ecloud.ECloudService.GetFirewallPolicyFirewallRulesPaginated),
			describeChildren[TaskCollection]("Tasks", ecloud.ECloudService.GetFirewallPolicyTasksPaginated),
		},
	},
	{
		Prefix: "fwr",
		Name:   "Firewall rule",
		Fetch:  describeFetch(ecloud.ECloudService.GetFirewallRule),
		Children: []describeChild{
			describeChildren[FirewallRulePortCollection]("Ports", ecloud.ECloudService.GetFirewallRuleFirewallRulePortsPaginated),
		},
	},
	{
		Prefix: "fwrp",
		Name:   "Firewall rule port",
		Fetch:  describeFetch(ecloud.ECloudService.GetFirewallRulePort),
	},
	{
		Prefix: "h",
		Name:   "Host",
		Fetch:  describeFetch(ecloud.ECloudService.GetHost),
		Children: []describeChild{
			describeChildren[TaskCollection]("Tasks", ecloud.ECloudService.GetHostTasksPaginated),
		},
	},
	{
		Prefix: "hg",
		Name:   "Host group",
		Fetch:  describeFetch(ecloud.ECloudService.GetHostGroup),
		Children: []describeChild{
			describeChildren[TaskCollection]("Tasks", ecloud.ECloudService.GetHostGroupTasksPaginated),
		},
	},
	{
		Prefix: "hs",
		Name:   "Host specification",
		Fetch:  describeFetch(ecloud.ECloudService.GetHostSpec),
	},
	{
		Prefix: "i",
		Name:   "Instance",
		Fetch:  describeFetch(ecloud.ECloudService.GetInstance),
		Children: []describeChild{
			describeChildren[NICCollection]("NICs", ecloud.ECloudService.GetInstanceNICsPaginated),
			describeChildren[VolumeCollection]("Volumes", ecloud.ECloudService.GetInstanceVolumesPaginated),
			describeChildren[FloatingIPCollection]("Floating IPs", ecloud.ECloudService.GetInstanceFloatingIPsPaginated),
			describeChildren[CredentialCollection]("Credentials", ecloud.ECloudService.GetInstanceCredentialsPaginated),
			describeChildren[TaskCollection]("Tasks", ecloud.ECloudService.GetInstanceTasksPaginated),
		},
	},
	{
		Prefix: "img",
		Name:   "Image",
		Fetch:  describeFetch(ecloud.ECloudService.GetImage),
		Children: []describeChild{
			describeChildren[ImageParameterCollection]("Parameters", ecloud.ECloudService.GetImageParametersPaginated),
			describeChildren[ImageMetadataCollection]("Metadata", ecloud.ECloudService.GetImageMetadataPaginated),
		},
	},
	{
		Prefix: "iops",
		Name:   "IOPS tier",
		Fetch:  describeFetch(ecloud.ECloudService.GetIOPSTier),
	},
	{
		Prefix: "ip",
		Name:   "IP address",
		Fetch:  describeFetch(ecloud.ECloudService.GetIPAddress),
	},
	{
		Prefix: "lb",
		Name:   "Load balancer",
		Fetch:  describeFetch(ecloud.ECloudService.GetLoadBalancer),
	},
	{
		Prefix: "lbs",
		Name:   "Load balancer specification",
		Fetch:  describeFetch(ecloud.ECloudService.GetLoadBalancerSpec),
	},
	{
		Prefix: "mgw",
		Name:   "Monitoring gateway",
		Fetch:  describeFetch(ecloud.ECloudService.GetMonitoringGateway),
	},
	{
		Prefix: "net",
		Name:   "Network",
		Fetch:  describeFetch(ecloud.ECloudService.GetNetwork),
		Children: []describeChild{
			describeChildren[NICCollection]("NICs", ecloud.ECloudService.GetNetworkNICsPaginated),
			describeChildren[TaskCollection]("Tasks", ecloud.ECloudService.GetNetworkTasksPaginated),
		},
	},
	{
		Prefix: "nic",
		Name:   "NIC",
		Fetch:  describeFetch(ecloud.ECloudService.GetNIC),
		Children: []describeChild{
			describeChildren[IPAddressCollection]("IP addresses", ecloud.ECloudService.GetNICIPAddressesPaginated),
			describeChildren[TaskCollection]("Tasks", ecloud.ECloudService.GetNICTasksPaginated),
		},
	},
	{
		Prefix: "nor",
		Name:   "NAT overload rule",
		Fetch:  describeFetch(ecloud.ECloudService.GetNATOverloadRule),
	},
	{
		Prefix: "np",
		Name:   "Network policy",
		Fetch:  describeFetch(ecloud.ECloudService.GetNetworkPolicy),
		Children: []describeChild{
			describeChildren[NetworkRuleCollection]("Network rules", ecloud.ECloudService.GetNetworkPolicyNetworkRulesPaginated),
			describeChildren[TaskCollection]("Tasks", ecloud.ECloudService.GetNetworkPolicyTasksPaginated),
		},
	},
	{
		Prefix: "nr",
		Name:   "Network rule",
		Fetch:  describeFetch(ecloud.ECloudService.GetNetworkRule),
		Children: []describeChild{
			describeChildren[NetworkRulePortCollection]("Ports", ecloud.ECloudService.GetNetworkRuleNetworkRulePortsPaginated),
		},
	},
	{
		Prefix: "nrp",
		Name:   "Network rule port",
		Fetch:  describeFetch(ecloud.ECloudService.GetNetworkRulePort),
	},
	{
		Prefix: "reg",
		Name:   "Region",
		Fetch:  describeFetch(ecloud.ECloudService.GetRegion),
	},
	{
		Prefix: "rt",
		Name:   "Resource tier",
		Fetch:  describeFetch(ecloud.ECloudService.GetResourceTier),
	},
	{
		Prefix: "rtp",
		Name:   "Router throughput",
		Fetch:  describeFetch(ecloud.ECloudService.GetRouterThroughput),
	},
	{
		Prefix: "rtr",
		Name:   "Router",
		Fetch:  describeFetch(ecloud.ECloudService.GetRouter),
		Children: []describeChild{
			describeChildren[NetworkCollection]("Networks", ecloud.ECloudService.GetRouterNetworksPaginated),
			describeChildren[FirewallPolicyCollection]("Firewall policies", ecloud.ECloudService.GetRouterFirewallPoliciesPaginated),
			describeChildren[[]ecloud.VPN]("VPNs", ecloud.ECloudService.GetRouterVPNsPaginated),
			describeChildren[TaskCollection]("Tasks", ecloud.ECloudService.GetRouterTasksPaginated),
		},
	},
	{
		Prefix: "ssh",
		Name:   "SSH key pair",
		Fetch:  describeFetch(ecloud.ECloudService.GetSSHKeyPair),
	},
	{
		Prefix: "tag",
		Name:   "Tag",
		Fetch:  describeFetch(ecloud.ECloudService.GetTag),
	},
	{
		Prefix: "task",
		Name:   "Task",
		Fetch:  describeFetch(ecloud.ECloudService.GetTask),
	},
	{
		Prefix: "vip",
		Name:   "VIP",
		Fetch:  describeFetch(ecloud.ECloudService.GetVIP),
	},
	{
		Prefix: "vol",
		Name:   "Volume",
		Fetch:  describeFetch(ecloud.ECloudService.GetVolume),
		Children: []describeChild{
			describeChildren[InstanceCollection]("Instances", ecloud.ECloudService.GetVolumeInstancesPaginated),
			describeChildren[TaskCollection]("Tasks", ecloud.ECloudService.GetVolumeTasksPaginated),
		},
	},
	{
		Prefix: "volgroup",
		Name:   "Volume group",
		Fetch:  describeFetch(ecloud.ECloudService.GetVolumeGroup),
		Children: []describeChild{
			describeChildren[VolumeCollection]("Volumes", ecloud.ECloudService.GetVolumeGroupVolumesPaginated),
		},
	},
	{
		Prefix: "vpc",
		Name:   "VPC",
		Fetch:  describeFetch(ecloud.ECloudService.GetVPC),
		Children: []describeChild{
			describeChildren[InstanceCollection]("Instances", ecloud.ECloudService.GetVPCInstancesPaginated),
			describeChildren[VolumeCollection]("Volumes", ecloud.ECloudService.GetVPCVolumesPaginated),
			describeChildren[TaskCollection]("Tasks", ecloud.ECloudService.GetVPCTasksPaginated),
		},
	},
	{
		Prefix: "vpn",
		Name:   "VPN service",
		Fetch:  describeFetch(ecloud.ECloudService.GetVPNService),
	},
	{
		Prefix: "vpne",
		Name:   "VPN endpoint",
		Fetch:  describeFetch(ecloud.ECloudService.GetVPNEndpoint),
	},
	{
		Prefix: "vpng",
		Name:   "VPN gateway",
		Fetch:  describeFetch(ecloud.ECloudService.GetVPNGateway),
		Children: []describeChild{
			describeChildren[TaskCollection]("Tasks", ecloud.ECloudService.GetVPNGatewayTasksPaginated),
		},
	},
	{
		Prefix: "vpngs",
		Name:   "VPN gateway specification",
		Fetch:  describeFetch(ecloud.ECloudService.GetVPNGatewaySpecification),
		Children: []describeChild{
			describeChildren[AvailabilityZoneCollection]("Availability zones", ecloud.ECloudService.GetVPNGatewaySpecificationAvailabilityZonesPaginated),
		},
	},
	{
		Prefix: "vpngu",
		Name:   "VPN gateway user",
		Fetch:  describeFetch(ecloud.ECloudService.GetVPNGatewayUser),
	},
	{
		Prefix: "vpnpg",
		Name:   "VPN profile group",
		Fetch:  describeFetch(ecloud.ECloudService.GetVPNProfileGroup),
	},
	{
		Prefix: "vpns",
		Name:   "VPN session",
		Fetch:  describeFetch(ecloud.ECloudService.GetVPNSession),
	},
}

// ecloudDescribeDescriptorsByPrefix indexes ecloudDescribeDescriptors by ID prefix.
var ecloudDescribeDescriptorsByPrefix = describeDescriptorsByPrefix(ecloudDescribeDescriptors)

func describeDescriptorsByPrefix(descriptors []describeDescriptor) map[string]describeDescriptor {
	index := make(map[string]describeDescriptor, len(descriptors))
	for _, descriptor := range descriptors {
		index[descriptor.Prefix] = descriptor
	}

	return index
}

// ecloudDescribePrefixes returns every supported ID prefix, sorted alphabetically.
func ecloudDescribePrefixes() []string {
	prefixes := make([]string, 0, len(ecloudDescribeDescriptors))
	for _, descriptor := range ecloudDescribeDescriptors {
		prefixes = append(prefixes, descriptor.Prefix)
	}
	slices.Sort(prefixes)

	return prefixes
}
