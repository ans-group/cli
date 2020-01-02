package resource

import (
	"fmt"
	"reflect"
)

type ResourceLocatorProvider interface {
	Locate(property string, value string) (interface{}, error)
	SupportedProperties() []string
}

type ResourceLocator struct {
	Provider ResourceLocatorProvider
}

func NewResourceLocator(provider ResourceLocatorProvider) *ResourceLocator {
	return &ResourceLocator{Provider: provider}
}

func (f *ResourceLocator) Invoke(filter string) (interface{}, error) {
	for _, property := range f.Provider.SupportedProperties() {
		items, err := f.Provider.Locate(property, filter)
		if err != nil {
			return nil, fmt.Errorf("Error retrieving items: %s", err)
		}

		if items != nil {
			kind := reflect.TypeOf(items).Kind()
			switch kind {
			case reflect.Slice:
				s := reflect.ValueOf(items)
				length := s.Len()

				if length > 1 {
					return nil, fmt.Errorf("More than one item found matching [%s] (%s)", filter, property)
				}

				if length == 1 {
					return s.Index(0).Interface(), nil
				}
			default:
				return nil, fmt.Errorf("Unsupported non-slice type [%s]", kind.String())
			}
		}
	}

	return nil, fmt.Errorf("No items found matching [%s]", filter)
}
