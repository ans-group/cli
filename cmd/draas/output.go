package draas

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/draas"
)

type SolutionCollection []draas.Solution

func (s SolutionCollection) DefaultColumns() []string {
	return []string{"id", "name"}
}

func (s SolutionCollection) Fields() []*output.OrderedFields {
	var data []*output.OrderedFields
	for _, solution := range s {
		fields := output.NewOrderedFields()
		fields.Set("id", solution.ID)
		fields.Set("name", solution.Name)
		fields.Set("iops_tier_id", solution.IOPSTierID)
		fields.Set("billing_type_id", solution.BillingTypeID)

		data = append(data, fields)
	}

	return data
}

type BackupResourceCollection []draas.BackupResource

func (b BackupResourceCollection) DefaultColumns() []string {
	return []string{"id", "name", "quota", "used_quota"}
}

func (b BackupResourceCollection) Fields() []*output.OrderedFields {
	var data []*output.OrderedFields
	for _, resource := range b {
		fields := output.NewOrderedFields()
		fields.Set("id", resource.ID)
		fields.Set("name", resource.Name)
		fields.Set("quota", strconv.Itoa(resource.Quota))
		fields.Set("used_quota", fmt.Sprintf("%f", resource.UsedQuota))

		data = append(data, fields)
	}

	return data
}

type IOPSTierCollection []draas.IOPSTier

func (t IOPSTierCollection) DefaultColumns() []string {
	return []string{"id", "iops_limit"}
}

func (t IOPSTierCollection) Fields() []*output.OrderedFields {
	var data []*output.OrderedFields
	for _, tier := range t {
		fields := output.NewOrderedFields()
		fields.Set("id", tier.ID)
		fields.Set("iops_limit", strconv.Itoa(tier.IOPSLimit))

		data = append(data, fields)
	}

	return data
}

type BackupServiceCollection []draas.BackupService

func (b BackupServiceCollection) DefaultColumns() []string {
	return []string{"service", "account_name", "gateway", "port"}
}

func (b BackupServiceCollection) Fields() []*output.OrderedFields {
	var data []*output.OrderedFields
	for _, service := range b {
		fields := output.NewOrderedFields()
		fields.Set("service", service.Service)
		fields.Set("account_name", service.AccountName)
		fields.Set("gateway", service.Gateway)
		fields.Set("port", strconv.Itoa(service.Port))

		data = append(data, fields)
	}

	return data
}

type FailoverPlanCollection []draas.FailoverPlan

func (f FailoverPlanCollection) DefaultColumns() []string {
	return []string{"id", "name", "status"}
}

func (f FailoverPlanCollection) Fields() []*output.OrderedFields {
	var data []*output.OrderedFields
	for _, plan := range f {
		vms := []string{}
		if len(plan.VMs) > 0 {
			for _, vm := range plan.VMs {
				vms = append(vms, vm.Name)
			}
		}

		fields := output.NewOrderedFields()
		fields.Set("id", plan.ID)
		fields.Set("name", plan.Name)
		fields.Set("description", plan.Description)
		fields.Set("status", plan.Status)
		fields.Set("vms", strings.Join(vms, ", "))

		data = append(data, fields)
	}

	return data
}

type ComputeResourceCollection []draas.ComputeResource

func (c ComputeResourceCollection) DefaultColumns() []string {
	return []string{"id", "hardware_plan_id", "memory_used", "memory_limit", "cpu_used"}
}

func (c ComputeResourceCollection) Fields() []*output.OrderedFields {
	var data []*output.OrderedFields
	for _, resource := range c {
		fields := output.NewOrderedFields()
		fields.Set("id", resource.ID)
		fields.Set("hardware_plan_id", resource.HardwarePlanID)
		fields.Set("memory_used", fmt.Sprintf("%f", resource.Memory.Used))
		fields.Set("memory_limit", fmt.Sprintf("%f", resource.Memory.Limit))
		fields.Set("cpu_used", strconv.Itoa(resource.CPU.Used))

		for i, storage := range resource.Storage {
			fields.Set(fmt.Sprintf("storage_#%d_name", i), storage.Name)
			fields.Set(fmt.Sprintf("storage_#%d_used", i), strconv.Itoa(storage.Used))
			fields.Set(fmt.Sprintf("storage_#%d_limit", i), strconv.Itoa(storage.Limit))
		}

		data = append(data, fields)
	}

	return data
}

type HardwarePlanCollection []draas.HardwarePlan

func (h HardwarePlanCollection) DefaultColumns() []string {
	return []string{"id", "name", "description"}
}

func (h HardwarePlanCollection) Fields() []*output.OrderedFields {
	var data []*output.OrderedFields
	for _, plan := range h {
		fields := output.NewOrderedFields()
		fields.Set("id", plan.ID)
		fields.Set("name", plan.Name)
		fields.Set("description", plan.Description)
		fields.Set("limits_processor", strconv.Itoa(plan.Limits.Processor))
		fields.Set("limits_memory", strconv.Itoa(plan.Limits.Memory))
		fields.Set("networks_public", strconv.Itoa(plan.Networks.Public))
		fields.Set("networks_private", strconv.Itoa(plan.Networks.Private))
		for i, storage := range plan.Storage {
			fields.Set(fmt.Sprintf("storage_#%d_id", i), storage.ID)
			fields.Set(fmt.Sprintf("storage_#%d_name", i), storage.Name)
			fields.Set(fmt.Sprintf("storage_#%d_type", i), storage.Type)
			fields.Set(fmt.Sprintf("storage_#%d_quota", i), strconv.Itoa(storage.Quota))
		}

		data = append(data, fields)
	}

	return data
}

type ReplicaCollection []draas.Replica

func (r ReplicaCollection) DefaultColumns() []string {
	return []string{"id", "name", "platform", "power"}
}

func (r ReplicaCollection) Fields() []*output.OrderedFields {
	var data []*output.OrderedFields
	for _, replica := range r {
		fields := output.NewOrderedFields()
		fields.Set("id", replica.ID)
		fields.Set("name", replica.Name)
		fields.Set("platform", replica.Platform)
		fields.Set("cpu", strconv.Itoa(replica.CPU))
		fields.Set("ram", strconv.Itoa(replica.RAM))
		fields.Set("disk", strconv.Itoa(replica.Disk))
		fields.Set("iops", strconv.Itoa(replica.IOPS))
		fields.Set("power", strconv.FormatBool(replica.Power))

		data = append(data, fields)
	}

	return data
}

type BillingTypeCollection []draas.BillingType

func (b BillingTypeCollection) DefaultColumns() []string {
	return []string{"id", "type"}
}

func (b BillingTypeCollection) Fields() []*output.OrderedFields {
	var data []*output.OrderedFields
	for _, t := range b {
		fields := output.NewOrderedFields()
		fields.Set("id", t.ID)
		fields.Set("type", t.Type)

		data = append(data, fields)
	}

	return data
}
