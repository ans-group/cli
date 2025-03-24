package safedns

import (
	"strconv"
	"strings"

	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/safedns"
)

type ZoneCollection []safedns.Zone

func (z ZoneCollection) DefaultColumns() []string {
	return []string{"name", "description"}
}

func (z ZoneCollection) Fields() []*output.OrderedFields {
	var data []*output.OrderedFields
	for _, zone := range z {
		fields := output.NewOrderedFields()
		fields.Set("name", zone.Name)
		fields.Set("description", zone.Description)

		data = append(data, fields)
	}

	return data
}

type RecordCollection []safedns.Record

func (r RecordCollection) DefaultColumns() []string {
	return []string{"id", "name", "type", "content", "updated_at", "priority", "ttl"}
}

func (r RecordCollection) Fields() []*output.OrderedFields {
	var data []*output.OrderedFields
	for _, record := range r {
		fields := output.NewOrderedFields()
		fields.Set("id", strconv.Itoa(record.ID))
		fields.Set("name", record.Name)
		fields.Set("type", record.Type.String())
		fields.Set("content", record.Content)
		fields.Set("updated_at", record.UpdatedAt.String())
		fields.Set("priority", strconv.Itoa(record.Priority))
		fields.Set("ttl", strconv.Itoa(int(record.TTL)))

		data = append(data, fields)
	}

	return data
}

type NoteCollection []safedns.Note

func (n NoteCollection) DefaultColumns() []string {
	return []string{"id", "notes", "ip"}
}

func (n NoteCollection) Fields() []*output.OrderedFields {
	var data []*output.OrderedFields
	for _, note := range n {
		fields := output.NewOrderedFields()
		fields.Set("id", strconv.Itoa(note.ID))
		fields.Set("notes", note.Notes)
		fields.Set("ip", note.IP.String())

		data = append(data, fields)
	}

	return data
}

type TemplateCollection []safedns.Template

func (t TemplateCollection) DefaultColumns() []string {
	return []string{"id", "name", "default", "created_at"}
}

func (t TemplateCollection) Fields() []*output.OrderedFields {
	var data []*output.OrderedFields
	for _, template := range t {
		fields := output.NewOrderedFields()
		fields.Set("id", strconv.Itoa(template.ID))
		fields.Set("name", template.Name)
		fields.Set("default", strconv.FormatBool(template.Default))
		fields.Set("created_at", template.CreatedAt.String())

		data = append(data, fields)
	}

	return data
}

type SettingsCollection []safedns.Settings

func (s SettingsCollection) DefaultColumns() []string {
	return []string{"id", "custom_soa_allowed", "custom_base_ns_allowed", "delegation_allowed", "product"}
}

func (s SettingsCollection) Fields() []*output.OrderedFields {
	var data []*output.OrderedFields
	for _, setting := range s {
		nameservers := []string{}
		for _, nameserver := range setting.Nameservers {
			nameservers = append(nameservers, nameserver.Name)
		}

		fields := output.NewOrderedFields()
		fields.Set("id", strconv.Itoa(setting.ID))
		fields.Set("email", setting.Email)
		fields.Set("nameservers", strings.Join(nameservers, ", "))
		fields.Set("custom_soa_allowed", strconv.FormatBool(setting.CustomSOAAllowed))
		fields.Set("custom_base_ns_allowed", strconv.FormatBool(setting.CustomBaseNSAllowed))
		fields.Set("delegation_allowed", strconv.FormatBool(setting.DelegationAllowed))
		fields.Set("product", setting.Product)

		data = append(data, fields)
	}

	return data
}
