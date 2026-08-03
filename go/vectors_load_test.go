package layrz

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// vectorCase represents a single test case from the JSON vectors.
type vectorCase struct {
	Name           string          `json:"name"`
	Field          string          `json:"field"`
	Kwargs         map[string]any  `json:"kwargs"`
	Value          any             `json:"value,omitempty"`
	ValueAbsent    bool            `json:"value_absent"`
	ExpectedErrors []expectedError `json:"expected_errors"`
}

// expectedError represents an expected error from a test case.
type expectedError struct {
	Code     string `json:"code"`
	Expected any    `json:"expected,omitempty"`
	Received any    `json:"received,omitempty"`
}

// loadVectorFiles loads all vector JSON files from ../vectors/fields/*.json
// relative to the test file.
func loadVectorFiles(t *testing.T) map[string][]vectorCase {
	t.Helper()

	pattern := filepath.Join("..", "vectors", "fields", "*.json")
	files, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatalf("failed to glob vector files: %v", err)
	}

	if len(files) == 0 {
		t.Fatal("no vector files found at ../vectors/fields/*.json; test vectors are mandatory")
	}

	result := make(map[string][]vectorCase)

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("failed to read %q: %v", file, err)
		}

		var cases []vectorCase
		if err := json.Unmarshal(data, &cases); err != nil {
			t.Fatalf("failed to unmarshal %q: %v", file, err)
		}

		fieldName := filepath.Base(file[:len(file)-5]) // Remove .json
		result[fieldName] = cases
	}

	return result
}

// sameScalar reports whether a and b represent the same scalar value,
// accounting for type differences between JSON and Go (e.g., float64 vs int, []interface{} vs []string).
func sameScalar(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Check for numeric comparison first
	aFloat, aIsFloat := toFloat64(a)
	bFloat, bIsFloat := toFloat64(b)

	if aIsFloat && bIsFloat {
		// For floats, use exact comparison (matching JSON semantics)
		return aFloat == bFloat
	}

	// For slices, compare elements
	aSlice, aIsSlice := a.([]any)
	bSlice, bIsSlice := b.([]string)
	if aIsSlice && bIsSlice {
		if len(aSlice) != len(bSlice) {
			return false
		}
		for i := range aSlice {
			if aSlice[i] != bSlice[i] {
				return false
			}
		}
		return true
	}

	// Check reverse
	aSliceStr, aIsSliceStr := a.([]string)
	bSliceIface, bIsSliceIface := b.([]any)
	if aIsSliceStr && bIsSliceIface {
		if len(aSliceStr) != len(bSliceIface) {
			return false
		}
		for i := range aSliceStr {
			if aSliceStr[i] != bSliceIface[i] {
				return false
			}
		}
		return true
	}

	// Fall back to direct equality
	return a == b
}

// toFloat64 converts a numeric value to float64 if possible.
func toFloat64(v any) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	default:
		return 0, false
	}
}

// convertJSONValue converts a raw JSON value to the appropriate Go type based on the
// expected datatype and context.
// Rules for numbers:
// - If the number looks like an integer (no decimal point), it's an int (unless datatype=float, then it's invalid)
// - If the number has a decimal point, it's a float
// - If datatype=float is explicitly set and the number is a whole number, return as int so validator rejects it
func convertJSONValue(val any, expectedDatatype string) (any, error) {
	if val == nil {
		return nil, nil
	}

	// Handle different JSON types
	switch v := val.(type) {
	case bool:
		return v, nil
	case string:
		return v, nil
	case float64:
		// JSON numbers come in as float64
		// Check if this looks like an integer (no decimal point)
		str := fmt.Sprintf("%v", val)
		hasDecimal := strings.Contains(str, ".")
		isWholeNumber := v == float64(int64(v))

		if isWholeNumber && !hasDecimal {
			// It's an integer (e.g., 42 not 42.0)
			// If datatype=float is requested, return as int to trigger "invalid"
			if expectedDatatype == "float" {
				return int(int64(v)), nil
			}
			return int(int64(v)), nil
		}
		// It has a decimal point, so it's a float
		return v, nil
	case map[string]any:
		return v, nil
	case []any:
		return v, nil
	default:
		return v, nil
	}
}
