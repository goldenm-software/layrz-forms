package layrz

import (
	"encoding/json"
	"strings"
	"testing"
)

// Address is a subform used in the example.
type Address struct {
	StreetName *string `layrz:"char,required,min_length=5"`
	ZipCode    *string `layrz:"char,required"`
}

// Item is a subform in the subform list.
type Item struct {
	Name  *string  `layrz:"char,required"`
	Price *float64 `layrz:"number,required,min_value=0"`
}

// ExampleForm demonstrates all field types, nested forms, and clean methods.
type ExampleForm struct {
	IDTest        *int            `layrz:"id,required"`
	EmailText     *string         `layrz:"email,required"`
	JSONListTest  *[]any          `layrz:"json,required,datatype=list"`
	JSONDictTest  *map[string]any `layrz:"json,required,datatype=dict"`
	IntTest       *int            `layrz:"number,required,min_value=0,max_value=5"`
	FloatTest     *float64        `layrz:"number,required,min_value=0,max_value=5"`
	BoolTest      *bool           `layrz:"bool,required"`
	PlainTextTest *string         `layrz:"char,required"`
	EmptyTextTest *string         `layrz:"char,required,empty"`
	RangeTextTest *string         `layrz:"char,required,min_length=5,max_length=10"`
	Address       *Address        `layrz:"subform"`
	Items         []*Item         `layrz:"subform_list"`
}

// Convention B clean methods
func (f *ExampleForm) CleanFunc1() Errors {
	return Errors{
		"clean1": {
			{Code: codeError1},
			{Code: "error2"},
		},
	}
}

func (f *ExampleForm) CleanFunc2() Errors {
	return Errors{
		"clean2": {
			{Code: codeError1},
		},
	}
}

// Convention A clean method
func (f *ExampleForm) CleanEmailText(value *string) *FieldError {
	// Example: reject blocked domain
	if value != nil && len(*value) > 0 {
		if len(*value) > 13 && (*value)[len(*value)-13:] == "@blocked.com" {
			return &FieldError{Code: "blockedDomain"}
		}
	}
	return nil
}

// TestExampleForm validates the ExampleForm as shown in the README.
func TestExampleForm(t *testing.T) {
	// Build the input
	idVal := 1
	emailVal := "example@goldenmcorp.com"
	jsonListVal := []any{"hola mundo"}
	jsonDictVal := map[string]any{"hola": "mundo"}
	intVal := 5
	floatVal := 4.5
	boolVal := true
	plainTextVal := "hola mundo"
	emptyTextVal := "hola"
	rangeTextVal := "hola" // 4 characters, less than min_length=5

	form := &ExampleForm{
		IDTest:        &idVal,
		EmailText:     &emailVal,
		JSONListTest:  &jsonListVal,
		JSONDictTest:  &jsonDictVal,
		IntTest:       &intVal,
		FloatTest:     &floatVal,
		BoolTest:      &boolVal,
		PlainTextTest: &plainTextVal,
		EmptyTextTest: &emptyTextVal,
		RangeTextTest: &rangeTextVal,
		// Address and Items are nil
	}

	// Validate
	errs := Validate(form)

	// Assert IsValid() is false
	if IsValid(form) {
		t.Error("expected IsValid() to be false")
	}

	// Assert keys are exactly what we expect
	expectedKeys := []string{"clean1", "clean2", "rangeTextTest"}
	actualKeys := errs.Keys()
	if len(actualKeys) != len(expectedKeys) {
		t.Errorf("expected %d keys, got %d: %v", len(expectedKeys), len(actualKeys), actualKeys)
	}
	for _, key := range expectedKeys {
		if _, ok := errs[key]; !ok {
			t.Errorf("missing expected key: %q", key)
		}
	}

	// Assert rangeTextTest has the right error
	rangeErrs, ok := errs["rangeTextTest"]
	if !ok {
		t.Fatal("missing rangeTextTest key")
	}
	if len(rangeErrs) != 1 {
		t.Errorf("expected 1 error for rangeTextTest, got %d", len(rangeErrs))
	}
	err := rangeErrs[0]
	if err.Code != codeMinLength {
		t.Errorf("expected code minLength, got %q", err.Code)
	}
	if err.Expected != 5 {
		t.Errorf("expected Expected=5, got %v", err.Expected)
	}
	if err.Received != 4 {
		t.Errorf("expected Received=4, got %v", err.Received)
	}

	// Assert clean1 has two errors
	clean1Errs, ok := errs["clean1"]
	if !ok {
		t.Fatal("missing clean1 key")
	}
	if len(clean1Errs) != 2 {
		t.Errorf("expected 2 errors for clean1, got %d", len(clean1Errs))
	}
	if clean1Errs[0].Code != codeError1 || clean1Errs[1].Code != "error2" {
		t.Errorf("clean1 error codes mismatch: expected [error1, error2], got [%s, %s]", clean1Errs[0].Code, clean1Errs[1].Code)
	}

	// Assert clean2 has one error
	clean2Errs, ok := errs["clean2"]
	if !ok {
		t.Fatal("missing clean2 key")
	}
	if len(clean2Errs) != 1 {
		t.Errorf("expected 1 error for clean2, got %d", len(clean2Errs))
	}
	if clean2Errs[0].Code != codeError1 {
		t.Errorf("expected clean2 code error1, got %q", clean2Errs[0].Code)
	}

	// Assert JSON serialization (prove wire-compat with Python)
	data, marshalErr := json.Marshal(errs)
	if marshalErr != nil {
		t.Fatalf("failed to marshal errors: %v", marshalErr)
	}

	var unmarshalled map[string][]map[string]any
	if unmarshalErr := json.Unmarshal(data, &unmarshalled); unmarshalErr != nil {
		t.Fatalf("failed to unmarshal errors: %v", unmarshalErr)
	}

	// Check rangeTextTest error serialization
	rangeEntry, ok := unmarshalled["rangeTextTest"]
	if !ok {
		t.Fatal("missing rangeTextTest in JSON")
	}
	if len(rangeEntry) != 1 {
		t.Errorf("expected 1 error entry, got %d", len(rangeEntry))
	}
	rangeErr := rangeEntry[0]
	if code, ok := rangeErr["code"]; !ok || code != codeMinLength {
		t.Errorf("expected code=minLength, got code=%v", code)
	}
	if expected, ok := rangeErr["expected"]; !ok || expected != float64(5) {
		// JSON numbers come back as float64
		t.Errorf("expected expected=5, got expected=%v (%T)", expected, expected)
	}
	if received, ok := rangeErr["received"]; !ok || received != float64(4) {
		t.Errorf("expected received=4, got received=%v (%T)", received, received)
	}
	// Assert no "extra" key and no "expected"/"received" with null values
	if _, ok := rangeErr["extra"]; ok && rangeErr["extra"] != nil {
		t.Errorf("expected no extra field, got %v", rangeErr["extra"])
	}
}

// TestExampleFormWithConventionA tests Convention A clean methods.
type EmailBlockerForm struct {
	EmailText *string `layrz:"email,required"`
}

func (f *EmailBlockerForm) CleanEmailText(value *string) *FieldError {
	if value != nil && len(*value) >= 12 && strings.HasSuffix(*value, "@blocked.com") {
		return &FieldError{Code: "blockedDomain"}
	}
	return nil
}

func TestConventionACleanMethod(t *testing.T) {
	// Test that Convention A correctly validates and rejects based on field value
	blockedEmail := "user@blocked.com"
	form := &EmailBlockerForm{
		EmailText: &blockedEmail,
	}

	// First, manually call the clean method to verify it works
	manualErr := form.CleanEmailText(form.EmailText)
	if manualErr == nil {
		t.Fatal("expected CleanEmailText to return an error, got nil")
	}
	if manualErr.Code != "blockedDomain" {
		t.Errorf("expected blockedDomain, got %q", manualErr.Code)
	}

	// Now test through validation
	errs := Validate(form)

	// Debug: print all errors
	t.Logf("All errors: %v", errs.Keys())
	for key, errList := range errs {
		for _, err := range errList {
			t.Logf("  %s: %v", key, err.Code)
		}
	}

	// Should have emailText error from Convention A
	emailErrs, ok := errs["emailText"]
	if !ok {
		t.Fatalf("missing emailText key; got keys: %v", errs.Keys())
	}
	if len(emailErrs) != 1 {
		t.Errorf("expected 1 error, got %d", len(emailErrs))
	}
	if emailErrs[0].Code != "blockedDomain" {
		t.Errorf("expected blockedDomain, got %q", emailErrs[0].Code)
	}

	// Test that a valid email passes
	validEmail := "user@gmail.com"
	form.EmailText = &validEmail
	errs = Validate(form)
	if len(errs) > 0 {
		t.Errorf("expected no errors for valid email, got %v", errs)
	}
}

// TestNestedSubform tests nested subform validation.
func TestNestedSubform(t *testing.T) {
	form := &ExampleForm{
		IDTest:    Ptr(1),
		EmailText: Ptr("test@example.com"),
		BoolTest:  Ptr(true),
	}

	// Add a non-nil Address with too-short StreetName
	streetName := "abc" // 3 chars, min_length=5
	form.Address = &Address{
		StreetName: &streetName,
		ZipCode:    Ptr("12345"),
	}

	errs := Validate(form)

	// Should have address.streetName error
	addrErrs, ok := errs["address.streetName"]
	if !ok {
		t.Fatal("missing address.streetName key")
	}
	if len(addrErrs) != 1 {
		t.Errorf("expected 1 error, got %d", len(addrErrs))
	}
	if addrErrs[0].Code != codeMinLength {
		t.Errorf("expected minLength, got %q", addrErrs[0].Code)
	}
}

// TestNestedSubformList tests subform_list validation.
func TestNestedSubformList(t *testing.T) {
	form := &ExampleForm{
		IDTest:    Ptr(1),
		EmailText: Ptr("test@example.com"),
		BoolTest:  Ptr(true),
	}

	// Add items: first with missing Name, second with negative Price
	form.Items = []*Item{
		{
			Name:  nil, // Missing required field
			Price: Ptr(10.0),
		},
		{
			Name:  Ptr("Item 2"),
			Price: Ptr(-5.0), // Negative, violates min_value=0
		},
	}

	errs := Validate(form)

	// Should have items.0.name and items.1.price errors
	item0NameErrs, ok := errs["items.0.name"]
	if !ok {
		t.Fatal("missing items.0.name key")
	}
	if len(item0NameErrs) != 1 || item0NameErrs[0].Code != codeRequired {
		t.Errorf("expected required error for items.0.name, got %v", item0NameErrs)
	}

	item1PriceErrs, ok := errs["items.1.price"]
	if !ok {
		t.Fatal("missing items.1.price key")
	}
	if len(item1PriceErrs) != 1 || item1PriceErrs[0].Code != codeMinValue {
		t.Errorf("expected minValue error for items.1.price, got %v", item1PriceErrs)
	}
}

// TestNilSubform tests that a nil subform is skipped (no errors).
func TestNilSubform(t *testing.T) {
	form := &ExampleForm{
		IDTest:    Ptr(1),
		EmailText: Ptr("test@example.com"),
		BoolTest:  Ptr(true),
		Address:   nil, // Nil subform should be skipped
	}

	errs := Validate(form)

	// Should not have any address.* errors
	for key := range errs {
		if len(key) > 8 && key[:7] == "address" {
			t.Errorf("unexpected error for nil subform: %q", key)
		}
	}
}

// TestDeeplyNestedForms tests 3-level nesting to verify prefix composition.
type Level2 struct {
	Name *string `layrz:"char,required"`
}

type Level1 struct {
	L2 *Level2 `layrz:"subform"`
}

type Level0 struct {
	L1 *Level1 `layrz:"subform"`
}

func TestDeeplyNestedForms(t *testing.T) {
	form := &Level0{
		L1: &Level1{
			L2: &Level2{
				Name: nil, // Missing required field
			},
		},
	}

	errs := Validate(form)

	// Should have l1.l2.name error
	if _, ok := errs["l1.l2.name"]; !ok {
		t.Errorf("missing l1.l2.name key; got keys: %v", errs.Keys())
	}
}
