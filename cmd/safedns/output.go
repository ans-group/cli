package safedns

import (
	"strconv"
	"strings"

	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/safedns"
)

func OutputSafeDNSZonesProvider(zones []safedns.Zone) output.OutputHandlerDataProvider {
	return output.NewGenericOutputHandlerDataProvider(
		output.WithData(zones),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, zone := range zones {
				fields := output.NewOrderedFields()
				fields.Set("name", output.NewFieldValue(zone.Name, true))
				fields.Set("description", output.NewFieldValue(zone.Description, true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputSafeDNSRecordsProvider(records []safedns.Record) output.OutputHandlerDataProvider {
	return output.NewGenericOutputHandlerDataProvider(
		output.WithData(records),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, record := range records {
				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(strconv.Itoa(record.ID), true))
				fields.Set("name", output.NewFieldValue(record.Name, true))
				fields.Set("type", output.NewFieldValue(record.Type.String(), true))
				fields.Set("content", output.NewFieldValue(record.Content, true))
				fields.Set("updated_at", output.NewFieldValue(record.UpdatedAt.String(), true))
				fields.Set("priority", output.NewFieldValue(strconv.Itoa(record.Priority), true))
				fields.Set("ttl", output.NewFieldValue(strconv.Itoa(int(record.TTL)), true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputSafeDNSNotesProvider(notes []safedns.Note) output.OutputHandlerDataProvider {
	return output.NewGenericOutputHandlerDataProvider(
		output.WithData(notes),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, note := range notes {
				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(strconv.Itoa(note.ID), true))
				fields.Set("notes", output.NewFieldValue(note.Notes, true))
				fields.Set("ip", output.NewFieldValue(note.IP.String(), true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputSafeDNSTemplatesProvider(templates []safedns.Template) output.OutputHandlerDataProvider {
	return output.NewGenericOutputHandlerDataProvider(
		output.WithData(templates),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, template := range templates {
				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(strconv.Itoa(template.ID), true))
				fields.Set("name", output.NewFieldValue(template.Name, true))
				fields.Set("default", output.NewFieldValue(strconv.FormatBool(template.Default), true))
				fields.Set("created_at", output.NewFieldValue(template.CreatedAt.String(), true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputSafeDNSSettingsProvider(settings []safedns.Settings) output.OutputHandlerDataProvider {
	return output.NewGenericOutputHandlerDataProvider(
		output.WithData(settings),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, setting := range settings {
				nameservers := []string{}
				for _, nameserver := range setting.Nameservers {
					nameservers = append(nameservers, nameserver.Name)
				}

				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(strconv.Itoa(setting.ID), true))
				fields.Set("email", output.NewFieldValue(setting.Email, false))
				fields.Set("nameservers", output.NewFieldValue(strings.Join(nameservers, ", "), false))
				fields.Set("custom_soa_allowed", output.NewFieldValue(strconv.FormatBool(setting.CustomSOAAllowed), true))
				fields.Set("custom_base_ns_allowed", output.NewFieldValue(strconv.FormatBool(setting.CustomBaseNSAllowed), true))
				fields.Set("delegation_allowed", output.NewFieldValue(strconv.FormatBool(setting.DelegationAllowed), true))
				fields.Set("product", output.NewFieldValue(setting.Product, true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}
