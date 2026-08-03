# Layrz Forms (Go)

A form validation library for Go with a simple API inspired by Django Forms. Validates struct fields based on tags, supports nested forms, custom validation methods, and accumulates all errors for a complete validation report.

## Install

```bash
go get github.com/goldenm-software/layrz-forms/go
```

## Usage

### Basic Example

```go
package main

import (
	"fmt"
	layrz "github.com/goldenm-software/layrz-forms/go"
)

type Address struct {
	StreetName *string `layrz:"char,required,min_length=5"`
	ZipCode    *string `layrz:"char,required"`
}

type Item struct {
	Name  *string  `layrz:"char,required"`
	Price *float64 `layrz:"number,required,min_value=0"`
}

type ExampleForm struct {
	IdTest        *int            `layrz:"id,required"`
	EmailText     *string         `layrz:"email,required"`
	JsonListTest  *[]any          `layrz:"json,required,datatype=list"`
	JsonDictTest  *map[string]any `layrz:"json,required,datatype=dict"`
	IntTest       *int            `layrz:"number,required,min_value=0,max_value=5"`
	FloatTest     *float64        `layrz:"number,required,min_value=0,max_value=5"`
	BoolTest      *bool           `layrz:"bool,required"`
	PlainTextTest *string         `layrz:"char,required"`
	EmptyTextTest *string         `layrz:"char,required,empty"`
	RangeTextTest *string         `layrz:"char,required,min_length=5,max_length=10"`
	Address       *Address        `layrz:"subform"`
	Items         []*Item         `layrz:"subform_list"`
}

// Convention B clean methods (cross-field validation)
func (f *ExampleForm) CleanFunc1() layrz.Errors {
	return layrz.Errors{
		"clean1": {
			{Code: "error1"},
			{Code: "error2"},
		},
	}
}

func (f *ExampleForm) CleanFunc2() layrz.Errors {
	return layrz.Errors{
		"clean2": {
			{Code: "error1"},
		},
	}
}

// Convention A clean method (per-field validation)
func (f *ExampleForm) CleanEmailText(value *string) *layrz.FieldError {
	if value != nil && len(*value) > 13 && (*value)[len(*value)-13:] == "@blocked.com" {
		return &layrz.FieldError{Code: "blockedDomain"}
	}
	return nil
}

func main() {
	form := &ExampleForm{
		IdTest:        layrz.Ptr(1),
		EmailText:     layrz.Ptr("example@goldenmcorp.com"),
		JsonDictTest:  layrz.Ptr(map[string]any{"hola": "mundo"}),
		JsonListTest:  layrz.Ptr([]any{"hola mundo"}),
		IntTest:       layrz.Ptr(5),
		FloatTest:     layrz.Ptr(4.5),
		BoolTest:      layrz.Ptr(true),
		PlainTextTest: layrz.Ptr("hola mundo"),
		EmptyTextTest: layrz.Ptr("hola"),
		RangeTextTest: layrz.Ptr("hola"), // 4 chars, min_length=5
	}

	errs := layrz.Validate(form)

	fmt.Println("IsValid:", layrz.IsValid(form))
	// Output: IsValid: false

	fmt.Println("Errors:")
	for _, key := range errs.Keys() {
		fmt.Printf("  %s: ", key)
		for i, err := range errs[key] {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Printf("{Code: %q, Expected: %v, Received: %v}", err.Code, err.Expected, err.Received)
		}
		fmt.Println()
	}
}
```

Output:
```
IsValid: false
Errors:
  clean1: {Code: "error1", Expected: <nil>, Received: <nil>}, {Code: "error2", Expected: <nil>, Received: <nil>}
  clean2: {Code: "error1", Expected: <nil>, Received: <nil>}
  rangeTextTest: {Code: "minLength", Expected: 5, Received: 4}
```

## Tag Grammar

Field validation rules are declared using the `layrz` struct tag. The grammar is:

```
layrz:"<kind>{,<rule>}{,<rule>}..."
```

Where `<kind>` is one of:
- `id` — Integer ID (must be > 0)
- `email` — Email address
- `uuid` — UUID
- `char` — String/Character field
- `number` — Numeric field (int or float)
- `bool` — Boolean field
- `json` — JSON object or array (list or dict)
- `subform` — Nested struct (pointer to struct)
- `subform_list` — List of nested structs (slice of struct or slice of pointer-to-struct)

Rules are comma-separated and include:

| Rule | Applies To | Description |
|------|-----------|-------------|
| `required` | All | Value must be present (non-nil) |
| `empty` | char, email, json | Empty string/container is valid |
| `min_length=<int>` | char, json(list) | Minimum length in characters (for char) or items (for list) |
| `max_length=<int>` | char, json(list) | Maximum length in characters (for char) or items (for list) |
| `min_value=<float>` | number, id | Minimum numeric value |
| `max_value=<float>` | number, id | Maximum numeric value |
| `regex=<pattern>` | char, email | Regex pattern to match (PCRE, must be last rule) |
| `choices=<a\|b\|c>` | char | Pipe-separated valid choices |
| `datatype=<type>` | number, json | Type specification: `int`/`float` for number, `list`/`dict` for json |

### Examples

```go
type User struct {
	// Required fields
	ID       *int    `layrz:"id,required"`
	Email    *string `layrz:"email,required"`
	Username *string `layrz:"char,required,min_length=3,max_length=20"`

	// Optional fields
	Bio           *string `layrz:"char,empty"` // Can be empty string
	Age           *int    `layrz:"number,min_value=0,max_value=150"`
	PreferredLang *string `layrz:"char,choices=en|es|pt"`

	// JSON fields
	Profile *map[string]any `layrz:"json,datatype=dict"`
	Tags    *[]any          `layrz:"json,datatype=list"`

	// Nested forms
	Address *Address `layrz:"subform"`
	Orders  []*Order `layrz:"subform_list"`
}
```

## Field Types & Validators

### ID Field (`id`)
Validates positive integers. Accepts:
- Go integers: `int`, `int8`–`int64`
- Strings parseable as positive integers
- Rejects: floats, negative numbers, zero, non-numeric strings, booleans

### Email Field (`email`)
Validates email addresses against a regex pattern. Empty strings rejected unless `empty:true`.
- Default regex: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-z]{2,63}$`
- Custom regex via `regex=<pattern>`

### UUID Field (`uuid`)
Validates UUID in any of these formats:
- Hyphenated: `8-4-4-4-12` (e.g., `550e8400-e29b-41d4-a716-446655440000`)
- Unhyphenated: 32 hex characters
- Braced: `{8-4-4-4-12}` or `{32-hex}`
- URN prefix: `urn:uuid:8-4-4-4-12`

### Char Field (`char`)
Validates strings. Length is measured in UTF-8 rune count, not bytes.
- `min_length=<int>` enforces minimum rune count
- `max_length=<int>` enforces maximum rune count
- `regex=<pattern>` validates against a regex
- `choices=<a|b|c>` restricts to specific values
- `empty:true` allows empty strings (default rejects)

### Number Field (`number`)
Validates numeric values. Requires `datatype=int` or `datatype=float`:
- `int`: Accepts `*int`, `*int8`–`*int64`; rejects floats
- `float`: Accepts `*float32`, `*float64`; rejects ints
- `min_value=<float>` and `max_value=<float>` enforce bounds
- Rejects booleans (Go differs from Python here)

### Bool Field (`bool`)
Validates boolean values. Simply checks that the value is a `*bool`.

### JSON Field (`json`)
Validates JSON objects (dicts) or arrays (lists). Requires `datatype=list` or `datatype=dict`:
- `list`: Accepts `*[]T`, rejects maps
- `dict`: Accepts `*map[K]V`, rejects slices
- `empty:true` allows empty containers
- `datatype` is inferred from the Go type if not specified

### Subform (`subform`)
Validates a nested struct. Must be a non-nil pointer to a struct. Nil subforms are **skipped** (divergent from Python, which validates against an empty dict).

### Subform List (`subform_list`)
Validates a slice of structs or pointers-to-structs. Errors are keyed with numeric indices (e.g., `items.0.name`, `items.1.price`).
- Nil slice elements are skipped
- Empty slices produce no errors

## Validation & Clean Methods

### Entry Points

```go
// Validate returns the complete error map
errs := layrz.Validate(&form)

// IsValid reports whether validation passed
if layrz.IsValid(&form) {
	fmt.Println("Form is valid!")
}
```

### Custom Validation with Clean Methods

**Convention A** (per-field): Validate a single field based on its current value.

```go
func (f *MyForm) Clean<FieldName>(value <FieldType>) *FieldError {
	if /* invalid */ {
		return &FieldError{Code: "customError", Expected: "...", Received: value}
	}
	return nil
}
```

Example: Reject emails from a blocked domain.

```go
func (f *UserForm) CleanEmail(value *string) *FieldError {
	if value != nil && strings.HasSuffix(*value, "@spam.com") {
		return &FieldError{Code: "blockedDomain"}
	}
	return nil
}
```

**Convention B** (cross-field): Validate relationships between fields.

```go
func (f *MyForm) Clean<Suffix>() Errors {
	errs := make(Errors)
	if /* password mismatch */ {
		errs.Add("password_confirm", &FieldError{Code: "mismatch"})
	}
	return errs
}
```

Example: Check that password and confirmation match.

```go
func (f *UserForm) CleanPasswords() Errors {
	if f.Password != nil && f.PasswordConfirm != nil && *f.Password != *f.PasswordConfirm {
		return Errors{
			"passwordConfirm": {
				{Code: "mustMatch"},
			},
		}
	}
	return nil
}
```

**Execution Order:**
1. Built-in tag validation (all fields and nested forms)
2. Convention A clean methods (alphabetical order)
3. Convention B clean methods (alphabetical order)

All errors are accumulated; validation does not fail-fast.

## Error Structure

Errors are returned as `map[string][]*FieldError` where each field name maps to a slice of errors.

```go
type FieldError struct {
	Code     string         `json:"code"`
	Expected any            `json:"expected,omitempty"`
	Received any            `json:"received,omitempty"`
	Extra    map[string]any `json:"extra,omitempty"`
}
```

- **Code**: Error code (e.g., `"required"`, `"minLength"`, `"invalid"`)
- **Expected**: The constraint that was violated (e.g., `5` for a min_length error)
- **Received**: The offending value
- **Extra**: Additional context (e.g., error metadata from custom validators)

### JSON Serialization

Errors marshal to JSON with `omitempty` on zero fields:

```go
import "encoding/json"

errs := layrz.Validate(&form)
data, _ := json.Marshal(errs)
// Output: {"field":{"code":"minLength","expected":5,"received":4}}
```

Fields with zero/nil values are omitted, ensuring clean JSON output.

## Nested Forms

Subforms and subform lists recursively validate their nested structs. Error keys are prefixed with the subform path using dot notation:

```go
type Company struct {
	Name *string `layrz:"char,required"`
}

type Employee struct {
	Name    *string `layrz:"char,required"`
	Company *Company `layrz:"subform"`
}

errs := layrz.Validate(&emp)
// Errors keyed as:
//   "name" (employee's name)
//   "company.name" (company's name)
```

### Subform Lists

List errors are keyed with numeric indices:

```go
type Order struct {
	Items []*Item `layrz:"subform_list"`
}

errs := layrz.Validate(&order)
// Errors keyed as:
//   "items.0.name" (first item's name)
//   "items.1.price" (second item's price)
//   etc.
```

## Differences from the Python Implementation

This Go implementation diverges from the Python version in several ways, mostly bug fixes:

1. **Absent optional field**: Go emits no errors; Python's `JsonField` wrongly emitted `invalid`.
2. **Wrong type**: Always `invalid`, regardless of `required` flag. Python's optional fields silently accepted wrong types.
3. **Empty string with `empty:false`**: Go emits `empty` error. Python's `EmailField` wrongly said `required`.
4. **Empty string with `empty:true`**: Go permits `""` but still regex-validates non-empty values. Python skipped regex entirely if `empty:true`.
5. **Boolean as number or ID**: Go always rejects `bool`. Python accepted `True` because `isinstance(True, int)`.
6. **Nil subform pointer**: Go skips it entirely (no errors). Python validated against an empty dict and emitted its `required` errors.
7. **Async validation**: Go has no equivalent of Python's `ais_valid()` (async); all methods are synchronous.
8. **Future alignment**: The Python implementation still has these bugs and will be corrected in a future 4.0.0 release.

## License

MIT License. See the repository for details.

---

Maintained by [Golden M](https://goldenm.com) with authorization of [Layrz LTD](https://layrz.com).
