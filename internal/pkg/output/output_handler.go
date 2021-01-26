package output

import (
	"fmt"
	"strings"

	"github.com/ryanuber/go-glob"
)

type OutputHandlerOpts map[string]interface{}

type OutputHandler struct {
	Format           string
	FormatArg        string
	Properties       []string
	SupportedFormats []string
	DataProvider     OutputHandlerDataProvider
}

func NewOutputHandler(dataProvider OutputHandlerDataProvider, format string, formatArg string) *OutputHandler {
	if format == "" {
		format = "table"
	}

	return &OutputHandler{
		DataProvider: dataProvider,
		Format:       format,
		FormatArg:    formatArg,
	}
}

func (o *OutputHandler) WithSupportedFormats(formats []string) *OutputHandler {
	o.SupportedFormats = formats
	return o
}

// Handle calls the relevant OutputProvider data retrieval methods for given value
// in struct property 'Format'
func (o *OutputHandler) Handle() error {
	if !o.supportedFormat() {
		return fmt.Errorf("Unsupported output format [%s], supported formats: %s", o.Format, strings.Join(o.SupportedFormats, ", "))
	}

	switch o.Format {
	case "json":
		return JSON(o.DataProvider.GetData())
	case "jsonpath":
		return JSONPath(o.FormatArg, o.DataProvider.GetData())
	case "template":
		return Template(o.FormatArg, o.DataProvider.GetData())
	case "value":
		d, err := o.getProcessedFieldData()
		if err != nil {
			return err
		}
		return Value(o.Properties, d)
	case "csv":
		d, err := o.getProcessedFieldData()
		if err != nil {
			return err
		}
		return CSV(o.Properties, d)
	case "list":
		d, err := o.getProcessedFieldData()
		if err != nil {
			return err
		}
		return List(o.Properties, d)
	default:
		Errorf("Invalid output format [%s], defaulting to 'table'", o.Format)
		fallthrough
	case "table":
		d, err := o.getProcessedFieldData()
		if err != nil {
			return err
		}
		return Table(d)
	}
}

func (o *OutputHandler) getProcessedFieldData() ([]*OrderedFields, error) {
	var filteredFieldsCollectionArray []*OrderedFields

	fieldsCollectionArray, err := o.DataProvider.GetFieldData()
	if err != nil {
		return nil, err
	}

	for _, fieldCollection := range fieldsCollectionArray {
		filteredFieldsCollection := NewOrderedFields()

		if len(o.Properties) > 0 {
			// For each property in o.Properties, add field to filteredFields array
			// if that property exists in fields
			for _, prop := range o.Properties {
				for _, fieldKey := range fieldCollection.Keys() {
					if glob.Glob(strings.ToLower(prop), fieldKey) {
						filteredFieldsCollection.Set(fieldKey, fieldCollection.Get(fieldKey))
					}
				}
			}

		} else {
			isDefaultField := func(key string) bool {
				for _, defaultFieldKey := range o.DataProvider.DefaultFields() {
					if key == defaultFieldKey {
						return true
					}
				}

				return false
			}

			// Use default fields
			for _, fieldKey := range fieldCollection.Keys() {
				if isDefaultField(fieldKey) {
					filteredFieldsCollection.Set(fieldKey, fieldCollection.Get(fieldKey))
				}
			}
		}

		filteredFieldsCollectionArray = append(filteredFieldsCollectionArray, filteredFieldsCollection)
	}

	return filteredFieldsCollectionArray, nil
}

func (o *OutputHandler) supportedFormat() bool {
	if o.SupportedFormats == nil {
		return true
	}

	for _, supportedFormat := range o.SupportedFormats {
		if strings.ToLower(supportedFormat) == o.Format {
			return true
		}
	}

	return false
}
