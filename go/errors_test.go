package layrz

import (
	"encoding/json"
	"testing"
)

func TestErrorsAdd(t *testing.T) {
	t.Run("add single error", func(t *testing.T) {
		e := make(Errors)
		e.Add("field", &FieldError{Code: "required"})
		if len(e["field"]) != 1 {
			t.Fatalf("expected 1 error, got %d", len(e["field"]))
		}
		if e["field"][0].Code != "required" {
			t.Errorf("expected code 'required', got %q", e["field"][0].Code)
		}
	})

	t.Run("add multiple errors", func(t *testing.T) {
		e := make(Errors)
		e.Add("field", &FieldError{Code: "invalid"}, &FieldError{Code: "minValue"})
		if len(e["field"]) != 2 {
			t.Fatalf("expected 2 errors, got %d", len(e["field"]))
		}
	})

	t.Run("skip nil entries", func(t *testing.T) {
		e := make(Errors)
		e.Add("field", &FieldError{Code: "invalid"}, nil, &FieldError{Code: "minValue"})
		if len(e["field"]) != 2 {
			t.Fatalf("expected 2 errors (nil skipped), got %d", len(e["field"]))
		}
	})

	t.Run("no-op on empty slice", func(t *testing.T) {
		e := make(Errors)
		e.Add("field")
		if len(e) != 0 {
			t.Errorf("expected empty errors, got %d entries", len(e))
		}
	})
}

func TestErrorsMerge(t *testing.T) {
	e1 := make(Errors)
	e1["userName"] = []*FieldError{{Code: "required"}}

	e2 := make(Errors)
	e2["user_name"] = []*FieldError{{Code: "invalid"}}

	e1.Merge(e2)

	// user_name should be converted to camelCase (userName)
	if len(e1["userName"]) != 2 {
		t.Fatalf("expected 2 errors after merge, got %d", len(e1["userName"]))
	}
}

func TestErrorsIsEmpty(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		e := make(Errors)
		if !e.IsEmpty() {
			t.Errorf("expected empty, got non-empty")
		}
	})

	t.Run("not empty", func(t *testing.T) {
		e := make(Errors)
		e.Add("field", &FieldError{Code: "required"})
		if e.IsEmpty() {
			t.Errorf("expected non-empty, got empty")
		}
	})
}

func TestErrorsKeys(t *testing.T) {
	e := make(Errors)
	e.Add("zebra", &FieldError{Code: "invalid"})
	e.Add("apple", &FieldError{Code: "required"})
	e.Add("banana", &FieldError{Code: "invalid"})

	keys := e.Keys()
	expected := []string{"apple", "banana", "zebra"}

	if len(keys) != len(expected) {
		t.Fatalf("expected %d keys, got %d", len(expected), len(keys))
	}

	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("keys[%d] = %q, want %q", i, k, expected[i])
		}
	}
}

func TestFieldErrorJSONMarshal(t *testing.T) {
	t.Run("expected and received with 0 values", func(t *testing.T) {
		fe := &FieldError{Code: "minValue", Expected: 0, Received: 0}
		data, err := json.Marshal(fe)
		if err != nil {
			t.Fatalf("failed to marshal: %v", err)
		}

		str := string(data)
		// The JSON should include "expected":0 and "received":0, not omit them
		if !json.Valid(data) {
			t.Errorf("invalid JSON: %s", str)
		}

		// Unmarshal to verify the fields are present
		var m map[string]any
		json.Unmarshal(data, &m)
		if _, hasExpected := m["expected"]; !hasExpected {
			t.Errorf("expected field not present in JSON: %s", str)
		}
		if _, hasReceived := m["received"]; !hasReceived {
			t.Errorf("received field not present in JSON: %s", str)
		}
	})

	t.Run("omit nil extra", func(t *testing.T) {
		fe := &FieldError{Code: "invalid", Extra: nil}
		data, err := json.Marshal(fe)
		if err != nil {
			t.Fatalf("failed to marshal: %v", err)
		}

		var m map[string]any
		json.Unmarshal(data, &m)
		if _, hasExtra := m["extra"]; hasExtra {
			t.Errorf("extra field should not be present when nil")
		}
	})
}

func TestErrorsJSONMarshal(t *testing.T) {
	e := make(Errors)
	e.Add("fieldName", &FieldError{Code: "minValue", Expected: 5, Received: 4})

	data, err := json.Marshal(e)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	if !json.Valid(data) {
		t.Errorf("invalid JSON: %s", string(data))
	}

	// Unmarshal and verify structure
	var m map[string]any
	json.Unmarshal(data, &m)

	fieldErrs, ok := m["fieldName"].([]any)
	if !ok || len(fieldErrs) != 1 {
		t.Errorf("fieldName not found or wrong type in JSON")
	}
}
