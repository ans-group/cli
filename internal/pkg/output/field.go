package output

// OrderedFields holds a string map with field values, and a slice of keys for
// maintaining order
type OrderedFields struct {
	m    map[string]FieldValue
	keys []string
}

// NewOrderedFields returns a pointer to an initialized OrderedFields struct
func NewOrderedFields() *OrderedFields {
	return &OrderedFields{
		m: make(map[string]FieldValue),
	}
}

// Set adds/updates given key k with FieldValue v
func (o *OrderedFields) Set(k string, v FieldValue) {
	exists := o.Exists(k)
	o.m[k] = v
	if !exists {
		o.keys = append(o.keys, k)
	}
}

// Get retrieves FieldValue for given key k
func (o *OrderedFields) Get(k string) FieldValue {
	return o.m[k]
}

// Exists returns true if given key k exists, otherwise false
func (o *OrderedFields) Exists(k string) bool {
	_, exists := o.m[k]
	return exists
}

// Keys returns a list of ordered keys
func (o *OrderedFields) Keys() []string {
	return o.keys
}

// FieldValue holds the value for a table field
type FieldValue struct {
	Value   string
	Default bool
}

// NewFieldValue returns a new, initialized FieldValue struct
func NewFieldValue(value string, def bool) FieldValue {
	return FieldValue{
		Value:   value,
		Default: def,
	}
}
