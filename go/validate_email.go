package layrz

import "regexp"

// ValidateEmail validates an Email field. Returns accumulated errors.
func ValidateEmail(value *string, r EmailRules) []*FieldError {
	if value == nil {
		if r.Required {
			return []*FieldError{{Code: "required"}}
		}
		return nil
	}

	// Empty string with Empty==false -> "empty" (Python wrongly says "required")
	if *value == "" {
		if !r.Empty {
			return []*FieldError{{Code: "empty"}}
		}
		// Empty is explicitly permitted, no regex validation
		return nil
	}

	// Non-empty value: always validate regex, even if Empty==true
	regex := r.Regex
	if regex == "" {
		regex = DefaultEmailRegex
	}

	// Replicate Python's re.match semantics: anchored at start only
	// If pattern doesn't start with ^, prefix it
	if regex[0] != '^' {
		regex = "^" + regex
	}

	compiled, _ := regexp.Compile(regex)
	if compiled == nil {
		// Shouldn't happen if ParseTag validated it, but be safe
		return []*FieldError{{Code: "invalid"}}
	}

	if !compiled.MatchString(*value) {
		return []*FieldError{{Code: "invalid"}}
	}

	return nil
}
