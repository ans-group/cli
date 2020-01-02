package output

import (
	"fmt"
	"strings"
)

type UnsupportedFormatHandler func() error

type OutputHandler struct {
	Format                   string
	Provider                 OutputHandlerProvider
	Properties               []string
	Template                 string
	SupportedFormats         []string
	UnsupportedFormatHandler UnsupportedFormatHandler
}

func NewOutputHandler(out OutputHandlerProvider, format string) *OutputHandler {
	if format == "" {
		format = "table"
	}

	return &OutputHandler{
		Provider: out,
		Format:   format,
	}
}

// Handle calls the relevant OutputProvider data retrieval methods for given value
// in struct property 'Format'
func (o *OutputHandler) Handle() error {
	if !o.supportedFormat() {
		if o.UnsupportedFormatHandler != nil {
			return o.UnsupportedFormatHandler()
		}

		return fmt.Errorf("Unsupported output format [%s], supported formats: %s", o.Format, strings.Join(o.SupportedFormats, ", "))
	}

	switch o.Format {
	case "json":
		return JSON(o.Provider.GetData())
	case "template":
		return Template(o.Template, o.Provider.GetData())
	case "value":
		d, err := o.Provider.GetFieldData()
		if err != nil {
			return err
		}
		return Value(o.Properties, d)
	case "csv":
		d, err := o.Provider.GetFieldData()
		if err != nil {
			return err
		}
		return CSV(o.Properties, d)
	case "list":
		d, err := o.Provider.GetFieldData()
		if err != nil {
			return err
		}
		return List(o.Properties, d)
	default:
		Errorf("Invalid output format [%s], defaulting to 'table'", o.Format)
		fallthrough
	case "table":
		d, err := o.Provider.GetFieldData()
		if err != nil {
			return err
		}
		return Table(o.Properties, d)
	}
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
