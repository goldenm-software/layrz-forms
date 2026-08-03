package layrz

import "reflect"

// ValidateJSON validates a JSON field. Checks container type (list vs dict) and
// optionally validates that it's not empty.
func ValidateJSON(value any, r JSONRules) []*FieldError {
	_, present := deref(value)
	if !present {
		if r.Required {
			return []*FieldError{{Code: "required"}}
		}
		return nil
	}

	// Check the container type
	v := reflect.ValueOf(value)

	// Dereference pointers
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			if r.Required {
				return []*FieldError{{Code: "required"}}
			}
			return nil
		}
		v = v.Elem()
	}

	var isDict, isList bool

	switch v.Kind() {
	case reflect.Map:
		isDict = true
	case reflect.Slice, reflect.Array:
		isList = true
	default:
		// Wrong container type
		return []*FieldError{{Code: "invalid"}}
	}

	// Check if the container type matches the expected datatype
	if r.Datatype == datatypeDict && !isDict {
		return []*FieldError{{Code: "invalid"}}
	}
	if r.Datatype == datatypeList && !isList {
		return []*FieldError{{Code: "invalid"}}
	}

	// Check if empty (and that's not allowed)
	if !r.Empty && v.Len() == 0 {
		return []*FieldError{{Code: "invalid"}}
	}

	return nil
}
