package layrz

import (
	"reflect"
	"testing"
)

// TestDerefNilInterface tests that deref(nil) returns (nil, false).
func TestDerefNilInterface(t *testing.T) {
	val, ok := deref(nil)
	if val != nil {
		t.Errorf("expected nil, got %v", val)
	}
	if ok {
		t.Errorf("expected false, got %v", ok)
	}
}

// TestDerefPointerTypes tests that deref correctly unwraps all supported pointer types
// and asserts both the dereferenced value and its dynamic type.
func TestDerefPointerTypes(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected any
		ok       bool
	}{
		// *int cases
		{
			name:     "non-nil *int",
			input:    Ptr(42),
			expected: 42,
			ok:       true,
		},
		{
			name:     "nil *int",
			input:    (*int)(nil),
			expected: nil,
			ok:       false,
		},

		// *int8 cases
		{
			name:     "non-nil *int8",
			input:    func() any { v := int8(8); return &v }(),
			expected: int8(8),
			ok:       true,
		},
		{
			name:     "nil *int8",
			input:    (*int8)(nil),
			expected: nil,
			ok:       false,
		},

		// *int16 cases
		{
			name:     "non-nil *int16",
			input:    func() any { v := int16(16); return &v }(),
			expected: int16(16),
			ok:       true,
		},
		{
			name:     "nil *int16",
			input:    (*int16)(nil),
			expected: nil,
			ok:       false,
		},

		// *int32 cases
		{
			name:     "non-nil *int32",
			input:    func() any { v := int32(32); return &v }(),
			expected: int32(32),
			ok:       true,
		},
		{
			name:     "nil *int32",
			input:    (*int32)(nil),
			expected: nil,
			ok:       false,
		},

		// *int64 cases
		{
			name:     "non-nil *int64",
			input:    func() any { v := int64(64); return &v }(),
			expected: int64(64),
			ok:       true,
		},
		{
			name:     "nil *int64",
			input:    (*int64)(nil),
			expected: nil,
			ok:       false,
		},

		// *float32 cases
		{
			name:     "non-nil *float32",
			input:    func() any { v := float32(3.14); return &v }(),
			expected: float32(3.14),
			ok:       true,
		},
		{
			name:     "nil *float32",
			input:    (*float32)(nil),
			expected: nil,
			ok:       false,
		},

		// *float64 cases
		{
			name:     "non-nil *float64",
			input:    Ptr(2.71),
			expected: 2.71,
			ok:       true,
		},
		{
			name:     "nil *float64",
			input:    (*float64)(nil),
			expected: nil,
			ok:       false,
		},

		// *string cases
		{
			name:     "non-nil *string",
			input:    Ptr("hello"),
			expected: "hello",
			ok:       true,
		},
		{
			name:     "nil *string",
			input:    (*string)(nil),
			expected: nil,
			ok:       false,
		},

		// *bool cases
		{
			name:     "non-nil *bool (true)",
			input:    Ptr(true),
			expected: true,
			ok:       true,
		},
		{
			name:     "non-nil *bool (false)",
			input:    Ptr(false),
			expected: false,
			ok:       true,
		},
		{
			name:     "nil *bool",
			input:    (*bool)(nil),
			expected: nil,
			ok:       false,
		},

		// *[]any cases
		{
			name:     "non-nil *[]any",
			input:    func() any { v := []any{1, "two", 3.0}; return &v }(),
			expected: []any{1, "two", 3.0},
			ok:       true,
		},
		{
			name:     "nil *[]any",
			input:    (*[]any)(nil),
			expected: nil,
			ok:       false,
		},

		// *map[string]any cases
		{
			name:     "non-nil *map[string]any",
			input:    func() any { v := map[string]any{"key": "value"}; return &v }(),
			expected: map[string]any{"key": "value"},
			ok:       true,
		},
		{
			name:     "nil *map[string]any",
			input:    (*map[string]any)(nil),
			expected: nil,
			ok:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := deref(tt.input)

			if ok != tt.ok {
				t.Errorf("expected ok=%v, got %v", tt.ok, ok)
			}

			if !ok && tt.expected == nil {
				if result != nil {
					t.Errorf("expected nil for ok=false, got %v", result)
				}
				return
			}

			if !ok {
				return
			}

			// Assert the dynamic type is the element type, not the pointer type
			resultType := reflect.TypeOf(result)
			expectedType := reflect.TypeOf(tt.expected)
			if resultType != expectedType {
				t.Errorf("expected type %v, got %v (pointer was not dereferenced)", expectedType, resultType)
			}

			// For non-nil results, compare values
			// Special handling for slices and maps which can't be compared with ==
			switch expectedType.Kind() {
			case reflect.Slice:
				resultSlice := reflect.ValueOf(result)
				expectedSlice := reflect.ValueOf(tt.expected)
				if resultSlice.Len() != expectedSlice.Len() {
					t.Errorf("expected slice len %d, got %d", expectedSlice.Len(), resultSlice.Len())
				} else {
					for i := 0; i < resultSlice.Len(); i++ {
						if !reflect.DeepEqual(resultSlice.Index(i).Interface(), expectedSlice.Index(i).Interface()) {
							t.Errorf("expected slice[%d]=%v, got %v", i, expectedSlice.Index(i).Interface(), resultSlice.Index(i).Interface())
						}
					}
				}
			case reflect.Map:
				resultMap := reflect.ValueOf(result)
				expectedMap := reflect.ValueOf(tt.expected)
				if resultMap.Len() != expectedMap.Len() {
					t.Errorf("expected map len %d, got %d", expectedMap.Len(), resultMap.Len())
				} else {
					for _, key := range expectedMap.MapKeys() {
						resultVal := resultMap.MapIndex(key)
						expectedVal := expectedMap.MapIndex(key)
						if !resultVal.IsValid() {
							t.Errorf("expected map key %v, not found in result", key.Interface())
						} else if !reflect.DeepEqual(resultVal.Interface(), expectedVal.Interface()) {
							t.Errorf("expected map[%v]=%v, got %v", key.Interface(), expectedVal.Interface(), resultVal.Interface())
						}
					}
				}
			default:
				if result != tt.expected {
					t.Errorf("expected value %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

// TestDerefNonPointerTypes tests that non-pointer values are returned unchanged.
func TestDerefNonPointerTypes(t *testing.T) {
	tests := []struct {
		name  string
		input any
	}{
		{
			name:  "int value",
			input: 42,
		},
		{
			name:  "string value",
			input: "hello",
		},
		{
			name:  "bool value",
			input: true,
		},
		{
			name:  "float64 value",
			input: 3.14,
		},
		{
			name:  "[]any slice",
			input: []any{1, 2, 3},
		},
		{
			name:  "map[string]any map",
			input: map[string]any{"a": 1},
		},
		{
			name: "struct value",
			input: struct {
				Field string
			}{Field: "test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := deref(tt.input)

			if !ok {
				t.Errorf("expected ok=true for non-pointer %v, got false", tt.input)
			}

			// Type should be unchanged
			resultType := reflect.TypeOf(result)
			inputType := reflect.TypeOf(tt.input)
			if resultType != inputType {
				t.Errorf("expected type %v, got %v", inputType, resultType)
			}

			// Compare values using reflect.DeepEqual to handle slices and maps
			if !reflect.DeepEqual(result, tt.input) {
				t.Errorf("expected %v, got %v", tt.input, result)
			}
		})
	}
}
