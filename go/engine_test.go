package layrz

import (
	"testing"
)

// TestNilFormInput tests that a nil form produces a config error.
func TestNilFormInput(t *testing.T) {
	errs := Validate(nil)

	configErrs, ok := errs[configErrorKey]
	if !ok {
		t.Fatal("expected _config error for nil form")
	}
	if len(configErrs) == 0 {
		t.Fatal("expected at least one config error")
	}
	if configErrs[0].Code != "internalError" {
		t.Errorf("expected internalError code, got %q", configErrs[0].Code)
	}
}

// TestNonPointerFormInput tests that a non-pointer form produces a config error.
func TestNonPointerFormInput(t *testing.T) {
	type TestForm struct {
		Name *string `layrz:"char,required"`
	}

	form := TestForm{
		Name: Ptr("test"),
	}

	// Pass by value (not a pointer)
	errs := Validate(form)

	configErrs, ok := errs[configErrorKey]
	if !ok {
		t.Fatalf("expected _config error for non-pointer form, got keys: %v", errs.Keys())
	}
	if len(configErrs) == 0 {
		t.Fatal("expected at least one config error")
	}
	// Error message should mention pointer receiver
	if configErrs[0].Code != "internalError" {
		t.Errorf("expected internalError code, got %q", configErrs[0].Code)
	}
	t.Logf("error message: %v", configErrs[0].Extra["reason"])
}

// TestNilPointerInput tests that a nil pointer produces a config error.
func TestNilPointerInput(t *testing.T) {
	type TestForm struct {
		Name *string `layrz:"char,required"`
	}

	var form *TestForm = nil

	errs := Validate(form)

	configErrs, ok := errs[configErrorKey]
	if !ok {
		t.Fatal("expected _config error for nil pointer")
	}
	if len(configErrs) == 0 {
		t.Fatal("expected at least one config error")
	}
}

// TestPointerToNonStructInput tests that a pointer to non-struct produces a config error.
func TestPointerToNonStructInput(t *testing.T) {
	val := "not a struct"
	errs := Validate(&val)

	configErrs, ok := errs[configErrorKey]
	if !ok {
		t.Fatal("expected _config error for pointer to non-struct")
	}
	if len(configErrs) == 0 {
		t.Fatal("expected at least one config error")
	}
}

// TestMalformedTag tests that a field with a malformed tag produces a config error.
func TestMalformedTag(t *testing.T) {
	type BadForm struct {
		Field *string `layrz:"invalid_kind"`
	}

	form := &BadForm{
		Field: Ptr("test"),
	}

	errs := Validate(form)

	configErrs, ok := errs[configErrorKey]
	if !ok {
		t.Fatalf("expected _config error for malformed tag, got keys: %v", errs.Keys())
	}
	if len(configErrs) == 0 {
		t.Fatal("expected at least one config error")
	}
}

// TestTagTypeMismatch tests that a field with a tag/type mismatch produces a config error.
func TestTagTypeMismatch(t *testing.T) {
	type MismatchForm struct {
		NumField *int `layrz:"char,required"` // char expects *string, not *int
	}

	form := &MismatchForm{
		NumField: Ptr(42),
	}

	errs := Validate(form)

	configErrs, ok := errs[configErrorKey]
	if !ok {
		t.Fatalf("expected _config error for tag/type mismatch, got keys: %v", errs.Keys())
	}
	if len(configErrs) == 0 {
		t.Fatal("expected at least one config error")
	}
}

// TestUntaggedFieldsSkipped tests that untagged fields are silently skipped.
func TestUntaggedFieldsSkipped(t *testing.T) {
	type SkipForm struct {
		TaggedField   *string `layrz:"char,required"`
		UntaggedField *string // No tag
	}

	form := &SkipForm{
		TaggedField:   Ptr("valid"),
		UntaggedField: nil, // Missing but should be skipped
	}

	errs := Validate(form)

	// Should have no errors (untaggedField is skipped)
	if len(errs) > 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

// TestSkipTagField tests that layrz:"-" fields are skipped.
func TestSkipTagField(t *testing.T) {
	type SkipDashForm struct {
		Field1 *string `layrz:"char,required"`
		Field2 *string `layrz:"-"` // Should skip
	}

	form := &SkipDashForm{
		Field1: Ptr("valid"),
		Field2: nil, // Should not be validated
	}

	errs := Validate(form)

	// Should have no errors (field2 is skipped)
	if len(errs) > 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

// TestUnexportedFieldsSkipped tests that unexported fields are skipped.
func TestUnexportedFieldsSkipped(t *testing.T) {
	type UnexportedForm struct {
		PublicField  *string `layrz:"char,required"`
		privateField *string `layrz:"char,required"` // Unexported, should skip
	}

	form := &UnexportedForm{
		PublicField:  Ptr("valid"),
		privateField: nil, // Should not cause error
	}

	errs := Validate(form)

	// Should have no errors (privateField is unexported and skipped)
	if len(errs) > 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

// TestEmbeddedAnonymousStruct tests that untagged anonymous structs promote their fields.
type EmbeddedAddress struct {
	StreetName *string `layrz:"char,required"`
}

type EmbeddingForm struct {
	EmbeddedAddress // Anonymous, no tag -> promote fields
}

func TestEmbeddedAnonymousStruct(t *testing.T) {
	form := &EmbeddingForm{
		EmbeddedAddress: EmbeddedAddress{
			StreetName: nil, // Missing required
		},
	}

	errs := Validate(form)

	// Should have streetName error (promoted from embedded struct)
	if _, ok := errs["streetName"]; !ok {
		t.Fatalf("expected streetName key (promoted field), got keys: %v", errs.Keys())
	}
}

// TestRecursionDepthGuard tests that deep recursion is guarded.
type SelfReferentialForm struct {
	Name  *string              `layrz:"char,required"`
	Child *SelfReferentialForm `layrz:"subform"`
}

func TestRecursionDepthGuard(t *testing.T) {
	// Create a deeply nested structure
	var form *SelfReferentialForm
	current := &SelfReferentialForm{Name: Ptr("test")}
	form = current

	// Create a chain exceeding maxDepth
	for i := 0; i < maxDepth+10; i++ {
		child := &SelfReferentialForm{Name: Ptr("test")}
		current.Child = child
		current = child
	}

	errs := Validate(form)

	// Should have _config error for depth exceeded
	configErrs, ok := errs[configErrorKey]
	if !ok {
		t.Fatalf("expected _config error for depth guard, got keys: %v", errs.Keys())
	}
	if len(configErrs) == 0 {
		t.Fatal("expected at least one config error from depth guard")
	}
	// Error should mention max depth
	found := false
	for _, err := range configErrs {
		if reason, ok := err.Extra["reason"].(string); ok && len(reason) > 0 {
			if reason[0:len("maximum")] == "maximum" {
				found = true
				break
			}
		}
	}
	if !found {
		t.Logf("config errors: %v", configErrs)
	}
}

// TestFieldValidationOrder tests that multiple validators are run in the right order.
type MultiValidatorForm struct {
	Text *string `layrz:"char,required,min_length=5,max_length=10"`
}

func TestFieldValidationOrder(t *testing.T) {
	// Test too short
	form := &MultiValidatorForm{
		Text: Ptr("hi"), // 2 chars, min_length=5
	}

	errs := Validate(form)

	textErrs, ok := errs["text"]
	if !ok {
		t.Fatal("missing text key")
	}
	if len(textErrs) != 1 || textErrs[0].Code != "minLength" {
		t.Errorf("expected minLength error, got %v", textErrs)
	}

	// Test too long
	form.Text = Ptr("verylongtext") // 12 chars, max_length=10
	errs = Validate(form)

	textErrs, ok = errs["text"]
	if !ok {
		t.Fatal("missing text key")
	}
	if len(textErrs) != 1 || textErrs[0].Code != "maxLength" {
		t.Errorf("expected maxLength error, got %v", textErrs)
	}

	// Test valid
	form.Text = Ptr("valid")
	errs = Validate(form)

	if len(errs) > 0 {
		t.Errorf("expected no errors for valid text, got %v", errs)
	}
}

// TestSubformWithoutTag tests that a struct field without a tag is not treated as a subform.
type SubformWithoutTagForm struct {
	Name     *string `layrz:"char,required"`
	Embedded struct {
		Field *string
	}
}

func TestSubformWithoutTag(t *testing.T) {
	form := &SubformWithoutTagForm{
		Name: Ptr("test"),
		Embedded: struct {
			Field *string
		}{
			Field: nil,
		},
	}

	errs := Validate(form)

	// Embedded struct should not be validated (no tag)
	if len(errs) > 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

// TestSubformListWithNilElements tests that nil elements in a subform list are skipped.
func TestSubformListWithNilElements(t *testing.T) {
	type Item struct {
		Name *string `layrz:"char,required"`
	}

	type ListForm struct {
		Items []*Item `layrz:"subform_list"`
	}

	form := &ListForm{
		Items: []*Item{
			{Name: Ptr("item1")},
			nil,         // Nil element
			{Name: nil}, // Required field missing
		},
	}

	errs := Validate(form)

	// Should have items.2.name error (for the third item's missing Name)
	if _, ok := errs["items.2.name"]; !ok {
		t.Fatalf("expected items.2.name key, got keys: %v", errs.Keys())
	}

	// Should NOT have items.1 error (nil element skipped)
	for key := range errs {
		if key == "items.1.name" {
			t.Errorf("unexpected error for nil element: %q", key)
		}
	}
}

// TestSubformListWithEmptySlice tests that an empty subform list produces no errors.
func TestSubformListWithEmptySlice(t *testing.T) {
	type Item struct {
		Name *string `layrz:"char,required"`
	}

	type ListForm struct {
		Items []*Item `layrz:"subform_list"`
	}

	form := &ListForm{
		Items: nil, // Empty/nil slice
	}

	errs := Validate(form)

	// Should have no errors
	if len(errs) > 0 {
		t.Errorf("expected no errors for empty list, got %v", errs)
	}
}

// TestNumericPrefixInKey tests that numeric indices in subform lists are not camelCased.
func TestNumericPrefixInKey(t *testing.T) {
	type Item struct {
		Name *string `layrz:"char,required"`
	}

	type ListForm struct {
		Items []*Item `layrz:"subform_list"`
	}

	form := &ListForm{
		Items: []*Item{
			nil,
			nil,
			{Name: nil},
		},
	}

	errs := Validate(form)

	// Should have items.2.name (index 2, not camelCased)
	if _, ok := errs["items.2.name"]; !ok {
		t.Fatalf("expected items.2.name key, got keys: %v", errs.Keys())
	}
}

// TestCustomErrorExtraFields tests that custom errors with extra fields are preserved.
type CustomErrorForm struct {
	Name *string `layrz:"char,required,min_length=3"`
}

func (f *CustomErrorForm) CleanName(value *string) *FieldError {
	if value != nil && *value == "reserved" {
		return &FieldError{
			Code: "reserved",
			Extra: map[string]any{
				"reason": "this name is reserved",
			},
		}
	}
	return nil
}

func TestCustomErrorExtraFields(t *testing.T) {
	form := &CustomErrorForm{
		Name: Ptr("reserved"),
	}

	errs := Validate(form)

	nameErrs, ok := errs["name"]
	if !ok {
		t.Fatal("missing name key")
	}
	if len(nameErrs) != 1 || nameErrs[0].Code != "reserved" {
		t.Errorf("expected reserved error, got %v", nameErrs)
	}
	if nameErrs[0].Extra == nil {
		t.Fatal("expected extra fields")
	}
	if reason, ok := nameErrs[0].Extra["reason"]; !ok || reason != "this name is reserved" {
		t.Errorf("expected extra['reason']='this name is reserved', got %v", nameErrs[0].Extra)
	}
}

// TestNumericFieldTypes tests all numeric field types (*int8, *int16, *int32, *int64, *float32).
func TestNumericFieldTypes(t *testing.T) {
	type NumericForm struct {
		Int8Test    *int8    `layrz:"number,required,datatype=int,min_value=0,max_value=10"`
		Int16Test   *int16   `layrz:"number,required,datatype=int,min_value=0,max_value=100"`
		Int32Test   *int32   `layrz:"number,required,datatype=int,min_value=0,max_value=1000"`
		Int64Test   *int64   `layrz:"number,required,datatype=int,min_value=0,max_value=10000"`
		Float32Test *float32 `layrz:"number,required,datatype=float,min_value=0,max_value=5.5"`
	}

	// Test valid values
	form := &NumericForm{
		Int8Test:    Ptr(int8(5)),
		Int16Test:   Ptr(int16(50)),
		Int32Test:   Ptr(int32(500)),
		Int64Test:   Ptr(int64(5000)),
		Float32Test: Ptr(float32(3.5)),
	}

	errs := Validate(form)
	if len(errs) > 0 {
		t.Errorf("expected no errors for valid numeric types, got %v", errs.Keys())
	}

	// Test out of bounds
	form.Int8Test = Ptr(int8(20)) // Exceeds max_value=10
	errs = Validate(form)
	if maxErr, ok := errs["int8Test"]; !ok || len(maxErr) == 0 {
		t.Errorf("expected maxValue error for int8Test, got %v", errs.Keys())
	}

	// Test float32 bounds
	form.Int8Test = Ptr(int8(5))          // Reset
	form.Float32Test = Ptr(float32(10.0)) // Exceeds max_value=5.5
	errs = Validate(form)
	if maxErr, ok := errs["float32Test"]; !ok || len(maxErr) == 0 {
		t.Errorf("expected maxValue error for float32Test, got %v", errs.Keys())
	}
}

// TestAllScalarFieldKinds tests that all scalar field kinds (id, email, uuid, char, number, bool, json)
// are dispatched correctly through validateScalarField and produce zero errors for valid values.
func TestAllScalarFieldKinds(t *testing.T) {
	type AllFieldsForm struct {
		IDIntField    *int            `layrz:"id,required"`
		IDStringField *string         `layrz:"id,required"`
		EmailField    *string         `layrz:"email,required"`
		UUIDField     *string         `layrz:"uuid,required"`
		CharField     *string         `layrz:"char,required"`
		NumberInt     *int            `layrz:"number,required,datatype=int"`
		NumberFloat64 *float64        `layrz:"number,required,datatype=float"`
		BoolField     *bool           `layrz:"bool,required"`
		JSONListField *[]any          `layrz:"json,required,datatype=list"`
		JSONDictField *map[string]any `layrz:"json,required,datatype=dict"`
	}

	form := &AllFieldsForm{
		IDIntField:    Ptr(123),
		IDStringField: Ptr("999"), // ID string must be parseable as integer
		EmailField:    Ptr("test@example.com"),
		UUIDField:     Ptr("550e8400-e29b-41d4-a716-446655440000"),
		CharField:     Ptr("hello"),
		NumberInt:     Ptr(42),
		NumberFloat64: Ptr(3.14),
		BoolField:     Ptr(true),
		JSONListField: Ptr([]any{1, "two", 3.0}),
		JSONDictField: Ptr(map[string]any{"key": "value"}),
	}

	errs := Validate(form)
	if len(errs) > 0 {
		codes := []string{}
		for key, fieldErrs := range errs {
			for _, e := range fieldErrs {
				codes = append(codes, key+":"+e.Code)
			}
		}
		t.Errorf("expected no errors for valid fields, got: %v", codes)
	}
}

// TestNilPointerScalarFields tests that nil pointers in scalar fields produce required errors.
func TestNilPointerScalarFields(t *testing.T) {
	type NilFieldsForm struct {
		IdField     *int    `layrz:"id,required"`
		EmailField  *string `layrz:"email,required"`
		UuidField   *string `layrz:"uuid,required"`
		CharField   *string `layrz:"char,required"`
		NumberField *int    `layrz:"number,required"`
		BoolField   *bool   `layrz:"bool,required"`
		JsonField   *[]any  `layrz:"json,required"`
	}

	form := &NilFieldsForm{
		// All fields are nil, should produce required errors
	}

	errs := Validate(form)
	if len(errs) == 0 {
		t.Fatal("expected errors for nil required fields")
	}

	expectedKeys := []string{"idField", "emailField", "uuidField", "charField", "numberField", "boolField", "jsonField"}
	for _, key := range expectedKeys {
		if fieldErrs, ok := errs[key]; !ok || len(fieldErrs) == 0 {
			t.Errorf("expected error for %s, got: %v", key, errs.Keys())
		} else if fieldErrs[0].Code != "required" {
			t.Errorf("expected required error for %s, got: %s", key, fieldErrs[0].Code)
		}
	}
}

// TestEmbeddedStructValueType tests that untagged embedded struct values (not pointers) are handled.
type EmbeddedValueAddress struct {
	StreetName *string `layrz:"char,required"`
}

type EmbeddingFormValue struct {
	EmbeddedValueAddress // Anonymous value, no tag -> promote fields
}

func TestEmbeddedStructValueType(t *testing.T) {
	form := &EmbeddingFormValue{
		EmbeddedValueAddress: EmbeddedValueAddress{
			StreetName: nil, // Missing required
		},
	}

	errs := Validate(form)

	// Should have streetName error (promoted from embedded struct)
	if _, ok := errs["streetName"]; !ok {
		t.Fatalf("expected streetName key (promoted field), got keys: %v", errs.Keys())
	}

	// Test with valid value
	form.EmbeddedValueAddress.StreetName = Ptr("Main St")
	errs = Validate(form)
	if len(errs) > 0 {
		t.Errorf("expected no errors for valid street name, got: %v", errs.Keys())
	}
}

// TestSubformListWithNonPointerElements tests subform_list with slice of non-pointer structs.
func TestSubformListWithNonPointerElements(t *testing.T) {
	type Item struct {
		Name *string `layrz:"char,required"`
	}

	type ListFormNonPtr struct {
		Items []Item `layrz:"subform_list"` // Non-pointer elements
	}

	form := &ListFormNonPtr{
		Items: []Item{
			{Name: Ptr("item1")},
			{Name: nil}, // Missing required field
		},
	}

	errs := Validate(form)

	// Should have items.1.name error
	if _, ok := errs["items.1.name"]; !ok {
		t.Fatalf("expected items.1.name key, got keys: %v", errs.Keys())
	}

	// Test with empty non-pointer slice
	form.Items = []Item{}
	errs = Validate(form)
	if len(errs) > 0 {
		t.Errorf("expected no errors for empty non-pointer list, got %v", errs.Keys())
	}
}

// TestSubformWithErrorsAndClean tests subform with errors and clean method discovery.
type SubformWithClean struct {
	Name *string `layrz:"char,required"`
}

func (s *SubformWithClean) CleanName(val *string) *FieldError {
	if val != nil && *val == "forbidden" {
		return &FieldError{Code: "forbidden"}
	}
	return nil
}

type FormWithSubformClean struct {
	Sub *SubformWithClean `layrz:"subform"`
}

func TestSubformWithErrorsAndClean(t *testing.T) {
	form := &FormWithSubformClean{
		Sub: &SubformWithClean{
			Name: Ptr("forbidden"),
		},
	}

	errs := Validate(form)

	// Should have sub.name error from clean method
	if subErrs, ok := errs["sub.name"]; !ok || len(subErrs) == 0 {
		t.Errorf("expected sub.name error, got keys: %v", errs.Keys())
	} else if subErrs[0].Code != "forbidden" {
		t.Errorf("expected forbidden code, got: %s", subErrs[0].Code)
	}
}

// TestNilSubformInEngine tests that a nil subform is skipped (Go divergence from Python).
func TestNilSubformInEngine(t *testing.T) {
	type SubForm struct {
		Name *string `layrz:"char,required"`
	}

	type FormWithNilSub struct {
		Sub *SubForm `layrz:"subform"`
	}

	form := &FormWithNilSub{
		Sub: nil, // Nil subform should be skipped
	}

	errs := Validate(form)
	if len(errs) > 0 {
		t.Errorf("expected no errors for nil subform, got: %v", errs.Keys())
	}
}

// TestNilSliceSubformList tests that a nil slice in subform_list produces no errors.
func TestNilSliceSubformList(t *testing.T) {
	type Item struct {
		Name *string `layrz:"char,required"`
	}

	type ListFormNil struct {
		Items []*Item `layrz:"subform_list"`
	}

	form := &ListFormNil{
		Items: nil, // Nil slice
	}

	errs := Validate(form)
	if len(errs) > 0 {
		t.Errorf("expected no errors for nil subform list, got: %v", errs.Keys())
	}
}

// TestEmbeddedStructWithTag tests that tagged anonymous structs are treated as configuration errors.
type TaggedEmbedded struct {
	Field *string `layrz:"char,required"`
}

type FormWithTaggedEmbedded struct {
	TaggedEmbedded `layrz:"char"` // Tagged anonymous struct (config error)
}

func TestEmbeddedStructWithTag(t *testing.T) {
	form := &FormWithTaggedEmbedded{
		TaggedEmbedded: TaggedEmbedded{
			Field: Ptr("test"),
		},
	}

	errs := Validate(form)

	// Tagged anonymous structs should produce a config error
	if configErrs, ok := errs[configErrorKey]; !ok {
		t.Fatalf("expected _config error for tagged anonymous struct, got keys: %v", errs.Keys())
	} else if len(configErrs) == 0 {
		t.Fatal("expected at least one config error")
	}
}

// TestJSONFieldWithNilPointer tests that JSON fields with nil pointers are handled correctly.
func TestJSONFieldWithNilPointer(t *testing.T) {
	type JSONForm struct {
		ListField *[]any          `layrz:"json,required,datatype=list"`
		DictField *map[string]any `layrz:"json,required,datatype=dict"`
	}

	form := &JSONForm{
		ListField: nil,
		DictField: nil,
	}

	errs := Validate(form)

	// Both should have required errors
	if _, ok := errs["listField"]; !ok {
		t.Errorf("expected listField error, got keys: %v", errs.Keys())
	}
	if _, ok := errs["dictField"]; !ok {
		t.Errorf("expected dictField error, got keys: %v", errs.Keys())
	}
}

// TestBoolFieldWithNilPointer tests that bool fields with nil pointers are handled correctly.
func TestBoolFieldWithNilPointer(t *testing.T) {
	type BoolForm struct {
		Flag *bool `layrz:"bool,required"`
	}

	form := &BoolForm{
		Flag: nil,
	}

	errs := Validate(form)

	if _, ok := errs["flag"]; !ok {
		t.Errorf("expected flag error, got keys: %v", errs.Keys())
	}
}

// TestIDFieldWithBothTypes tests that id field works with both *int and *string types.
func TestIDFieldWithBothTypes(t *testing.T) {
	type IDMixedForm struct {
		IntId    *int    `layrz:"id,required"`
		StringId *string `layrz:"id,required"`
	}

	form := &IDMixedForm{
		IntId:    Ptr(99),
		StringId: Ptr("456"), // String ID must be parseable as integer
	}

	errs := Validate(form)
	if len(errs) > 0 {
		t.Errorf("expected no errors for valid id fields, got: %v", errs.Keys())
	}

	// Test with nil values
	form.IntId = nil
	form.StringId = nil
	errs = Validate(form)

	if _, ok := errs["intId"]; !ok {
		t.Errorf("expected intId error")
	}
	if _, ok := errs["stringId"]; !ok {
		t.Errorf("expected stringId error")
	}
}

// TestValueFieldChar tests char validation with non-pointer string fields.
func TestValueFieldChar(t *testing.T) {
	t.Run("empty string with required and empty=false", func(t *testing.T) {
		type ValueCharForm struct {
			Name string `layrz:"char,required"`
		}
		form := &ValueCharForm{Name: ""}
		errs := Validate(form)

		// Should have "empty" error, NOT "required"
		if fieldErrs, ok := errs["name"]; !ok {
			t.Errorf("expected error for name field")
		} else if len(fieldErrs) > 0 && fieldErrs[0].Code == "required" {
			t.Errorf("value field should not emit required; got code=%q", fieldErrs[0].Code)
		} else if len(fieldErrs) == 0 || fieldErrs[0].Code != "empty" {
			t.Errorf("empty string with empty=false should emit empty error, got %v", fieldErrs)
		}
	})

	t.Run("empty string with required and empty=true", func(t *testing.T) {
		type ValueCharForm struct {
			Name string `layrz:"char,required,empty"`
		}
		form := &ValueCharForm{Name: ""}
		errs := Validate(form)

		// Should have no errors
		if len(errs) > 0 {
			t.Errorf("expected no errors with empty=true, got %v", errs)
		}
	})

	t.Run("non-empty string with min_length constraint", func(t *testing.T) {
		type ValueCharForm struct {
			Name string `layrz:"char,required,min_length=5"`
		}
		form := &ValueCharForm{Name: "abc"}
		errs := Validate(form)

		if fieldErrs, ok := errs["name"]; !ok {
			t.Errorf("expected error for name field")
		} else if len(fieldErrs) == 0 || fieldErrs[0].Code != "minLength" {
			t.Errorf("expected minLength error, got %v", fieldErrs)
		}
	})

	t.Run("valid string with required and constraints", func(t *testing.T) {
		type ValueCharForm struct {
			Name string `layrz:"char,required,min_length=3,max_length=10"`
		}
		form := &ValueCharForm{Name: "valid"}
		errs := Validate(form)

		if len(errs) > 0 {
			t.Errorf("expected no errors, got %v", errs)
		}
	})
}

// TestValueFieldEmail tests email validation with non-pointer string fields.
func TestValueFieldEmail(t *testing.T) {
	t.Run("empty string with required", func(t *testing.T) {
		type ValueEmailForm struct {
			Email string `layrz:"email,required"`
		}
		form := &ValueEmailForm{Email: ""}
		errs := Validate(form)

		// Should have "empty" error, NOT "required"
		if fieldErrs, ok := errs["email"]; !ok {
			t.Errorf("expected error for email field")
		} else if len(fieldErrs) > 0 && fieldErrs[0].Code == "required" {
			t.Errorf("value field should not emit required; got code=%q", fieldErrs[0].Code)
		}
	})

	t.Run("invalid email format", func(t *testing.T) {
		type ValueEmailForm struct {
			Email string `layrz:"email,required"`
		}
		form := &ValueEmailForm{Email: "nope"}
		errs := Validate(form)

		if fieldErrs, ok := errs["email"]; !ok {
			t.Errorf("expected error for email field")
		} else if len(fieldErrs) == 0 || fieldErrs[0].Code != "invalid" {
			t.Errorf("expected invalid error, got %v", fieldErrs)
		}
	})

	t.Run("valid email", func(t *testing.T) {
		type ValueEmailForm struct {
			Email string `layrz:"email,required"`
		}
		form := &ValueEmailForm{Email: "test@example.com"}
		errs := Validate(form)

		if len(errs) > 0 {
			t.Errorf("expected no errors, got %v", errs)
		}
	})
}

// TestValueFieldNumber tests number validation with non-pointer numeric fields.
func TestValueFieldNumber(t *testing.T) {
	t.Run("int with required and min_value constraint", func(t *testing.T) {
		type ValueNumberForm struct {
			Count int `layrz:"number,required,min_value=1"`
		}
		form := &ValueNumberForm{Count: 0}
		errs := Validate(form)

		// Should have "minValue" error, NOT "required"
		if fieldErrs, ok := errs["count"]; !ok {
			t.Errorf("expected error for count field")
		} else if len(fieldErrs) > 0 && fieldErrs[0].Code == "required" {
			t.Errorf("value field should not emit required; got code=%q", fieldErrs[0].Code)
		} else if len(fieldErrs) == 0 || fieldErrs[0].Code != "minValue" {
			t.Errorf("expected minValue error, got %v", fieldErrs)
		}
	})

	t.Run("int within valid range", func(t *testing.T) {
		type ValueNumberForm struct {
			Count int `layrz:"number,required,min_value=1,max_value=10"`
		}
		form := &ValueNumberForm{Count: 5}
		errs := Validate(form)

		if len(errs) > 0 {
			t.Errorf("expected no errors, got %v", errs)
		}
	})

	t.Run("float64 with required", func(t *testing.T) {
		type ValueNumberForm struct {
			Price float64 `layrz:"number,required,datatype=float,min_value=0"`
		}
		form := &ValueNumberForm{Price: 19.99}
		errs := Validate(form)

		if len(errs) > 0 {
			t.Errorf("expected no errors, got %v", errs)
		}
	})
}

// TestValueFieldBool tests bool validation with non-pointer bool fields.
func TestValueFieldBool(t *testing.T) {
	t.Run("false bool with required", func(t *testing.T) {
		type ValueBoolForm struct {
			Active bool `layrz:"bool,required"`
		}
		form := &ValueBoolForm{Active: false}
		errs := Validate(form)

		// Should have NO errors (false is a valid present bool)
		if len(errs) > 0 {
			t.Errorf("expected no errors for false bool, got %v", errs)
		}
	})

	t.Run("true bool with required", func(t *testing.T) {
		type ValueBoolForm struct {
			Active bool `layrz:"bool,required"`
		}
		form := &ValueBoolForm{Active: true}
		errs := Validate(form)

		// Should have NO errors
		if len(errs) > 0 {
			t.Errorf("expected no errors for true bool, got %v", errs)
		}
	})
}

// TestValueFieldJSON tests JSON validation with non-pointer slice/map fields.
func TestValueFieldJSON(t *testing.T) {
	t.Run("nil slice with datatype=list and empty=false", func(t *testing.T) {
		type ValueJSONForm struct {
			Tags []any `layrz:"json,datatype=list"`
		}
		form := &ValueJSONForm{Tags: nil}
		errs := Validate(form)

		// A nil slice (zero value) in a value field is present but empty
		// empty=false should emit "invalid"
		if fieldErrs, ok := errs["tags"]; !ok {
			t.Errorf("expected error for tags field")
		} else if len(fieldErrs) == 0 || fieldErrs[0].Code != "invalid" {
			t.Errorf("expected invalid error for nil slice, got %v", fieldErrs)
		}
	})

	t.Run("empty slice with datatype=list and empty=false", func(t *testing.T) {
		type ValueJSONForm struct {
			Tags []any `layrz:"json,datatype=list"`
		}
		form := &ValueJSONForm{Tags: []any{}}
		errs := Validate(form)

		// empty slice with empty=false should emit "invalid"
		if fieldErrs, ok := errs["tags"]; !ok {
			t.Errorf("expected error for tags field")
		} else if len(fieldErrs) == 0 || fieldErrs[0].Code != "invalid" {
			t.Errorf("expected invalid error for empty slice, got %v", fieldErrs)
		}
	})

	t.Run("non-empty slice with datatype=list", func(t *testing.T) {
		type ValueJSONForm struct {
			Tags []any `layrz:"json,datatype=list"`
		}
		form := &ValueJSONForm{Tags: []any{"x", "y"}}
		errs := Validate(form)

		if len(errs) > 0 {
			t.Errorf("expected no errors for non-empty slice, got %v", errs)
		}
	})

	t.Run("nil map with datatype=dict and empty=false", func(t *testing.T) {
		type ValueJSONForm struct {
			Metadata map[string]any `layrz:"json,datatype=dict"`
		}
		form := &ValueJSONForm{Metadata: nil}
		errs := Validate(form)

		// A nil map (zero value) in a value field is present but empty
		// empty=false should emit "invalid"
		if fieldErrs, ok := errs["metadata"]; !ok {
			t.Errorf("expected error for metadata field")
		} else if len(fieldErrs) == 0 || fieldErrs[0].Code != "invalid" {
			t.Errorf("expected invalid error for nil map, got %v", fieldErrs)
		}
	})

	t.Run("non-empty map with datatype=dict", func(t *testing.T) {
		type ValueJSONForm struct {
			Metadata map[string]any `layrz:"json,datatype=dict"`
		}
		form := &ValueJSONForm{Metadata: map[string]any{"key": "value"}}
		errs := Validate(form)

		if len(errs) > 0 {
			t.Errorf("expected no errors for non-empty map, got %v", errs)
		}
	})
}

// TestMixedPointerValueFields tests a struct with both pointer and value scalar fields.
func TestMixedPointerValueFields(t *testing.T) {
	type MixedForm struct {
		// Value fields (always present)
		Name  string `layrz:"char,required,min_length=1"`
		Count int    `layrz:"number,required,min_value=1"`

		// Pointer fields (can be absent)
		OptionalEmail *string `layrz:"email,required"`
		OptionalFlag  *bool   `layrz:"bool,required"`
	}

	form := &MixedForm{
		Name:          "", // zero value, not absent
		Count:         0,  // zero value, not absent
		OptionalEmail: nil,
		OptionalFlag:  nil,
	}

	errs := Validate(form)

	// Value fields should emit validation errors (not required, but actual validation)
	if _, ok := errs["name"]; !ok {
		t.Errorf("expected error for name value field (empty string)")
	}
	if _, ok := errs["count"]; !ok {
		t.Errorf("expected error for count value field (min_value violation)")
	}

	// Pointer fields should emit required errors
	if _, ok := errs["optionalEmail"]; !ok {
		t.Errorf("expected required error for optionalEmail pointer field")
	}
	if _, ok := errs["optionalFlag"]; !ok {
		t.Errorf("expected required error for optionalFlag pointer field")
	}
}

// TestValueFieldInSubform tests value scalar fields inside a subform.
func TestValueFieldInSubform(t *testing.T) {
	type Address struct {
		Street string `layrz:"char,required,min_length=5"`
	}

	type AddressForm struct {
		Addr *Address `layrz:"subform"`
	}

	form := &AddressForm{
		Addr: &Address{Street: "Main"},
	}

	errs := Validate(form)

	// Should have error keyed as "addr.street"
	if fieldErrs, ok := errs["addr.street"]; !ok {
		t.Errorf("expected error for addr.street, got keys: %v", errs.Keys())
	} else if len(fieldErrs) == 0 || fieldErrs[0].Code != "minLength" {
		t.Errorf("expected minLength error in subform, got %v", fieldErrs)
	}
}

// TestValueFieldInSubformList tests value scalar fields inside a subform_list.
func TestValueFieldInSubformList(t *testing.T) {
	type Item struct {
		Name string `layrz:"char,required,min_length=3"`
	}

	type OrderForm struct {
		Items []Item `layrz:"subform_list"`
	}

	form := &OrderForm{
		Items: []Item{
			{Name: "Valid Item"},
			{Name: "X"}, // too short
		},
	}

	errs := Validate(form)

	// Should have error keyed as "items.1.name"
	if fieldErrs, ok := errs["items.1.name"]; !ok {
		t.Errorf("expected error for items.1.name, got keys: %v", errs.Keys())
	} else if len(fieldErrs) == 0 || fieldErrs[0].Code != "minLength" {
		t.Errorf("expected minLength error in subform_list element, got %v", fieldErrs)
	}
}

// ValueFieldFormWithClean has a value field and a Convention A clean method.
type ValueFieldFormWithClean struct {
	Code string `layrz:"char,required"`
}

// CleanCode implements Convention A clean method for Code field (value field).
func (f *ValueFieldFormWithClean) CleanCode(value string) *FieldError {
	if value != "ALLOWED" {
		return &FieldError{Code: "notAllowed"}
	}
	return nil
}

// TestValueFieldCleanMethod tests Convention A clean methods with value fields.
func TestValueFieldCleanMethod(t *testing.T) {
	form := &ValueFieldFormWithClean{Code: "DENIED"}
	errs := Validate(form)

	if fieldErrs, ok := errs["code"]; !ok {
		t.Errorf("expected error for code field")
	} else if len(fieldErrs) == 0 || fieldErrs[0].Code != "notAllowed" {
		t.Errorf("expected notAllowed error from clean method, got %v", fieldErrs)
	}
}
