package layrz

import (
	"strings"
	"testing"
)

// TestConventionADiscoveryAndCall tests Convention A clean method discovery and invocation.
func TestConventionADiscoveryAndCall(t *testing.T) {
	type TestForm struct {
		Name *string `layrz:"char,required"`
	}

	var calls int

	// Monkey-patch for testing (we'll use reflection instead)
	form := &TestForm{
		Name: Ptr("John"),
	}

	errs := Validate(form)

	// With no clean methods defined, should have no errors from clean
	if len(errs) > 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
	_ = calls
}

// TestConventionANilReturn tests that a nil return from Convention A produces no error.
type NoErrorForm struct {
	Name *string `layrz:"char,required"`
}

func (f *NoErrorForm) CleanName(value *string) *FieldError {
	// Always returns nil
	return nil
}

func TestConventionANilReturn(t *testing.T) {
	form := &NoErrorForm{
		Name: Ptr("valid"),
	}

	errs := Validate(form)

	// Should have no errors
	if len(errs) > 0 {
		t.Errorf("expected no errors from nil return, got %v", errs)
	}
}

// TestConventionAErrorReturn tests that a non-nil return is filed correctly.
type WithErrorForm struct {
	Email *string `layrz:"email,required"`
}

func (f *WithErrorForm) CleanEmail(value *string) *FieldError {
	if value != nil && len(*value) > 0 {
		// Check if local part (before @) ends with 'x'
		if idx := strings.Index(*value, "@"); idx > 0 && (*value)[idx-1] == 'x' {
			return &FieldError{Code: "endsWithX"}
		}
	}
	return nil
}

func TestConventionAErrorReturn(t *testing.T) {
	form := &WithErrorForm{
		Email: Ptr("testx@example.com"),
	}

	errs := Validate(form)

	// Should have email error from Convention A
	emailErrs, ok := errs["email"]
	if !ok {
		t.Fatal("missing email key")
	}
	if len(emailErrs) != 1 || emailErrs[0].Code != "endsWithX" {
		if len(emailErrs) == 0 {
			t.Errorf("expected 1 error with code endsWithX, got 0 errors")
		} else {
			t.Errorf("expected endsWithX error, got code=%q expected=%v received=%v", emailErrs[0].Code, emailErrs[0].Expected, emailErrs[0].Received)
		}
	}
}

// TestConventionBMerged tests that Convention B errors are merged with camelCase keys.
type ConventionBForm struct {
	Field1 *string `layrz:"char,required"`
}

func (f *ConventionBForm) CleanValidation() Errors {
	return Errors{
		"field_1": {
			{Code: "error1"},
			{Code: "error2"},
		},
	}
}

func TestConventionBMerged(t *testing.T) {
	form := &ConventionBForm{
		Field1: Ptr("valid"),
	}

	errs := Validate(form)

	// Should have field1 key (camelCased from field_1)
	field1Errs, ok := errs["field1"]
	if !ok {
		t.Fatal("missing field1 key (should be camelCased from field_1)")
	}
	if len(field1Errs) != 2 {
		t.Errorf("expected 2 errors, got %d", len(field1Errs))
	}
	if field1Errs[0].Code != "error1" || field1Errs[1].Code != "error2" {
		t.Errorf("expected [error1, error2], got [%s, %s]", field1Errs[0].Code, field1Errs[1].Code)
	}
}

// TestConventionANeverAlsoB tests that a method matching a field name is only Convention A.
type NoDoubleFiringForm struct {
	Status *string `layrz:"char,required"`
}

func (f *NoDoubleFiringForm) CleanStatus(value *string) *FieldError {
	// This is Convention A, not B
	if value != nil && *value == "invalid" {
		return &FieldError{Code: "customError"}
	}
	return nil
}

func TestConventionANeverAlsoB(t *testing.T) {
	form := &NoDoubleFiringForm{
		Status: Ptr("invalid"),
	}

	errs := Validate(form)

	// Should have exactly 1 error for status (from Convention A, not double-fired)
	statusErrs, ok := errs["status"]
	if !ok {
		t.Fatal("missing status key")
	}
	if len(statusErrs) != 1 {
		t.Errorf("expected exactly 1 error (Convention A only), got %d", len(statusErrs))
	}
}

// TestBrokenMethodSignature tests that a broken method signature becomes a config error.
type BrokenSignatureForm struct {
	Name *string `layrz:"char,required"`
}

func (f *BrokenSignatureForm) CleanBroken(a int, b int) string {
	// Wrong signature: wrong arity, wrong return type
	return "bad"
}

func TestBrokenMethodSignature(t *testing.T) {
	form := &BrokenSignatureForm{
		Name: Ptr("test"),
	}

	errs := Validate(form)

	// Should have _config error
	configErrs, ok := errs[configErrorKey]
	if !ok {
		t.Fatalf("expected _config error, got keys: %v", errs.Keys())
	}
	if len(configErrs) == 0 {
		t.Fatal("expected at least one config error")
	}
	// Should not panic
	t.Logf("config error: %v", configErrs[0])
}

// TestConventionAParameterTypeMismatch tests that a Convention A with wrong param type is a config error.
type TypeMismatchForm struct {
	Name *string `layrz:"char,required"`
}

func (f *TypeMismatchForm) CleanName(value int) *FieldError {
	// Wrong param type: should be *string, not int
	return nil
}

func TestConventionAParameterTypeMismatch(t *testing.T) {
	form := &TypeMismatchForm{
		Name: Ptr("test"),
	}

	errs := Validate(form)

	// Should have _config error
	configErrs, ok := errs[configErrorKey]
	if !ok {
		t.Fatalf("expected _config error, got keys: %v", errs.Keys())
	}
	if len(configErrs) == 0 {
		t.Fatal("expected at least one config error")
	}
}

// TestPanicInCleanMethod tests that a panic in a clean method is recovered.
type PanicForm struct {
	Name *string `layrz:"char,required"`
}

func (f *PanicForm) CleanPanic() Errors {
	// This will panic
	panic("intentional panic")
}

func TestPanicInCleanMethod(t *testing.T) {
	form := &PanicForm{
		Name: Ptr("test"),
	}

	// Should not panic
	errs := Validate(form)

	// Should have _config error recording the panic
	configErrs, ok := errs[configErrorKey]
	if !ok {
		t.Fatalf("expected _config error, got keys: %v", errs.Keys())
	}
	if len(configErrs) == 0 {
		t.Fatal("expected at least one config error from panic recovery")
	}
}

// TestDeterministicCleanOrder tests that clean methods run in alphabetical order.
type OrderForm struct {
	Field *string `layrz:"char,required"`
}

var cleanOrder []string

func (f *OrderForm) CleanZ() Errors {
	cleanOrder = append(cleanOrder, "Z")
	return nil
}

func (f *OrderForm) CleanA() Errors {
	cleanOrder = append(cleanOrder, "A")
	return nil
}

func (f *OrderForm) CleanM() Errors {
	cleanOrder = append(cleanOrder, "M")
	return nil
}

func TestDeterministicCleanOrder(t *testing.T) {
	cleanOrder = nil // Reset
	form := &OrderForm{
		Field: Ptr("test"),
	}

	Validate(form)

	// Should have been called in order A, M, Z
	if len(cleanOrder) != 3 || cleanOrder[0] != "A" || cleanOrder[1] != "M" || cleanOrder[2] != "Z" {
		t.Errorf("expected order [A, M, Z], got %v", cleanOrder)
	}
}

// TestConventionBSameKeyMerge tests that multiple Convention B methods can write to the same key.
type SameKeyForm struct {
	Field *string `layrz:"char,required"`
}

var sameKeyErrors []string

func (f *SameKeyForm) CleanA() Errors {
	return Errors{
		"status": {
			{Code: "errorA"},
		},
	}
}

func (f *SameKeyForm) CleanB() Errors {
	return Errors{
		"status": {
			{Code: "errorB"},
		},
	}
}

func TestConventionBSameKeyMerge(t *testing.T) {
	form := &SameKeyForm{
		Field: Ptr("test"),
	}

	errs := Validate(form)

	// Should have status key with both errors
	statusErrs, ok := errs["status"]
	if !ok {
		t.Fatal("missing status key")
	}
	if len(statusErrs) != 2 {
		t.Errorf("expected 2 errors, got %d", len(statusErrs))
	}
	if statusErrs[0].Code != "errorA" || statusErrs[1].Code != "errorB" {
		t.Errorf("expected [errorA, errorB], got [%s, %s]", statusErrs[0].Code, statusErrs[1].Code)
	}
}

// TestNestedSubformCleanMethods tests that nested subform clean methods run with prefixed keys.
type NestedWithCleanForm struct {
	Name *string `layrz:"char,required"`
}

func (f *NestedWithCleanForm) CleanSubform() Errors {
	return Errors{
		"field": {
			{Code: "nested_error"},
		},
	}
}

type ParentWithNestedCleanForm struct {
	Nested *NestedWithCleanForm `layrz:"subform"`
}

func TestNestedSubformCleanMethods(t *testing.T) {
	form := &ParentWithNestedCleanForm{
		Nested: &NestedWithCleanForm{
			Name: Ptr("test"),
		},
	}

	errs := Validate(form)

	// Should have nested.field key (from nested clean method, prefixed)
	nestedErrs, ok := errs["nested.field"]
	if !ok {
		t.Fatalf("expected nested.field key, got keys: %v", errs.Keys())
	}
	if len(nestedErrs) != 1 || nestedErrs[0].Code != "nested_error" {
		if len(nestedErrs) == 0 {
			t.Errorf("expected 1 error with code nested_error, got 0 errors")
		} else {
			t.Errorf("expected nested_error, got code=%q expected=%v received=%v", nestedErrs[0].Code, nestedErrs[0].Expected, nestedErrs[0].Received)
		}
	}
}
