package layrz

import (
	"encoding/json"
	"sort"
)

// FieldError represents a single validation error on a form field.
type FieldError struct {
	Code     string         `json:"code"`
	Expected any            `json:"expected,omitempty"`
	Received any            `json:"received,omitempty"`
	Extra    map[string]any `json:"extra,omitempty"`
}

// Errors maps field names (in camelCase) to slices of FieldError.
// The zero value is ready to use.
type Errors map[string][]*FieldError

// Add appends errors to the field key. It creates the slice lazily and skips nil entries.
// Must not be called on a nil Errors value.
func (e Errors) Add(key string, errs ...*FieldError) {
	if len(errs) == 0 {
		return
	}

	// Filter out nil entries
	var nonNil []*FieldError
	for _, err := range errs {
		if err != nil {
			nonNil = append(nonNil, err)
		}
	}

	if len(nonNil) == 0 {
		return
	}

	e[key] = append(e[key], nonNil...)
}

// Merge appends all errors from other into e, converting field keys to camelCase.
func (e Errors) Merge(other Errors) {
	for key, errs := range other {
		camelKey := ToCamelCase(key)
		e[camelKey] = append(e[camelKey], errs...)
	}
}

// IsEmpty reports whether e contains no errors.
func (e Errors) IsEmpty() bool {
	return len(e) == 0
}

// Keys returns all field keys in e, sorted in ascending order for determinism.
func (e Errors) Keys() []string {
	keys := make([]string, 0, len(e))
	for k := range e {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// MarshalJSON encodes e to JSON. This satisfies json.Marshaler and ensures
// that fields like Expected:0 and Received:"" are included in the JSON output.
func (e Errors) MarshalJSON() ([]byte, error) {
	// Convert to a plain map[string]interface{} to avoid custom marshaling issues
	m := make(map[string]any)
	for k, errs := range e {
		m[k] = errs
	}
	return json.Marshal(m)
}
