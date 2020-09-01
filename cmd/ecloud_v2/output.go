package ecloud_v2

import (
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func OutputECloudV2VPCsProvider(vpcs []ecloud.VPC) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(vpcs).WithDefaultFields([]string{"id", "name"})
}
