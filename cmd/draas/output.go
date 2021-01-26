package draas

import (
	"fmt"
	"strconv"

	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/draas"
)

func OutputDRaaSSolutionsProvider(solutions []draas.Solution) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(solutions).
		WithDefaultFields([]string{"id", "name"})
}

func OutputDRaaSBackupResourcesProvider(resources []draas.BackupResource) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(resources).
		WithDefaultFields([]string{"id", "name", "quota", "used_quota"})
}

func OutputDRaaSIOPSTiersProvider(tiers []draas.IOPSTier) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(tiers).
		WithDefaultFields([]string{"id", "iops_limit"})
}

func OutputDRaaSBackupServicesProvider(services []draas.BackupService) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(services).
		WithDefaultFields([]string{"service", "account_name", "gateway", "port"})
}

func OutputDRaaSFailoverPlansProvider(plans []draas.FailoverPlan) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(plans).
		WithDefaultFields([]string{"id", "name", "status"})
}

func OutputDRaaSComputeResourcesProvider(resources []draas.ComputeResource) output.OutputHandlerDataProvider {
	return output.NewGenericOutputHandlerDataProvider(
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
					fields.Set(fmt.Sprintf("storage_#%d_name", i), output.NewFieldValue(storage.Name, true))
					fields.Set(fmt.Sprintf("storage_#%d_used", i), output.NewFieldValue(strconv.Itoa(storage.Used), true))
					fields.Set(fmt.Sprintf("storage_#%d_limit", i), output.NewFieldValue(strconv.Itoa(storage.Limit), true))
				}

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputDRaaSHardwarePlansProvider(plans []draas.HardwarePlan) output.OutputHandlerDataProvider {
	return output.NewGenericOutputHandlerDataProvider(
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
				for i, storage := range plan.Storage {
					fields.Set(fmt.Sprintf("storage_#%d_id", i), output.NewFieldValue(storage.ID, true))
					fields.Set(fmt.Sprintf("storage_#%d_name", i), output.NewFieldValue(storage.Name, true))
					fields.Set(fmt.Sprintf("storage_#%d_type", i), output.NewFieldValue(storage.Type, false))
					fields.Set(fmt.Sprintf("storage_#%d_quota", i), output.NewFieldValue(strconv.Itoa(storage.Quota), false))
				}

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputDRaaSReplicasProvider(replicas []draas.Replica) output.OutputHandlerDataProvider {
	return output.NewGenericOutputHandlerDataProvider(
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
				fields.Set("disk", output.NewFieldValue(strconv.Itoa(replica.Disk), false))
				fields.Set("iops", output.NewFieldValue(strconv.Itoa(replica.IOPS), false))
				fields.Set("power", output.NewFieldValue(strconv.FormatBool(replica.Power), true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputDRaaSBillingTypesProvider(types []draas.BillingType) output.OutputHandlerDataProvider {
	return output.NewGenericOutputHandlerDataProvider(
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
