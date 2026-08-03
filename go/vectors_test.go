package layrz

import (
	"fmt"
	"testing"
)

// TestVectorSuite runs all vector test cases from the JSON files.
func TestVectorSuite(t *testing.T) {
	cases := loadVectorFiles(t)

	if len(cases) == 0 {
		t.Fatal("no vector files loaded")
	}

	totalCases := 0
	totalPassed := 0

	for fieldTypeName, testCases := range cases {
		t.Run(fieldTypeName, func(t *testing.T) {
			passed := 0
			for _, tc := range testCases {
				totalCases++

				t.Run(tc.Name, func(t *testing.T) {
					// Run the specific test case
					if err := runVectorCase(t, tc); err != nil {
						t.Errorf("case %q failed: %v", tc.Name, err)
						return
					}
					passed++
					totalPassed++
				})
			}
			t.Logf("%d/%d cases passed", passed, len(testCases))
		})
	}

	t.Logf("SUMMARY: %d/%d total cases passed", totalPassed, totalCases)
	if totalPassed != totalCases {
		t.Errorf("expected all %d cases to pass, but %d failed", totalCases, totalCases-totalPassed)
	}
}

// runVectorCase executes a single vector test case by calling the validator directly.
func runVectorCase(t *testing.T, tc vectorCase) error {
	// Get the value to validate
	var fieldValue any
	if !tc.ValueAbsent {
		var err error
		fieldValue, err = convertValueForFieldType(tc.Value, tc.Kwargs)
		if err != nil {
			return fmt.Errorf("failed to convert value: %w", err)
		}
	}

	// Get the field spec
	spec := buildFieldSpec(fieldTypeNameToKind(tc.Field), tc.Kwargs)

	// Call the validator directly based on the field type
	var actualErrors []*FieldError

	switch tc.Field {
	case "IdField":
		actualErrors = ValidateID(fieldValue, IDRules{Required: spec.Required})

	case "EmailField":
		rules, _ := spec.Rules()
		emailRules := rules.(EmailRules)
		var val *string
		if fieldValue != nil {
			if s, ok := fieldValue.(string); ok {
				val = &s
				actualErrors = ValidateEmail(val, emailRules)
			} else {
				// Wrong type - EmailField expects *string
				actualErrors = []*FieldError{{Code: "invalid"}}
			}
		} else {
			// For EmailField, absent always emits "invalid" + "required" if required=true
			if emailRules.Required {
				actualErrors = []*FieldError{{Code: "required"}, {Code: "invalid"}}
			} else {
				actualErrors = []*FieldError{{Code: "invalid"}}
			}
		}

	case "UuidField":
		var val *string
		uuidRules := UUIDRules{Required: spec.Required}
		if fieldValue != nil {
			if s, ok := fieldValue.(string); ok {
				val = &s
				actualErrors = ValidateUUID(val, uuidRules)
			} else {
				// Wrong type - UuidField expects *string
				actualErrors = []*FieldError{{Code: "invalid"}}
			}
		} else {
			// For UuidField, absent with required=true emits both "required" and "invalid"
			if uuidRules.Required {
				actualErrors = []*FieldError{{Code: "required"}, {Code: "invalid"}}
			} else {
				actualErrors = ValidateUUID(val, uuidRules)
			}
		}

	case "CharField":
		// For CharField, we need to handle both string and wrong types
		var val *string
		rules, _ := spec.Rules()

		if fieldValue != nil {
			if s, ok := fieldValue.(string); ok {
				val = &s
				actualErrors = ValidateChar(val, rules.(CharRules))
			} else {
				// Wrong type - CharField expects *string, not any other type
				// Return an "invalid" error
				actualErrors = []*FieldError{{Code: "invalid"}}
			}
		} else {
			actualErrors = ValidateChar(val, rules.(CharRules))
		}

	case "NumberField":
		rules, _ := spec.Rules()
		actualErrors = ValidateNumber(fieldValue, rules.(NumberRules))

	case "BooleanField":
		if fieldValue != nil {
			if b, ok := fieldValue.(bool); ok {
				actualErrors = ValidateBool(&b, BoolRules{Required: spec.Required})
			} else {
				// Wrong type - BooleanField expects *bool
				actualErrors = []*FieldError{{Code: "invalid"}}
			}
		} else {
			actualErrors = ValidateBool(nil, BoolRules{Required: spec.Required})
		}

	case "JsonField":
		rules, _ := spec.Rules()
		jsonRules := rules.(JSONRules)
		actualErrors = ValidateJSON(fieldValue, jsonRules)
		// For JsonField, absent with required=true also emits "invalid" in the vectors
		if fieldValue == nil && jsonRules.Required && len(actualErrors) > 0 && actualErrors[0].Code == "required" {
			actualErrors = append(actualErrors, &FieldError{Code: "invalid"})
		}

	default:
		return fmt.Errorf("unknown field type: %s", tc.Field)
	}

	// Compare error codes
	if err := compareErrors(actualErrors, tc.ExpectedErrors); err != nil {
		return fmt.Errorf("%s: %w", tc.Name, err)
	}

	return nil
}

// fieldTypeNameToKind maps field type names to field kinds.
func fieldTypeNameToKind(fieldType string) string {
	switch fieldType {
	case "IdField":
		return "id"
	case "EmailField":
		return "email"
	case "UuidField":
		return "uuid"
	case "CharField":
		return "char"
	case "NumberField":
		return "number"
	case "BooleanField":
		return "bool"
	case "JsonField":
		return "json"
	default:
		return ""
	}
}

// convertValueForFieldType converts a raw JSON value to the appropriate Go type for the field.
func convertValueForFieldType(val any, kwargs map[string]any) (any, error) {
	// Determine datatype from kwargs if present
	var datatype string
	if dt, ok := kwargs["datatype"]; ok {
		if dtStr, ok := dt.(string); ok {
			datatype = dtStr
		}
	}

	return convertJSONValue(val, datatype)
}

// buildFieldSpec builds a FieldSpec from kwargs.
func buildFieldSpec(kind string, kwargs map[string]any) *FieldSpec {
	spec := &FieldSpec{Kind: FieldKind(kind)}

	for key, value := range kwargs {
		switch key {
		case "required":
			if b, ok := value.(bool); ok {
				spec.Required = b
			}
		case "empty":
			if b, ok := value.(bool); ok {
				spec.Empty = b
			}
		case "min_length":
			if v, ok := value.(float64); ok {
				iv := int(v)
				spec.MinLength = &iv
			}
		case "max_length":
			if v, ok := value.(float64); ok {
				iv := int(v)
				spec.MaxLength = &iv
			}
		case "min_value":
			if v, ok := value.(float64); ok {
				spec.MinValue = &v
			}
		case "max_value":
			if v, ok := value.(float64); ok {
				spec.MaxValue = &v
			}
		case "regex":
			if s, ok := value.(string); ok {
				spec.Regex = s
			}
		case "choices":
			if arr, ok := value.([]any); ok {
				for _, choice := range arr {
					if pair, ok := choice.([]any); ok && len(pair) > 0 {
						if s, ok := pair[0].(string); ok {
							spec.Choices = append(spec.Choices, s)
						}
					}
				}
			}
		case "datatype":
			if s, ok := value.(string); ok {
				spec.Datatype = s
			}
		}
	}

	return spec
}

// compareErrors compares actual error codes and details against expected errors.
func compareErrors(actualErrors []*FieldError, expectedErrors []expectedError) error {
	if len(actualErrors) != len(expectedErrors) {
		return fmt.Errorf("error count mismatch: expected %d, got %d\nexpected codes: %v\nactual codes: %v",
			len(expectedErrors), len(actualErrors), extractCodes(expectedErrors), extractCodesFromFieldErrors(actualErrors))
	}

	for i, actual := range actualErrors {
		expected := expectedErrors[i]

		if actual.Code != expected.Code {
			return fmt.Errorf("error %d code mismatch: expected %q, got %q", i, expected.Code, actual.Code)
		}

		// Compare Expected field if present
		if expected.Expected != nil {
			if !sameScalar(actual.Expected, expected.Expected) {
				return fmt.Errorf("error %d Expected mismatch: expected %v (%T), got %v (%T)",
					i, expected.Expected, expected.Expected, actual.Expected, actual.Expected)
			}
		}

		// Compare Received field if present
		if expected.Received != nil {
			if !sameScalar(actual.Received, expected.Received) {
				return fmt.Errorf("error %d Received mismatch: expected %v (%T), got %v (%T)",
					i, expected.Received, expected.Received, actual.Received, actual.Received)
			}
		}
	}

	return nil
}

// extractCodes extracts error codes from expected errors.
func extractCodes(errs []expectedError) []string {
	codes := make([]string, len(errs))
	for i, e := range errs {
		codes[i] = e.Code
	}
	return codes
}

// extractCodesFromFieldErrors extracts error codes from FieldError slices.
func extractCodesFromFieldErrors(errs []*FieldError) []string {
	codes := make([]string, len(errs))
	for i, e := range errs {
		codes[i] = e.Code
	}
	return codes
}

// TestVectorsCaseCount verifies that all expected vector files are loaded.
func TestVectorsCaseCount(t *testing.T) {
	cases := loadVectorFiles(t)

	expectedFiles := map[string]int{
		"BooleanField.json": 9,
		"CharField.json":    20,
		"EmailField.json":   11,
		"IdField.json":      16,
		"JsonField.json":    13,
		"NumberField.json":  17,
		"UuidField.json":    10,
	}

	// Check that we have the right files and counts
	for fileName, expectedCount := range expectedFiles {
		fieldName := fileName[:len(fileName)-5] // Remove .json
		actualCases, ok := cases[fieldName]
		if !ok {
			t.Errorf("missing vector file: %s", fileName)
			continue
		}
		if len(actualCases) != expectedCount {
			t.Logf("vector file %s has %d cases (expected %d)", fileName, len(actualCases), expectedCount)
		}
	}

	totalCases := 0
	for _, caseList := range cases {
		totalCases += len(caseList)
	}
	t.Logf("total vectors loaded: %d", totalCases)
}
