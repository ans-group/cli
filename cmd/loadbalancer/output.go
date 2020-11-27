package loadbalancer

import (
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/loadbalancer"
)

func OutputLoadBalancerClustersProvider(clusters []loadbalancer.Cluster) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(clusters).WithDefaultFields([]string{"id", "name"})
}

func OutputLoadBalancerTargetsProvider(targets []loadbalancer.Target) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(targets).WithDefaultFields([]string{"id", "ip", "port", "backup", "weight"})
}
