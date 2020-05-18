package draas

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/draas"
)

func OutputDRaaSSolutionsProvider(solutions []draas.Solution) output.OutputHandlerProvider {
	return output.NewGenericOutputHandlerProvider(
		output.WithData(solutions),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, solution := range solutions {
				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(solution.ID, true))
				fields.Set("name", output.NewFieldValue(solution.Name, true))
				fields.Set("iops_tier_id", output.NewFieldValue(solution.IOPSTierID, false))
				fields.Set("billing_type_id", output.NewFieldValue(solution.BillingTypeID, false))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputDRaaSBackupResourcesProvider(resources []draas.BackupResource) output.OutputHandlerProvider {
	return output.NewGenericOutputHandlerProvider(
		output.WithData(resources),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, resource := range resources {
				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(resource.ID, true))
				fields.Set("name", output.NewFieldValue(resource.Name, true))
				fields.Set("quota", output.NewFieldValue(strconv.Itoa(resource.Quota), true))
				fields.Set("used_quota", output.NewFieldValue(fmt.Sprintf("%f", resource.UsedQuota), true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputDRaaSIOPSTiersProvider(tiers []draas.IOPSTier) output.OutputHandlerProvider {
	return output.NewGenericOutputHandlerProvider(
		output.WithData(tiers),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, tier := range tiers {
				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(tier.ID, true))
				fields.Set("iops_limit", output.NewFieldValue(strconv.Itoa(tier.IOPSLimit), true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputDRaaSBackupServicesProvider(services []draas.BackupService) output.OutputHandlerProvider {
	return output.NewGenericOutputHandlerProvider(
		output.WithData(services),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, service := range services {
				fields := output.NewOrderedFields()
				fields.Set("service", output.NewFieldValue(service.Service, true))
				fields.Set("account_name", output.NewFieldValue(service.AccountName, true))
				fields.Set("gateway", output.NewFieldValue(service.Gateway, true))
				fields.Set("port", output.NewFieldValue(strconv.Itoa(service.Port), true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputDRaaSFailoverPlansProvider(plans []draas.FailoverPlan) output.OutputHandlerProvider {
	return output.NewGenericOutputHandlerProvider(
		output.WithData(plans),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, plan := range plans {
				vms := []string{}
				if len(plan.VMs) > 0 {
					for _, vm := range plan.VMs {
						vms = append(vms, vm.Name)
					}
				}

				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(plan.ID, true))
				fields.Set("name", output.NewFieldValue(plan.Name, true))
				fields.Set("description", output.NewFieldValue(plan.Description, false))
				fields.Set("status", output.NewFieldValue(plan.Status, true))
				fields.Set("vms", output.NewFieldValue(strings.Join(vms, ", "), false))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputDRaaSComputeResourcesProvider(resources []draas.ComputeResource) output.OutputHandlerProvider {
	return output.NewGenericOutputHandlerProvider(
		output.WithData(resources),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, resource := range resources {
				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(resource.ID, true))
				fields.Set("hardware_plan_id", output.NewFieldValue(resource.HardwarePlanID, true))
				fields.Set("memory_used", output.NewFieldValue(fmt.Sprintf("%f", resource.Memory.Used), true))
				fields.Set("memory_limit", output.NewFieldValue(fmt.Sprintf("%f", resource.Memory.Limit), true))
				fields.Set("cpu_used", output.NewFieldValue(strconv.Itoa(resource.CPU.Used), true))
				for i, storage := range resource.Storage {
					fields.Set(fmt.Sprintf("storage_%d_name", i), output.NewFieldValue(storage.Name, true))
					fields.Set(fmt.Sprintf("storage_%d_used", i), output.NewFieldValue(strconv.Itoa(storage.Used), true))
					fields.Set(fmt.Sprintf("storage_%d_limit", i), output.NewFieldValue(strconv.Itoa(storage.Limit), true))
				}

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputDRaaSHardwarePlansProvider(plans []draas.HardwarePlan) output.OutputHandlerProvider {
	return output.NewGenericOutputHandlerProvider(
		output.WithData(plans),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, plan := range plans {
				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(plan.ID, true))
				fields.Set("name", output.NewFieldValue(plan.Name, true))
				fields.Set("description", output.NewFieldValue(plan.Description, true))
				fields.Set("limits_processor", output.NewFieldValue(strconv.Itoa(plan.Limits.Processor), false))
				fields.Set("limits_memory", output.NewFieldValue(strconv.Itoa(plan.Limits.Memory), false))
				fields.Set("networks_public", output.NewFieldValue(strconv.Itoa(plan.Networks.Public), false))
				fields.Set("networks_private", output.NewFieldValue(strconv.Itoa(plan.Networks.Private), false))
				fields.Set("storage_id", output.NewFieldValue(plan.Storage.ID, false))
				fields.Set("storage_name", output.NewFieldValue(plan.Storage.Name, false))
				fields.Set("storage_type", output.NewFieldValue(plan.Storage.Type, false))
				fields.Set("storage_quota", output.NewFieldValue(strconv.Itoa(plan.Storage.Quota), false))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputDRaaSReplicasProvider(replicas []draas.Replica) output.OutputHandlerProvider {
	return output.NewGenericOutputHandlerProvider(
		output.WithData(replicas),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, replica := range replicas {
				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(replica.ID, true))
				fields.Set("name", output.NewFieldValue(replica.Name, true))
				fields.Set("platform", output.NewFieldValue(replica.Platform, true))
				fields.Set("cpu", output.NewFieldValue(strconv.Itoa(replica.CPU), false))
				fields.Set("ram", output.NewFieldValue(strconv.Itoa(replica.RAM), false))
				fields.Set("hdd", output.NewFieldValue(strconv.Itoa(replica.HDD), false))
				fields.Set("iops", output.NewFieldValue(strconv.Itoa(replica.IOPS), false))
				fields.Set("power", output.NewFieldValue(strconv.FormatBool(replica.Power), true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputDRaaSBillingTypesProvider(types []draas.BillingType) output.OutputHandlerProvider {
	return output.NewGenericOutputHandlerProvider(
		output.WithData(types),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, t := range types {
				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(t.ID, true))
				fields.Set("type", output.NewFieldValue(t.Type, true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}
