package layrz

import "strconv"

// ValidateID validates an ID field. Accepts an int-kind value or a string parseable
// as an integer. Value must be > 0. Returns accumulated errors.
func ValidateID(value any, r IDRules) []*FieldError {
	_, present := deref(value)
	if !present {
		if r.Required {
			return []*FieldError{{Code: "required"}}
		}
		return nil
	}

	// Wrong type -> invalid regardless of Required (fixes the "optional silently accepts wrong type" bug)
	// Explicitly reject bool (in case it sneaks through via any)
	if _, isBool := value.(bool); isBool {
		return []*FieldError{{Code: "invalid"}}
	}

	var numVal int64

	switch v := value.(type) {
	case int:
		numVal = int64(v)
	case int64:
		numVal = v
	case int32:
		numVal = int64(v)
	case int16:
		numVal = int64(v)
	case int8:
		numVal = int64(v)
	case string:
		// Try to parse string as integer
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return []*FieldError{{Code: "invalid"}}
		}
		numVal = n
	default:
		return []*FieldError{{Code: "invalid"}}
	}

	if numVal <= 0 {
		return []*FieldError{{Code: "invalid"}}
	}

	return nil
}
