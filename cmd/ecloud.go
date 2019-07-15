package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ecloud",
		Short: "Commands relating to eCloud service",
	}

	// Child root commands
	cmd.AddCommand(ecloudVirtualMachineRootCmd())
	cmd.AddCommand(ecloudSolutionRootCmd())
	cmd.AddCommand(ecloudSiteRootCmd())
	cmd.AddCommand(ecloudHostRootCmd())
	cmd.AddCommand(ecloudFirewallRootCmd())
	cmd.AddCommand(ecloudPodRootCmd())
	cmd.AddCommand(ecloudDatastoreRootCmd())
	cmd.AddCommand(ecloudApplianceRootCmd())
	cmd.AddCommand(ecloudCreditRootCmd())

	return cmd
}

// OutputECloudVirtualMachines implements OutputDataProvider for outputting an array of virtual machines
type OutputECloudVirtualMachines struct {
	VirtualMachines []ecloud.VirtualMachine
}

func outputECloudVirtualMachines(vms []ecloud.VirtualMachine) {
	err := Output(&OutputECloudVirtualMachines{VirtualMachines: vms})
	if err != nil {
		output.Fatalf("Failed to output virtual machines: %s", err)
	}
}

func (o *OutputECloudVirtualMachines) GetData() interface{} {
	return o.VirtualMachines
}

func (o *OutputECloudVirtualMachines) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, vm := range o.VirtualMachines {
		fields := o.getOrderedFields(vm)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputECloudVirtualMachines) getOrderedFields(vm ecloud.VirtualMachine) *output.OrderedFields {
	fields := output.NewOrderedFields()
	fields.Set("id", output.NewFieldValue(strconv.Itoa(vm.ID), true))
	fields.Set("name", output.NewFieldValue(vm.Name, true))
	fields.Set("hostname", output.NewFieldValue(vm.Hostname, false))
	fields.Set("computername", output.NewFieldValue(vm.ComputerName, false))
	fields.Set("cpu", output.NewFieldValue(strconv.Itoa(vm.CPU), true))
	fields.Set("ram_gb", output.NewFieldValue(strconv.Itoa(vm.RAM), true))
	fields.Set("hdd_gb", output.NewFieldValue(strconv.Itoa(vm.HDD), true))
	fields.Set("ip_internal", output.NewFieldValue(vm.IPInternal.String(), true))
	fields.Set("ip_external", output.NewFieldValue(vm.IPExternal.String(), true))
	fields.Set("platform", output.NewFieldValue(vm.Platform, false))
	fields.Set("template", output.NewFieldValue(vm.Template, false))
	fields.Set("backup", output.NewFieldValue(strconv.FormatBool(vm.Backup), false))
	fields.Set("support", output.NewFieldValue(strconv.FormatBool(vm.Support), false))
	fields.Set("environment", output.NewFieldValue(vm.Environment, false))
	fields.Set("solution_id", output.NewFieldValue(strconv.Itoa(vm.SolutionID), false))
	fields.Set("status", output.NewFieldValue(vm.Status.String(), true))
	fields.Set("power_status", output.NewFieldValue(vm.PowerStatus, true))
	fields.Set("tools_status", output.NewFieldValue(vm.ToolsStatus, false))
	fields.Set("encrypted", output.NewFieldValue(strconv.FormatBool(vm.Encrypted), false))
	fields.Set("role", output.NewFieldValue(vm.Role, false))

	return fields
}

// OutputECloudVirtualMachineDisks implements OutputDataProvider for outputting an array of virtual machine disks
type OutputECloudVirtualMachineDisks struct {
	VirtualMachineDisks []ecloud.VirtualMachineDisk
}

func outputECloudVirtualMachineDisks(disks []ecloud.VirtualMachineDisk) {
	err := Output(&OutputECloudVirtualMachineDisks{VirtualMachineDisks: disks})
	if err != nil {
		output.Fatalf("Failed to output virtual machine disks: %s", err)
	}
}

func (o *OutputECloudVirtualMachineDisks) GetData() interface{} {
	return o.VirtualMachineDisks
}

func (o *OutputECloudVirtualMachineDisks) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, disk := range o.VirtualMachineDisks {
		fields := o.getOrderedFields(disk)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputECloudVirtualMachineDisks) getOrderedFields(disk ecloud.VirtualMachineDisk) *output.OrderedFields {
	fields := output.NewOrderedFields()
	fields.Set("name", output.NewFieldValue(disk.Name, true))
	fields.Set("capacity", output.NewFieldValue(strconv.Itoa(disk.Capacity), true))
	fields.Set("uuid", output.NewFieldValue(disk.UUID, true))
	fields.Set("type", output.NewFieldValue(disk.Type.String(), true))
	fields.Set("key", output.NewFieldValue(strconv.Itoa(disk.Key), false))

	return fields
}

// OutputECloudTags implements OutputDataProvider for outputting an array of tags
type OutputECloudTags struct {
	Tags []ecloud.Tag
}

func outputECloudTags(tags []ecloud.Tag) {
	err := Output(&OutputECloudTags{Tags: tags})
	if err != nil {
		output.Fatalf("Failed to output tags: %s", err)
	}
}

func (o *OutputECloudTags) GetData() interface{} {
	return o.Tags
}

func (o *OutputECloudTags) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, tag := range o.Tags {
		fields := o.getOrderedFields(tag)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputECloudTags) getOrderedFields(tag ecloud.Tag) *output.OrderedFields {
	fields := output.NewOrderedFields()
	fields.Set("key", output.NewFieldValue(tag.Key, true))
	fields.Set("value", output.NewFieldValue(tag.Value, true))
	fields.Set("created_at", output.NewFieldValue(tag.CreatedAt.String(), true))

	return fields
}

// OutputECloudSolutions implements OutputDataProvider for outputting an array of tags
type OutputECloudSolutions struct {
	Solutions []ecloud.Solution
}

func outputECloudSolutions(solutions []ecloud.Solution) {
	err := Output(&OutputECloudSolutions{Solutions: solutions})
	if err != nil {
		output.Fatalf("Failed to output solutions: %s", err)
	}
}

func (o *OutputECloudSolutions) GetData() interface{} {
	return o.Solutions
}

func (o *OutputECloudSolutions) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, solution := range o.Solutions {
		fields := o.getOrderedFields(solution)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputECloudSolutions) getOrderedFields(solution ecloud.Solution) *output.OrderedFields {
	fields := output.NewOrderedFields()
	fields.Set("id", output.NewFieldValue(strconv.Itoa(solution.ID), true))
	fields.Set("name", output.NewFieldValue(solution.Name, true))
	fields.Set("environment", output.NewFieldValue(solution.Environment.String(), true))

	return fields
}

// OutputECloudSites implements OutputDataProvider for outputting an array of sites
type OutputECloudSites struct {
	Sites []ecloud.Site
}

func outputECloudSites(sites []ecloud.Site) {
	err := Output(&OutputECloudSites{Sites: sites})
	if err != nil {
		output.Fatalf("Failed to output sites: %s", err)
	}
}

func (o *OutputECloudSites) GetData() interface{} {
	return o.Sites
}

func (o *OutputECloudSites) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, site := range o.Sites {
		fields := o.getOrderedFields(site)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputECloudSites) getOrderedFields(site ecloud.Site) *output.OrderedFields {
	fields := output.NewOrderedFields()
	fields.Set("id", output.NewFieldValue(strconv.Itoa(site.ID), true))
	fields.Set("state", output.NewFieldValue(site.State, true))
	fields.Set("solution_id", output.NewFieldValue(strconv.Itoa(site.SolutionID), true))
	fields.Set("pod_id", output.NewFieldValue(strconv.Itoa(site.PodID), true))

	return fields
}

// OutputECloudHosts implements OutputDataProvider for outputting an array of hosts
type OutputECloudHosts struct {
	Hosts []ecloud.Host
}

func outputECloudHosts(hosts []ecloud.Host) {
	err := Output(&OutputECloudHosts{Hosts: hosts})
	if err != nil {
		output.Fatalf("Failed to output hosts: %s", err)
	}
}

func (o *OutputECloudHosts) GetData() interface{} {
	return o.Hosts
}

func (o *OutputECloudHosts) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, host := range o.Hosts {
		fields := o.getOrderedFields(host)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputECloudHosts) getOrderedFields(host ecloud.Host) *output.OrderedFields {
	fields := output.NewOrderedFields()
	fields.Set("id", output.NewFieldValue(strconv.Itoa(host.ID), true))
	fields.Set("solution_id", output.NewFieldValue(strconv.Itoa(host.SolutionID), true))
	fields.Set("pod_id", output.NewFieldValue(strconv.Itoa(host.PodID), true))
	fields.Set("name", output.NewFieldValue(host.Name, true))
	fields.Set("cpu_quantity", output.NewFieldValue(strconv.Itoa(host.CPU.Quantity), true))
	fields.Set("cpu_cores", output.NewFieldValue(strconv.Itoa(host.CPU.Cores), true))
	fields.Set("cpu_speed", output.NewFieldValue(host.CPU.Speed, false))
	fields.Set("ram_capacity", output.NewFieldValue(strconv.Itoa(host.RAM.Capacity), true))
	fields.Set("ram_reserved", output.NewFieldValue(strconv.Itoa(host.RAM.Reserved), false))
	fields.Set("ram_allocated", output.NewFieldValue(strconv.Itoa(host.RAM.Allocated), false))
	fields.Set("ram_available", output.NewFieldValue(strconv.Itoa(host.RAM.Available), false))

	return fields
}

// OutputECloudDatastores implements OutputDataProvider for outputting an array of hosts
type OutputECloudDatastores struct {
	Datastores []ecloud.Datastore
}

func outputECloudDatastores(datastores []ecloud.Datastore) {
	err := Output(&OutputECloudDatastores{Datastores: datastores})
	if err != nil {
		output.Fatalf("Failed to output datastores: %s", err)
	}
}

func (o *OutputECloudDatastores) GetData() interface{} {
	return o.Datastores
}

func (o *OutputECloudDatastores) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, datastore := range o.Datastores {
		fields := o.getOrderedFields(datastore)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputECloudDatastores) getOrderedFields(datastore ecloud.Datastore) *output.OrderedFields {
	fields := output.NewOrderedFields()
	fields.Set("id", output.NewFieldValue(strconv.Itoa(datastore.ID), true))
	fields.Set("solution_id", output.NewFieldValue(strconv.Itoa(datastore.SolutionID), true))
	fields.Set("site_id", output.NewFieldValue(strconv.Itoa(datastore.SiteID), true))
	fields.Set("name", output.NewFieldValue(datastore.Name, true))
	fields.Set("status", output.NewFieldValue(datastore.Status.String(), true))
	fields.Set("capacity", output.NewFieldValue(strconv.Itoa(datastore.Capacity), true))
	fields.Set("allocated", output.NewFieldValue(strconv.Itoa(datastore.Allocated), false))
	fields.Set("available", output.NewFieldValue(strconv.Itoa(datastore.Available), false))

	return fields
}

// OutputECloudTemplates implements OutputDataProvider for outputting an array of hosts
type OutputECloudTemplates struct {
	Templates []ecloud.Template
}

func outputECloudTemplates(templates []ecloud.Template) {
	err := Output(&OutputECloudTemplates{Templates: templates})
	if err != nil {
		output.Fatalf("Failed to output templates: %s", err)
	}
}

func (o *OutputECloudTemplates) GetData() interface{} {
	return o.Templates
}

func (o *OutputECloudTemplates) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, template := range o.Templates {
		fields := o.getOrderedFields(template)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputECloudTemplates) getOrderedFields(template ecloud.Template) *output.OrderedFields {
	fields := output.NewOrderedFields()
	fields.Set("name", output.NewFieldValue(template.Name, true))
	fields.Set("cpu", output.NewFieldValue(strconv.Itoa(template.CPU), true))
	fields.Set("ram_gb", output.NewFieldValue(strconv.Itoa(template.RAM), true))
	fields.Set("hdd_gb", output.NewFieldValue(strconv.Itoa(template.HDD), true))
	fields.Set("platform", output.NewFieldValue(template.Platform, true))
	fields.Set("operating_system", output.NewFieldValue(template.OperatingSystem, true))
	fields.Set("solution_id", output.NewFieldValue(strconv.Itoa(template.SolutionID), true))

	return fields
}

// OutputECloudNetworks implements OutputDataProvider for outputting an array of hosts
type OutputECloudNetworks struct {
	Networks []ecloud.Network
}

func outputECloudNetworks(networks []ecloud.Network) {
	err := Output(&OutputECloudNetworks{Networks: networks})
	if err != nil {
		output.Fatalf("Failed to output networks: %s", err)
	}
}

func (o *OutputECloudNetworks) GetData() interface{} {
	return o.Networks
}

func (o *OutputECloudNetworks) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, network := range o.Networks {
		fields := o.getOrderedFields(network)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputECloudNetworks) getOrderedFields(network ecloud.Network) *output.OrderedFields {
	fields := output.NewOrderedFields()
	fields.Set("id", output.NewFieldValue(strconv.Itoa(network.ID), true))
	fields.Set("name", output.NewFieldValue(network.Name, true))

	return fields
}

// OutputECloudFirewalls implements OutputDataProvider for outputting an array of hosts
type OutputECloudFirewalls struct {
	Firewalls []ecloud.Firewall
}

func outputECloudFirewalls(firewalls []ecloud.Firewall) {
	err := Output(&OutputECloudFirewalls{Firewalls: firewalls})
	if err != nil {
		output.Fatalf("Failed to output firewalls: %s", err)
	}
}

func (o *OutputECloudFirewalls) GetData() interface{} {
	return o.Firewalls
}

func (o *OutputECloudFirewalls) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, firewall := range o.Firewalls {
		fields := o.getOrderedFields(firewall)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputECloudFirewalls) getOrderedFields(firewall ecloud.Firewall) *output.OrderedFields {
	fields := output.NewOrderedFields()
	fields.Set("id", output.NewFieldValue(strconv.Itoa(firewall.ID), true))
	fields.Set("name", output.NewFieldValue(firewall.Name, true))
	fields.Set("hostname", output.NewFieldValue(firewall.Hostname, true))
	fields.Set("ip", output.NewFieldValue(firewall.IP.String(), true))
	fields.Set("role", output.NewFieldValue(firewall.Role.String(), true))

	return fields
}

// OutputECloudPods implements OutputDataProvider for outputting an array of pods
type OutputECloudPods struct {
	Pods []ecloud.Pod
}

func outputECloudPods(pods []ecloud.Pod) {
	err := Output(&OutputECloudPods{Pods: pods})
	if err != nil {
		output.Fatalf("Failed to output pods: %s", err)
	}
}

func (o *OutputECloudPods) GetData() interface{} {
	return o.Pods
}

func (o *OutputECloudPods) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, pod := range o.Pods {
		fields := o.getOrderedFields(pod)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputECloudPods) getOrderedFields(pod ecloud.Pod) *output.OrderedFields {
	fields := output.NewOrderedFields()
	fields.Set("id", output.NewFieldValue(strconv.Itoa(pod.ID), true))
	fields.Set("name", output.NewFieldValue(pod.Name, true))

	return fields
}

// OutputECloudAppliances implements OutputDataProvider for outputting an array of appliances
type OutputECloudAppliances struct {
	Appliances []ecloud.Appliance
}

func outputECloudAppliances(appliances []ecloud.Appliance) {
	err := Output(&OutputECloudAppliances{Appliances: appliances})
	if err != nil {
		output.Fatalf("Failed to output appliances: %s", err)
	}
}

func (o *OutputECloudAppliances) GetData() interface{} {
	return o.Appliances
}

func (o *OutputECloudAppliances) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, appliance := range o.Appliances {
		fields := o.getOrderedFields(appliance)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputECloudAppliances) getOrderedFields(appliance ecloud.Appliance) *output.OrderedFields {
	fields := output.NewOrderedFields()
	fields.Set("id", output.NewFieldValue(appliance.ID, true))
	fields.Set("name", output.NewFieldValue(appliance.Name, true))
	fields.Set("logo_uri", output.NewFieldValue(appliance.LogoURI, false))
	fields.Set("description", output.NewFieldValue(appliance.Description, false))
	fields.Set("documentation_uri", output.NewFieldValue(appliance.DocumentationURI, false))
	fields.Set("publisher", output.NewFieldValue(appliance.Publisher, true))
	fields.Set("created_at", output.NewFieldValue(appliance.CreatedAt.String(), true))

	return fields
}

// OutputECloudApplianceParameters implements OutputDataProvider for outputting an array of appliance parameters
type OutputECloudApplianceParameters struct {
	ApplianceParameters []ecloud.ApplianceParameter
}

func outputECloudApplianceParameters(parameters []ecloud.ApplianceParameter) {
	err := Output(&OutputECloudApplianceParameters{ApplianceParameters: parameters})
	if err != nil {
		output.Fatalf("Failed to output appliance parameters: %s", err)
	}
}

func (o *OutputECloudApplianceParameters) GetData() interface{} {
	return o.ApplianceParameters
}

func (o *OutputECloudApplianceParameters) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, parameter := range o.ApplianceParameters {
		fields := o.getOrderedFields(parameter)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputECloudApplianceParameters) getOrderedFields(parameter ecloud.ApplianceParameter) *output.OrderedFields {
	fields := output.NewOrderedFields()
	fields.Set("id", output.NewFieldValue(parameter.ID, true))
	fields.Set("name", output.NewFieldValue(parameter.Name, true))
	fields.Set("key", output.NewFieldValue(parameter.Key, true))
	fields.Set("type", output.NewFieldValue(parameter.Type, true))
	fields.Set("description", output.NewFieldValue(parameter.Description, true))
	fields.Set("required", output.NewFieldValue(strconv.FormatBool(parameter.Required), true))
	fields.Set("validation_rule", output.NewFieldValue(parameter.ValidationRule, false))

	return fields
}

// GetCreateTagRequestFromStringArrayFlag returns an array of CreateTagRequest structs from given tag string array flag
func GetCreateTagRequestFromStringArrayFlag(tagsFlag []string) ([]ecloud.CreateTagRequest, error) {
	var tags []ecloud.CreateTagRequest
	for _, tagFlag := range tagsFlag {
		key, value, err := GetKeyValueFromStringFlag(tagFlag)
		if err != nil {
			return tags, err
		}

		tags = append(tags, ecloud.CreateTagRequest{Key: key, Value: value})
	}

	return tags, nil
}

// GetCreateVirtualMachineRequestParameterFromStringArrayFlag returns an array of CreateVirtualMachineRequestParameter structs from given string array flag
func GetCreateVirtualMachineRequestParameterFromStringArrayFlag(parametersFlag []string) ([]ecloud.CreateVirtualMachineRequestParameter, error) {
	var parameters []ecloud.CreateVirtualMachineRequestParameter
	for _, parameterFlag := range parametersFlag {
		key, value, err := GetKeyValueFromStringFlag(parameterFlag)
		if err != nil {
			return parameters, err
		}

		parameters = append(parameters, ecloud.CreateVirtualMachineRequestParameter{Key: key, Value: value})
	}

	return parameters, nil
}

// GetKeyValueFromStringFlag returns a string map from given string flag. Expects format 'key=value'
func GetKeyValueFromStringFlag(flag string) (key, value string, err error) {
	if flag == "" {
		return key, value, errors.New("Missing key/value")
	}

	parts := strings.Split(flag, "=")
	if len(parts) < 2 || len(parts) > 2 {
		return key, value, errors.New("Invalid format, expecting: key=value")
	}
	if parts[0] == "" {
		return key, value, errors.New("Missing key")
	}
	if parts[1] == "" {
		return key, value, errors.New("Missing value")
	}

	return parts[0], parts[1], nil
}

// SolutionTemplateExistsWaitFunc returns WaitFunc for waiting for a template to exist
func SolutionTemplateExistsWaitFunc(service ecloud.ECloudService, solutionID int, templateName string, exists bool) WaitFunc {
	return func() (finished bool, err error) {
		_, err = service.GetSolutionTemplate(solutionID, templateName)
		if err != nil {
			if _, ok := err.(*ecloud.TemplateNotFoundError); ok {
				return (exists == false), nil
			}

			return false, fmt.Errorf("Failed to retrieve solution template [%s]: %s", templateName, err.Error())
		}

		return (exists == true), nil
	}
}

// PodTemplateExistsWaitFunc returns WaitFunc for waiting for a template to exist
func PodTemplateExistsWaitFunc(service ecloud.ECloudService, podID int, templateName string, exists bool) WaitFunc {
	return func() (finished bool, err error) {
		_, err = service.GetPodTemplate(podID, templateName)
		if err != nil {
			if _, ok := err.(*ecloud.TemplateNotFoundError); ok {
				return (exists == false), nil
			}

			return false, fmt.Errorf("Failed to retrieve pod template [%s]: %s", templateName, err.Error())
		}

		return (exists == true), nil
	}
}
