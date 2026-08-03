# Form Validation Vectors

This directory contains cross-language test vectors for form field validation. These vectors are designed to be consumed by test suites in any language (Python, Go, etc.) to ensure consistent behavior across implementations.

## Schema

Each field type has a JSON file containing an array of test cases. The schema for each test case is:

```json
{
  "name": "unique_case_identifier",
  "field": "FieldTypeName",
  "kwargs": {
    "required": true,
    "other_param": "value"
  },
  "value": "the_test_value",
  "expected_errors": [
    {"code": "errorCode", "expected": 5, "received": 4}
  ]
}
```

### Key Fields

- **name** (string): Unique identifier for this test case (snake_case)
- **field** (string): The field type name (e.g., "CharField", "NumberField")
- **kwargs** (object): Constructor parameters for the field. Use Python types:
  - `true`/`false` for booleans
  - Numbers for ints/floats
  - Strings for string types
  - `"dict"` or `"list"` for datatype parameters (not JSON object/array literals)
- **value** (any): The value to validate. Omit this key entirely if testing the absent case (use `value_absent` instead).
- **value_absent** (boolean): If `true`, the field is not present in the input dict. Mutually exclusive with `value`.
- **expected_errors** (array): Array of error objects. Each error object:
  - **code** (string): The error code (e.g., "required", "minLength", "invalid")
  - Other keys are extra fields like `expected`, `received`, `message`, etc.
- Empty array `[]` means validation passes with no errors.

## Special Cases

### Datatype Parameters

For fields like `JsonField` and `NumberField` that take a Python type, the JSON vectors use string representations:
- `"dict"` for `dict` type
- `"list"` for `list` type
- `"int"` for `int` type
- `"float"` for `float` type

The test suite must convert these to the actual Python types when constructing the field.

### Choices Parameter

Choices are represented as a 2D array (list of tuples):
```json
"choices": [["value1", "Label1"], ["value2", "Label2"]]
```

### Absent vs. None

- **Absent**: Field key is not in the input dict at all. Use `"value_absent": true`.
- **None**: Field key is present but has a `null` value. Use `"value": null`.

These are different and must be tested separately.

## Consuming Vectors

A language-specific test suite should:

1. Load each JSON file from this directory
2. For each test case:
   - Create a Form with the field type and kwargs
   - If `value_absent` is true, don't include the field in the input dict
   - If `value` is present (including `null`), include it in the input dict
   - Call `form.is_valid()` (or equivalent)
   - Assert that `form.errors()` matches `expected_errors` exactly

Example (Python):

```python
with open('CharField.json') as f:
    cases = json.load(f)

for case in cases:
    field_class = getattr(forms, case['field'])
    
    class TestForm(Form):
        field_name = field_class(**case['kwargs'])
    
    input_dict = {}
    if 'value' in case:
        input_dict['field_name'] = case['value']
    # If 'value_absent' is true, don't add the key
    
    form = TestForm(input_dict)
    form.is_valid()
    assert form.errors() == {'fieldName': case['expected_errors']}
```

## Files

- `BooleanField.json` - BooleanField test cases
- `CharField.json` - CharField test cases
- `EmailField.json` - EmailField test cases
- `IdField.json` - IdField test cases
- `JsonField.json` - JsonField test cases
- `NumberField.json` - NumberField test cases
- `UuidField.json` - UuidField test cases
