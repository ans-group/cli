package loadbalancer

import (
	"github.com/ans-group/sdk-go/pkg/service/loadbalancer"
)

type ClusterCollection []loadbalancer.Cluster

func (m ClusterCollection) DefaultColumns() []string {
	return []string{"id", "name", "deployed", "deployed_at"}
}

type ListenerCollection []loadbalancer.Listener

func (m ListenerCollection) DefaultColumns() []string {
	return []string{"id", "name", "cluster_id"}
}

type TargetGroupCollection []loadbalancer.TargetGroup

func (m TargetGroupCollection) DefaultColumns() []string {
	return []string{"id", "name", "cluster_id", "mode"}
}

type TargetCollection []loadbalancer.Target

func (m TargetCollection) DefaultColumns() []string {
	return []string{"id", "name", "targetgroup_id", "ip", "port", "weight", "backup", "active"}
}

type BindCollection []loadbalancer.Bind

func (m BindCollection) DefaultColumns() []string {
	return []string{"id", "listener_id", "vip_id", "port"}
}

type CertificateCollection []loadbalancer.Certificate

func (m CertificateCollection) DefaultColumns() []string {
	return []string{"id", "listener_id", "name", "expires_at"}
}

type AccessIPCollection []loadbalancer.AccessIP

func (m AccessIPCollection) DefaultColumns() []string {
	return []string{"id", "ip"}
}

type ACLCollection []loadbalancer.ACL

func (m ACLCollection) DefaultColumns() []string {
	return []string{"id", "name", "conditions", "actions"}
}

// ACLCondition represents an ACL condition
type ACLCondition struct {
	loadbalancer.ACLCondition
	Index int `json:"index"`
}

type ACLConditionCollection []ACLCondition

func (m ACLConditionCollection) DefaultColumns() []string {
	return []string{"index", "name", "inverted", "arguments"}
}

// ACLAction represents an ACL action
type ACLAction struct {
	loadbalancer.ACLAction
	Index int `json:"index"`
}

type ACLActionCollection []ACLAction

func (m ACLActionCollection) DefaultColumns() []string {
	return []string{"index", "name", "arguments"}
}

type ACLTemplatesCollection []loadbalancer.ACLTemplates

func (m ACLTemplatesCollection) DefaultColumns() []string {
	return []string{"id", "conditions", "actions"}
}

type DeploymentCollection []loadbalancer.Deployment

func (m DeploymentCollection) DefaultColumns() []string {
	return []string{"id", "cluster_id", "successful", "created_at"}
}

type VIPCollection []loadbalancer.VIP

func (m VIPCollection) DefaultColumns() []string {
	return []string{"id", "cluster_id", "internal_cidr", "external_cidr"}
}
