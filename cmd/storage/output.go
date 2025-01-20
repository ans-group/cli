package storage

import (
	"strconv"

	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/storage"
)

type SolutionCollection []storage.Solution

func (s SolutionCollection) DefaultColumns() []string {
	return []string{"id", "name", "san_id", "created_at", "updated_at"}
}

func (s SolutionCollection) Fields() []*output.OrderedFields {
	var data []*output.OrderedFields
	for _, solution := range s {
		fields := output.NewOrderedFields()
		fields.Set("id", strconv.Itoa(solution.ID))
		fields.Set("name", solution.Name)
		fields.Set("san_id", strconv.Itoa(solution.SanID))
		fields.Set("created_at", solution.CreatedAt.String())
		fields.Set("updated_at", solution.UpdatedAt.String())

		data = append(data, fields)
	}

	return data
}

type VolumeCollection []storage.Volume

func (v VolumeCollection) DefaultColumns() []string {
	return []string{"id", "name", "size_gb", "status", "solution_id", "created_at", "updated_at"}
}

func (v VolumeCollection) Fields() []*output.OrderedFields {
	var data []*output.OrderedFields
	for _, volume := range v {
		fields := output.NewOrderedFields()
		fields.Set("id", strconv.Itoa(volume.ID))
		fields.Set("name", volume.Name)
		fields.Set("wwn", volume.WWN)
		fields.Set("size_gb", strconv.Itoa(volume.SizeGB))
		fields.Set("status", volume.Status)
		fields.Set("solution_id", strconv.Itoa(volume.SolutionID))
		fields.Set("created_at", volume.CreatedAt.String())
		fields.Set("updated_at", volume.UpdatedAt.String())

		data = append(data, fields)
	}

	return data
}

type HostCollection []storage.Host

func (h HostCollection) DefaultColumns() []string {
	return []string{"id", "name", "server_id", "status", "solution_id", "created_at", "updated_at"}
}

func (h HostCollection) Fields() []*output.OrderedFields {
	var data []*output.OrderedFields
	for _, host := range h {
		fields := output.NewOrderedFields()
		fields.Set("id", strconv.Itoa(host.ID))
		fields.Set("name", host.Name)
		fields.Set("os_type", host.OSType)
		fields.Set("iqn", host.IQN)
		fields.Set("server_id", strconv.Itoa(host.ServerID))
		fields.Set("status", host.Status)
		fields.Set("solution_id", strconv.Itoa(host.SolutionID))
		fields.Set("created_at", host.CreatedAt.String())
		fields.Set("updated_at", host.UpdatedAt.String())

		data = append(data, fields)
	}

	return data
}
