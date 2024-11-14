package output

// OrderedFields holds a string map with field values, and a slice of keys for
// maintaining order
type OrderedFields struct {
	m    map[string]string
	keys []string
}

// NewOrderedFields returns a pointer to an initialized OrderedFields struct
func NewOrderedFields() *OrderedFields {
	return &OrderedFields{
		m: make(map[string]string),
	}
}

// Set adds/updates given key k with FieldValue v
func (o *OrderedFields) Set(k string, v string) {
	exists := o.Exists(k)
	o.m[k] = v
	if !exists {
		o.keys = append(o.keys, k)
	}
}

// Get retrieves FieldValue for given key k
func (o *OrderedFields) Get(k string) string {
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
