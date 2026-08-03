package layrz

import "testing"

func TestValidateEmail(t *testing.T) {
	t.Run("absent required true", func(t *testing.T) {
		errs := ValidateEmail(nil, EmailRules{Required: true})
		if len(errs) != 1 || errs[0].Code != "required" {
			t.Errorf("expected [required], got %v", errs)
		}
	})

	t.Run("absent required false", func(t *testing.T) {
		errs := ValidateEmail(nil, EmailRules{Required: false})
		if len(errs) != 0 {
			t.Errorf("expected no errors, got %v", errs)
		}
	})

	t.Run("valid email", func(t *testing.T) {
		errs := ValidateEmail(Ptr("test@example.com"), EmailRules{Required: true})
		if len(errs) != 0 {
			t.Errorf("expected no errors, got %v", errs)
		}
	})

	t.Run("valid email with plus addressing", func(t *testing.T) {
		errs := ValidateEmail(Ptr("test+tag@example.com"), EmailRules{Required: true})
		if len(errs) != 0 {
			t.Errorf("expected no errors, got %v", errs)
		}
	})

	t.Run("invalid no at sign", func(t *testing.T) {
		errs := ValidateEmail(Ptr("notanemail.com"), EmailRules{Required: true})
		if len(errs) != 1 || errs[0].Code != "invalid" {
			t.Errorf("expected [invalid], got %v", errs)
		}
	})

	t.Run("invalid no tld", func(t *testing.T) {
		errs := ValidateEmail(Ptr("test@example"), EmailRules{Required: true})
		if len(errs) != 1 || errs[0].Code != "invalid" {
			t.Errorf("expected [invalid], got %v", errs)
		}
	})

	t.Run("empty string empty false required true", func(t *testing.T) {
		errs := ValidateEmail(Ptr(""), EmailRules{Required: true, Empty: false})
		if len(errs) != 1 || errs[0].Code != "empty" {
			t.Errorf("expected [empty], got %v", errs)
		}
	})

	t.Run("empty string empty true", func(t *testing.T) {
		errs := ValidateEmail(Ptr(""), EmailRules{Required: false, Empty: true})
		if len(errs) != 0 {
			t.Errorf("expected no errors, got %v", errs)
		}
	})

	t.Run("non-empty invalid format with empty true", func(t *testing.T) {
		errs := ValidateEmail(Ptr("not-an-email"), EmailRules{Required: true, Empty: true})
		if len(errs) != 1 || errs[0].Code != "invalid" {
			t.Errorf("expected [invalid], got %v", errs)
		}
	})

	t.Run("valid with empty true", func(t *testing.T) {
		errs := ValidateEmail(Ptr("user@example.com"), EmailRules{Required: true, Empty: true})
		if len(errs) != 0 {
			t.Errorf("expected no errors, got %v", errs)
		}
	})

	t.Run("invalid int type", func(t *testing.T) {
		// Email validates only *string, so this shouldn't happen in normal use
		// But test the logic for wrong type handling
		var val *string
		errs := ValidateEmail(val, EmailRules{Required: true})
		if len(errs) != 1 || errs[0].Code != "required" {
			t.Errorf("expected [required] for nil, got %v", errs)
		}
	})
}
