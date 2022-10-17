package loadbalancer

import (
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/loadbalancer"
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

func OutputLoadBalancerTargetsProvider(targets []loadbalancer.Target) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(targets).
		WithDefaultFields([]string{"id", "name", "targetgroup_id", "ip", "port", "weight", "backup", "active"})
}

func OutputLoadBalancerBindsProvider(binds []loadbalancer.Bind) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(binds).
		WithDefaultFields([]string{"id", "listener_id", "vip_id", "port"})
}

func OutputLoadBalancerCertificatesProvider(certs []loadbalancer.Certificate) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(certs).
		WithDefaultFields([]string{"id", "listener_id", "name", "expires_at"})
}

func OutputLoadBalancerAccessIPsProvider(accessIPs []loadbalancer.AccessIP) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(accessIPs).
		WithDefaultFields([]string{"id", "ip"})
}

func OutputLoadBalancerACLsProvider(acls []loadbalancer.ACL) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(acls).
		WithDefaultFields([]string{"id", "name", "conditions", "actions"})
}

// ACLCondition represents an ACL condition
type ACLCondition struct {
	loadbalancer.ACLCondition
	Index int `json:"index"`
}

func OutputLoadBalancerACLConditionsProvider(conditions []ACLCondition) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(conditions).
		WithDefaultFields([]string{"index", "name", "inverted", "arguments"})
}

// ACLAction represents an ACL action
type ACLAction struct {
	loadbalancer.ACLAction
	Index int `json:"index"`
}

func OutputLoadBalancerACLActionsProvider(actions []ACLAction) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(actions).
		WithDefaultFields([]string{"index", "name", "arguments"})
}

func OutputLoadBalancerACLTemplatesProvider(templates []loadbalancer.ACLTemplates) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(templates).
		WithDefaultFields([]string{"id", "conditions", "actions"})
}

func OutputLoadBalancerDeploymentsProvider(deployments []loadbalancer.Deployment) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(deployments).
		WithDefaultFields([]string{"id", "cluster_id", "successful", "created_at"})
}
