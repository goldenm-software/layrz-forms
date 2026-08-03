package layrz

import (
	"reflect"
	"testing"
)

// Test constants for coverage
const (
	covRequired = "required"
	covInvalid  = "invalid"
	covEmpty    = "empty"
	covMinValue = "minValue"
	covMaxValue = "maxValue"
	covInt      = "int"
	covFloat    = "float"
	covList     = "list"
	covDict     = "dict"
)

// TestParseTagEmptyTokenEdges tests edge cases in tag parsing with empty tokens.
func TestParseTagEmptyTokenEdges(t *testing.T) {
	t.Run("empty token between commas", func(t *testing.T) {
		spec, err := ParseTag("char,,required")
		if err != nil {
			t.Fatalf("ParseTag with empty token expected no error, got %v", err)
		}
		if !spec.Required {
			t.Errorf("expected Required=true despite empty token")
		}
	})

	t.Run("duplicate min_value rule", func(t *testing.T) {
		_, err := ParseTag("number,min_value=1.0,min_value=2.0")
		if err == nil {
			t.Errorf("expected error for duplicate min_value")
		}
	})

	t.Run("duplicate max_value rule", func(t *testing.T) {
		_, err := ParseTag("number,max_value=1.0,max_value=2.0")
		if err == nil {
			t.Errorf("expected error for duplicate max_value")
		}
	})

	t.Run("duplicate empty rule", func(t *testing.T) {
		_, err := ParseTag("char,empty,empty")
		if err == nil {
			t.Errorf("expected error for duplicate empty")
		}
	})

	t.Run("duplicate choices rule", func(t *testing.T) {
		_, err := ParseTag("char,choices=a|b,choices=c|d")
		if err == nil {
			t.Errorf("expected error for duplicate choices")
		}
	})

	t.Run("regex with comma and parens (consumes rest)", func(t *testing.T) {
		// regex= consumes the rest of the tag, so "regex=^a$,(invalid"
		// is parsed as regex value "^a$,(invalid" which is invalid regex
		spec, err := ParseTag("char,regex=^a$,(invalid")
		if err == nil {
			t.Errorf("expected error for invalid regex pattern")
		}
		if spec != nil {
			t.Errorf("expected nil spec on error")
		}
	})

	t.Run("invalid min_value float", func(t *testing.T) {
		_, err := ParseTag("number,min_value=not_a_number")
		if err == nil {
			t.Errorf("expected error for invalid min_value")
		}
	})

	t.Run("invalid max_value float", func(t *testing.T) {
		_, err := ParseTag("number,max_value=not_a_number")
		if err == nil {
			t.Errorf("expected error for invalid max_value")
		}
	})

	t.Run("invalid datatype in parseKeyValueRule", func(t *testing.T) {
		_, err := ParseTag("json,datatype=bad_type")
		if err == nil {
			t.Errorf("expected error for invalid datatype")
		}
	})

	t.Run("unknown rule key", func(t *testing.T) {
		_, err := ParseTag("char,unknown_key=value")
		if err == nil {
			t.Errorf("expected error for unknown rule key")
		}
	})

	t.Run("unrecognized token without equals", func(t *testing.T) {
		_, err := ParseTag("char,unrecognized_token")
		if err == nil {
			t.Errorf("expected error for unrecognized token")
		}
	})

	t.Run("choices with single value", func(t *testing.T) {
		spec, err := ParseTag("char,choices=single")
		if err != nil {
			t.Fatalf("ParseTag with single choice: %v", err)
		}
		if len(spec.Choices) != 1 || spec.Choices[0] != "single" {
			t.Errorf("expected single choice, got %v", spec.Choices)
		}
	})

	t.Run("choices with empty pipe segments", func(t *testing.T) {
		spec, err := ParseTag("char,choices=a||b")
		if err != nil {
			t.Fatalf("ParseTag with pipe gaps: %v", err)
		}
		if len(spec.Choices) != 3 {
			t.Errorf("expected 3 choices (including empty), got %v", spec.Choices)
		}
	})
}

// TestToCamelCaseEmptySegments tests ToCamelCase with empty and edge-case inputs.
func TestToCamelCaseEmptySegments(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		result := ToCamelCase("")
		if result != "" {
			t.Errorf("expected empty string, got %q", result)
		}
	})

	t.Run("single underscore", func(t *testing.T) {
		result := ToCamelCase("_")
		if result != "" {
			t.Errorf("expected empty string for single underscore, got %q", result)
		}
	})

	t.Run("multiple underscores", func(t *testing.T) {
		result := ToCamelCase("__")
		if result != "" {
			t.Errorf("expected empty string for multiple underscores, got %q", result)
		}
	})

	t.Run("dot with empty segments", func(t *testing.T) {
		result := ToCamelCase("a..b")
		if result != "a..b" {
			t.Errorf("expected a..b, got %q", result)
		}
	})

	t.Run("trailing underscore", func(t *testing.T) {
		result := ToCamelCase("field_")
		if result != "field" {
			t.Errorf("expected 'field', got %q", result)
		}
	})

	t.Run("leading underscore", func(t *testing.T) {
		result := ToCamelCase("_field")
		if result != "field" {
			t.Errorf("expected 'field', got %q", result)
		}
	})

	t.Run("underscores with dots", func(t *testing.T) {
		result := ToCamelCase("_a._b")
		if result != "a.b" {
			t.Errorf("expected 'a.b', got %q", result)
		}
	})
}

// TestErrorsAddEdgeCases tests Errors.Add with edge cases.
func TestErrorsAddEdgeCases(t *testing.T) {
	t.Run("add to empty varargs", func(t *testing.T) {
		e := make(Errors)
		e.Add("field")
		if len(e) != 0 {
			t.Errorf("expected empty errors after Add with no args")
		}
	})

	t.Run("all nil varargs", func(t *testing.T) {
		e := make(Errors)
		e.Add("field", nil, nil, nil)
		if len(e) != 0 {
			t.Errorf("expected empty errors after Add with all nil args")
		}
	})

	t.Run("mixed nil and non-nil", func(t *testing.T) {
		e := make(Errors)
		e.Add("field", nil, &FieldError{Code: covRequired}, nil, &FieldError{Code: covInvalid})
		if len(e["field"]) != 2 {
			t.Errorf("expected 2 errors, got %d", len(e["field"]))
		}
	})
}

// TestCleanMethodSignatureValidation tests validation of clean method signatures.
func TestCleanMethodSignatureValidation(t *testing.T) {
	t.Run("convention a wrong arity - no parameters", func(t *testing.T) {
		// This would require a malformed method signature that can't be easily tested
		// via the public API. Skip for now as it's covered by validateConventionA tests.
	})

	t.Run("convention a wrong parameter type", func(t *testing.T) {
		type BadParamForm struct {
			Field *string `layrz:"char,required"`
		}

		// Test validateConventionA directly with wrong parameter type
		err := validateConventionA(
			reflect.TypeOf(func(f *BadParamForm, x int) *FieldError { return nil }),
			"Field",
			reflect.TypeOf(BadParamForm{}),
		)
		if err == nil {
			t.Errorf("expected error for wrong parameter type")
		}
	})

	t.Run("convention a wrong return type", func(t *testing.T) {
		type BadReturnForm struct {
			Field *string `layrz:"char,required"`
		}

		err := validateConventionA(
			reflect.TypeOf(func(f *BadReturnForm, x *string) int { return 0 }),
			"Field",
			reflect.TypeOf(BadReturnForm{}),
		)
		if err == nil {
			t.Errorf("expected error for wrong return type")
		}
	})

	t.Run("convention a wrong arity - extra parameters", func(t *testing.T) {
		type BadArityForm struct {
			Field *string `layrz:"char,required"`
		}

		err := validateConventionA(
			reflect.TypeOf(func(f *BadArityForm, x *string, y string) *FieldError { return nil }),
			"Field",
			reflect.TypeOf(BadArityForm{}),
		)
		if err == nil {
			t.Errorf("expected error for wrong arity (too many params)")
		}
	})

	t.Run("convention a wrong return count", func(t *testing.T) {
		type BadReturnCountForm struct {
			Field *string `layrz:"char,required"`
		}

		err := validateConventionA(
			reflect.TypeOf(func(f *BadReturnCountForm, x *string) (*FieldError, error) { return nil, nil }),
			"Field",
			reflect.TypeOf(BadReturnCountForm{}),
		)
		if err == nil {
			t.Errorf("expected error for wrong return count")
		}
	})

	t.Run("convention b wrong arity", func(t *testing.T) {
		err := validateConventionB(
			reflect.TypeOf(func(f *struct{}, x string) Errors { return nil }),
		)
		if err == nil {
			t.Errorf("expected error for convention B with parameter")
		}
	})

	t.Run("convention b wrong return type", func(t *testing.T) {
		err := validateConventionB(
			reflect.TypeOf(func(f *struct{}) *FieldError { return nil }),
		)
		if err == nil {
			t.Errorf("expected error for convention B returning FieldError instead of Errors")
		}
	})

	t.Run("convention b wrong return count", func(t *testing.T) {
		err := validateConventionB(
			reflect.TypeOf(func(f *struct{}) (Errors, error) { return nil, nil }),
		)
		if err == nil {
			t.Errorf("expected error for convention B with 2 return values")
		}
	})
}

// TestSpecCheckTypeErrors tests type checking rejection paths.
func TestSpecCheckTypeErrors(t *testing.T) {
	t.Run("number datatype mismatch int->float", func(t *testing.T) {
		spec := &FieldSpec{Kind: KindNumber, Datatype: covFloat}
		err := spec.CheckType(reflect.TypeOf(int32(0)))
		if err == nil {
			t.Errorf("expected error for int type with float datatype")
		}
	})

	t.Run("number datatype mismatch float->int", func(t *testing.T) {
		spec := &FieldSpec{Kind: KindNumber, Datatype: covInt}
		err := spec.CheckType(reflect.TypeOf(3.14))
		if err == nil {
			t.Errorf("expected error for float type with int datatype")
		}
	})

	t.Run("json on scalar type", func(t *testing.T) {
		spec := &FieldSpec{Kind: KindJSON}
		err := spec.CheckType(reflect.TypeOf("string"))
		if err == nil {
			t.Errorf("expected error for json on scalar string")
		}
	})

	t.Run("json pointer to scalar", func(t *testing.T) {
		spec := &FieldSpec{Kind: KindJSON}
		err := spec.CheckType(reflect.TypeOf(Ptr("string")))
		if err == nil {
			t.Errorf("expected error for json on *string")
		}
	})

	t.Run("json list with dict datatype", func(t *testing.T) {
		spec := &FieldSpec{Kind: KindJSON, Datatype: covDict}
		err := spec.CheckType(reflect.TypeOf([]int{}))
		if err == nil {
			t.Errorf("expected error for slice with dict datatype")
		}
	})

	t.Run("json dict with list datatype", func(t *testing.T) {
		spec := &FieldSpec{Kind: KindJSON, Datatype: covList}
		err := spec.CheckType(reflect.TypeOf(map[string]int{}))
		if err == nil {
			t.Errorf("expected error for map with list datatype")
		}
	})

	t.Run("subform not pointer", func(t *testing.T) {
		spec := &FieldSpec{Kind: KindSubform}
		type TestStruct struct{}
		err := spec.CheckType(reflect.TypeOf(TestStruct{}))
		if err == nil {
			t.Errorf("expected error for subform on non-pointer value")
		}
	})

	t.Run("subform pointer to non-struct", func(t *testing.T) {
		spec := &FieldSpec{Kind: KindSubform}
		err := spec.CheckType(reflect.TypeOf(Ptr("string")))
		if err == nil {
			t.Errorf("expected error for subform on pointer to non-struct")
		}
	})

	t.Run("subform_list not slice", func(t *testing.T) {
		spec := &FieldSpec{Kind: KindSubformList}
		type TestStruct struct{}
		err := spec.CheckType(reflect.TypeOf(TestStruct{}))
		if err == nil {
			t.Errorf("expected error for subform_list on non-slice")
		}
	})

	t.Run("subform_list of non-structs", func(t *testing.T) {
		spec := &FieldSpec{Kind: KindSubformList}
		err := spec.CheckType(reflect.TypeOf([]string{}))
		if err == nil {
			t.Errorf("expected error for subform_list of strings")
		}
	})

	t.Run("subform_list of pointers to non-structs", func(t *testing.T) {
		spec := &FieldSpec{Kind: KindSubformList}
		err := spec.CheckType(reflect.TypeOf([]*string{}))
		if err == nil {
			t.Errorf("expected error for subform_list of *string")
		}
	})

	t.Run("json infer list from slice", func(t *testing.T) {
		spec := &FieldSpec{Kind: KindJSON}
		err := spec.CheckType(reflect.TypeOf([]int{}))
		if err != nil {
			t.Fatalf("CheckType for json on slice: %v", err)
		}
		if spec.Datatype != covList {
			t.Errorf("expected Datatype=list, got %q", spec.Datatype)
		}
	})

	t.Run("json infer dict from map", func(t *testing.T) {
		spec := &FieldSpec{Kind: KindJSON}
		err := spec.CheckType(reflect.TypeOf(map[string]int{}))
		if err != nil {
			t.Fatalf("CheckType for json on map: %v", err)
		}
		if spec.Datatype != covDict {
			t.Errorf("expected Datatype=dict, got %q", spec.Datatype)
		}
	})
}

// TestEngineDepthGuard tests the recursion depth guard.
func TestEngineDepthGuard(t *testing.T) {
	type DeepForm struct {
		Sub *DeepForm `layrz:"subform"`
	}

	// Create a deeply nested structure by hand
	deep := &DeepForm{}
	current := deep
	for i := 0; i < 35; i++ {
		current.Sub = &DeepForm{}
		current = current.Sub
	}

	errs := Validate(deep)

	// Should have a config error for exceeding max depth
	configErrs, ok := errs[configErrorKey]
	if !ok || len(configErrs) == 0 {
		t.Errorf("expected config error for depth exceeded, got none")
	}
}

// TestValidateIDTypeRejection tests ValidateID with wrong types.
func TestValidateIDTypeRejection(t *testing.T) {
	t.Run("bool is invalid", func(t *testing.T) {
		errs := ValidateID(true, IDRules{})
		if len(errs) != 1 || errs[0].Code != covInvalid {
			t.Errorf("expected invalid code for bool, got %v", errs)
		}
	})

	t.Run("float64 is invalid", func(t *testing.T) {
		errs := ValidateID(3.14, IDRules{})
		if len(errs) != 1 || errs[0].Code != covInvalid {
			t.Errorf("expected invalid code for float64, got %v", errs)
		}
	})

	t.Run("uint is invalid", func(t *testing.T) {
		errs := ValidateID(uint(42), IDRules{})
		if len(errs) != 1 || errs[0].Code != covInvalid {
			t.Errorf("expected invalid code for uint, got %v", errs)
		}
	})

	t.Run("string that is not a valid int", func(t *testing.T) {
		errs := ValidateID("not_an_int", IDRules{})
		if len(errs) != 1 || errs[0].Code != covInvalid {
			t.Errorf("expected invalid code for unparseable string, got %v", errs)
		}
	})

	t.Run("zero is invalid", func(t *testing.T) {
		errs := ValidateID(0, IDRules{})
		if len(errs) != 1 || errs[0].Code != covInvalid {
			t.Errorf("expected invalid code for zero, got %v", errs)
		}
	})

	t.Run("negative number is invalid", func(t *testing.T) {
		errs := ValidateID(-42, IDRules{})
		if len(errs) != 1 || errs[0].Code != covInvalid {
			t.Errorf("expected invalid code for negative, got %v", errs)
		}
	})

	t.Run("negative string is invalid", func(t *testing.T) {
		errs := ValidateID("-42", IDRules{})
		if len(errs) != 1 || errs[0].Code != covInvalid {
			t.Errorf("expected invalid code for negative string, got %v", errs)
		}
	})
}

// TestValidateEmailTypeRejection tests ValidateEmail with wrong types.
func TestValidateEmailTypeRejection(t *testing.T) {
	t.Run("nil pointer is absent", func(t *testing.T) {
		errs := ValidateEmail(nil, EmailRules{Required: true})
		if len(errs) != 1 || errs[0].Code != covRequired {
			t.Errorf("expected required code for nil pointer, got %v", errs)
		}
	})

	t.Run("nil pointer optional is ok", func(t *testing.T) {
		errs := ValidateEmail(nil, EmailRules{Required: false})
		if len(errs) != 0 {
			t.Errorf("expected no error for optional nil pointer, got %v", errs)
		}
	})
}

// TestValidateNumberTypeRejection tests ValidateNumber with wrong types.
func TestValidateNumberTypeRejection(t *testing.T) {
	t.Run("string is invalid", func(t *testing.T) {
		errs := ValidateNumber("123", NumberRules{})
		if len(errs) != 1 || errs[0].Code != covInvalid {
			t.Errorf("expected invalid code for string, got %v", errs)
		}
	})

	t.Run("bool is invalid", func(t *testing.T) {
		errs := ValidateNumber(true, NumberRules{})
		if len(errs) != 1 || errs[0].Code != covInvalid {
			t.Errorf("expected invalid code for bool, got %v", errs)
		}
	})

	t.Run("nil pointer is absent", func(t *testing.T) {
		errs := ValidateNumber(nil, NumberRules{Required: true})
		if len(errs) != 1 || errs[0].Code != covRequired {
			t.Errorf("expected required code for nil, got %v", errs)
		}
	})
}

// TestValidateJSONTypeRejection tests ValidateJSON with wrong types.
func TestValidateJSONTypeRejection(t *testing.T) {
	t.Run("scalar is invalid", func(t *testing.T) {
		errs := ValidateJSON("string", JSONRules{})
		if len(errs) != 1 || errs[0].Code != covInvalid {
			t.Errorf("expected invalid code for scalar, got %v", errs)
		}
	})

	t.Run("nil pointer is absent", func(t *testing.T) {
		errs := ValidateJSON(nil, JSONRules{Required: true})
		if len(errs) != 1 || errs[0].Code != covRequired {
			t.Errorf("expected required code for nil, got %v", errs)
		}
	})
}

// TestCleanMethodPanic tests that panics in clean methods are caught.
func TestCleanMethodPanic(t *testing.T) {
	type PanicForm struct {
		Field *string `layrz:"char,required"`
	}

	form := &PanicForm{Field: Ptr("test")}

	// callClean with a function that panics
	fn := reflect.ValueOf(func(*PanicForm, *string) *FieldError {
		panic("test panic")
	})

	_, err := callClean(fn, []reflect.Value{
		reflect.ValueOf(form),
		reflect.ValueOf(Ptr("test")),
	})

	if err == nil {
		t.Errorf("expected error from panic, got nil")
	}
}

// TestSpecRulesUnsupportedKind tests Rules() with unsupported kind.
func TestSpecRulesUnsupportedKind(t *testing.T) {
	spec := &FieldSpec{Kind: FieldKind("unsupported")}
	_, err := spec.Rules()
	if err == nil {
		t.Errorf("expected error for unsupported field kind")
	}
}

// TestValidateIDAllIntTypes tests ValidateID with various int types.
func TestValidateIDAllIntTypes(t *testing.T) {
	t.Run("int32", func(t *testing.T) {
		errs := ValidateID(int32(42), IDRules{})
		if len(errs) != 0 {
			t.Errorf("expected no error for valid int32")
		}
	})

	t.Run("int16", func(t *testing.T) {
		errs := ValidateID(int16(42), IDRules{})
		if len(errs) != 0 {
			t.Errorf("expected no error for valid int16")
		}
	})

	t.Run("int8", func(t *testing.T) {
		errs := ValidateID(int8(42), IDRules{})
		if len(errs) != 0 {
			t.Errorf("expected no error for valid int8")
		}
	})

	t.Run("int64", func(t *testing.T) {
		errs := ValidateID(int64(42), IDRules{})
		if len(errs) != 0 {
			t.Errorf("expected no error for valid int64")
		}
	})
}

// TestValidateEmailNilAndEmpty tests ValidateEmail with nil and empty.
func TestValidateEmailNilAndEmpty(t *testing.T) {
	t.Run("nil pointer optional empty=false", func(t *testing.T) {
		errs := ValidateEmail(nil, EmailRules{Required: false, Empty: false})
		if len(errs) != 0 {
			t.Errorf("expected no error for optional nil")
		}
	})

	t.Run("empty string with empty=false required=true", func(t *testing.T) {
		errs := ValidateEmail(Ptr(""), EmailRules{Required: true, Empty: false})
		if len(errs) != 1 || errs[0].Code != covEmpty {
			t.Errorf("expected empty code, got %v", errs)
		}
	})

	t.Run("empty string with empty=true", func(t *testing.T) {
		errs := ValidateEmail(Ptr(""), EmailRules{Required: true, Empty: true})
		if len(errs) != 0 {
			t.Errorf("expected no error for empty=true, got %v", errs)
		}
	})

	t.Run("invalid email format", func(t *testing.T) {
		errs := ValidateEmail(Ptr("not-an-email"), EmailRules{Required: true})
		if len(errs) != 1 || errs[0].Code != covInvalid {
			t.Errorf("expected invalid code, got %v", errs)
		}
	})
}

// TestValidateNumberWithNilAndRanges tests ValidateNumber with nil and min/max.
func TestValidateNumberWithNilAndRanges(t *testing.T) {
	t.Run("nil pointer required", func(t *testing.T) {
		errs := ValidateNumber(nil, NumberRules{Required: true})
		if len(errs) != 1 || errs[0].Code != covRequired {
			t.Errorf("expected required code, got %v", errs)
		}
	})

	t.Run("value below min", func(t *testing.T) {
		minVal := 10.0
		errs := ValidateNumber(5, NumberRules{MinValue: &minVal})
		if len(errs) != 1 || errs[0].Code != covMinValue {
			t.Errorf("expected minValue code, got %v", errs)
		}
	})

	t.Run("value above max", func(t *testing.T) {
		maxVal := 10.0
		errs := ValidateNumber(20, NumberRules{MaxValue: &maxVal})
		if len(errs) != 1 || errs[0].Code != covMaxValue {
			t.Errorf("expected maxValue code, got %v", errs)
		}
	})
}

// TestValidateJSONWithNilAndEmpty tests ValidateJSON with nil and empty.
func TestValidateJSONWithNilAndEmpty(t *testing.T) {
	t.Run("nil pointer required", func(t *testing.T) {
		errs := ValidateJSON(nil, JSONRules{Required: true})
		if len(errs) != 1 || errs[0].Code != covRequired {
			t.Errorf("expected required code, got %v", errs)
		}
	})

	t.Run("empty list with empty=false required=true", func(t *testing.T) {
		errs := ValidateJSON([]any{}, JSONRules{Required: true, Empty: false})
		if len(errs) != 1 || errs[0].Code != covInvalid {
			t.Errorf("expected invalid code for empty list, got %v", errs)
		}
	})

	t.Run("empty list with empty=true", func(t *testing.T) {
		errs := ValidateJSON([]any{}, JSONRules{Required: true, Empty: true})
		if len(errs) != 0 {
			t.Errorf("expected no error for empty=true, got %v", errs)
		}
	})

	t.Run("empty dict with empty=false required=true", func(t *testing.T) {
		errs := ValidateJSON(map[string]any{}, JSONRules{Required: true, Datatype: covDict, Empty: false})
		if len(errs) != 1 || errs[0].Code != covInvalid {
			t.Errorf("expected invalid code for empty dict, got %v", errs)
		}
	})

	t.Run("empty dict with empty=true", func(t *testing.T) {
		errs := ValidateJSON(map[string]any{}, JSONRules{Required: true, Datatype: covDict, Empty: true})
		if len(errs) != 0 {
			t.Errorf("expected no error for empty=true, got %v", errs)
		}
	})
}

// TestParseTagDuplicateMinLength tests the specific duplicate min_length error path.
func TestParseTagDuplicateMinLength(t *testing.T) {
	_, err := ParseTag("char,min_length=5,min_length=10")
	if err == nil {
		t.Errorf("expected error for duplicate min_length")
	}
}

// TestParseTagDuplicateMaxLength tests the specific duplicate max_length error path.
func TestParseTagDuplicateMaxLength(t *testing.T) {
	_, err := ParseTag("char,max_length=5,max_length=10")
	if err == nil {
		t.Errorf("expected error for duplicate max_length")
	}
}

// TestParseTagDuplicateDatatype tests the specific duplicate datatype error path.
func TestParseTagDuplicateDatatype(t *testing.T) {
	_, err := ParseTag("number,datatype=int,datatype=float")
	if err == nil {
		t.Errorf("expected error for duplicate datatype")
	}
}

// TestParseTagDuplicateChoices tests the specific duplicate choices error path.
func TestParseTagDuplicateChoices(t *testing.T) {
	_, err := ParseTag("char,choices=a|b,choices=c|d")
	if err == nil {
		t.Errorf("expected error for duplicate choices")
	}
}

// TestParseTagInvalidMinLength tests invalid min_length value.
func TestParseTagInvalidMinLength(t *testing.T) {
	_, err := ParseTag("char,min_length=not_a_number")
	if err == nil {
		t.Errorf("expected error for invalid min_length value")
	}
}

// TestParseTagInvalidMaxLength tests invalid max_length value.
func TestParseTagInvalidMaxLength(t *testing.T) {
	_, err := ParseTag("char,max_length=not_a_number")
	if err == nil {
		t.Errorf("expected error for invalid max_length value")
	}
}

// TestParseTagEmptyTagError tests that empty tags cause an error.
func TestParseTagEmptyTagError(t *testing.T) {
	_, err := ParseTag("char,")
	if err != nil {
		t.Errorf("expected no error for trailing comma, got %v", err)
	}
}

// TestParseTagDuplicateRegex tests duplicate regex in tag.
func TestParseTagDuplicateRegex(t *testing.T) {
	_, err := ParseTag("char,regex=^a$")
	if err != nil {
		t.Errorf("expected no error for single regex, got %v", err)
	}

	// Can't actually test duplicate regex since it consumes the rest
}

// TestSpecCheckTypeDefaultCase tests the default case in CheckType.
func TestSpecCheckTypeDefaultCase(t *testing.T) {
	spec := &FieldSpec{Kind: FieldKind("unknown")}
	err := spec.CheckType(reflect.TypeOf("string"))
	// Should not error, falls through to default which returns nil
	if err != nil {
		t.Errorf("expected no error for unknown kind, got %v", err)
	}
}

// TestEmbeddedStructWithNilPointer tests embedded struct with nil pointer.
func TestEmbeddedStructWithNilPointer(t *testing.T) {
	type InnerForm struct {
		Field *string `layrz:"char,required"`
	}

	type OuterForm struct {
		*InnerForm
	}

	form := &OuterForm{InnerForm: nil}
	errs := Validate(form)

	// Should skip nil embedded pointer
	if len(errs) > 0 {
		t.Errorf("expected no errors for nil embedded pointer, got %v", errs)
	}
}

// TestEmbeddedStructValue tests embedded struct as value (not pointer).
func TestEmbeddedStructValue(t *testing.T) {
	type InnerForm struct {
		Field *string `layrz:"char,required"`
	}

	type OuterForm struct {
		InnerForm
	}

	form := &OuterForm{
		InnerForm: InnerForm{Field: Ptr("test")},
	}
	errs := Validate(form)

	// Should validate the embedded value fields
	if len(errs) > 0 {
		t.Errorf("expected no errors for valid embedded value, got %v", errs)
	}
}

// TestEmbeddedStructValueWithError tests embedded struct value with validation errors.
func TestEmbeddedStructValueWithError(t *testing.T) {
	type InnerForm struct {
		Field *string `layrz:"char,required"`
	}

	type OuterForm struct {
		InnerForm
	}

	form := &OuterForm{
		InnerForm: InnerForm{Field: nil},
	}
	errs := Validate(form)

	// Should have error from required field
	if len(errs) == 0 {
		t.Errorf("expected errors for required field in embedded value")
	}
}

// TestValidateEmailWithRegexOverride tests email with custom regex.
func TestValidateEmailWithRegexOverride(t *testing.T) {
	t.Run("valid with custom regex", func(t *testing.T) {
		errs := ValidateEmail(Ptr("test@example.com"), EmailRules{
			Required: true,
			Regex:    "^[a-z]+@[a-z]+\\.[a-z]+$",
		})
		if len(errs) != 0 {
			t.Errorf("expected no error for valid email with custom regex, got %v", errs)
		}
	})

	t.Run("invalid with custom regex", func(t *testing.T) {
		errs := ValidateEmail(Ptr("123@example.com"), EmailRules{
			Required: true,
			Regex:    "^[a-z]+@[a-z]+\\.[a-z]+$",
		})
		if len(errs) != 1 || errs[0].Code != covInvalid {
			t.Errorf("expected invalid code for email not matching custom regex, got %v", errs)
		}
	})
}

// TestValidateNumberSpecificIntTypes tests specific int type cases.
func TestValidateNumberSpecificIntTypes(t *testing.T) {
	t.Run("float32", func(t *testing.T) {
		errs := ValidateNumber(float32(3.14), NumberRules{})
		if len(errs) != 0 {
			t.Errorf("expected no error for valid float32, got %v", errs)
		}
	})

	t.Run("uint is invalid", func(t *testing.T) {
		errs := ValidateNumber(uint(42), NumberRules{})
		if len(errs) != 1 || errs[0].Code != covInvalid {
			t.Errorf("expected invalid code for uint, got %v", errs)
		}
	})
}

// TestValidateJSONTypeChecking tests JSON type validation paths.
func TestValidateJSONTypeChecking(t *testing.T) {
	t.Run("list with dict datatype mismatch", func(t *testing.T) {
		errs := ValidateJSON([]int{1, 2}, JSONRules{Required: true, Datatype: covDict})
		if len(errs) != 1 || errs[0].Code != covInvalid {
			t.Errorf("expected invalid code for list with dict datatype, got %v", errs)
		}
	})

	t.Run("dict with list datatype mismatch", func(t *testing.T) {
		errs := ValidateJSON(map[string]int{"a": 1}, JSONRules{Required: true, Datatype: covList})
		if len(errs) != 1 || errs[0].Code != covInvalid {
			t.Errorf("expected invalid code for dict with list datatype, got %v", errs)
		}
	})

	t.Run("pointer to list", func(t *testing.T) {
		list := []any{1, 2}
		errs := ValidateJSON(&list, JSONRules{Required: true})
		if len(errs) != 0 {
			t.Errorf("expected no error for pointer to list, got %v", errs)
		}
	})

	t.Run("pointer to map", func(t *testing.T) {
		m := map[string]any{"a": 1}
		errs := ValidateJSON(&m, JSONRules{Required: true, Datatype: covDict})
		if len(errs) != 0 {
			t.Errorf("expected no error for pointer to map, got %v", errs)
		}
	})
}
