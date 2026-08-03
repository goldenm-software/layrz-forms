package layrz

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// FieldKind identifies which built-in validator a struct field uses.
type FieldKind string

const (
	KindID          FieldKind = "id"
	KindEmail       FieldKind = "email"
	KindUUID        FieldKind = "uuid"
	KindChar        FieldKind = "char"
	KindNumber      FieldKind = "number"
	KindBool        FieldKind = "bool"
	KindJSON        FieldKind = "json"
	KindSubform     FieldKind = "subform"
	KindSubformList FieldKind = "subform_list"
)

// FieldSpec is a parsed `layrz` struct tag.
type FieldSpec struct {
	Kind      FieldKind
	Required  bool
	Empty     bool
	MinLength *int
	MaxLength *int
	MinValue  *float64
	MaxValue  *float64
	Regex     string
	Choices   []string
	Datatype  string // "int"|"float"|"list"|"dict"|""
}

// ParseTag parses a `layrz` tag value. Returns (nil, nil) for "" or "-" (skip).
// Grammar: `fieldtype{,rule}` where rule is one of:
//   - bare `required`, bare `empty`
//   - `min_length=<int>`, `max_length=<int>`, `min_value=<float>`, `max_value=<float>`
//   - `regex=<string>` (must be last; consumes remainder of tag)
//   - `choices=<a|b|c>`
//   - `datatype=<int|float|list|dict>`
//
// Returns an error if the tag is malformed, the first token is not a valid FieldKind,
// a rule is duplicated, or a regex fails to compile.
func ParseTag(tag string) (*FieldSpec, error) {
	tag = strings.TrimSpace(tag)
	if tag == "" || tag == "-" {
		return nil, nil
	}

	// Split on comma, but regex= is special: it consumes the rest
	var tokens []string
	var regexValue string
	parts := strings.Split(tag, ",")

	for i, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "regex=") {
			// regex= must be the last rule; consume the rest
			regexValue = strings.Join(parts[i:], ",")
			regexValue = strings.TrimPrefix(regexValue, "regex=")
			break
		}
		tokens = append(tokens, part)
	}

	if len(tokens) == 0 {
		return nil, fmt.Errorf("layrz tag: empty or missing field kind")
	}

	// First token must be a valid FieldKind
	kindStr := tokens[0]
	var kind FieldKind
	switch FieldKind(kindStr) {
	case KindID, KindEmail, KindUUID, KindChar, KindNumber, KindBool, KindJSON, KindSubform, KindSubformList:
		kind = FieldKind(kindStr)
	default:
		return nil, fmt.Errorf("layrz tag: invalid field kind %q", kindStr)
	}

	spec := &FieldSpec{Kind: kind}
	seen := make(map[string]bool)

	// Parse rules (tokens[1:])
	for _, token := range tokens[1:] {
		if token == "" {
			continue
		}

		if token == "required" {
			if seen["required"] {
				return nil, fmt.Errorf("layrz tag: duplicate rule 'required'")
			}
			spec.Required = true
			seen["required"] = true
			continue
		}

		if token == "empty" {
			if seen["empty"] {
				return nil, fmt.Errorf("layrz tag: duplicate rule 'empty'")
			}
			spec.Empty = true
			seen["empty"] = true
			continue
		}

		// Parse key=value rules
		if idx := strings.Index(token, "="); idx != -1 {
			key := token[:idx]
			value := token[idx+1:]

			switch key {
			case "min_length":
				if seen["min_length"] {
					return nil, fmt.Errorf("layrz tag: duplicate rule 'min_length'")
				}
				v, err := strconv.Atoi(value)
				if err != nil {
					return nil, fmt.Errorf("layrz tag: min_length value must be an integer: %w", err)
				}
				spec.MinLength = &v
				seen["min_length"] = true

			case "max_length":
				if seen["max_length"] {
					return nil, fmt.Errorf("layrz tag: duplicate rule 'max_length'")
				}
				v, err := strconv.Atoi(value)
				if err != nil {
					return nil, fmt.Errorf("layrz tag: max_length value must be an integer: %w", err)
				}
				spec.MaxLength = &v
				seen["max_length"] = true

			case "min_value":
				if seen["min_value"] {
					return nil, fmt.Errorf("layrz tag: duplicate rule 'min_value'")
				}
				v, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return nil, fmt.Errorf("layrz tag: min_value must be a number: %w", err)
				}
				spec.MinValue = &v
				seen["min_value"] = true

			case "max_value":
				if seen["max_value"] {
					return nil, fmt.Errorf("layrz tag: duplicate rule 'max_value'")
				}
				v, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return nil, fmt.Errorf("layrz tag: max_value must be a number: %w", err)
				}
				spec.MaxValue = &v
				seen["max_value"] = true

			case "choices":
				if seen["choices"] {
					return nil, fmt.Errorf("layrz tag: duplicate rule 'choices'")
				}
				spec.Choices = strings.Split(value, "|")
				seen["choices"] = true

			case "datatype":
				if seen["datatype"] {
					return nil, fmt.Errorf("layrz tag: duplicate rule 'datatype'")
				}
				if value != "int" && value != "float" && value != "list" && value != "dict" {
					return nil, fmt.Errorf("layrz tag: datatype must be one of 'int', 'float', 'list', 'dict', got %q", value)
				}
				spec.Datatype = value
				seen["datatype"] = true

			default:
				return nil, fmt.Errorf("layrz tag: unknown rule %q", key)
			}
			continue
		}

		return nil, fmt.Errorf("layrz tag: unrecognized token %q", token)
	}

	// Handle regex= (must be last)
	if regexValue != "" {
		if seen["regex"] {
			return nil, fmt.Errorf("layrz tag: duplicate rule 'regex'")
		}
		// Validate that the regex compiles
		if _, err := regexp.Compile(regexValue); err != nil {
			return nil, fmt.Errorf("layrz tag: regex value is invalid: %w", err)
		}
		spec.Regex = regexValue
		seen["regex"] = true
	}

	return spec, nil
}
