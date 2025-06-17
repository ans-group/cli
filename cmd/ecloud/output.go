package ecloud

import (
	"strconv"

	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
)

type VirtualMachineCollection []ecloud.VirtualMachine

func (m VirtualMachineCollection) DefaultColumns() []string {
	return []string{"id", "name", "cpu", "ram_gb", "hdd_gb", "ip_internal", "ip_external", "status", "power_status"}
}

func (m VirtualMachineCollection) Fields() []*output.OrderedFields {
	var data []*output.OrderedFields
	for _, vm := range m {
		fields := output.NewOrderedFields()
		fields.Set("id", strconv.Itoa(vm.ID))
		fields.Set("name", vm.Name)
		fields.Set("hostname", vm.Hostname)
		fields.Set("computername", vm.ComputerName)
		fields.Set("cpu", strconv.Itoa(vm.CPU))
		fields.Set("ram_gb", strconv.Itoa(vm.RAM))
		fields.Set("hdd_gb", strconv.Itoa(vm.HDD))
		fields.Set("ip_internal", vm.IPInternal.String())
		fields.Set("ip_external", vm.IPExternal.String())
		fields.Set("platform", vm.Platform)
		fields.Set("template", vm.Template)
		fields.Set("backup", strconv.FormatBool(vm.Backup))
		fields.Set("support", strconv.FormatBool(vm.Support))
		fields.Set("environment", vm.Environment)
		fields.Set("solution_id", strconv.Itoa(vm.SolutionID))
		fields.Set("status", string(vm.Status))
		fields.Set("power_status", vm.PowerStatus)
		fields.Set("tools_status", vm.ToolsStatus)
		fields.Set("encrypted", strconv.FormatBool(vm.Encrypted))
		fields.Set("role", vm.Role)
		fields.Set("gpu_profile", vm.GPUProfile)

		data = append(data, fields)
	}

	return data
}

type VirtualMachineDiskCollection []ecloud.VirtualMachineDisk

func (m VirtualMachineDiskCollection) DefaultColumns() []string {
	return []string{"name", "capacity", "uuid", "type"}
}

func (m VirtualMachineDiskCollection) Fields() []*output.OrderedFields {
	var data []*output.OrderedFields
	for _, disk := range m {
		fields := output.NewOrderedFields()
		fields.Set("name", disk.Name)
		fields.Set("capacity", strconv.Itoa(disk.Capacity))
		fields.Set("uuid", disk.UUID)
		fields.Set("type", string(disk.Type))
		fields.Set("key", strconv.Itoa(disk.Key))

		data = append(data, fields)
	}

	return data
}

type TagCollection []ecloud.Tag

func (m TagCollection) DefaultColumns() []string {
	return []string{"key", "value", "created_at"}
}

func (m TagCollection) Fields() []*output.OrderedFields {
	var data []*output.OrderedFields
	for _, tag := range m {
		fields := output.NewOrderedFields()
		fields.Set("key", tag.Key)
		fields.Set("value", tag.Value)
		fields.Set("created_at", tag.CreatedAt.String())

		data = append(data, fields)
	}

	return data
}

type SolutionCollection []ecloud.Solution

func (m SolutionCollection) DefaultColumns() []string {
	return []string{"id", "name", "environment", "pod_id"}
}

func (m SolutionCollection) Fields() []*output.OrderedFields {
	var data []*output.OrderedFields
	for _, solution := range m {
		fields := output.NewOrderedFields()
		fields.Set("id", strconv.Itoa(solution.ID))
		fields.Set("name", solution.Name)
		fields.Set("environment", string(solution.Environment))
		fields.Set("pod_id", strconv.Itoa(solution.PodID))
		fields.Set("encryption_enabled", strconv.FormatBool(solution.EncryptionEnabled))
		fields.Set("encryption_default", strconv.FormatBool(solution.EncryptionDefault))

		data = append(data, fields)
	}

	return data
}

type SiteCollection []ecloud.Site

func (m SiteCollection) DefaultColumns() []string {
	return []string{"id", "state", "solution_id", "pod_id"}
}

func (m SiteCollection) Fields() []*output.OrderedFields {
	var data []*output.OrderedFields
	for _, site := range m {
		fields := output.NewOrderedFields()
		fields.Set("id", strconv.Itoa(site.ID))
		fields.Set("state", site.State)
		fields.Set("solution_id", strconv.Itoa(site.SolutionID))
		fields.Set("pod_id", strconv.Itoa(site.PodID))

		data = append(data, fields)
	}

	return data
}

type V1HostCollection []ecloud.V1Host

func (m V1HostCollection) DefaultColumns() []string {
	return []string{"id", "solution_id", "pod_id", "name", "cpu_quantity", "cpu_cores", "ram_capacity"}
}

func (m V1HostCollection) Fields() []*output.OrderedFields {
	var data []*output.OrderedFields
	for _, host := range m {
		fields := output.NewOrderedFields()
		fields.Set("id", strconv.Itoa(host.ID))
		fields.Set("solution_id", strconv.Itoa(host.SolutionID))
		fields.Set("pod_id", strconv.Itoa(host.PodID))
		fields.Set("name", host.Name)
		fields.Set("cpu_quantity", strconv.Itoa(host.CPU.Quantity))
		fields.Set("cpu_cores", strconv.Itoa(host.CPU.Cores))
		fields.Set("cpu_speed", host.CPU.Speed)
		fields.Set("ram_capacity", strconv.Itoa(host.RAM.Capacity))
		fields.Set("ram_reserved", strconv.Itoa(host.RAM.Reserved))
		fields.Set("ram_allocated", strconv.Itoa(host.RAM.Allocated))
		fields.Set("ram_available", strconv.Itoa(host.RAM.Available))

		data = append(data, fields)
	}

	return data
}

type DatastoreCollection []ecloud.Datastore

func (m DatastoreCollection) DefaultColumns() []string {
	return []string{"id", "solution_id", "site_id", "name", "status", "capacity"}
}

func (m DatastoreCollection) Fields() []*output.OrderedFields {
	var data []*output.OrderedFields
	for _, datastore := range m {
		fields := output.NewOrderedFields()
		fields.Set("id", strconv.Itoa(datastore.ID))
		fields.Set("solution_id", strconv.Itoa(datastore.SolutionID))
		fields.Set("site_id", strconv.Itoa(datastore.SiteID))
		fields.Set("name", datastore.Name)
		fields.Set("status", string(datastore.Status))
		fields.Set("capacity", strconv.Itoa(datastore.Capacity))
		fields.Set("allocated", strconv.Itoa(datastore.Allocated))
		fields.Set("available", strconv.Itoa(datastore.Available))

		data = append(data, fields)
	}

	return data
}

type TemplateCollection []ecloud.Template

func (m TemplateCollection) DefaultColumns() []string {
	return []string{"name", "cpu", "ram_gb", "hdd_gb", "platform", "operating_system", "solution_id"}
}

func (m TemplateCollection) Fields() []*output.OrderedFields {
	var data []*output.OrderedFields
	for _, template := range m {
		fields := output.NewOrderedFields()
		fields.Set("name", template.Name)
		fields.Set("cpu", strconv.Itoa(template.CPU))
		fields.Set("ram_gb", strconv.Itoa(template.RAM))
		fields.Set("hdd_gb", strconv.Itoa(template.HDD))
		fields.Set("platform", template.Platform)
		fields.Set("operating_system", template.OperatingSystem)
		fields.Set("solution_id", strconv.Itoa(template.SolutionID))

		data = append(data, fields)
	}

	return data
}

type V1NetworkCollection []ecloud.V1Network

func (m V1NetworkCollection) DefaultColumns() []string {
	return []string{"id", "name"}
}

func (m V1NetworkCollection) Fields() []*output.OrderedFields {
	var data []*output.OrderedFields
	for _, network := range m {
		fields := output.NewOrderedFields()
		fields.Set("id", strconv.Itoa(network.ID))
		fields.Set("name", network.Name)

		data = append(data, fields)
	}

	return data
}

type FirewallCollection []ecloud.Firewall

func (m FirewallCollection) DefaultColumns() []string {
	return []string{"id", "name", "hostname", "ip", "role"}
}

func (m FirewallCollection) Fields() []*output.OrderedFields {
	var data []*output.OrderedFields
	for _, firewall := range m {
		fields := output.NewOrderedFields()
		fields.Set("id", strconv.Itoa(firewall.ID))
		fields.Set("name", firewall.Name)
		fields.Set("hostname", firewall.Hostname)
		fields.Set("ip", firewall.IP.String())
		fields.Set("role", string(firewall.Role))

		data = append(data, fields)
	}

	return data
}

type PodCollection []ecloud.Pod

func (m PodCollection) DefaultColumns() []string {
	return []string{"id", "name"}
}

func (m PodCollection) Fields() []*output.OrderedFields {
	var data []*output.OrderedFields
	for _, pod := range m {
		fields := output.NewOrderedFields()
		fields.Set("id", strconv.Itoa(pod.ID))
		fields.Set("name", pod.Name)

		data = append(data, fields)
	}

	return data
}

type ApplianceCollection []ecloud.Appliance

func (m ApplianceCollection) DefaultColumns() []string {
	return []string{"id", "name", "publisher", "created_at"}
}

func (m ApplianceCollection) Fields() []*output.OrderedFields {
	var data []*output.OrderedFields
	for _, appliance := range m {
		fields := output.NewOrderedFields()
		fields.Set("id", appliance.ID)
		fields.Set("name", appliance.Name)
		fields.Set("logo_uri", appliance.LogoURI)
		fields.Set("description", appliance.Description)
		fields.Set("documentation_uri", appliance.DocumentationURI)
		fields.Set("publisher", appliance.Publisher)
		fields.Set("created_at", appliance.CreatedAt.String())

		data = append(data, fields)
	}

	return data
}

type ApplianceParameterCollection []ecloud.ApplianceParameter

func (m ApplianceParameterCollection) DefaultColumns() []string {
	return []string{"id", "name", "key", "type", "description", "required"}
}

func (m ApplianceParameterCollection) Fields() []*output.OrderedFields {
	var data []*output.OrderedFields
	for _, parameter := range m {
		fields := output.NewOrderedFields()
		fields.Set("id", parameter.ID)
		fields.Set("name", parameter.Name)
		fields.Set("key", parameter.Key)
		fields.Set("type", parameter.Type)
		fields.Set("description", parameter.Description)
		fields.Set("required", strconv.FormatBool(parameter.Required))
		fields.Set("validation_rule", parameter.ValidationRule)

		data = append(data, fields)
	}

	return data
}

type ConsoleSessionCollection []ecloud.ConsoleSession

func (m ConsoleSessionCollection) DefaultColumns() []string {
	return []string{"url"}
}

func (m ConsoleSessionCollection) Fields() []*output.OrderedFields {
	var data []*output.OrderedFields
	for _, session := range m {
		fields := output.NewOrderedFields()
		fields.Set("url", session.URL)

		data = append(data, fields)
	}

	return data
}

type VPCCollection []ecloud.VPC

func (m VPCCollection) DefaultColumns() []string {
	return []string{"id", "name", "region_id", "sync_status"}
}

type InstanceCollection []ecloud.Instance

func (m InstanceCollection) DefaultColumns() []string {
	return []string{"id", "name", "vpc_id", "vcpu_sockets", "vcpu_cores_per_socket", "ram_capacity", "sync_status"}
}

type FloatingIPCollection []ecloud.FloatingIP

func (m FloatingIPCollection) DefaultColumns() []string {
	return []string{"id", "name", "ip_address", "sync_status"}
}

type FirewallPolicyCollection []ecloud.FirewallPolicy

func (m FirewallPolicyCollection) DefaultColumns() []string {
	return []string{"id", "name", "router_id", "sync_status"}
}

type FirewallRuleCollection []ecloud.FirewallRule

func (m FirewallRuleCollection) DefaultColumns() []string {
	return []string{"id", "name", "firewall_policy_id", "source", "destination", "action", "direction", "enabled"}
}

type FirewallRulePortCollection []ecloud.FirewallRulePort

func (m FirewallRulePortCollection) DefaultColumns() []string {
	return []string{"id", "name", "firewall_rule_id", "protocol", "source", "destination"}
}

type RegionCollection []ecloud.Region

func (m RegionCollection) DefaultColumns() []string {
	return []string{"id", "name"}
}

type VolumeCollection []ecloud.Volume

func (m VolumeCollection) DefaultColumns() []string {
	return []string{"id", "name", "type", "capacity", "sync_status"}
}

type CredentialCollection []ecloud.Credential

func (m CredentialCollection) DefaultColumns() []string {
	return []string{"id", "name", "username", "password"}
}

type NICCollection []ecloud.NIC

func (m NICCollection) DefaultColumns() []string {
	return []string{"id", "mac_address", "instance_id", "network_id", "ip_address"}
}

type RouterCollection []ecloud.Router

func (m RouterCollection) DefaultColumns() []string {
	return []string{"id", "name", "vpc_id", "availability_zone_id", "sync_status"}
}

type NetworkCollection []ecloud.Network

func (m NetworkCollection) DefaultColumns() []string {
	return []string{"id", "name", "router_id", "subnet", "sync_status"}
}

type DHCPCollection []ecloud.DHCP

func (m DHCPCollection) DefaultColumns() []string {
	return []string{"id", "vpc_id", "availability_zone_id", "sync_status"}
}

type RouterThroughputCollection []ecloud.RouterThroughput

func (m RouterThroughputCollection) DefaultColumns() []string {
	return []string{"id", "availability_zone_id", "name"}
}

type ImageCollection []ecloud.Image

func (m ImageCollection) DefaultColumns() []string {
	return []string{"id", "name"}
}

type ImageParameterCollection []ecloud.ImageParameter

func (m ImageParameterCollection) DefaultColumns() []string {
	return []string{"key", "type", "required"}
}

type ImageMetadataCollection []ecloud.ImageMetadata

func (m ImageMetadataCollection) DefaultColumns() []string {
	return []string{"key", "value"}
}

type SSHKeyPairCollection []ecloud.SSHKeyPair

func (m SSHKeyPairCollection) DefaultColumns() []string {
	return []string{"id", "name"}
}

type TaskCollection []ecloud.Task

func (m TaskCollection) DefaultColumns() []string {
	return []string{"id", "resource_id", "name", "status", "created_at", "updated_at"}
}

type NetworkPolicyCollection []ecloud.NetworkPolicy

func (m NetworkPolicyCollection) DefaultColumns() []string {
	return []string{"id", "name", "router_id", "sync_status"}
}

type NetworkRuleCollection []ecloud.NetworkRule

func (m NetworkRuleCollection) DefaultColumns() []string {
	return []string{"id", "name", "network_policy_id", "source", "destination", "action", "direction", "enabled"}
}

type NetworkRulePortCollection []ecloud.NetworkRulePort

func (m NetworkRulePortCollection) DefaultColumns() []string {
	return []string{"id", "name", "network_rule_id", "protocol", "source", "destination"}
}

type HostGroupCollection []ecloud.HostGroup

func (m HostGroupCollection) DefaultColumns() []string {
	return []string{"id", "name", "vpc_id", "sync_status"}
}

type HostCollection []ecloud.Host

func (m HostCollection) DefaultColumns() []string {
	return []string{"id", "name", "host_group_id", "sync_status"}
}

type HostSpecCollection []ecloud.HostSpec

func (m HostSpecCollection) DefaultColumns() []string {
	return []string{"id", "name", "cpu_sockets", "cpu_cores", "cpu_type", "cpu_clock_speed", "ram_capacity"}
}

type AvailabilityZoneCollection []ecloud.AvailabilityZone

func (m AvailabilityZoneCollection) DefaultColumns() []string {
	return []string{"id", "name", "region_id"}
}

type VPNServiceCollection []ecloud.VPNService

func (m VPNServiceCollection) DefaultColumns() []string {
	return []string{"id", "name", "router_id", "vpc_id", "sync_status"}
}

type VPNEndpointCollection []ecloud.VPNEndpoint

func (m VPNEndpointCollection) DefaultColumns() []string {
	return []string{"id", "name", "vpn_service_id", "floating_ip_id", "sync_status"}
}

type VPNSessionCollection []ecloud.VPNSession

func (m VPNSessionCollection) DefaultColumns() []string {
	return []string{"id", "name", "vpn_service_id", "vpn_endpoint_id", "remote_ip", "sync_status", "tunnel_details_session_state"}
}

type VPNProfileGroupCollection []ecloud.VPNProfileGroup

func (m VPNProfileGroupCollection) DefaultColumns() []string {
	return []string{"id", "name", "availability_zone_id"}
}

type VPNSessionPreSharedKeyCollection []ecloud.VPNSessionPreSharedKey

func (m VPNSessionPreSharedKeyCollection) DefaultColumns() []string {
	return []string{"psk"}
}

type VolumeGroupCollection []ecloud.VolumeGroup

func (m VolumeGroupCollection) DefaultColumns() []string {
	return []string{"id", "name", "vpc_id", "availability_zone_id", "sync_status"}
}

type LoadBalancerCollection []ecloud.LoadBalancer

func (m LoadBalancerCollection) DefaultColumns() []string {
	return []string{"id", "name", "vpc_id", "availability_zone_id", "sync_status"}
}

type LoadBalancerSpecCollection []ecloud.LoadBalancerSpec

func (m LoadBalancerSpecCollection) DefaultColumns() []string {
	return []string{"id", "name", "description"}
}

type VIPCollection []ecloud.VIP

func (m VIPCollection) DefaultColumns() []string {
	return []string{"id", "name", "load_balancer_id", "ip_address_id", "config_id", "sync_status"}
}

type IPAddressCollection []ecloud.IPAddress

func (m IPAddressCollection) DefaultColumns() []string {
	return []string{"id", "name", "ip_address", "sync_status"}
}

type AffinityRuleCollection []ecloud.AffinityRule

func (m AffinityRuleCollection) DefaultColumns() []string {
	return []string{"id", "name", "vpc_id", "availability_zone_id", "type", "sync_status"}
}

type AffinityRuleMemberCollection []ecloud.AffinityRuleMember

func (m AffinityRuleMemberCollection) DefaultColumns() []string {
	return []string{"id", "instance_id", "affinity_rule_id", "sync_status"}
}

type ResourceTierCollection []ecloud.ResourceTier

func (m ResourceTierCollection) DefaultColumns() []string {
	return []string{"id", "name", "availability_zone_id"}
}

type NATOverloadRuleCollection []ecloud.NATOverloadRule

func (m NATOverloadRuleCollection) DefaultColumns() []string {
	return []string{"id", "name", "network_id", "subnet", "floating_ip_id", "action"}
}

type IOPSTierCollection []ecloud.IOPSTier

func (m IOPSTierCollection) DefaultColumns() []string {
	return []string{"id", "name", "level"}
}

type VPNGatewayCollection []ecloud.VPNGateway

func (m VPNGatewayCollection) DefaultColumns() []string {
	return []string{"id", "name", "fqdn", "router_id", "specification_id", "sync_status"}
}

type VPNGatewaySpecificationCollection []ecloud.VPNGatewaySpecification

func (m VPNGatewaySpecificationCollection) DefaultColumns() []string {
	return []string{"id", "name", "description"}
}

type VPNGatewayUserCollection []ecloud.VPNGatewayUser

func (m VPNGatewayUserCollection) DefaultColumns() []string {
	return []string{"id", "name", "username", "vpn_gateway_id", "sync_status"}
}

type BackupGatewaySpecificationCollection []ecloud.BackupGatewaySpecification

func (m BackupGatewaySpecificationCollection) DefaultColumns() []string {
	return []string{"id", "name", "description"}
}

type MonitoringGatewayCollection []ecloud.MonitoringGateway

func (m MonitoringGatewayCollection) DefaultColumns() []string {
	return []string{"id", "name", "vpc_id", "router_id", "specification_id", "sync_status"}
}

type BackupGatewayCollection []ecloud.BackupGateway

func (m BackupGatewayCollection) DefaultColumns() []string {
	return []string{"id", "name", "vpc_id", "availability_zone", "gateway_spec_id", "sync_status"}
}
