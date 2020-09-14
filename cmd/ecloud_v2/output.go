package ecloud_v2

import (
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func OutputECloudVPCsProvider(vpcs []ecloud.VPC) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(vpcs).WithDefaultFields([]string{"id", "name"})
}

func OutputECloudAvailabilityZonesProvider(azs []ecloud.AvailabilityZone) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(azs).WithDefaultFields([]string{"id", "name", "code"})
}

func OutputECloudNetworksProvider(networks []ecloud.Network) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(networks).WithDefaultFields([]string{"id", "name", "router_id"})
}
