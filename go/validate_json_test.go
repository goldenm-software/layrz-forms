package layrz

import "testing"

func TestValidateJSON(t *testing.T) {
	t.Run("dict absent required true", func(t *testing.T) {
		errs := ValidateJSON(nil, JSONRules{Required: true, Datatype: "dict"})
		if len(errs) != 1 || errs[0].Code != "required" {
			t.Errorf("expected [required], got %v", errs)
		}
	})

	t.Run("dict absent required false", func(t *testing.T) {
		errs := ValidateJSON(nil, JSONRules{Required: false, Datatype: "dict"})
		if len(errs) != 0 {
			t.Errorf("expected no errors, got %v", errs)
		}
	})

	t.Run("dict populated", func(t *testing.T) {
		m := map[string]any{"key": "value"}
		errs := ValidateJSON(&m, JSONRules{Required: true, Datatype: "dict"})
		if len(errs) != 0 {
			t.Errorf("expected no errors, got %v", errs)
		}
	})

	t.Run("dict empty empty true", func(t *testing.T) {
		m := map[string]any{}
		errs := ValidateJSON(&m, JSONRules{Required: false, Datatype: "dict", Empty: true})
		if len(errs) != 0 {
			t.Errorf("expected no errors, got %v", errs)
		}
	})

	t.Run("dict empty empty false", func(t *testing.T) {
		m := map[string]any{}
		errs := ValidateJSON(&m, JSONRules{Required: false, Datatype: "dict", Empty: false})
		if len(errs) != 1 || errs[0].Code != "invalid" {
			t.Errorf("expected [invalid], got %v", errs)
		}
	})

	t.Run("dict wrong type list", func(t *testing.T) {
		l := []any{}
		errs := ValidateJSON(&l, JSONRules{Required: true, Datatype: "dict"})
		if len(errs) != 1 || errs[0].Code != "invalid" {
			t.Errorf("expected [invalid], got %v", errs)
		}
	})

	t.Run("dict wrong type string", func(t *testing.T) {
		s := "not a dict"
		errs := ValidateJSON(&s, JSONRules{Required: true, Datatype: "dict"})
		if len(errs) != 1 || errs[0].Code != "invalid" {
			t.Errorf("expected [invalid], got %v", errs)
		}
	})

	t.Run("list absent required true", func(t *testing.T) {
		errs := ValidateJSON(nil, JSONRules{Required: true, Datatype: "list"})
		if len(errs) != 1 || errs[0].Code != "required" {
			t.Errorf("expected [required], got %v", errs)
		}
	})

	t.Run("list absent required false", func(t *testing.T) {
		errs := ValidateJSON(nil, JSONRules{Required: false, Datatype: "list"})
		if len(errs) != 0 {
			t.Errorf("expected no errors, got %v", errs)
		}
	})

	t.Run("list populated", func(t *testing.T) {
		l := []any{1, 2, 3}
		errs := ValidateJSON(&l, JSONRules{Required: true, Datatype: "list"})
		if len(errs) != 0 {
			t.Errorf("expected no errors, got %v", errs)
		}
	})

	t.Run("list empty empty true", func(t *testing.T) {
		l := []any{}
		errs := ValidateJSON(&l, JSONRules{Required: false, Datatype: "list", Empty: true})
		if len(errs) != 0 {
			t.Errorf("expected no errors, got %v", errs)
		}
	})

	t.Run("list empty empty false", func(t *testing.T) {
		l := []any{}
		errs := ValidateJSON(&l, JSONRules{Required: false, Datatype: "list", Empty: false})
		if len(errs) != 1 || errs[0].Code != "invalid" {
			t.Errorf("expected [invalid], got %v", errs)
		}
	})

	t.Run("list wrong type dict", func(t *testing.T) {
		m := map[string]any{}
		errs := ValidateJSON(&m, JSONRules{Required: true, Datatype: "list"})
		if len(errs) != 1 || errs[0].Code != "invalid" {
			t.Errorf("expected [invalid], got %v", errs)
		}
	})
}
