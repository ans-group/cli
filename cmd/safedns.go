package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/cli/internal/pkg/resource"
	"github.com/ukfast/sdk-go/pkg/connection"
	"github.com/ukfast/sdk-go/pkg/service/safedns"
)

func safednsRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "safedns",
		Short: "Commands relating to SafeDNS service",
	}

	// Child root commands
	cmd.AddCommand(safednsZoneRootCmd())
	cmd.AddCommand(safednsZoneRecordRootCmd())
	cmd.AddCommand(safednsZoneNoteRootCmd())
	cmd.AddCommand(safednsTemplateRootCmd())
	cmd.AddCommand(safednsSettingsRootCmd())

	return cmd
}

type SafeDNSTemplateLocatorProvider struct {
	service safedns.SafeDNSService
}

func NewSafeDNSTemplateLocatorProvider(service safedns.SafeDNSService) *SafeDNSTemplateLocatorProvider {
	return &SafeDNSTemplateLocatorProvider{service: service}
}

func (p *SafeDNSTemplateLocatorProvider) SupportedProperties() []string {
	return []string{"name"}
}

func (p *SafeDNSTemplateLocatorProvider) Locate(property string, value string) (interface{}, error) {
	params := connection.APIRequestParameters{}
	params.WithFilter(connection.APIRequestFiltering{Property: property, Operator: connection.EQOperator, Value: []string{value}})

	return p.service.GetTemplates(params)
}

func getSafeDNSTemplateByNameOrID(service safedns.SafeDNSService, nameOrID string) (safedns.Template, error) {
	templateID, err := strconv.Atoi(nameOrID)
	if err != nil {
		locator := resource.NewResourceLocator(NewSafeDNSTemplateLocatorProvider(service))

		template, err := locator.Invoke(nameOrID)
		if err != nil {
			return safedns.Template{}, fmt.Errorf("Error locating template [%s]: %s", nameOrID, err)
		}

		return template.(safedns.Template), nil
	}

	template, err := service.GetTemplate(templateID)
	if err != nil {
		return safedns.Template{}, fmt.Errorf("Error retrieving template by ID [%d]: %s", templateID, err)
	}

	return template, nil
}

func getSafeDNSTemplateIDByNameOrID(service safedns.SafeDNSService, nameOrID string) (int, error) {
	templateID, err := strconv.Atoi(nameOrID)
	if err != nil {
		locator := resource.NewResourceLocator(NewSafeDNSTemplateLocatorProvider(service))

		template, err := locator.Invoke(nameOrID)
		if err != nil {
			return 0, fmt.Errorf("Error locating template [%s]: %s", nameOrID, err)
		}

		return template.(safedns.Template).ID, nil
	}

	return templateID, nil
}

// OutputSafeDNSZones implements OutputDataProvider for outputting an array of Zones
type OutputSafeDNSZones struct {
	Zones []safedns.Zone
}

func outputSafeDNSZones(zones []safedns.Zone) {
	err := Output(&OutputSafeDNSZones{Zones: zones})
	if err != nil {
		output.Fatalf("Failed to output zones: %s", err)
	}
}

func (o *OutputSafeDNSZones) GetData() interface{} {
	return o.Zones
}

func (o *OutputSafeDNSZones) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, zone := range o.Zones {
		fields := o.getOrderedFields(zone)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputSafeDNSZones) getOrderedFields(zone safedns.Zone) *output.OrderedFields {
	fields := output.NewOrderedFields()
	fields.Set("name", output.NewFieldValue(zone.Name, true))
	fields.Set("description", output.NewFieldValue(zone.Description, true))

	return fields
}

// OutputSafeDNSRecords implements OutputDataProvider for outputting an array of Records
type OutputSafeDNSRecords struct {
	Records []safedns.Record
}

func outputSafeDNSRecords(records []safedns.Record) {
	err := Output(&OutputSafeDNSRecords{Records: records})
	if err != nil {
		output.Fatalf("Failed to output records: %s", err)
	}
}

func (o *OutputSafeDNSRecords) GetData() interface{} {
	return o.Records
}

func (o *OutputSafeDNSRecords) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, record := range o.Records {
		OutputSafeDNSRecords := o.getOrderedFields(record)
		data = append(data, OutputSafeDNSRecords)
	}

	return data, nil
}

func (o *OutputSafeDNSRecords) getOrderedFields(record safedns.Record) *output.OrderedFields {
	fields := output.NewOrderedFields()
	fields.Set("id", output.NewFieldValue(strconv.Itoa(record.ID), true))
	fields.Set("name", output.NewFieldValue(record.Name, true))
	fields.Set("type", output.NewFieldValue(record.Type.String(), true))
	fields.Set("content", output.NewFieldValue(record.Content, true))
	fields.Set("updated_at", output.NewFieldValue(record.UpdatedAt.String(), true))
	fields.Set("priority", output.NewFieldValue(strconv.Itoa(record.Priority), true))
	fields.Set("ttl", output.NewFieldValue(strconv.Itoa(int(record.TTL)), true))

	return fields
}

// OutputSafeDNSNotes implements OutputDataProvider for outputting an array of Notes
type OutputSafeDNSNotes struct {
	Notes []safedns.Note
}

func outputSafeDNSNotes(notes []safedns.Note) {
	err := Output(&OutputSafeDNSNotes{Notes: notes})
	if err != nil {
		output.Fatalf("Failed to output notes: %s", err)
	}
}

func (o *OutputSafeDNSNotes) GetData() interface{} {
	return o.Notes
}

func (o *OutputSafeDNSNotes) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, note := range o.Notes {
		fields := o.getOrderedFields(note)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputSafeDNSNotes) getOrderedFields(note safedns.Note) *output.OrderedFields {
	fields := output.NewOrderedFields()
	fields.Set("id", output.NewFieldValue(strconv.Itoa(note.ID), true))
	fields.Set("notes", output.NewFieldValue(note.Notes, true))
	fields.Set("ip", output.NewFieldValue(note.IP.String(), true))

	return fields
}

// OutputSafeDNSTemplates implements OutputDataProvider for outputting an array of Templates
type OutputSafeDNSTemplates struct {
	Templates []safedns.Template
}

func outputSafeDNSTemplates(templates []safedns.Template) {
	err := Output(&OutputSafeDNSTemplates{Templates: templates})
	if err != nil {
		output.Fatalf("Failed to output templates: %s", err)
	}
}

func (o *OutputSafeDNSTemplates) GetData() interface{} {
	return o.Templates
}

func (o *OutputSafeDNSTemplates) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, template := range o.Templates {
		fields := o.getOrderedFields(template)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputSafeDNSTemplates) getOrderedFields(template safedns.Template) *output.OrderedFields {
	fields := output.NewOrderedFields()
	fields.Set("id", output.NewFieldValue(strconv.Itoa(template.ID), true))
	fields.Set("name", output.NewFieldValue(template.Name, true))
	fields.Set("default", output.NewFieldValue(strconv.FormatBool(template.Default), true))
	fields.Set("created_at", output.NewFieldValue(template.CreatedAt.String(), true))

	return fields
}

// OutputSafeDNSSettings implements OutputDataProvider for outputting an array of Settings
type OutputSafeDNSSettings struct {
	Settings []safedns.Settings
}

func outputSafeDNSSettings(settings []safedns.Settings) {
	err := Output(&OutputSafeDNSSettings{Settings: settings})
	if err != nil {
		output.Fatalf("Failed to output settings: %s", err)
	}
}

func (o *OutputSafeDNSSettings) GetData() interface{} {
	return o.Settings
}

func (o *OutputSafeDNSSettings) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, settings := range o.Settings {
		fields := o.getOrderedFields(settings)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputSafeDNSSettings) getOrderedFields(settings safedns.Settings) *output.OrderedFields {
	nameservers := []string{}
	for _, nameserver := range settings.Nameservers {
		nameservers = append(nameservers, nameserver.Name)
	}

	fields := output.NewOrderedFields()
	fields.Set("id", output.NewFieldValue(strconv.Itoa(settings.ID), true))
	fields.Set("email", output.NewFieldValue(settings.Email, false))
	fields.Set("nameservers", output.NewFieldValue(strings.Join(nameservers, ", "), false))
	fields.Set("custom_soa_allowed", output.NewFieldValue(strconv.FormatBool(settings.CustomSOAAllowed), true))
	fields.Set("custom_base_ns_allowed", output.NewFieldValue(strconv.FormatBool(settings.CustomBaseNSAllowed), true))
	fields.Set("delegation_allowed", output.NewFieldValue(strconv.FormatBool(settings.DelegationAllowed), true))
	fields.Set("product", output.NewFieldValue(settings.Product, true))

	return fields
}
