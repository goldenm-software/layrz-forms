package layrz

import "testing"

func TestValidateBool(t *testing.T) {
	t.Run("absent required true", func(t *testing.T) {
		errs := ValidateBool(nil, BoolRules{Required: true})
		if len(errs) != 1 || errs[0].Code != "required" {
			t.Errorf("expected [required], got %v", errs)
		}
	})

	t.Run("absent required false", func(t *testing.T) {
		errs := ValidateBool(nil, BoolRules{Required: false})
		if len(errs) != 0 {
			t.Errorf("expected no errors, got %v", errs)
		}
	})

	t.Run("valid true", func(t *testing.T) {
		errs := ValidateBool(Ptr(true), BoolRules{Required: true})
		if len(errs) != 0 {
			t.Errorf("expected no errors, got %v", errs)
		}
	})

	t.Run("valid false", func(t *testing.T) {
		errs := ValidateBool(Ptr(false), BoolRules{Required: true})
		if len(errs) != 0 {
			t.Errorf("expected no errors, got %v", errs)
		}
	})
}
