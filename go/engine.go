package layrz

import (
	"fmt"
	"reflect"
	"strconv"
)

const maxDepth = 32

// validateStruct walks ptrVal (a non-nil pointer to a struct) and accumulates
// errors into out, prefixing every key with prefix (empty at top level).
// Recursion depth is guarded; exceeding maxDepth adds a config error.
func validateStruct(ptrVal reflect.Value, prefix string, out Errors, depth int) {
	// Guard against infinite recursion
	if depth > maxDepth {
		out.Add(configErrorKey, &FieldError{
			Code: "internalError",
			Extra: map[string]any{
				"reason": fmt.Sprintf("maximum recursion depth (%d) exceeded", maxDepth),
			},
		})
		return
	}

	structType := ptrVal.Elem().Type()

	// Iterate over each field
	for i := 0; i < structType.NumField(); i++ {
		fieldDef := structType.Field(i)

		// Skip unexported fields (can't read them)
		if !fieldDef.IsExported() {
			continue
		}

		// Skip embedded/anonymous structs for now; we'll handle them below
		// For now, only process fields with tags
		tag := fieldDef.Tag.Get("layrz")

		// Handle anonymous embedded structs: if the field is unnamed and is a struct,
		// walk its fields inline with the SAME prefix (promoted fields).
		// Only do this if the field is NOT tagged (tagged anonymous structs are errors).
		if fieldDef.Anonymous && tag == "" {
			if fieldDef.Type.Kind() == reflect.Struct {
				// This is an embedded struct without a tag; promote its fields
				fieldVal := ptrVal.Elem().Field(i)
				if fieldVal.Kind() == reflect.Ptr {
					if !fieldVal.IsNil() {
						validateStruct(fieldVal, prefix, out, depth+1)
					}
				} else {
					// Embedded struct value; take its address
					validateStruct(fieldVal.Addr(), prefix, out, depth+1)
				}
				continue
			}
		}

		// Parse the tag
		spec, err := ParseTag(tag)
		if err != nil {
			out.Add(configErrorKey, &FieldError{
				Code: "internalError",
				Extra: map[string]any{
					"reason": err.Error(),
					"field":  fieldDef.Name,
				},
			})
			continue
		}

		// Skip if tag says to skip
		if spec == nil {
			continue
		}

		fieldVal := ptrVal.Elem().Field(i)

		// Check type compatibility
		if err := spec.CheckType(fieldDef.Type); err != nil {
			out.Add(configErrorKey, &FieldError{
				Code: "internalError",
				Extra: map[string]any{
					"reason": err.Error(),
					"field":  fieldDef.Name,
				},
			})
			continue
		}

		key := joinKey(prefix, ToCamelCase(fieldDef.Name))

		// Handle by kind
		switch spec.Kind {
		case KindID, KindEmail, KindUUID, KindChar, KindNumber, KindBool, KindJSON:
			// Scalar field validator
			validateScalarField(spec, fieldVal, key, out)

		case KindSubform:
			// Subform: must be non-nil pointer to struct
			if fieldVal.Kind() == reflect.Ptr && !fieldVal.IsNil() {
				validateStruct(fieldVal, key, out, depth+1)
				// Also run clean methods on the subform
				subStructType := fieldVal.Type().Elem()
				conventionA, conventionB, configErrors := discoverCleanMethods(fieldVal, subStructType)
				for _, cfgErr := range configErrors {
					out.Add(key, cfgErr)
				}
				runConventionACleanWithPrefix(fieldVal, key, conventionA, out)
				runConventionBCleanWithPrefix(fieldVal, key, conventionB, out)
			}
			// nil subform is skipped (deliberate Go divergence from Python)

		case KindSubformList:
			// Subform list: slice of struct or slice of pointer-to-struct
			if fieldVal.Kind() == reflect.Slice || fieldVal.Kind() == reflect.Array {
				for j := 0; j < fieldVal.Len(); j++ {
					elemVal := fieldVal.Index(j)

					// Handle pointer elements
					if elemVal.Kind() == reflect.Ptr {
						if elemVal.IsNil() {
							continue
						}
						// For pointer elements, use them directly
						validateStruct(elemVal, key+"."+strconv.Itoa(j), out, depth+1)
					} else {
						// For non-pointer elements, take the address
						validateStruct(elemVal.Addr(), key+"."+strconv.Itoa(j), out, depth+1)
					}
				}
			}
		}
	}
}

// validateScalarField dispatches to the appropriate validator based on spec.Kind.
func validateScalarField(spec *FieldSpec, fieldVal reflect.Value, key string, out Errors) {
	rules, err := spec.Rules()
	if err != nil {
		out.Add(configErrorKey, &FieldError{
			Code: "internalError",
			Extra: map[string]any{
				"reason": err.Error(),
			},
		})
		return
	}

	var errs []*FieldError

	switch spec.Kind {
	case KindID:
		r := rules.(IDRules)
		// Unwrap pointer before passing to validator
		var val any
		if fieldVal.Kind() == reflect.Ptr {
			if !fieldVal.IsNil() {
				val = fieldVal.Elem().Interface()
			}
		} else {
			val = fieldVal.Interface()
		}
		errs = ValidateID(val, r)

	case KindEmail:
		// Type assert to *string
		var val *string
		if fieldVal.Kind() == reflect.Ptr {
			if fieldVal.IsNil() {
				val = nil
			} else {
				ptrVal := fieldVal.Interface().(*string)
				val = ptrVal
			}
		}
		r := rules.(EmailRules)
		errs = ValidateEmail(val, r)

	case KindUUID:
		// Type assert to *string
		var val *string
		if fieldVal.Kind() == reflect.Ptr {
			if fieldVal.IsNil() {
				val = nil
			} else {
				ptrVal := fieldVal.Interface().(*string)
				val = ptrVal
			}
		}
		r := rules.(UUIDRules)
		errs = ValidateUUID(val, r)

	case KindChar:
		// Type assert to *string
		var val *string
		if fieldVal.Kind() == reflect.Ptr {
			if fieldVal.IsNil() {
				val = nil
			} else {
				ptrVal := fieldVal.Interface().(*string)
				val = ptrVal
			}
		}
		r := rules.(CharRules)
		errs = ValidateChar(val, r)

	case KindNumber:
		r := rules.(NumberRules)
		// Unwrap pointer before passing to validator
		var val any
		if fieldVal.Kind() == reflect.Ptr {
			if !fieldVal.IsNil() {
				val = fieldVal.Elem().Interface()
			}
		} else {
			val = fieldVal.Interface()
		}
		errs = ValidateNumber(val, r)

	case KindBool:
		// Type assert to *bool
		var val *bool
		if fieldVal.Kind() == reflect.Ptr {
			if fieldVal.IsNil() {
				val = nil
			} else {
				ptrVal := fieldVal.Interface().(*bool)
				val = ptrVal
			}
		}
		r := rules.(BoolRules)
		errs = ValidateBool(val, r)

	case KindJSON:
		r := rules.(JSONRules)
		// Unwrap pointer before passing to validator
		var val any
		if fieldVal.Kind() == reflect.Ptr {
			if !fieldVal.IsNil() {
				val = fieldVal.Elem().Interface()
			}
		} else {
			val = fieldVal.Interface()
		}
		errs = ValidateJSON(val, r)
	}

	out.Add(key, errs...)
}

// joinKey combines a prefix and a key with a dot, or returns key if prefix is empty.
func joinKey(prefix, key string) string {
	if prefix == "" {
		return key
	}
	return prefix + "." + key
}

// validateForm validates the tag rules on a form struct, discovering and calling
// clean methods on the top-level form only.
// Returns a new Errors map with all accumulated errors.
func validateForm(form any) Errors {
	out := make(Errors)

	ptrVal := reflect.ValueOf(form)

	// Input validation
	if form == nil {
		out.Add(configErrorKey, &FieldError{
			Code: "internalError",
			Extra: map[string]any{
				"reason": "form is nil",
			},
		})
		return out
	}

	if ptrVal.Kind() != reflect.Ptr {
		out.Add(configErrorKey, &FieldError{
			Code: "internalError",
			Extra: map[string]any{
				"reason": "form must be a pointer (pointer receiver for clean methods and to modify struct fields)",
			},
		})
		return out
	}

	if ptrVal.IsNil() {
		out.Add(configErrorKey, &FieldError{
			Code: "internalError",
			Extra: map[string]any{
				"reason": "form pointer is nil",
			},
		})
		return out
	}

	elemType := ptrVal.Type().Elem()
	if elemType.Kind() != reflect.Struct {
		out.Add(configErrorKey, &FieldError{
			Code: "internalError",
			Extra: map[string]any{
				"reason": "form must point to a struct",
			},
		})
		return out
	}

	// Phase 1: validate built-in tag rules
	validateStruct(ptrVal, "", out, 0)

	// Discover clean methods on the top-level form
	structType := elemType
	conventionA, conventionB, configErrors := discoverCleanMethods(ptrVal, structType)

	// File any config errors from discovery
	for _, cfgErr := range configErrors {
		out.Add(configErrorKey, cfgErr)
	}

	// Phase 2: run Convention A clean methods
	runConventionAClean(ptrVal, conventionA, out)

	// Phase 3: run Convention B clean methods
	runConventionBClean(ptrVal, conventionB, out)

	return out
}
