package layrz

import "testing"

func TestValidateChar(t *testing.T) {
	t.Run("absent required true", func(t *testing.T) {
		errs := ValidateChar(nil, CharRules{Required: true})
		if len(errs) != 1 || errs[0].Code != "required" {
			t.Errorf("expected [required], got %v", errs)
		}
	})

	t.Run("absent required false", func(t *testing.T) {
		errs := ValidateChar(nil, CharRules{Required: false})
		if len(errs) != 0 {
			t.Errorf("expected no errors, got %v", errs)
		}
	})

	t.Run("valid string", func(t *testing.T) {
		errs := ValidateChar(Ptr("John"), CharRules{Required: true})
		if len(errs) != 0 {
			t.Errorf("expected no errors, got %v", errs)
		}
	})

	t.Run("empty string empty false", func(t *testing.T) {
		errs := ValidateChar(Ptr(""), CharRules{Required: true, Empty: false})
		if len(errs) != 1 || errs[0].Code != "empty" {
			t.Errorf("expected [empty], got %v", errs)
		}
	})

	t.Run("empty string empty true", func(t *testing.T) {
		errs := ValidateChar(Ptr(""), CharRules{Required: false, Empty: true})
		if len(errs) != 0 {
			t.Errorf("expected no errors, got %v", errs)
		}
	})

	t.Run("min_length exact", func(t *testing.T) {
		errs := ValidateChar(Ptr("hello"), CharRules{Required: true, MinLength: Ptr(5)})
		if len(errs) != 0 {
			t.Errorf("expected no errors, got %v", errs)
		}
	})

	t.Run("min_length under", func(t *testing.T) {
		errs := ValidateChar(Ptr("hola"), CharRules{Required: true, MinLength: Ptr(5)})
		if len(errs) != 1 || errs[0].Code != "minLength" {
			t.Errorf("expected [minLength], got %v", errs)
		}
		if errs[0].Expected != 5 || errs[0].Received != 4 {
			t.Errorf("expected Expected=5, Received=4, got %v, %v", errs[0].Expected, errs[0].Received)
		}
	})

	t.Run("max_length exact", func(t *testing.T) {
		errs := ValidateChar(Ptr("hello"), CharRules{Required: true, MaxLength: Ptr(5)})
		if len(errs) != 0 {
			t.Errorf("expected no errors, got %v", errs)
		}
	})

	t.Run("max_length over", func(t *testing.T) {
		errs := ValidateChar(Ptr("hello world"), CharRules{Required: true, MaxLength: Ptr(5)})
		if len(errs) != 1 || errs[0].Code != "maxLength" {
			t.Errorf("expected [maxLength], got %v", errs)
		}
		if errs[0].Expected != 5 || errs[0].Received != 11 {
			t.Errorf("expected Expected=5, Received=11, got %v, %v", errs[0].Expected, errs[0].Received)
		}
	})

	t.Run("choices valid", func(t *testing.T) {
		errs := ValidateChar(Ptr("a"), CharRules{Required: true, Choices: []string{"a", "b", "c"}})
		if len(errs) != 0 {
			t.Errorf("expected no errors, got %v", errs)
		}
	})

	t.Run("choices invalid", func(t *testing.T) {
		errs := ValidateChar(Ptr("d"), CharRules{Required: true, Choices: []string{"a", "b", "c"}})
		if len(errs) != 1 || errs[0].Code != "invalidChoice" {
			t.Errorf("expected [invalidChoice], got %v", errs)
		}
	})

	t.Run("regex valid", func(t *testing.T) {
		errs := ValidateChar(Ptr("abc"), CharRules{Required: true, Regex: "^[a-z]+$"})
		if len(errs) != 0 {
			t.Errorf("expected no errors, got %v", errs)
		}
	})

	t.Run("regex invalid", func(t *testing.T) {
		errs := ValidateChar(Ptr("abc123"), CharRules{Required: true, Regex: "^[a-z]+$"})
		if len(errs) != 1 || errs[0].Code != "invalidFormat" {
			t.Errorf("expected [invalidFormat], got %v", errs)
		}
	})

	t.Run("accumulate multiple errors", func(t *testing.T) {
		errs := ValidateChar(Ptr("a"), CharRules{
			Required:  true,
			MinLength: Ptr(5),
			MaxLength: Ptr(10),
		})
		if len(errs) != 1 || errs[0].Code != "minLength" {
			t.Errorf("expected [minLength], got %v", errs)
		}
	})

	t.Run("utf8 rune counting", func(t *testing.T) {
		// "Café" is 4 runes (C-a-f-é) even though it's more bytes
		errs := ValidateChar(Ptr("Café"), CharRules{Required: true, MinLength: Ptr(4)})
		if len(errs) != 0 {
			t.Errorf("expected no errors for 4-rune string, got %v", errs)
		}

		// But 5 min_length should fail
		errs = ValidateChar(Ptr("Café"), CharRules{Required: true, MinLength: Ptr(5)})
		if len(errs) != 1 || errs[0].Code != "minLength" {
			t.Errorf("expected [minLength], got %v", errs)
		}
	})
}
