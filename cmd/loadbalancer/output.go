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
		WithDefaultFields([]string{"id", "name", "cluster_id"})
}

func OutputLoadBalancerTargetGroupsProvider(groups []loadbalancer.TargetGroup) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(groups).
		WithDefaultFields([]string{"id", "name", "cluster_id", "mode"})
}

func OutputLoadBalancerBindsProvider(binds []loadbalancer.Bind) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(binds).
		WithDefaultFields([]string{"id", "listener_id", "vip_id", "port"})
}

func OutputLoadBalancerCertificatesProvider(certs []loadbalancer.Certificate) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(certs).
		WithDefaultFields([]string{"id", "listener_id", "name"})
}

func OutputLoadBalancerAccessIPsProvider(accessIPs []loadbalancer.AccessIP) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(accessIPs).
		WithDefaultFields([]string{"id", "ip"})
}
