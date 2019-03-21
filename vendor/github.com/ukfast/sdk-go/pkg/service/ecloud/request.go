package ecloud

import (
	"github.com/ukfast/sdk-go/pkg/connection"
)

// PatchTagRequest represents an eCloud tag patch request
type PatchTagRequest struct {
	Value string `json:"value,omitempty"`
}

// Validate returns an error if struct properties are missing/invalid
func (c *PatchTagRequest) Validate() *connection.ValidationError {
	return nil
}

// CreateTagRequest represents a request to create an eCloud tag
type CreateTagRequest struct {
	connection.APIRequestBodyDefaultValidator

	Key   string `json:"key" validate:"required"`
	Value string `json:"value" validate:"required"`
}

// Validate returns an error if struct properties are missing/invalid
func (c *CreateTagRequest) Validate() *connection.ValidationError {
	return c.APIRequestBodyDefaultValidator.Validate(c)
}

// CreateVirtualMachineRequest represents a request to create an eCloud virtual machine
type CreateVirtualMachineRequest struct {
	connection.APIRequestBodyDefaultValidator

	Environment      string `json:"environment" validate:"required"`
	Template         string `json:"template" validate:"required"`
	TemplatePassword string `json:"template_password,omitempty"`
	// Count in Cores
	CPU int `json:"cpu" validate:"required"`
	// Size in GB
	RAM int `json:"ram" validate:"required"`
	// Size in GB
	HDD                int                               `json:"hdd,omitempty"`
	Disks              []CreateVirtualMachineRequestDisk `json:"hdd_disks,omitempty"`
	Name               string                            `json:"name,omitempty"`
	ComputerName       string                            `json:"computername,omitempty"`
	Tags               []CreateTagRequest                `json:"tags,omitempty"`
	Backup             bool                              `json:"backup"`
	Support            bool                              `json:"support"`
	Monitoring         bool                              `json:"monitoring"`
	MonitoringContacts []int                             `json:"monitoring_contacts"`
	SolutionID         int                               `json:"solution_id,omitempty"`
	DatastoreID        int                               `json:"datastore_id,omitempty"`
	SiteID             int                               `json:"site_id,omitempty"`
	NetworkID          int                               `json:"network_id,omitempty"`
	ExternalIPRequired bool                              `json:"external_ip_required"`
	SSHKeys            []string                          `json:"ssh_keys,omitempty"`
}

// CreateVirtualMachineRequestDisk represents a request to create an eCloud virtual machine disk
type CreateVirtualMachineRequestDisk struct {
	Name string `json:"name,omitempty"`
	// Size in GB
	Capacity int `json:"capacity" validate:"required"`
}

// Validate returns an error if struct properties are missing/invalid
func (c *CreateVirtualMachineRequest) Validate() *connection.ValidationError {
	if c.HDD == 0 && (c.Disks == nil || len(c.Disks) < 1) {
		return connection.NewValidationError("HDD or Disks must be provided")
	}

	return c.APIRequestBodyDefaultValidator.Validate(c)
}

// PatchSolutionRequest represents an eCloud solution patch request
type PatchSolutionRequest struct {
	Name *string `json:"name,omitempty"`
}

// Validate returns an error if struct properties are missing/invalid
func (c *PatchSolutionRequest) Validate() *connection.ValidationError {
	return nil
}

// RenameTemplateRequest represents an eCloud template rename request
type RenameTemplateRequest struct {
	connection.APIRequestBodyDefaultValidator

	NewTemplateName string `json:"new_template_name" validate:"required"`
	SolutionID      *int   `json:"solution_id,omitempty"`
}

// Validate returns an error if struct properties are missing/invalid
func (c *RenameTemplateRequest) Validate() *connection.ValidationError {
	return c.APIRequestBodyDefaultValidator.Validate(c)
}

// CloneVirtualMachineRequest represents a request to clone an eCloud virtual machine
type CloneVirtualMachineRequest struct {
	connection.APIRequestBodyDefaultValidator

	Name string `json:"name" validate:"required"`
}

// Validate returns an error if struct properties are missing/invalid
func (c *CloneVirtualMachineRequest) Validate() *connection.ValidationError {
	return c.APIRequestBodyDefaultValidator.Validate(c)
}

// PatchVirtualMachineRequest represents an eCloud virtual machine patch request
type PatchVirtualMachineRequest struct {
	Name *string `json:"name,omitempty"`
	// Count in Cores
	CPU int `json:"cpu,omitempty"`
	// Size in GB
	RAM int `json:"ram,omitempty"`
	// KV map of hard disks, key being hard disk name, value being size in GB
	Disks []PatchVirtualMachineRequestDisk `json:"hdd_disks,omitempty"`
}

// Validate returns an error if struct properties are missing/invalid
func (c *PatchVirtualMachineRequest) Validate() *connection.ValidationError {
	return nil
}

type PatchVirtualMachineRequestDiskState string

const (
	PatchVirtualMachineRequestDiskStatePresent PatchVirtualMachineRequestDiskState = "present"
	PatchVirtualMachineRequestDiskStateAbsent  PatchVirtualMachineRequestDiskState = "absent"
)

func (s PatchVirtualMachineRequestDiskState) String() string {
	return string(s)
}

// PatchVirtualMachineRequestDisk represents an eCloud virtual machine patch request disk
type PatchVirtualMachineRequestDisk struct {
	UUID string `json:"uuid" validate:"required"`
	// Size in GB
	Capacity int                                 `json:"capacity,omitempty"`
	State    PatchVirtualMachineRequestDiskState `json:"state,omitempty"`
}
