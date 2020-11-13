package ecloudv2

import (
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func OutputECloudVPCsProvider(vpcs []ecloud.VPC) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(vpcs).WithDefaultFields([]string{"id", "name", "region_id", "created_at", "updated_at"})
}

func OutputECloudInstancesProvider(instances []ecloud.Instance) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(instances).WithDefaultFields([]string{"id", "name", "vpc_id", "status", "vcpu_cores", "ram_capacity"})
}

func OutputECloudFloatingIPsProvider(fips []ecloud.FloatingIP) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(fips).WithDefaultFields([]string{"id", "created_at", "updated_at"})
}

func OutputECloudFirewallRulesProvider(rules []ecloud.FirewallRule) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(rules).WithDefaultFields([]string{"id", "router_id", "created_at", "updated_at"})
}

func OutputECloudRegionsProvider(regions []ecloud.Region) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(regions).WithDefaultFields([]string{"id", "name"})
}

func OutputECloudLoadBalancerClustersProvider(lbcs []ecloud.LoadBalancerCluster) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(lbcs).WithDefaultFields([]string{"id", "name"})
}

func OutputECloudVolumesProvider(volumes []ecloud.Volume) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(volumes).WithDefaultFields([]string{"id", "name", "capacity"})
}

func OutputECloudCredentialsProvider(credentials []ecloud.Credential) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(credentials).WithDefaultFields([]string{"id", "name", "username", "password"})
}

func OutputECloudNICsProvider(nics []ecloud.NIC) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(nics).WithDefaultFields([]string{"id", "mac_address", "instance", "ip_address"})
}
