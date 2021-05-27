package ecloud

import (
	"strconv"

	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func OutputECloudVirtualMachinesProvider(vms []ecloud.VirtualMachine) output.OutputHandlerDataProvider {
	return output.NewGenericOutputHandlerDataProvider(
		output.WithData(vms),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, vm := range vms {
				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(strconv.Itoa(vm.ID), true))
				fields.Set("name", output.NewFieldValue(vm.Name, true))
				fields.Set("hostname", output.NewFieldValue(vm.Hostname, false))
				fields.Set("computername", output.NewFieldValue(vm.ComputerName, false))
				fields.Set("cpu", output.NewFieldValue(strconv.Itoa(vm.CPU), true))
				fields.Set("ram_gb", output.NewFieldValue(strconv.Itoa(vm.RAM), true))
				fields.Set("hdd_gb", output.NewFieldValue(strconv.Itoa(vm.HDD), true))
				fields.Set("ip_internal", output.NewFieldValue(vm.IPInternal.String(), true))
				fields.Set("ip_external", output.NewFieldValue(vm.IPExternal.String(), true))
				fields.Set("platform", output.NewFieldValue(vm.Platform, false))
				fields.Set("template", output.NewFieldValue(vm.Template, false))
				fields.Set("backup", output.NewFieldValue(strconv.FormatBool(vm.Backup), false))
				fields.Set("support", output.NewFieldValue(strconv.FormatBool(vm.Support), false))
				fields.Set("environment", output.NewFieldValue(vm.Environment, false))
				fields.Set("solution_id", output.NewFieldValue(strconv.Itoa(vm.SolutionID), false))
				fields.Set("status", output.NewFieldValue(vm.Status.String(), true))
				fields.Set("power_status", output.NewFieldValue(vm.PowerStatus, true))
				fields.Set("tools_status", output.NewFieldValue(vm.ToolsStatus, false))
				fields.Set("encrypted", output.NewFieldValue(strconv.FormatBool(vm.Encrypted), false))
				fields.Set("role", output.NewFieldValue(vm.Role, false))
				fields.Set("gpu_profile", output.NewFieldValue(vm.GPUProfile, false))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputECloudVirtualMachineDisksProvider(disks []ecloud.VirtualMachineDisk) output.OutputHandlerDataProvider {
	return output.NewGenericOutputHandlerDataProvider(
		output.WithData(disks),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, disk := range disks {
				fields := output.NewOrderedFields()
				fields.Set("name", output.NewFieldValue(disk.Name, true))
				fields.Set("capacity", output.NewFieldValue(strconv.Itoa(disk.Capacity), true))
				fields.Set("uuid", output.NewFieldValue(disk.UUID, true))
				fields.Set("type", output.NewFieldValue(disk.Type.String(), true))
				fields.Set("key", output.NewFieldValue(strconv.Itoa(disk.Key), false))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputECloudTagsProvider(tags []ecloud.Tag) output.OutputHandlerDataProvider {
	return output.NewGenericOutputHandlerDataProvider(
		output.WithData(tags),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, tag := range tags {
				fields := output.NewOrderedFields()
				fields.Set("key", output.NewFieldValue(tag.Key, true))
				fields.Set("value", output.NewFieldValue(tag.Value, true))
				fields.Set("created_at", output.NewFieldValue(tag.CreatedAt.String(), true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputECloudSolutionsProvider(solutions []ecloud.Solution) output.OutputHandlerDataProvider {
	return output.NewGenericOutputHandlerDataProvider(
		output.WithData(solutions),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, solution := range solutions {
				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(strconv.Itoa(solution.ID), true))
				fields.Set("name", output.NewFieldValue(solution.Name, true))
				fields.Set("environment", output.NewFieldValue(solution.Environment.String(), true))
				fields.Set("pod_id", output.NewFieldValue(strconv.Itoa(solution.PodID), true))
				fields.Set("encryption_enabled", output.NewFieldValue(strconv.FormatBool(solution.EncryptionEnabled), false))
				fields.Set("encryption_default", output.NewFieldValue(strconv.FormatBool(solution.EncryptionDefault), false))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputECloudSitesProvider(sites []ecloud.Site) output.OutputHandlerDataProvider {
	return output.NewGenericOutputHandlerDataProvider(
		output.WithData(sites),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, site := range sites {
				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(strconv.Itoa(site.ID), true))
				fields.Set("state", output.NewFieldValue(site.State, true))
				fields.Set("solution_id", output.NewFieldValue(strconv.Itoa(site.SolutionID), true))
				fields.Set("pod_id", output.NewFieldValue(strconv.Itoa(site.PodID), true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputECloudV1HostsProvider(hosts []ecloud.V1Host) output.OutputHandlerDataProvider {
	return output.NewGenericOutputHandlerDataProvider(
		output.WithData(hosts),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, host := range hosts {
				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(strconv.Itoa(host.ID), true))
				fields.Set("solution_id", output.NewFieldValue(strconv.Itoa(host.SolutionID), true))
				fields.Set("pod_id", output.NewFieldValue(strconv.Itoa(host.PodID), true))
				fields.Set("name", output.NewFieldValue(host.Name, true))
				fields.Set("cpu_quantity", output.NewFieldValue(strconv.Itoa(host.CPU.Quantity), true))
				fields.Set("cpu_cores", output.NewFieldValue(strconv.Itoa(host.CPU.Cores), true))
				fields.Set("cpu_speed", output.NewFieldValue(host.CPU.Speed, false))
				fields.Set("ram_capacity", output.NewFieldValue(strconv.Itoa(host.RAM.Capacity), true))
				fields.Set("ram_reserved", output.NewFieldValue(strconv.Itoa(host.RAM.Reserved), false))
				fields.Set("ram_allocated", output.NewFieldValue(strconv.Itoa(host.RAM.Allocated), false))
				fields.Set("ram_available", output.NewFieldValue(strconv.Itoa(host.RAM.Available), false))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputECloudDatastoresProvider(datastores []ecloud.Datastore) output.OutputHandlerDataProvider {
	return output.NewGenericOutputHandlerDataProvider(
		output.WithData(datastores),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, datastore := range datastores {
				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(strconv.Itoa(datastore.ID), true))
				fields.Set("solution_id", output.NewFieldValue(strconv.Itoa(datastore.SolutionID), true))
				fields.Set("site_id", output.NewFieldValue(strconv.Itoa(datastore.SiteID), true))
				fields.Set("name", output.NewFieldValue(datastore.Name, true))
				fields.Set("status", output.NewFieldValue(datastore.Status.String(), true))
				fields.Set("capacity", output.NewFieldValue(strconv.Itoa(datastore.Capacity), true))
				fields.Set("allocated", output.NewFieldValue(strconv.Itoa(datastore.Allocated), false))
				fields.Set("available", output.NewFieldValue(strconv.Itoa(datastore.Available), false))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputECloudTemplatesProvider(templates []ecloud.Template) output.OutputHandlerDataProvider {
	return output.NewGenericOutputHandlerDataProvider(
		output.WithData(templates),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, template := range templates {
				fields := output.NewOrderedFields()
				fields.Set("name", output.NewFieldValue(template.Name, true))
				fields.Set("cpu", output.NewFieldValue(strconv.Itoa(template.CPU), true))
				fields.Set("ram_gb", output.NewFieldValue(strconv.Itoa(template.RAM), true))
				fields.Set("hdd_gb", output.NewFieldValue(strconv.Itoa(template.HDD), true))
				fields.Set("platform", output.NewFieldValue(template.Platform, true))
				fields.Set("operating_system", output.NewFieldValue(template.OperatingSystem, true))
				fields.Set("solution_id", output.NewFieldValue(strconv.Itoa(template.SolutionID), true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputECloudV1NetworksProvider(networks []ecloud.V1Network) output.OutputHandlerDataProvider {
	return output.NewGenericOutputHandlerDataProvider(
		output.WithData(networks),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, network := range networks {
				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(strconv.Itoa(network.ID), true))
				fields.Set("name", output.NewFieldValue(network.Name, true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputECloudFirewallsProvider(firewalls []ecloud.Firewall) output.OutputHandlerDataProvider {
	return output.NewGenericOutputHandlerDataProvider(
		output.WithData(firewalls),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, firewall := range firewalls {
				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(strconv.Itoa(firewall.ID), true))
				fields.Set("name", output.NewFieldValue(firewall.Name, true))
				fields.Set("hostname", output.NewFieldValue(firewall.Hostname, true))
				fields.Set("ip", output.NewFieldValue(firewall.IP.String(), true))
				fields.Set("role", output.NewFieldValue(firewall.Role.String(), true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputECloudPodsProvider(pods []ecloud.Pod) output.OutputHandlerDataProvider {
	return output.NewGenericOutputHandlerDataProvider(
		output.WithData(pods),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, pod := range pods {
				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(strconv.Itoa(pod.ID), true))
				fields.Set("name", output.NewFieldValue(pod.Name, true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputECloudAppliancesProvider(appliances []ecloud.Appliance) output.OutputHandlerDataProvider {
	return output.NewGenericOutputHandlerDataProvider(
		output.WithData(appliances),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, appliance := range appliances {
				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(appliance.ID, true))
				fields.Set("name", output.NewFieldValue(appliance.Name, true))
				fields.Set("logo_uri", output.NewFieldValue(appliance.LogoURI, false))
				fields.Set("description", output.NewFieldValue(appliance.Description, false))
				fields.Set("documentation_uri", output.NewFieldValue(appliance.DocumentationURI, false))
				fields.Set("publisher", output.NewFieldValue(appliance.Publisher, true))
				fields.Set("created_at", output.NewFieldValue(appliance.CreatedAt.String(), true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputECloudApplianceParametersProvider(parameters []ecloud.ApplianceParameter) output.OutputHandlerDataProvider {
	return output.NewGenericOutputHandlerDataProvider(
		output.WithData(parameters),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, parameter := range parameters {
				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(parameter.ID, true))
				fields.Set("name", output.NewFieldValue(parameter.Name, true))
				fields.Set("key", output.NewFieldValue(parameter.Key, true))
				fields.Set("type", output.NewFieldValue(parameter.Type, true))
				fields.Set("description", output.NewFieldValue(parameter.Description, true))
				fields.Set("required", output.NewFieldValue(strconv.FormatBool(parameter.Required), true))
				fields.Set("validation_rule", output.NewFieldValue(parameter.ValidationRule, false))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputECloudConsoleSessionsProvider(sessions []ecloud.ConsoleSession) output.OutputHandlerDataProvider {
	return output.NewGenericOutputHandlerDataProvider(
		output.WithData(sessions),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, session := range sessions {
				fields := output.NewOrderedFields()
				fields.Set("url", output.NewFieldValue(session.URL, true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputECloudVPCsProvider(vpcs []ecloud.VPC) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(vpcs).WithDefaultFields([]string{"id", "name", "region_id", "sync_status"})
}

func OutputECloudInstancesProvider(instances []ecloud.Instance) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(instances).WithDefaultFields([]string{"id", "name", "vpc_id", "vcpu_cores", "ram_capacity", "sync_status"})
}

func OutputECloudFloatingIPsProvider(fips []ecloud.FloatingIP) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(fips).WithDefaultFields([]string{"id", "name", "ip_address", "sync_status"})
}

func OutputECloudFirewallPoliciesProvider(policies []ecloud.FirewallPolicy) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(policies).WithDefaultFields([]string{"id", "name", "router_id", "sync_status"})
}

func OutputECloudFirewallRulesProvider(rules []ecloud.FirewallRule) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(rules).WithDefaultFields([]string{"id", "name", "firewall_policy_id", "source", "destination", "action", "direction", "enabled"})
}

func OutputECloudFirewallRulePortsProvider(rules []ecloud.FirewallRulePort) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(rules).WithDefaultFields([]string{"id", "name", "firewall_rule_id", "protocol", "source", "destination"})
}

func OutputECloudRegionsProvider(regions []ecloud.Region) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(regions).WithDefaultFields([]string{"id", "name"})
}

func OutputECloudVolumesProvider(volumes []ecloud.Volume) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(volumes).WithDefaultFields([]string{"id", "name", "type", "capacity", "sync_status"})
}

func OutputECloudCredentialsProvider(credentials []ecloud.Credential) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(credentials).WithDefaultFields([]string{"id", "name", "username", "password"})
}

func OutputECloudNICsProvider(nics []ecloud.NIC) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(nics).WithDefaultFields([]string{"id", "mac_address", "instance_id", "network_id", "ip_address"})
}

func OutputECloudRoutersProvider(routers []ecloud.Router) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(routers).WithDefaultFields([]string{"id", "name", "vpc_id", "sync_status"})
}

func OutputECloudNetworksProvider(networks []ecloud.Network) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(networks).WithDefaultFields([]string{"id", "name", "router_id", "subnet", "sync_status"})
}

func OutputECloudDHCPsProvider(dhcps []ecloud.DHCP) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(dhcps).WithDefaultFields([]string{"id", "vpc_id", "availability_zone_id", "sync_status"})
}

func OutputECloudRouterThroughputsProvider(throughputs []ecloud.RouterThroughput) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(throughputs).WithDefaultFields([]string{"id", "availability_zone_id", "name"})
}

func OutputECloudImagesProvider(images []ecloud.Image) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(images).WithDefaultFields([]string{"id", "name"})
}

func OutputECloudImageParametersProvider(parameters []ecloud.ImageParameter) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(parameters).WithDefaultFields([]string{"key", "type", "required"})
}

func OutputECloudImageMetadataProvider(metadata []ecloud.ImageMetadata) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(metadata).WithDefaultFields([]string{"key", "value"})
}

func OutputECloudSSHKeyPairsProvider(keypairs []ecloud.SSHKeyPair) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(keypairs).WithDefaultFields([]string{"id", "name"})
}

func OutputECloudTasksProvider(tasks []ecloud.Task) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(tasks).WithDefaultFields([]string{"id", "name", "status", "created_at", "updated_at"})
}
