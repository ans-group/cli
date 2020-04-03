package storage

import (
	"strconv"

	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/storage"
)

func OutputStorageSolutionsProvider(solutions []storage.Solution) output.OutputHandlerProvider {
	return output.NewGenericOutputHandlerProvider(
		output.WithData(solutions),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, solution := range solutions {
				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(strconv.Itoa(solution.ID), true))
				fields.Set("name", output.NewFieldValue(solution.Name, true))
				fields.Set("san_id", output.NewFieldValue(strconv.Itoa(solution.SanID), true))
				fields.Set("created_at", output.NewFieldValue(solution.CreatedAt.String(), true))
				fields.Set("updated_at", output.NewFieldValue(solution.UpdatedAt.String(), true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputStorageVolumesProvider(volumes []storage.Volume) output.OutputHandlerProvider {
	return output.NewGenericOutputHandlerProvider(
		output.WithData(volumes),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, volume := range volumes {
				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(strconv.Itoa(volume.ID), true))
				fields.Set("name", output.NewFieldValue(volume.Name, true))
				fields.Set("wwn", output.NewFieldValue(volume.WWN, false))
				fields.Set("size_gb", output.NewFieldValue(strconv.Itoa(volume.SizeGB), true))
				fields.Set("status", output.NewFieldValue(volume.Status, true))
				fields.Set("solution_id", output.NewFieldValue(strconv.Itoa(volume.SolutionID), true))
				fields.Set("created_at", output.NewFieldValue(volume.CreatedAt.String(), true))
				fields.Set("updated_at", output.NewFieldValue(volume.UpdatedAt.String(), true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputStorageHostsProvider(hosts []storage.Host) output.OutputHandlerProvider {
	return output.NewGenericOutputHandlerProvider(
		output.WithData(hosts),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, host := range hosts {
				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(strconv.Itoa(host.ID), true))
				fields.Set("name", output.NewFieldValue(host.Name, true))
				fields.Set("os_type", output.NewFieldValue(host.OSType, false))
				fields.Set("iqn", output.NewFieldValue(host.IQN, false))
				fields.Set("server_id", output.NewFieldValue(strconv.Itoa(host.ServerID), true))
				fields.Set("status", output.NewFieldValue(host.Status, true))
				fields.Set("solution_id", output.NewFieldValue(strconv.Itoa(host.SolutionID), true))
				fields.Set("created_at", output.NewFieldValue(host.CreatedAt.String(), true))
				fields.Set("updated_at", output.NewFieldValue(host.UpdatedAt.String(), true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}
