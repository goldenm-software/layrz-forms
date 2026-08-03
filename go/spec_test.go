package layrz

import (
	"reflect"
	"testing"
)

func TestCheckType(t *testing.T) {
	t.Run("id kind", func(t *testing.T) {
		tests := []struct {
			name    string
			t       reflect.Type
			wantErr bool
		}{
			{"*int", reflect.TypeOf(Ptr(0)), false},
			{"*int64", reflect.TypeOf(Ptr(int64(0))), false},
			{"*string", reflect.TypeOf(Ptr("")), false},
			{"int (non-ptr)", reflect.TypeOf(0), true},
			{"*float64", reflect.TypeOf(Ptr(0.0)), true},
		}

		for _, tt := range tests {
			spec := &FieldSpec{Kind: KindID}
			err := spec.CheckType(tt.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckType(KindID, %v): wantErr=%v, got error=%v", tt.t, tt.wantErr, err)
			}
		}
	})

	t.Run("email kind", func(t *testing.T) {
		spec := &FieldSpec{Kind: KindEmail}

		err := spec.CheckType(reflect.TypeOf(Ptr("")))
		if err != nil {
			t.Errorf("CheckType(KindEmail, *string): expected no error, got %v", err)
		}

		err = spec.CheckType(reflect.TypeOf(Ptr(0)))
		if err == nil {
			t.Errorf("CheckType(KindEmail, *int): expected error, got nil")
		}
	})

	t.Run("uuid kind", func(t *testing.T) {
		spec := &FieldSpec{Kind: KindUUID}

		err := spec.CheckType(reflect.TypeOf(Ptr("")))
		if err != nil {
			t.Errorf("CheckType(KindUUID, *string): expected no error, got %v", err)
		}

		err = spec.CheckType(reflect.TypeOf(Ptr(0)))
		if err == nil {
			t.Errorf("CheckType(KindUUID, *int): expected error, got nil")
		}
	})

	t.Run("char kind", func(t *testing.T) {
		spec := &FieldSpec{Kind: KindChar}

		err := spec.CheckType(reflect.TypeOf(Ptr("")))
		if err != nil {
			t.Errorf("CheckType(KindChar, *string): expected no error, got %v", err)
		}

		err = spec.CheckType(reflect.TypeOf(Ptr(0)))
		if err == nil {
			t.Errorf("CheckType(KindChar, *int): expected error, got nil")
		}
	})

	t.Run("number kind with int inference", func(t *testing.T) {
		spec := &FieldSpec{Kind: KindNumber}

		err := spec.CheckType(reflect.TypeOf(Ptr(int32(0))))
		if err != nil {
			t.Fatalf("CheckType(KindNumber, *int32): expected no error, got %v", err)
		}

		if spec.Datatype != "int" {
			t.Errorf("expected Datatype=int, got %q", spec.Datatype)
		}
	})

	t.Run("number kind with float inference", func(t *testing.T) {
		spec := &FieldSpec{Kind: KindNumber}

		err := spec.CheckType(reflect.TypeOf(Ptr(0.0)))
		if err != nil {
			t.Fatalf("CheckType(KindNumber, *float64): expected no error, got %v", err)
		}

		if spec.Datatype != "float" {
			t.Errorf("expected Datatype=float, got %q", spec.Datatype)
		}
	})

	t.Run("number kind with explicit datatype mismatch", func(t *testing.T) {
		spec := &FieldSpec{Kind: KindNumber, Datatype: "float"}

		err := spec.CheckType(reflect.TypeOf(Ptr(int32(0))))
		if err == nil {
			t.Errorf("CheckType with datatype=float and *int32: expected error, got nil")
		}
	})

	t.Run("bool kind", func(t *testing.T) {
		spec := &FieldSpec{Kind: KindBool}

		err := spec.CheckType(reflect.TypeOf(Ptr(false)))
		if err != nil {
			t.Errorf("CheckType(KindBool, *bool): expected no error, got %v", err)
		}

		err = spec.CheckType(reflect.TypeOf(Ptr("")))
		if err == nil {
			t.Errorf("CheckType(KindBool, *string): expected error, got nil")
		}
	})

	t.Run("json kind with slice (list inference)", func(t *testing.T) {
		spec := &FieldSpec{Kind: KindJSON}
		var val *[]any

		err := spec.CheckType(reflect.TypeOf(val))
		if err != nil {
			t.Fatalf("CheckType(KindJSON, *[]any): expected no error, got %v", err)
		}

		if spec.Datatype != "list" {
			t.Errorf("expected Datatype=list, got %q", spec.Datatype)
		}
	})

	t.Run("json kind with map (dict inference)", func(t *testing.T) {
		spec := &FieldSpec{Kind: KindJSON}
		var val *map[string]any

		err := spec.CheckType(reflect.TypeOf(val))
		if err != nil {
			t.Fatalf("CheckType(KindJSON, *map[string]any): expected no error, got %v", err)
		}

		if spec.Datatype != "dict" {
			t.Errorf("expected Datatype=dict, got %q", spec.Datatype)
		}
	})

	t.Run("json kind with datatype mismatch", func(t *testing.T) {
		spec := &FieldSpec{Kind: KindJSON, Datatype: "dict"}
		var val *[]any

		err := spec.CheckType(reflect.TypeOf(val))
		if err == nil {
			t.Errorf("CheckType with datatype=dict and *[]any: expected error, got nil")
		}
	})

	t.Run("subform kind", func(t *testing.T) {
		type TestForm struct{}
		spec := &FieldSpec{Kind: KindSubform}
		var val *TestForm

		err := spec.CheckType(reflect.TypeOf(val))
		if err != nil {
			t.Errorf("CheckType(KindSubform, *TestForm): expected no error, got %v", err)
		}

		err = spec.CheckType(reflect.TypeOf(Ptr("not a struct")))
		if err == nil {
			t.Errorf("CheckType(KindSubform, *string): expected error, got nil")
		}
	})

	t.Run("subform_list kind", func(t *testing.T) {
		type TestForm struct{}
		spec := &FieldSpec{Kind: KindSubformList}

		// []TestForm
		var val1 []TestForm
		err := spec.CheckType(reflect.TypeOf(val1))
		if err != nil {
			t.Errorf("CheckType(KindSubformList, []TestForm): expected no error, got %v", err)
		}

		// []*TestForm
		var val2 []*TestForm
		err = spec.CheckType(reflect.TypeOf(val2))
		if err != nil {
			t.Errorf("CheckType(KindSubformList, []*TestForm): expected no error, got %v", err)
		}

		// Not a slice
		err = spec.CheckType(reflect.TypeOf(Ptr(TestForm{})))
		if err == nil {
			t.Errorf("CheckType(KindSubformList, *TestForm): expected error, got nil")
		}
	})
}

func TestFieldSpecRules(t *testing.T) {
	t.Run("id rules", func(t *testing.T) {
		spec := &FieldSpec{Kind: KindID, Required: true}
		rules, err := spec.Rules()
		if err != nil {
			t.Fatalf("Rules() error: %v", err)
		}

		idRules, ok := rules.(IDRules)
		if !ok {
			t.Fatalf("expected IDRules, got %T", rules)
		}
		if !idRules.Required {
			t.Errorf("expected Required=true")
		}
	})

	t.Run("email rules with default regex", func(t *testing.T) {
		spec := &FieldSpec{Kind: KindEmail, Empty: false}
		rules, err := spec.Rules()
		if err != nil {
			t.Fatalf("Rules() error: %v", err)
		}

		emailRules, ok := rules.(EmailRules)
		if !ok {
			t.Fatalf("expected EmailRules, got %T", rules)
		}
		if emailRules.Regex != DefaultEmailRegex {
			t.Errorf("expected default regex, got %q", emailRules.Regex)
		}
	})

	t.Run("email rules with custom regex", func(t *testing.T) {
		spec := &FieldSpec{Kind: KindEmail, Regex: "^[a-z]+@[a-z]+$"}
		rules, err := spec.Rules()
		if err != nil {
			t.Fatalf("Rules() error: %v", err)
		}

		emailRules := rules.(EmailRules)
		if emailRules.Regex != "^[a-z]+@[a-z]+$" {
			t.Errorf("expected custom regex")
		}
	})
}
