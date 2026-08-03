package layrz

import "testing"

//nolint:gocyclo // table-driven test legitimately exceeds complexity 30
func TestParseTag(t *testing.T) {
	t.Run("empty and skip cases", func(t *testing.T) {
		spec, err := ParseTag("")
		if spec != nil || err != nil {
			t.Errorf("ParseTag(\"\") = (%v, %v), want (nil, nil)", spec, err)
		}

		spec, err = ParseTag("-")
		if spec != nil || err != nil {
			t.Errorf("ParseTag(\"-\") = (%v, %v), want (nil, nil)", spec, err)
		}
	})

	t.Run("valid field kinds", func(t *testing.T) {
		kinds := []string{"id", "email", "uuid", "char", "number", "bool", "json", "subform", "subform_list"}
		for _, kind := range kinds {
			spec, err := ParseTag(kind)
			if err != nil {
				t.Errorf("ParseTag(%q) error: %v", kind, err)
				continue
			}
			if spec == nil {
				t.Errorf("ParseTag(%q) = nil, want *FieldSpec", kind)
				continue
			}
			if spec.Kind != FieldKind(kind) {
				t.Errorf("ParseTag(%q).Kind = %q, want %q", kind, spec.Kind, kind)
			}
		}
	})

	t.Run("invalid field kind", func(t *testing.T) {
		_, err := ParseTag("invalid_kind")
		if err == nil {
			t.Errorf("ParseTag(\"invalid_kind\") expected error, got nil")
		}
	})

	t.Run("required flag", func(t *testing.T) {
		spec, err := ParseTag("char,required")
		if err != nil {
			t.Fatalf("ParseTag error: %v", err)
		}
		if !spec.Required {
			t.Errorf("expected Required=true")
		}
	})

	t.Run("empty flag", func(t *testing.T) {
		spec, err := ParseTag("char,empty")
		if err != nil {
			t.Fatalf("ParseTag error: %v", err)
		}
		if !spec.Empty {
			t.Errorf("expected Empty=true")
		}
	})

	t.Run("min_length and max_length", func(t *testing.T) {
		spec, err := ParseTag("char,min_length=5,max_length=10")
		if err != nil {
			t.Fatalf("ParseTag error: %v", err)
		}
		if spec.MinLength == nil || *spec.MinLength != 5 {
			t.Errorf("expected MinLength=5, got %v", spec.MinLength)
		}
		if spec.MaxLength == nil || *spec.MaxLength != 10 {
			t.Errorf("expected MaxLength=10, got %v", spec.MaxLength)
		}
	})

	t.Run("min_value and max_value", func(t *testing.T) {
		spec, err := ParseTag("number,min_value=1.5,max_value=99.99")
		if err != nil {
			t.Fatalf("ParseTag error: %v", err)
		}
		if spec.MinValue == nil || *spec.MinValue != 1.5 {
			t.Errorf("expected MinValue=1.5, got %v", spec.MinValue)
		}
		if spec.MaxValue == nil || *spec.MaxValue != 99.99 {
			t.Errorf("expected MaxValue=99.99, got %v", spec.MaxValue)
		}
	})

	t.Run("choices", func(t *testing.T) {
		spec, err := ParseTag("char,choices=a|b|c")
		if err != nil {
			t.Fatalf("ParseTag error: %v", err)
		}
		if len(spec.Choices) != 3 || spec.Choices[0] != "a" || spec.Choices[1] != "b" || spec.Choices[2] != "c" {
			t.Errorf("expected choices [a, b, c], got %v", spec.Choices)
		}
	})

	t.Run("datatype", func(t *testing.T) {
		spec, err := ParseTag("number,datatype=int")
		if err != nil {
			t.Fatalf("ParseTag error: %v", err)
		}
		if spec.Datatype != testSpecDatatypeInt {
			t.Errorf("expected Datatype=%s, got %q", testSpecDatatypeInt, spec.Datatype)
		}
	})

	t.Run("regex as last rule", func(t *testing.T) {
		spec, err := ParseTag("char,min_length=5,regex=^[a-z]+$")
		if err != nil {
			t.Fatalf("ParseTag error: %v", err)
		}
		if spec.Regex != "^[a-z]+$" {
			t.Errorf("expected Regex=^[a-z]+$, got %q", spec.Regex)
		}
		if spec.MinLength == nil || *spec.MinLength != 5 {
			t.Errorf("expected MinLength=5 before regex")
		}
	})

	t.Run("regex with comma in pattern", func(t *testing.T) {
		spec, err := ParseTag("char,regex=^[A-Z]{2,3}$")
		if err != nil {
			t.Fatalf("ParseTag error: %v", err)
		}
		if spec.Regex != "^[A-Z]{2,3}$" {
			t.Errorf("expected Regex=^[A-Z]{2,3}$, got %q", spec.Regex)
		}
	})

	t.Run("invalid datatype value", func(t *testing.T) {
		_, err := ParseTag("number,datatype=bad")
		if err == nil {
			t.Errorf("ParseTag with invalid datatype expected error, got nil")
		}
	})

	t.Run("invalid min_length value", func(t *testing.T) {
		_, err := ParseTag("char,min_length=abc")
		if err == nil {
			t.Errorf("ParseTag with invalid min_length expected error, got nil")
		}
	})

	t.Run("duplicate rule", func(t *testing.T) {
		_, err := ParseTag("char,required,required")
		if err == nil {
			t.Errorf("ParseTag with duplicate required expected error, got nil")
		}
	})

	t.Run("invalid regex", func(t *testing.T) {
		_, err := ParseTag("char,regex=[invalid(")
		if err == nil {
			t.Errorf("ParseTag with invalid regex expected error, got nil")
		}
	})

	t.Run("unknown rule", func(t *testing.T) {
		_, err := ParseTag("char,unknown_rule=value")
		if err == nil {
			t.Errorf("ParseTag with unknown rule expected error, got nil")
		}
	})
}
