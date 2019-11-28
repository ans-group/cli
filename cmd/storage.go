package cmd

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/storage"
)

func storageRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "storage",
		Short: "Commands relating to Storage service",
	}

	// Child root commands
	cmd.AddCommand(storageSolutionRootCmd())
	cmd.AddCommand(storageHostRootCmd())
	cmd.AddCommand(storageVolumeRootCmd())

	return cmd
}

func outputStorageSolutions(solutions []storage.Solution) {
	err := Output(&OutputStorageSolutions{Solutions: solutions})
	if err != nil {
		output.Fatalf("Failed to output solutions: %s", err)
	}
}

// OutputStorageSolutions implements OutputDataProvider for outputting an array of Solutions
type OutputStorageSolutions struct {
	Solutions []storage.Solution
}

func (o *OutputStorageSolutions) GetData() interface{} {
	return o.Solutions
}

func (o *OutputStorageSolutions) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, solution := range o.Solutions {
		fields := o.getOrderedFields(solution)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputStorageSolutions) getOrderedFields(solution storage.Solution) *output.OrderedFields {
	fields := output.NewOrderedFields()
	fields.Set("id", output.NewFieldValue(strconv.Itoa(solution.ID), true))
	fields.Set("name", output.NewFieldValue(solution.Name, true))
	fields.Set("san_id", output.NewFieldValue(strconv.Itoa(solution.SanID), true))
	fields.Set("created_at", output.NewFieldValue(solution.CreatedAt.String(), true))
	fields.Set("updated_at", output.NewFieldValue(solution.UpdatedAt.String(), true))

	return fields
}

func outputStorageVolumes(volumes []storage.Volume) {
	err := Output(&OutputStorageVolumes{Volumes: volumes})
	if err != nil {
		output.Fatalf("Failed to output volumes: %s", err)
	}
}

// OutputStorageVolumes implements OutputDataProvider for outputting an array of Volumes
type OutputStorageVolumes struct {
	Volumes []storage.Volume
}

func (o *OutputStorageVolumes) GetData() interface{} {
	return o.Volumes
}

func (o *OutputStorageVolumes) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, volume := range o.Volumes {
		fields := o.getOrderedFields(volume)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputStorageVolumes) getOrderedFields(volume storage.Volume) *output.OrderedFields {
	fields := output.NewOrderedFields()
	fields.Set("id", output.NewFieldValue(strconv.Itoa(volume.ID), true))
	fields.Set("name", output.NewFieldValue(volume.Name, true))
	fields.Set("wwn", output.NewFieldValue(volume.WWN, false))
	fields.Set("size_gb", output.NewFieldValue(strconv.Itoa(volume.SizeGB), true))
	fields.Set("status", output.NewFieldValue(volume.Status, true))
	fields.Set("solution_id", output.NewFieldValue(strconv.Itoa(volume.SolutionID), true))
	fields.Set("created_at", output.NewFieldValue(volume.CreatedAt.String(), true))
	fields.Set("updated_at", output.NewFieldValue(volume.UpdatedAt.String(), true))

	return fields
}

func outputStorageHosts(hosts []storage.Host) {
	err := Output(&OutputStorageHosts{Hosts: hosts})
	if err != nil {
		output.Fatalf("Failed to output hosts: %s", err)
	}
}

// OutputStorageHosts implements OutputDataProvider for outputting an array of Hosts
type OutputStorageHosts struct {
	Hosts []storage.Host
}

func (o *OutputStorageHosts) GetData() interface{} {
	return o.Hosts
}

func (o *OutputStorageHosts) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, host := range o.Hosts {
		fields := o.getOrderedFields(host)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputStorageHosts) getOrderedFields(host storage.Host) *output.OrderedFields {
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

	return fields
}
