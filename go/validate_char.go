package layrz

import (
	"regexp"
	"unicode/utf8"
)

// ValidateChar validates a Char field. Performs all checks in order and accumulates errors.
// Length is measured in UTF-8 rune count, not bytes.
func ValidateChar(value *string, r CharRules) []*FieldError {
	if value == nil {
		if r.Required {
			return []*FieldError{{Code: "required"}}
		}
		return nil
	}

	var errs []*FieldError

	// Length in runes (code points), not bytes
	runeLen := utf8.RuneCountInString(*value)

	// 1. Empty==false && len==0 -> "empty"
	if !r.Empty && runeLen == 0 {
		errs = append(errs, &FieldError{Code: "empty"})
	}

	// 2. MaxLength set && len > max
	if r.MaxLength != nil && runeLen > *r.MaxLength {
		errs = append(errs, &FieldError{
			Code:     "maxLength",
			Expected: *r.MaxLength,
			Received: runeLen,
		})
	}

	// 3. MinLength set && len < min
	if r.MinLength != nil && runeLen < *r.MinLength {
		errs = append(errs, &FieldError{
			Code:     "minLength",
			Expected: *r.MinLength,
			Received: runeLen,
		})
	}

	// 4. Choices set && value not in choices
	if len(r.Choices) > 0 {
		found := false
		for _, choice := range r.Choices {
			if choice == *value {
				found = true
				break
			}
		}
		if !found {
			errs = append(errs, &FieldError{
				Code:     "invalidChoice",
				Expected: r.Choices,
				Received: *value,
			})
		}
	}

	// 5. Regex set && no match
	if r.Regex != "" {
		compiled, _ := regexp.Compile(r.Regex)
		if compiled != nil && !compiled.MatchString(*value) {
			errs = append(errs, &FieldError{
				Code:     "invalidFormat",
				Expected: r.Regex,
				Received: *value,
			})
		}
	}

	return errs
}
