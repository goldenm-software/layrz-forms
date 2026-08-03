package layrz

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// configErrorKey is the reserved key for internal configuration errors.
const configErrorKey = "_config"

// cleanMethodInfo describes a clean method that was discovered and validated.
type cleanMethodInfo struct {
	name      string
	typ       reflect.Type
	value     reflect.Value
	fieldName string // non-empty only for Convention A
}

// discoverCleanMethods walks the pointer type's method set and classifies
// each clean method as Convention A (per-field) or Convention B (cross-field).
// Returns a tuple of (conventionA, conventionB, configErrors).
func discoverCleanMethods(ptrVal reflect.Value, structType reflect.Type) ([]cleanMethodInfo, []cleanMethodInfo, []*FieldError) {
	ptrType := ptrVal.Type()

	// Get the set of exported field names
	fieldNames := make(map[string]bool)
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if field.IsExported() {
			fieldNames[field.Name] = true
		}
	}

	var conventionA []cleanMethodInfo
	var conventionB []cleanMethodInfo
	var configErrors []*FieldError

	// Iterate over all methods on the pointer type
	for i := 0; i < ptrType.NumMethod(); i++ {
		m := ptrType.Method(i)

		// Skip if it doesn't start with "Clean" or is literally "Clean"
		if !strings.HasPrefix(m.Name, "Clean") || m.Name == "Clean" {
			continue
		}

		suffix := m.Name[len("Clean"):]
		mt := m.Type

		// Try to classify as Convention A first
		if fieldNames[suffix] {
			// Convention A: Clean<FieldName>(value <FieldType>) *FieldError
			if err := validateConventionA(mt, suffix, structType); err != nil {
				configErrors = append(configErrors, &FieldError{
					Code: "internalError",
					Extra: map[string]any{
						"reason":  err.Error(),
						"method":  m.Name,
						"context": "method signature validation",
					},
				})
				continue
			}

			conventionA = append(conventionA, cleanMethodInfo{
				name:      m.Name,
				typ:       mt,
				value:     m.Func,
				fieldName: suffix,
			})
		} else {
			// Convention B: Clean<Suffix>() Errors
			if err := validateConventionB(mt); err != nil {
				configErrors = append(configErrors, &FieldError{
					Code: "internalError",
					Extra: map[string]any{
						"reason":  err.Error(),
						"method":  m.Name,
						"context": "method signature validation",
					},
				})
				continue
			}

			conventionB = append(conventionB, cleanMethodInfo{
				name:  m.Name,
				typ:   mt,
				value: m.Func,
			})
		}
	}

	// Sort both slices by method name for determinism
	sort.Slice(conventionA, func(i, j int) bool {
		return conventionA[i].name < conventionA[j].name
	})
	sort.Slice(conventionB, func(i, j int) bool {
		return conventionB[i].name < conventionB[j].name
	})

	return conventionA, conventionB, configErrors
}

// validateConventionA checks that the method signature matches Convention A.
// Signature: func (f *T) CleanFieldName(value <FieldType>) *FieldError
func validateConventionA(mt reflect.Type, fieldName string, structType reflect.Type) error {
	// Must have exactly 2 parameters: receiver + value
	if mt.NumIn() != 2 {
		return fmt.Errorf("wrong arity: expected 2 params (receiver + value), got %d", mt.NumIn())
	}

	// Must have exactly 1 return value
	if mt.NumOut() != 1 {
		return fmt.Errorf("wrong return count: expected 1, got %d", mt.NumOut())
	}

	// Return type must be *FieldError
	if mt.Out(0) != reflect.TypeOf((*FieldError)(nil)) {
		return fmt.Errorf("return type must be *FieldError, got %v", mt.Out(0))
	}

	// Parameter 1 (index 1, since 0 is receiver) must match the field's type exactly
	field, ok := structType.FieldByName(fieldName)
	if !ok {
		return fmt.Errorf("field %q not found", fieldName)
	}

	expectedParamType := field.Type
	actualParamType := mt.In(1)
	if actualParamType != expectedParamType {
		return fmt.Errorf("parameter type mismatch: expected %v, got %v", expectedParamType, actualParamType)
	}

	return nil
}

// validateConventionB checks that the method signature matches Convention B.
// Signature: func (f *T) CleanSuffix() Errors
func validateConventionB(mt reflect.Type) error {
	// Must have exactly 1 parameter: receiver only
	if mt.NumIn() != 1 {
		return fmt.Errorf("wrong arity: expected 1 param (receiver), got %d", mt.NumIn())
	}

	// Must have exactly 1 return value
	if mt.NumOut() != 1 {
		return fmt.Errorf("wrong return count: expected 1, got %d", mt.NumOut())
	}

	// Return type must be Errors
	if mt.Out(0) != reflect.TypeOf(Errors(nil)) {
		return fmt.Errorf("return type must be Errors, got %v", mt.Out(0))
	}

	return nil
}

// callClean invokes a clean method, catching panics and returning them as errors.
func callClean(fn reflect.Value, args []reflect.Value) (out []reflect.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic in clean method: %v", r)
		}
	}()
	out = fn.Call(args)
	return out, nil
}

// runConventionAClean calls all Convention A methods and files results into out.
func runConventionAClean(ptrVal reflect.Value, methods []cleanMethodInfo, out Errors) {
	runConventionACleanWithPrefix(ptrVal, "", methods, out)
}

// runConventionACleanWithPrefix calls all Convention A methods and files results into out,
// prefixing all produced keys with the given prefix.
func runConventionACleanWithPrefix(ptrVal reflect.Value, prefix string, methods []cleanMethodInfo, out Errors) {
	for _, info := range methods {
		// Get the field value
		field := ptrVal.Elem().FieldByName(info.fieldName)

		// Call the method
		result, err := callClean(info.value, []reflect.Value{ptrVal, field})
		if err != nil {
			// Add config error
			out.Add(configErrorKey, &FieldError{
				Code: "internalError",
				Extra: map[string]any{
					"reason":  err.Error(),
					"method":  info.name,
					"context": "method execution",
				},
			})
			continue
		}

		// Extract the returned *FieldError
		if len(result) != 1 {
			out.Add(configErrorKey, &FieldError{
				Code: "internalError",
				Extra: map[string]any{
					"reason": "unexpected return value count",
					"method": info.name,
				},
			})
			continue
		}

		returnedErr := result[0].Interface().(*FieldError)
		if returnedErr != nil {
			// File under camelCased field name, with prefix
			key := ToCamelCase(info.fieldName)
			prefixedKey := joinKey(prefix, key)
			out.Add(prefixedKey, returnedErr)
		}
	}
}

// runConventionBClean calls all Convention B methods and merges results into out.
func runConventionBClean(ptrVal reflect.Value, methods []cleanMethodInfo, out Errors) {
	runConventionBCleanWithPrefix(ptrVal, "", methods, out)
}

// runConventionBCleanWithPrefix calls all Convention B methods and merges results into out,
// prefixing all produced keys with the given prefix.
func runConventionBCleanWithPrefix(ptrVal reflect.Value, prefix string, methods []cleanMethodInfo, out Errors) {
	for _, info := range methods {
		// Call the method
		result, err := callClean(info.value, []reflect.Value{ptrVal})
		if err != nil {
			// Add config error
			out.Add(configErrorKey, &FieldError{
				Code: "internalError",
				Extra: map[string]any{
					"reason":  err.Error(),
					"method":  info.name,
					"context": "method execution",
				},
			})
			continue
		}

		// Extract the returned Errors
		if len(result) != 1 {
			out.Add(configErrorKey, &FieldError{
				Code: "internalError",
				Extra: map[string]any{
					"reason": "unexpected return value count",
					"method": info.name,
				},
			})
			continue
		}

		returnedErrors := result[0].Interface().(Errors)
		if returnedErrors != nil && !returnedErrors.IsEmpty() {
			// Prefix the keys if needed
			if prefix == "" {
				out.Merge(returnedErrors)
			} else {
				for key, errs := range returnedErrors {
					camelKey := ToCamelCase(key)
					prefixedKey := prefix + "." + camelKey
					out.Add(prefixedKey, errs...)
				}
			}
		}
	}
}
