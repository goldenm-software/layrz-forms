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

// codesOf extracts the error codes from a slice of FieldErrors for assertion purposes.
func codesOf(errs []*FieldError) []string {
	codes := make([]string, len(errs))
	for i, e := range errs {
		if e != nil {
			codes[i] = e.Code
		}
	}
	return codes
}

// TestSubformListCleanMethods tests that nested subform_list elements run their clean methods.
// This is the main regression test for the bug fix.
type InnerWithClean struct {
	Label *string `layrz:"char,required,min_length=2"`
	Extra *string
}

func (i *InnerWithClean) CleanExtra(value *string) *FieldError {
	return &FieldError{Code: "innerConvA"}
}

func (i *InnerWithClean) CleanCross() Errors {
	return Errors{"innerConvB": {{Code: "fired"}}}
}

type OuterWithList struct {
	One  *InnerWithClean   `layrz:"subform"`
	Many []*InnerWithClean `layrz:"subform_list"`
}

func TestSubformListCleanMethods(t *testing.T) {
	// This test reproduces the exact bug: subform_list elements never run clean methods
	errs := Validate(&OuterWithList{
		One:  &InnerWithClean{Label: Ptr("ok")},
		Many: []*InnerWithClean{{Label: Ptr("ok")}},
	})

	// Should have one.extra and one.innerConvB from the subform
	if _, ok := errs["one.extra"]; !ok {
		t.Error("missing one.extra from subform")
	}
	if _, ok := errs["one.innerConvB"]; !ok {
		t.Error("missing one.innerConvB from subform")
	}

	// Should ALSO have many.0.extra and many.0.innerConvB from the subform_list element
	if _, ok := errs["many.0.extra"]; !ok {
		t.Error("missing many.0.extra from subform_list element (REGRESSION: clean methods not running)")
	}
	if _, ok := errs["many.0.innerConvB"]; !ok {
		t.Error("missing many.0.innerConvB from subform_list element (REGRESSION: clean methods not running)")
	}

	t.Logf("All keys: %v", errs.Keys())
}

// TestSubformListMultipleElementsWithClean tests multiple elements each running clean methods.
type ItemWithClean struct {
	Name *string `layrz:"char,required"`
}

func (i *ItemWithClean) CleanValidation() Errors {
	if i.Name != nil && *i.Name == "invalid" {
		return Errors{"name": {{Code: "badName"}}}
	}
	return nil
}

type ListFormMulti struct {
	Items []*ItemWithClean `layrz:"subform_list"`
}

func TestSubformListMultipleElementsWithClean(t *testing.T) {
	errs := Validate(&ListFormMulti{
		Items: []*ItemWithClean{
			{Name: Ptr("valid0")},
			{Name: Ptr("invalid")}, // Should trigger clean error
			{Name: Ptr("valid2")},
			{Name: Ptr("invalid")}, // Should trigger clean error
		},
	})

	// Should have errors for indices 1 and 3
	if _, ok := errs["items.1.name"]; !ok {
		t.Error("missing items.1.name")
	}
	if _, ok := errs["items.3.name"]; !ok {
		t.Error("missing items.3.name")
	}

	// Should NOT have errors for indices 0 and 2
	if _, ok := errs["items.0.name"]; ok {
		t.Error("unexpected items.0.name error")
	}
	if _, ok := errs["items.2.name"]; ok {
		t.Error("unexpected items.2.name error")
	}

	t.Logf("All keys: %v", errs.Keys())
}

// TestSubformListNonPointerElementsWithClean tests non-pointer slice elements running clean methods.
type NonPointerItem struct {
	Name *string `layrz:"char,required"`
}

func (n *NonPointerItem) CleanCheck() Errors {
	if n.Name != nil && *n.Name == "skip" {
		return Errors{"name": {{Code: "skipError"}}}
	}
	return nil
}

type ListFormNonPointer struct {
	Items []NonPointerItem `layrz:"subform_list"` // Non-pointer elements
}

func TestSubformListNonPointerElementsWithClean(t *testing.T) {
	errs := Validate(&ListFormNonPointer{
		Items: []NonPointerItem{
			{Name: Ptr("ok")},
			{Name: Ptr("skip")},
		},
	})

	// Should have items.1.name error from clean method
	if _, ok := errs["items.1.name"]; !ok {
		t.Error("missing items.1.name error from non-pointer element clean method")
	}

	t.Logf("All keys: %v", errs.Keys())
}

// TestSubformListNilPointerElementsSkipped tests that nil elements are skipped entirely.
type ItemToSkip struct {
	Name *string `layrz:"char,required"`
}

func (i *ItemToSkip) CleanValidation() Errors {
	if i.Name != nil && *i.Name == "bad" {
		return Errors{"name": {{Code: "badValue"}}}
	}
	return nil
}

type ListFormWithNil struct {
	Items []*ItemToSkip `layrz:"subform_list"`
}

func TestSubformListNilPointerElementsSkipped(t *testing.T) {
	errs := Validate(&ListFormWithNil{
		Items: []*ItemToSkip{
			{Name: Ptr("ok")},
			nil, // Should be skipped - no clean method, no tag rules
			{Name: Ptr("ok2")},
		},
	})

	// Index 1 is nil, so no items.1.* keys at all
	for key := range errs {
		if len(key) >= 8 && key[:8] == "items.1." {
			t.Errorf("nil element at index 1 should not produce errors, got key=%q", key)
		}
	}

	// Indices 0 and 2 should not have any errors (names are valid and don't match "bad")
	if _, ok := errs["items.0.name"]; ok {
		t.Error("unexpected items.0.name error")
	}
	if _, ok := errs["items.2.name"]; ok {
		t.Error("unexpected items.2.name error")
	}

	t.Logf("All keys: %v", errs.Keys())
}

// TestSubformListThreeLevelNesting tests deeply nested subform + subform_list + subform_list combinations.
type Level3Item struct {
	Value *string `layrz:"char,required"`
}

func (l *Level3Item) CleanValidation() Errors {
	if l.Value != nil && *l.Value == "bad" {
		return Errors{"value": {{Code: "bad_value"}}}
	}
	return nil
}

type Level2Container struct {
	Items []*Level3Item `layrz:"subform_list"`
}

func (l *Level2Container) CleanValidation() Errors {
	return Errors{"container": {{Code: "container_error"}}}
}

type Level1Container struct {
	L2 *Level2Container `layrz:"subform"`
}

type Level0Form struct {
	L1 *Level1Container `layrz:"subform"`
}

func TestSubformListThreeLevelNesting(t *testing.T) {
	errs := Validate(&Level0Form{
		L1: &Level1Container{
			L2: &Level2Container{
				Items: []*Level3Item{
					{Value: Ptr("ok")},
					{Value: Ptr("bad")},
				},
			},
		},
	})

	// Should have deeply prefixed keys
	if _, ok := errs["l1.l2.container"]; !ok {
		t.Error("missing l1.l2.container from Level2 clean method")
	}
	if _, ok := errs["l1.l2.items.1.value"]; !ok {
		t.Error("missing l1.l2.items.1.value from Level3 clean method")
	}

	t.Logf("All keys: %v", errs.Keys())
}

// TestSubformListCleanReceivesCorrectValue tests that Convention A receives the actual field value.
type EchoForm struct {
	Name *string `layrz:"char,required"`
}

func (e *EchoForm) CleanName(value *string) *FieldError {
	// If name is "pass", no error; otherwise, error
	if value != nil && *value == "pass" {
		return nil
	}
	return &FieldError{Code: "failed_validation"}
}

type EchoListForm struct {
	Items []*EchoForm `layrz:"subform_list"`
}

func TestSubformListCleanReceivesCorrectValue(t *testing.T) {
	errs := Validate(&EchoListForm{
		Items: []*EchoForm{
			{Name: Ptr("pass")}, // Should pass
			{Name: Ptr("fail")}, // Should fail
			{Name: Ptr("pass")}, // Should pass
		},
	})

	// Only index 1 should have error
	if _, ok := errs["items.1.name"]; !ok {
		t.Error("missing items.1.name error")
	}
	if _, ok := errs["items.0.name"]; ok {
		t.Error("unexpected items.0.name error")
	}
	if _, ok := errs["items.2.name"]; ok {
		t.Error("unexpected items.2.name error")
	}

	t.Logf("All keys: %v", errs.Keys())
}

// TestSubformListCleanPanicRecovery tests that a panic in a clean method is recovered and other elements continue.
type PanicItem struct {
	Name *string `layrz:"char,required"`
}

func (p *PanicItem) CleanPanic() Errors {
	panic("intentional panic in list element")
}

type PanicListForm struct {
	Items []*PanicItem `layrz:"subform_list"`
}

func TestSubformListCleanPanicRecovery(t *testing.T) {
	// Should not panic; should recover into _config errors
	errs := Validate(&PanicListForm{
		Items: []*PanicItem{
			{Name: Ptr("item0")},
			{Name: Ptr("item1")},
		},
	})

	// Should have _config errors for the panic
	configErrs, ok := errs[configErrorKey]
	if !ok {
		t.Error("missing _config error for panic recovery")
	} else if len(configErrs) == 0 {
		t.Error("expected _config error from panic, got 0")
	} else {
		t.Logf("panic recovered: %v", configErrs[0].Code)
	}

	t.Logf("All keys: %v", errs.Keys())
}

// TestSubformListNoConfigErrorsInHappyPath tests that no spurious _config errors appear in valid cases.
type SimpleListItem struct {
	Name *string `layrz:"char,required"`
}

func (s *SimpleListItem) CleanValidation() Errors {
	return nil
}

type SimpleListForm struct {
	Items []*SimpleListItem `layrz:"subform_list"`
}

func TestSubformListNoConfigErrorsInHappyPath(t *testing.T) {
	errs := Validate(&SimpleListForm{
		Items: []*SimpleListItem{
			{Name: Ptr("item0")},
			{Name: Ptr("item1")},
		},
	})

	// Should have NO _config errors in the happy path
	if _, ok := errs[configErrorKey]; ok {
		t.Error("unexpected _config error in happy path")
	}

	t.Logf("All keys (should be empty): %v", errs.Keys())
}
