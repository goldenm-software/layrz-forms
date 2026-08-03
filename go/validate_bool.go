package layrz

// ValidateBool validates a Bool field. Rejects any non-bool type.
func ValidateBool(value *bool, r BoolRules) []*FieldError {
	if value == nil {
		if r.Required {
			return []*FieldError{{Code: "required"}}
		}
		return nil
	}

	// If we got here, value is a valid *bool, so no errors
	return nil
}
