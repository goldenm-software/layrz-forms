package layrz

// Validate validates the form and returns every accumulated error, keyed by camelCase field name.
// form must be a non-nil pointer to a struct. Returns a map of field names to error slices;
// an empty map means validation passed.
func Validate(form any) Errors {
	return validateForm(form)
}

// IsValid reports whether the form passes validation (i.e., has no errors).
func IsValid(form any) bool {
	return Validate(form).IsEmpty()
}
