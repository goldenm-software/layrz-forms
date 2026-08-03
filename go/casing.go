package layrz

import (
	"strings"
	"unicode"
)

// ToCamelCase converts a snake_case key (with optional dot notation) to camelCase.
// For example:
//   - "user_name" → "userName"
//   - "a_b.c_d" → "aB.cD"
//   - "range_text_test" → "rangeTextTest"
//   - "address.street_name" → "address.streetName"
//   - "" → ""
//   - "_foo" → "foo" (leading underscore becomes part of the first word)
//   - "foo_" → "foo" (trailing underscore is dropped)
//
// Each dot-separated segment is processed independently, with underscores
// converted to camelCase within that segment.
func ToCamelCase(key string) string {
	if key == "" {
		return ""
	}

	segments := strings.Split(key, ".")
	result := make([]string, 0, len(segments))

	for _, segment := range segments {
		parts := strings.Split(segment, "_")
		if len(parts) == 0 {
			result = append(result, "")
			continue
		}

		// First part stays lowercase
		init := parts[0]

		// Remaining parts are title-cased (uppercase first letter, lowercase rest)
		var camel string
		for i, part := range parts {
			if i == 0 {
				camel = init
			} else {
				camel += pythonTitle(part)
			}
		}

		// Lowercase the first character of the final camelCase word
		if camel == "" {
			result = append(result, "")
		} else {
			result = append(result, string(unicode.ToLower(rune(camel[0])))+camel[1:])
		}
	}

	return strings.Join(result, ".")
}

// pythonTitle replicates Python's str.title() behavior:
// - Uppercase a letter if the previous char is not a letter
// - Lowercase the first letter
// - Lowercase all remaining letters
//
// Examples:
//   - "foo" → "Foo"
//   - "bC" → "Bc" (the C is lowercased)
//   - "a1b" → "A1B" (both letters uppercased, digit treated as non-letter separator)
func pythonTitle(s string) string {
	if s == "" {
		return ""
	}

	result := []rune(s)
	for i, ch := range result {
		if i == 0 {
			// First character: always uppercase if it's a letter
			result[i] = unicode.ToUpper(ch)
		} else {
			// Check if previous character is a letter
			prevChar := result[i-1]
			if !unicode.IsLetter(prevChar) {
				// Previous is not a letter: uppercase this letter
				result[i] = unicode.ToUpper(ch)
			} else {
				// Previous is a letter: lowercase this letter
				result[i] = unicode.ToLower(ch)
			}
		}
	}

	return string(result)
}
