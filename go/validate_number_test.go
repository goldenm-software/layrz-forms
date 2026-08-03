package layrz

import "testing"

func TestValidateNumber(t *testing.T) {
	t.Run("int absent required true", func(t *testing.T) {
		errs := ValidateNumber(nil, NumberRules{Required: true, Datatype: "int"})
		if len(errs) != 1 || errs[0].Code != "required" {
			t.Errorf("expected [required], got %v", errs)
		}
	})

	t.Run("int absent required false", func(t *testing.T) {
		errs := ValidateNumber(nil, NumberRules{Required: false, Datatype: "int"})
		if len(errs) != 0 {
			t.Errorf("expected no errors, got %v", errs)
		}
	})

	t.Run("int valid positive", func(t *testing.T) {
		errs := ValidateNumber(42, NumberRules{Required: true, Datatype: "int"})
		if len(errs) != 0 {
			t.Errorf("expected no errors, got %v", errs)
		}
	})

	t.Run("int zero", func(t *testing.T) {
		errs := ValidateNumber(0, NumberRules{Required: true, Datatype: "int"})
		if len(errs) != 0 {
			t.Errorf("expected no errors, got %v", errs)
		}
	})

	t.Run("int negative", func(t *testing.T) {
		errs := ValidateNumber(-10, NumberRules{Required: true, Datatype: "int"})
		if len(errs) != 0 {
			t.Errorf("expected no errors, got %v", errs)
		}
	})

	t.Run("int min_value exact", func(t *testing.T) {
		errs := ValidateNumber(5, NumberRules{Required: true, Datatype: "int", MinValue: Ptr(5.0)})
		if len(errs) != 0 {
			t.Errorf("expected no errors, got %v", errs)
		}
	})

	t.Run("int min_value under", func(t *testing.T) {
		errs := ValidateNumber(4, NumberRules{Required: true, Datatype: "int", MinValue: Ptr(5.0)})
		if len(errs) != 1 || errs[0].Code != "minValue" {
			t.Errorf("expected [minValue], got %v", errs)
		}
		if errs[0].Expected != int64(5) || errs[0].Received != int64(4) {
			t.Errorf("expected Expected=5, Received=4, got %v, %v", errs[0].Expected, errs[0].Received)
		}
	})

	t.Run("int max_value exact", func(t *testing.T) {
		errs := ValidateNumber(10, NumberRules{Required: true, Datatype: "int", MaxValue: Ptr(10.0)})
		if len(errs) != 0 {
			t.Errorf("expected no errors, got %v", errs)
		}
	})

	t.Run("int max_value over", func(t *testing.T) {
		errs := ValidateNumber(15, NumberRules{Required: true, Datatype: "int", MaxValue: Ptr(10.0)})
		if len(errs) != 1 || errs[0].Code != "maxValue" {
			t.Errorf("expected [maxValue], got %v", errs)
		}
		if errs[0].Expected != int64(10) || errs[0].Received != int64(15) {
			t.Errorf("expected Expected=10, Received=15, got %v, %v", errs[0].Expected, errs[0].Received)
		}
	})

	t.Run("float valid", func(t *testing.T) {
		errs := ValidateNumber(3.14, NumberRules{Required: true, Datatype: "float"})
		if len(errs) != 0 {
			t.Errorf("expected no errors, got %v", errs)
		}
	})

	t.Run("float rejects int", func(t *testing.T) {
		errs := ValidateNumber(42, NumberRules{Required: true, Datatype: "float"})
		if len(errs) != 1 || errs[0].Code != "invalid" {
			t.Errorf("expected [invalid], got %v", errs)
		}
	})

	t.Run("float min_value", func(t *testing.T) {
		errs := ValidateNumber(1.5, NumberRules{Required: true, Datatype: "float", MinValue: Ptr(1.5)})
		if len(errs) != 0 {
			t.Errorf("expected no errors, got %v", errs)
		}
	})

	t.Run("float max_value", func(t *testing.T) {
		errs := ValidateNumber(99.99, NumberRules{Required: true, Datatype: "float", MaxValue: Ptr(99.99)})
		if len(errs) != 0 {
			t.Errorf("expected no errors, got %v", errs)
		}
	})

	t.Run("int rejects string", func(t *testing.T) {
		errs := ValidateNumber("not a number", NumberRules{Required: true, Datatype: "int"})
		if len(errs) != 1 || errs[0].Code != "invalid" {
			t.Errorf("expected [invalid], got %v", errs)
		}
	})

	t.Run("int rejects bool true", func(t *testing.T) {
		errs := ValidateNumber(true, NumberRules{Required: true, Datatype: "int"})
		if len(errs) != 1 || errs[0].Code != "invalid" {
			t.Errorf("expected [invalid] for bool, got %v", errs)
		}
	})

	t.Run("int rejects bool false", func(t *testing.T) {
		errs := ValidateNumber(false, NumberRules{Required: true, Datatype: "int"})
		if len(errs) != 1 || errs[0].Code != "invalid" {
			t.Errorf("expected [invalid] for bool, got %v", errs)
		}
	})

	t.Run("float rejects bool", func(t *testing.T) {
		errs := ValidateNumber(true, NumberRules{Required: false, Datatype: "float"})
		if len(errs) != 1 || errs[0].Code != "invalid" {
			t.Errorf("expected [invalid] for bool, got %v", errs)
		}
	})

	t.Run("accumulate min and max errors", func(t *testing.T) {
		errs := ValidateNumber(200, NumberRules{
			Required: true,
			Datatype: "int",
			MinValue: Ptr(100.0),
			MaxValue: Ptr(150.0),
		})
		if len(errs) != 1 {
			t.Errorf("expected 1 error (maxValue), got %d: %v", len(errs), errs)
		}
		if errs[0].Code != "maxValue" {
			t.Errorf("expected maxValue error, got %q", errs[0].Code)
		}
	})

	t.Run("accumulate both min and max errors", func(t *testing.T) {
		errs := ValidateNumber(50, NumberRules{
			Required: true,
			Datatype: "int",
			MinValue: Ptr(100.0),
			MaxValue: Ptr(40.0),
		})
		if len(errs) != 2 {
			t.Errorf("expected 2 errors, got %d: %v", len(errs), errs)
		}
		codes := map[string]bool{}
		for _, e := range errs {
			codes[e.Code] = true
		}
		if !codes["minValue"] || !codes["maxValue"] {
			t.Errorf("expected minValue and maxValue errors, got %v", codes)
		}
	})
}
