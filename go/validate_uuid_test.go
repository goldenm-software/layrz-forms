package layrz

import "testing"

func TestValidateUUID(t *testing.T) {
	t.Run("absent required true", func(t *testing.T) {
		errs := ValidateUUID(nil, UUIDRules{Required: true})
		if len(errs) != 1 || errs[0].Code != "required" {
			t.Errorf("expected [required], got %v", errs)
		}
	})

	t.Run("absent required false", func(t *testing.T) {
		errs := ValidateUUID(nil, UUIDRules{Required: false})
		if len(errs) != 0 {
			t.Errorf("expected no errors, got %v", errs)
		}
	})

	t.Run("valid hyphenated uuid", func(t *testing.T) {
		errs := ValidateUUID(Ptr("550e8400-e29b-41d4-a716-446655440000"), UUIDRules{Required: true})
		if len(errs) != 0 {
			t.Errorf("expected no errors, got %v", errs)
		}
	})

	t.Run("valid unhyphenated uuid", func(t *testing.T) {
		errs := ValidateUUID(Ptr("550e8400e29b41d4a716446655440000"), UUIDRules{Required: true})
		if len(errs) != 0 {
			t.Errorf("expected no errors, got %v", errs)
		}
	})

	t.Run("valid braced uuid", func(t *testing.T) {
		errs := ValidateUUID(Ptr("{550e8400-e29b-41d4-a716-446655440000}"), UUIDRules{Required: true})
		if len(errs) != 0 {
			t.Errorf("expected no errors for braced uuid, got %v", errs)
		}
	})

	t.Run("valid urn:uuid prefix", func(t *testing.T) {
		errs := ValidateUUID(Ptr("urn:uuid:550e8400-e29b-41d4-a716-446655440000"), UUIDRules{Required: true})
		if len(errs) != 0 {
			t.Errorf("expected no errors for urn:uuid, got %v", errs)
		}
	})

	t.Run("valid case insensitive", func(t *testing.T) {
		errs := ValidateUUID(Ptr("550E8400-E29B-41D4-A716-446655440000"), UUIDRules{Required: true})
		if len(errs) != 0 {
			t.Errorf("expected no errors for uppercase, got %v", errs)
		}
	})

	t.Run("invalid format", func(t *testing.T) {
		errs := ValidateUUID(Ptr("not-a-uuid"), UUIDRules{Required: true})
		if len(errs) != 1 || errs[0].Code != "invalid" {
			t.Errorf("expected [invalid], got %v", errs)
		}
	})

	t.Run("invalid empty string", func(t *testing.T) {
		errs := ValidateUUID(Ptr(""), UUIDRules{Required: true})
		if len(errs) != 1 || errs[0].Code != "invalid" {
			t.Errorf("expected [invalid], got %v", errs)
		}
	})

	t.Run("invalid malformed hyphenated", func(t *testing.T) {
		errs := ValidateUUID(Ptr("550e8400-e29b-41d4-a716"), UUIDRules{Required: true})
		if len(errs) != 1 || errs[0].Code != "invalid" {
			t.Errorf("expected [invalid], got %v", errs)
		}
	})
}

func TestIsValidUUID(t *testing.T) {
	tests := []struct {
		name     string
		uuid     string
		expected bool
	}{
		// Hyphenated
		{"hyphenated standard", "550e8400-e29b-41d4-a716-446655440000", true},
		{"hyphenated uppercase", "550E8400-E29B-41D4-A716-446655440000", true},
		{"hyphenated mixed case", "550e8400-E29B-41d4-a716-446655440000", true},

		// Unhyphenated
		{"unhyphenated", "550e8400e29b41d4a716446655440000", true},
		{"unhyphenated uppercase", "550E8400E29B41D4A716446655440000", true},

		// Braced
		{"braced", "{550e8400-e29b-41d4-a716-446655440000}", true},

		// URN
		{"urn:uuid", "urn:uuid:550e8400-e29b-41d4-a716-446655440000", true},
		{"urn:uuid uppercase", "URN:UUID:550e8400-e29b-41d4-a716-446655440000", true},

		// Invalid
		{"empty", "", false},
		{"not uuid", "not-a-uuid", false},
		{"short", "550e8400", false},
		{"too long", "550e8400-e29b-41d4-a716-446655440000-extra", false},
		{"missing digits", "550e8400-e29b-41d4-a716-44665544000", false},
		{"non-hex", "g50e8400-e29b-41d4-a716-446655440000", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidUUID(tt.uuid)
			if got != tt.expected {
				t.Errorf("isValidUUID(%q) = %v, want %v", tt.uuid, got, tt.expected)
			}
		})
	}
}
