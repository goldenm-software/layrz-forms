package layrz

import (
	"testing"
)

func TestValidateID(t *testing.T) {
	t.Run("absent required true", func(t *testing.T) {
		errs := ValidateID(nil, IDRules{Required: true})
		if len(errs) != 1 || errs[0].Code != "required" {
			t.Errorf("expected [required], got %v", errs)
		}
	})

	t.Run("absent required false", func(t *testing.T) {
		errs := ValidateID(nil, IDRules{Required: false})
		if len(errs) != 0 {
			t.Errorf("expected no errors, got %v", errs)
		}
	})

	t.Run("valid positive int", func(t *testing.T) {
		errs := ValidateID(42, IDRules{Required: true})
		if len(errs) != 0 {
			t.Errorf("expected no errors, got %v", errs)
		}
	})

	t.Run("valid positive string", func(t *testing.T) {
		errs := ValidateID("123", IDRules{Required: true})
		if len(errs) != 0 {
			t.Errorf("expected no errors, got %v", errs)
		}
	})

	t.Run("invalid zero", func(t *testing.T) {
		errs := ValidateID(0, IDRules{Required: true})
		if len(errs) != 1 || errs[0].Code != "invalid" {
			t.Errorf("expected [invalid], got %v", errs)
		}
	})

	t.Run("invalid negative", func(t *testing.T) {
		errs := ValidateID(-5, IDRules{Required: true})
		if len(errs) != 1 || errs[0].Code != "invalid" {
			t.Errorf("expected [invalid], got %v", errs)
		}
	})

	t.Run("invalid negative string", func(t *testing.T) {
		errs := ValidateID("-42", IDRules{Required: true})
		if len(errs) != 1 || errs[0].Code != "invalid" {
			t.Errorf("expected [invalid], got %v", errs)
		}
	})

	t.Run("invalid non-numeric string required true", func(t *testing.T) {
		errs := ValidateID("abc", IDRules{Required: true})
		if len(errs) != 1 || errs[0].Code != "invalid" {
			t.Errorf("expected [invalid], got %v", errs)
		}
	})

	t.Run("invalid non-numeric string required false", func(t *testing.T) {
		errs := ValidateID("abc", IDRules{Required: false})
		if len(errs) != 1 || errs[0].Code != "invalid" {
			t.Errorf("expected [invalid], got %v", errs)
		}
	})

	t.Run("invalid float required true", func(t *testing.T) {
		errs := ValidateID(3.14, IDRules{Required: true})
		if len(errs) != 1 || errs[0].Code != "invalid" {
			t.Errorf("expected [invalid], got %v", errs)
		}
	})

	t.Run("invalid float required false", func(t *testing.T) {
		errs := ValidateID(1.5, IDRules{Required: false})
		if len(errs) != 1 || errs[0].Code != "invalid" {
			t.Errorf("expected [invalid], got %v", errs)
		}
	})

	t.Run("invalid bool", func(t *testing.T) {
		errs := ValidateID(true, IDRules{Required: true})
		if len(errs) != 1 || errs[0].Code != "invalid" {
			t.Errorf("expected [invalid], got %v", errs)
		}
	})

	t.Run("invalid list", func(t *testing.T) {
		errs := ValidateID([]int{}, IDRules{Required: true})
		if len(errs) != 1 || errs[0].Code != "invalid" {
			t.Errorf("expected [invalid], got %v", errs)
		}
	})

	t.Run("invalid dict", func(t *testing.T) {
		errs := ValidateID(map[string]int{}, IDRules{Required: true})
		if len(errs) != 1 || errs[0].Code != "invalid" {
			t.Errorf("expected [invalid], got %v", errs)
		}
	})
}
