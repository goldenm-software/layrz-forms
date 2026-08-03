package layrz

import "testing"

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Basic cases
		{"user_name", "userName"},
		{"id_test", "idTest"},
		{"range_text_test", "rangeTextTest"},

		// Dot notation
		{"a_b.c_d", "aB.cD"},
		{"address.street_name", "address.streetName"},

		// Numeric segments
		{"items.0.name", "items.0.name"},

		// Empty and edge cases
		{"", ""},
		{"_foo", "foo"},
		{"foo_", "foo"},

		// Single word (no underscores)
		{"user", "user"},
		{"address", "address"},

		// str.title quirks (uppercase first, lowercase rest)
		{"a_bC", "aBc"},
		{"a_BC", "aBc"},

		// Non-alpha treated as word separators (a1b.title() == A1B)
		{"a_1b", "a1B"},
		{"a_b_1c", "aB1C"},

		// Underscores at various positions
		{"_", ""},
		{"__", ""},
		{"_a_b", "aB"},
		{"a__b", "aB"},

		// Complex cases
		{"field_name_test", "fieldNameTest"},
		{"outer.inner_field", "outer.innerField"},
		{"a.b_c.d_e_f", "a.bC.dEF"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ToCamelCase(tt.input)
			if got != tt.expected {
				t.Errorf("ToCamelCase(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestPythonTitle(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Basic
		{"foo", "Foo"},
		{"bar", "Bar"},

		// str.title quirk: uppercase first, lowercase rest
		{"bC", "Bc"},
		{"BC", "Bc"},
		{"Bc", "Bc"},

		// Non-alpha as separators (previous char not a letter)
		{"a1b", "A1B"},
		{"a1B2c", "A1B2C"},

		// Single character
		{"a", "A"},
		{"1", "1"},

		// Empty
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := pythonTitle(tt.input)
			if got != tt.expected {
				t.Errorf("pythonTitle(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
