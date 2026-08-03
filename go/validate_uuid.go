package layrz

import (
	"regexp"
	"strings"
)

// ValidateUUID validates a UUID field. Accepts 8-4-4-4-12 hyphenated form,
// 32-hex-char unhyphenated form, braced forms, and urn:uuid: prefix.
// Returns accumulated errors.
func ValidateUUID(value *string, r UUIDRules) []*FieldError {
	if value == nil {
		if r.Required {
			return []*FieldError{{Code: "required"}}
		}
		return nil
	}

	if !isValidUUID(*value) {
		return []*FieldError{{Code: "invalid"}}
	}

	return nil
}

// isValidUUID checks if a string is a valid UUID in any of the accepted formats:
//   - Hyphenated: 8-4-4-4-12 (e.g., 550e8400-e29b-41d4-a716-446655440000)
//   - Unhyphenated: 32 hex chars (e.g., 550e8400e29b41d4a716446655440000)
//   - Braced: {8-4-4-4-12} or {32-hex}
//   - URN prefix: urn:uuid:8-4-4-4-12
func isValidUUID(s string) bool {
	if s == "" {
		return false
	}

	// Remove urn:uuid: prefix if present (case-insensitive)
	if strings.HasPrefix(strings.ToLower(s), "urn:uuid:") {
		s = s[9:]
	}

	// Remove braces if present
	if len(s) > 0 && s[0] == '{' && s[len(s)-1] == '}' {
		s = s[1 : len(s)-1]
	}

	// Check hyphenated format: 8-4-4-4-12
	hyphenatedPattern := `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`
	if matched, _ := regexp.MatchString(hyphenatedPattern, s); matched {
		return true
	}

	// Check unhyphenated format: 32 hex chars
	unhyphenatedPattern := `^[0-9a-fA-F]{32}$`
	if matched, _ := regexp.MatchString(unhyphenatedPattern, s); matched {
		return true
	}

	return false
}
