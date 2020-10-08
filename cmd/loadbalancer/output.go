package loadbalancer

import (
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/loadbalancer"
)

func OutputLoadBalancerGroupsProvider(groups []loadbalancer.Group) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(groups).WithDefaultFields([]string{"id", "name"})
}

func OutputLoadBalancerConfigurationsProvider(groups []loadbalancer.Configuration) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(groups).WithDefaultFields([]string{"id", "name"})
}
