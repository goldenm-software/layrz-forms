package layrz

import (
	"fmt"
	"reflect"
)

// IDRules holds validation rules for ID fields.
type IDRules struct {
	Required bool
}

// EmailRules holds validation rules for Email fields.
type EmailRules struct {
	Required bool
	Empty    bool
	Regex    string
}

// UUIDRules holds validation rules for UUID fields.
type UUIDRules struct {
	Required bool
}

// CharRules holds validation rules for Char fields.
type CharRules struct {
	Required  bool
	Empty     bool
	MinLength *int
	MaxLength *int
	Regex     string
	Choices   []string
}

// NumberRules holds validation rules for Number fields.
type NumberRules struct {
	Required bool
	Datatype string // "int"|"float"
	MinValue *float64
	MaxValue *float64
}

// BoolRules holds validation rules for Bool fields.
type BoolRules struct {
	Required bool
}

// JSONRules holds validation rules for JSON fields.
type JSONRules struct {
	Required bool
	Empty    bool
	Datatype string // "list"|"dict"
}

// DefaultEmailRegex is the email pattern used when EmailRules.Regex is empty.
const DefaultEmailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-z]{2,63}$`

// Rules converts a spec into the concrete rules value for its Kind.
// Returns the rules as `any` plus an error if the spec is incoherent.
func (s *FieldSpec) Rules() (any, error) {
	switch s.Kind {
	case KindID:
		return IDRules{Required: s.Required}, nil

	case KindEmail:
		regex := s.Regex
		if regex == "" {
			regex = DefaultEmailRegex
		}
		return EmailRules{
			Required: s.Required,
			Empty:    s.Empty,
			Regex:    regex,
		}, nil

	case KindUUID:
		return UUIDRules{Required: s.Required}, nil

	case KindChar:
		return CharRules{
			Required:  s.Required,
			Empty:     s.Empty,
			MinLength: s.MinLength,
			MaxLength: s.MaxLength,
			Regex:     s.Regex,
			Choices:   s.Choices,
		}, nil

	case KindNumber:
		return NumberRules{
			Required: s.Required,
			Datatype: s.Datatype,
			MinValue: s.MinValue,
			MaxValue: s.MaxValue,
		}, nil

	case KindBool:
		return BoolRules{Required: s.Required}, nil

	case KindJSON:
		return JSONRules{
			Required: s.Required,
			Empty:    s.Empty,
			Datatype: s.Datatype,
		}, nil

	default:
		return nil, fmt.Errorf("spec: unsupported field kind %q", s.Kind)
	}
}

// CheckType reports whether the spec's Kind is compatible with the Go field type,
// and infers Datatype from the type when the tag omitted it.
// For scalar kinds, t can be either a pointer (which models absence as nil) or a value field
// (which is always present). A non-pointer scalar field satisfies required trivially.
// t is the struct field's type.
func (s *FieldSpec) CheckType(t reflect.Type) error {
	// elemType is the underlying scalar type: t itself for a value field, or t.Elem() for a pointer.
	elemType := t
	isPointer := t.Kind() == reflect.Ptr
	if isPointer {
		elemType = t.Elem()
	}

	// Helper to check if elemType is one of the given kinds
	elemIsKind := func(kinds ...reflect.Kind) bool {
		for _, k := range kinds {
			if elemType.Kind() == k {
				return true
			}
		}
		return false
	}

	switch s.Kind {
	case KindID:
		// id: int, int64, string, *int, *int64, or *string
		if !elemIsKind(reflect.Int, reflect.Int64, reflect.String) {
			return fmt.Errorf("layrz: kind \"id\" requires an integer or string (or pointer to one), got %v", t)
		}

	case KindEmail, KindUUID, KindChar:
		// email/uuid/char: string or *string
		if elemType.Kind() != reflect.String {
			return fmt.Errorf("layrz: kind %q requires a string or *string, got %v", s.Kind, t)
		}

	case KindNumber:
		// number: int, int8, int16, int32, int64, float32, float64, or pointers to them
		if !elemIsKind(reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64) {
			return fmt.Errorf("layrz: kind \"number\" requires an int or float type (or pointer to one), got %v", t)
		}

		// Infer Datatype from the element type if not explicitly set
		if s.Datatype == "" {
			switch elemType.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				s.Datatype = "int"
			case reflect.Float32, reflect.Float64:
				s.Datatype = "float"
			}
		} else {
			// Check that explicit Datatype matches the Go type
			switch s.Datatype {
			case "int":
				if !elemIsKind(reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64) {
					return fmt.Errorf("layrz: datatype \"int\" incompatible with %v", t)
				}
			case "float":
				if !elemIsKind(reflect.Float32, reflect.Float64) {
					return fmt.Errorf("layrz: datatype \"float\" incompatible with %v", t)
				}
			}
		}

	case KindBool:
		// bool: bool or *bool
		if elemType.Kind() != reflect.Bool {
			return fmt.Errorf("layrz: kind \"bool\" requires a bool or *bool, got %v", t)
		}

	case KindJSON:
		// json: []any, map[string]any, *[]any, *map[string]any, or more general slice/map types (pointer or value)
		switch elemType.Kind() {
		case reflect.Slice, reflect.Array:
			// It's a slice/array, infer as "list"
			if s.Datatype == "" {
				s.Datatype = "list"
			} else if s.Datatype != "list" {
				return fmt.Errorf("layrz: datatype %q incompatible with slice/array type %v", s.Datatype, t)
			}

		case reflect.Map:
			// It's a map, infer as "dict"
			if s.Datatype == "" {
				s.Datatype = "dict"
			} else if s.Datatype != "dict" {
				return fmt.Errorf("layrz: datatype %q incompatible with map type %v", s.Datatype, t)
			}

		default:
			return fmt.Errorf("layrz: kind \"json\" requires a slice/array or map type (or pointer to one), got %v", t)
		}

	case KindSubform:
		// subform: pointer to struct (value struct not yet supported)
		if t.Kind() != reflect.Ptr {
			return fmt.Errorf("layrz: kind \"subform\" requires a pointer to struct, got %v", t)
		}
		if t.Elem().Kind() != reflect.Struct {
			return fmt.Errorf("layrz: kind \"subform\" requires a pointer to struct, got %v", t)
		}

	case KindSubformList:
		// subform_list: slice of struct or slice of pointer-to-struct
		if t.Kind() != reflect.Slice && t.Kind() != reflect.Array {
			return fmt.Errorf("layrz: kind \"subform_list\" requires a slice type, got %v", t)
		}

		elemType := t.Elem()
		if elemType.Kind() == reflect.Ptr {
			if elemType.Elem().Kind() != reflect.Struct {
				return fmt.Errorf("layrz: kind \"subform_list\" requires []struct or []*struct, got %v", t)
			}
		} else if elemType.Kind() != reflect.Struct {
			return fmt.Errorf("layrz: kind \"subform_list\" requires []struct or []*struct, got %v", t)
		}
	}

	return nil
}
