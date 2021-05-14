package loadbalancer

import (
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/loadbalancer"
)

func OutputLoadBalancerClustersProvider(clusters []loadbalancer.Cluster) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(clusters).
		WithDefaultFields([]string{"id", "name", "deployed", "deployed_at"})
}

func OutputLoadBalancerListenersProvider(listeners []loadbalancer.Listener) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(listeners).
		WithDefaultFields([]string{"id", "name", "deployed", "deployed_at"})
}
