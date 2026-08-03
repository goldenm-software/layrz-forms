package layrz

import (
	"fmt"
	"reflect"
)

// Datatype constants for tag parsing and validation.
const (
	datatypeInt   = "int"
	datatypeFloat = "float"
	datatypeList  = "list"
	datatypeDict  = "dict"
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

// checkTypeID validates and infers types for KindID.
func (s *FieldSpec) checkTypeID(elemType reflect.Type, t reflect.Type) error {
	if elemType.Kind() != reflect.Int && elemType.Kind() != reflect.Int64 && elemType.Kind() != reflect.String {
		return fmt.Errorf("layrz: kind \"id\" requires an integer or string (or pointer to one), got %v", t)
	}
	return nil
}

// checkTypeString validates types for string-based kinds (email, uuid, char).
func (s *FieldSpec) checkTypeString(elemType reflect.Type, t reflect.Type) error {
	if elemType.Kind() != reflect.String {
		return fmt.Errorf("layrz: kind %q requires a string or *string, got %v", s.Kind, t)
	}
	return nil
}

// checkTypeNumber validates and infers types for KindNumber.
func (s *FieldSpec) checkTypeNumber(elemType reflect.Type, t reflect.Type) error {
	isNumeric := elemType.Kind() == reflect.Int || elemType.Kind() == reflect.Int8 || elemType.Kind() == reflect.Int16 ||
		elemType.Kind() == reflect.Int32 || elemType.Kind() == reflect.Int64 ||
		elemType.Kind() == reflect.Float32 || elemType.Kind() == reflect.Float64

	if !isNumeric {
		return fmt.Errorf("layrz: kind \"number\" requires an int or float type (or pointer to one), got %v", t)
	}

	// Infer Datatype from the element type if not explicitly set
	if s.Datatype == "" {
		switch elemType.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			s.Datatype = datatypeInt
		case reflect.Float32, reflect.Float64:
			s.Datatype = datatypeFloat
		}
	} else {
		// Check that explicit Datatype matches the Go type
		isInt := elemType.Kind() == reflect.Int || elemType.Kind() == reflect.Int8 ||
			elemType.Kind() == reflect.Int16 || elemType.Kind() == reflect.Int32 || elemType.Kind() == reflect.Int64
		isFloat := elemType.Kind() == reflect.Float32 || elemType.Kind() == reflect.Float64

		switch s.Datatype {
		case datatypeInt:
			if !isInt {
				return fmt.Errorf("layrz: datatype \"int\" incompatible with %v", t)
			}
		case datatypeFloat:
			if !isFloat {
				return fmt.Errorf("layrz: datatype \"float\" incompatible with %v", t)
			}
		}
	}
	return nil
}

// checkTypeJSON validates and infers types for KindJSON.
func (s *FieldSpec) checkTypeJSON(elemType reflect.Type, t reflect.Type) error {
	switch elemType.Kind() {
	case reflect.Slice, reflect.Array:
		// It's a slice/array, infer as "list"
		if s.Datatype == "" {
			s.Datatype = datatypeList
		} else if s.Datatype != datatypeList {
			return fmt.Errorf("layrz: datatype %q incompatible with slice/array type %v", s.Datatype, t)
		}

	case reflect.Map:
		// It's a map, infer as "dict"
		if s.Datatype == "" {
			s.Datatype = datatypeDict
		} else if s.Datatype != datatypeDict {
			return fmt.Errorf("layrz: datatype %q incompatible with map type %v", s.Datatype, t)
		}

	default:
		return fmt.Errorf("layrz: kind \"json\" requires a slice/array or map type (or pointer to one), got %v", t)
	}
	return nil
}

// checkTypeSubform validates types for KindSubform.
func (s *FieldSpec) checkTypeSubform(t reflect.Type) error {
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("layrz: kind \"subform\" requires a pointer to struct, got %v", t)
	}
	if t.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("layrz: kind \"subform\" requires a pointer to struct, got %v", t)
	}
	return nil
}

// checkTypeSubformList validates types for KindSubformList.
func (s *FieldSpec) checkTypeSubformList(t reflect.Type) error {
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
	return nil
}

// CheckType reports whether the spec's Kind is compatible with the Go field type,
// and infers Datatype from the type when the tag omitted it.
// For scalar kinds, t can be either a pointer (which models absence as nil) or a value field
// (which is always present). A non-pointer scalar field satisfies required trivially.
// t is the struct field's type.
func (s *FieldSpec) CheckType(t reflect.Type) error {
	// elemType is the underlying scalar type: t itself for a value field, or t.Elem() for a pointer.
	elemType := t
	if t.Kind() == reflect.Ptr {
		elemType = t.Elem()
	}

	switch s.Kind {
	case KindID:
		return s.checkTypeID(elemType, t)
	case KindEmail, KindUUID, KindChar:
		return s.checkTypeString(elemType, t)
	case KindNumber:
		return s.checkTypeNumber(elemType, t)
	case KindBool:
		if elemType.Kind() != reflect.Bool {
			return fmt.Errorf("layrz: kind \"bool\" requires a bool or *bool, got %v", t)
		}
	case KindJSON:
		return s.checkTypeJSON(elemType, t)
	case KindSubform:
		return s.checkTypeSubform(t)
	case KindSubformList:
		return s.checkTypeSubformList(t)
	default:
		return nil
	}

	return nil
}
