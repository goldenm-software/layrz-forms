---
name: go-form-builder
description: Build and validate forms with the layrz-forms Go library. Use for struct tags, Validate(), IsValid(), Clean<Field> hooks, FieldError, error codes, pointer vs. value fields, GraphQL inputs, and subform recursion. GO-SPECIFIC only — see form-builder router for Python.
---

# Go form builder

layrz-forms is a struct validation library. A form is a pointer to a struct with exported fields tagged with `layrz:`. Call `Validate(&form)` to get an `Errors` map (field → error slices), keyed in camelCase. This guide covers the Go API and the exact semantics that make pointer and value fields work differently.

## Import and entry points

```go
import "github.com/goldenm-software/layrz-forms/go/v3"
```

**Key exports:**
- `Validate(form any) Errors` — validates and returns all errors, keyed by camelCase field name. Form MUST be a non-nil pointer to a struct; passing a value or nil yields a `_config` error.
- `IsValid(form any) bool` — shorthand for `Validate(form).IsEmpty()`.
- `FieldError` — single error with `Code` (string), `Expected` (any), `Received` (any), `Extra` (map[string]any).
- `Errors` — map[string][]*FieldError, the result type.

The form argument must always be a **NON-NIL POINTER TO A STRUCT**. This requirement exists because Clean methods have pointer receivers — without a pointer, receivers are invisible to reflection, so the engine cannot find or invoke them. Bad input yields a `_config` error instead of a panic.

```go
// WRONG: passing a value
errs := Validate(form)  // _config error: "form must be a pointer"

// CORRECT: pass the address
errs := Validate(&form)
```

## Struct tag grammar

Tags are declared on struct fields as `layrz:"<rules>"`. The first token is the field kind; the rest are optional rules separated by commas.

**Grammar:**
```
layrz:"<kind>[,rule][,rule]..."
```

**Field kinds (mutually exclusive):**
- `id` — int or string ID; required to be > 0
- `email` — email string
- `uuid` — UUID string (hyphenated, unhyphenated, braced, or urn: prefix formats accepted)
- `char` — arbitrary string
- `number` — int or float numeric value
- `bool` — boolean
- `json` — slice or map (JSON-like container)
- `subform` — pointer to a nested struct (validated recursively)
- `subform_list` — slice of struct or slice of pointer-to-struct

**Rules (applied to the kind):**
- `required` — field is absent if nil (pointer fields) or missing (value fields never absent). No value.
- `empty` — empty string or container is OK. Bare keyword. Default: false (empty rejected).
- `min_length=<int>` — minimum string length in runes.
- `max_length=<int>` — maximum string length in runes.
- `min_value=<float>` — minimum numeric value.
- `max_value=<float>` — maximum numeric value.
- `choices=<a|b|c>` — string must be one of the pipe-separated values.
- `regex=<pattern>` — string must match pattern. **MUST BE LAST** in the tag because its value may contain commas. Pattern is compiled; invalid patterns produce a `_config` error.
- `datatype=<type>` — disambiguate numeric or JSON type. Values: `int` (integers only), `float` (floats only), `list` (arrays), `dict` (maps). Mostly inferred from Go type; use to override.
- `-` — skip this field entirely.

**Common mistakes:**
- Using `:` instead of `=`: `required:true` is wrong; use bare `required`.
- Forgetting `=`: `min_length5` is wrong; use `min_length=5`.
- `regex=` not last: `layrz:"char,regex=[0-9]+,required"` fails because regex consumes the rest of the string. Move it to the end.

**Examples:**
```go
type User struct {
    ID          *int    `layrz:"id,required"`
    Email       *string `layrz:"email,required"`
    Name        *string `layrz:"char,required,min_length=2,max_length=50"`
    Age         *int    `layrz:"number,min_value=0,max_value=150"`
    Status      *string `layrz:"char,choices=active|inactive|pending"`
    Description *string `layrz:"char,regex=^[a-z0-9]+$"`  // regex last
    Active      *bool   `layrz:"bool,required"`
    Tags        *[]any  `layrz:"json,datatype=list"`
    Metadata    *map[string]any `layrz:"json,datatype=dict"`
    Address     *Address `layrz:"subform"`  // nested
    Items       []*Item  `layrz:"subform_list"`  // nested list
    Ignored     string   `layrz:"-"`  // not validated
    Untagged    string   // skipped; no tag
}
```

## Pointer vs. value fields: the key concept

This is the single most important semantic difference from Python. In Go, a POINTER field models **absence**, and a VALUE field is **always present**.

### Pointer fields (absence is nil)

A pointer field is nil when absent.

```go
type Form struct {
    Email *string `layrz:"email,required"`
    Age   *int    `layrz:"number,required"`
}

// nil Email → required error
// nil Age → required error
```

- Nil → absent → `required` fires if rule is set, otherwise no error
- Non-nil → present → validate the pointed value

### Value fields (always present, even at zero)

A value field has a Go zero value that counts as PRESENT. A zero string `""` is present. A zero int `0` is present.

```go
type Form struct {
    Name  string `layrz:"char,required"`      // "" is PRESENT, not required
    Count int    `layrz:"number,required,min_value=1"`  // 0 is PRESENT
}

// Name = "" with required → NO error (it's present)
// Name = "" with empty:false → empty error (present but blank)
// Count = 0 with required → NO error (it's present)
// Count = 0 with min_value=1 → minValue error (value too low)
```

- `required` **never fires** on a value field because the zero value is always present
- `empty` **does fire** on a zero-length string or empty slice/map
- Other rules (min/max, regex, choices) apply to the zero value

### Consequences

**Example 1: a required string field**
```go
// Pointer: absence is error
type F1 struct {
    Name *string `layrz:"char,required"`
}
f := &F1{Name: nil}         // error: required
f = &F1{Name: Ptr("")}      // error: empty (present but blank)
f = &F1{Name: Ptr("Bob")}   // OK

// Value: zero value is OK, but empty is rejected
type F2 struct {
    Name string `layrz:"char,required"`
}
f := &F2{Name: ""}          // error: empty (required does NOT fire)
f = &F2{Name: "Bob"}        // OK
```

**Example 2: a number with min bound**
```go
// Pointer: absence skips min_value check
type F1 struct {
    Count *int64 `layrz:"number,min_value=1"`
}
f := &F1{Count: nil}        // OK (absent, no bound check)
f = &F1{Count: Ptr(0)}      // error: minValue

// Value: zero passes required but fails min_value
type F2 struct {
    Count int64 `layrz:"number,min_value=1"`
}
f := &F2{Count: 0}          // error: minValue (zero violates bound)
f = &F2{Count: 1}           // OK
```

**Example 3: a bool field**
```go
// Pointer: nil is absent
type F1 struct {
    Active *bool `layrz:"bool,required"`
}
f := &F1{Active: nil}       // error: required
f = &F1{Active: Ptr(false)} // OK (false is a valid bool)

// Value: false is present and always valid (no min/max for bool)
type F2 struct {
    Active bool `layrz:"bool,required"`
}
f := &F2{Active: false}     // OK (false is present)
f = &F2{Active: true}       // OK
```

### Accepted types per kind

| Kind | Pointer Forms | Value Forms |
|------|---|---|
| `id` | `*int`, `*int64`, `*string` | `int`, `int64`, `string` |
| `email` | `*string` | `string` |
| `uuid` | `*string` | `string` |
| `char` | `*string` | `string` |
| `number` | `*int`, `*int8...64`, `*float32`, `*float64` | `int`, `int8...64`, `float32`, `float64` |
| `bool` | `*bool` | `bool` |
| `json` | `*[]any`, `*map[string]any`, `*[]T`, `*map[string]T` | `[]any`, `map[string]any`, `[]T`, `map[string]T` |
| `subform` | `*struct` | (not supported; must be pointer) |
| `subform_list` | `[]*struct`, `[]struct` | N/A |

## GraphQL interop

layrz-forms integrates cleanly with `graph-gophers/graphql-go` input types. GraphQL matches fields to Go by name (case-insensitive, underscores stripped in the matching), NOT by `json:` or `layrz:` tags. This means validation tags never collide with GraphQL.

**Key rule:** GraphQL rejects a pointer for a non-null (`!`) input field. A `String!` requires a value field; a nullable `String` maps to a pointer field.

**Example GraphQL schema:**
```graphql
input CreateUserInput {
    id: Int!          # non-null, must be value field
    email: String!    # non-null, must be value field
    bio: String       # nullable, must be pointer field
    age: Int          # nullable, must be pointer field
}
```

**Corresponding Go struct:**
```go
type CreateUserInput struct {
    ID    int     `layrz:"id,required"`           // value: non-null
    Email string  `layrz:"email,required"`        // value: non-null
    Bio   *string `layrz:"char"`                  // pointer: nullable
    Age   *int    `layrz:"number,min_value=0"`   // pointer: nullable
}
```

GraphQL sets fields by name matching and populates the struct. Then call `Validate(&input)`. Custom scalars (like a UUID struct or Status enum) should be left untagged and validated in a Convention A hook:

```go
// Custom scalar: a UUID struct
type UUID struct {
    Value string
    Valid bool
}

type Form struct {
    UserID UUID `layrz:"-"`  // untagged; no built-in validation
}

// Convention A hook: validate by exact type
func (f *Form) CleanUserID(value UUID) *FieldError {
    if !value.Valid {
        return &FieldError{Code: "invalid"}
    }
    // Further validation if needed
    return nil
}
```

## Clean hooks: per-field and cross-field

After built-in field rules, the engine discovers and calls clean methods. Two conventions:

### Convention A: per-field

**Signature:** `func (f *MyForm) Clean<FieldName>(value <FieldType>) *FieldError`

- Method name is `Clean` + the Go struct field name exactly (e.g., `CleanEmail` for field `Email`).
- Parameter type **must match the field's type exactly** (pointer or value).
- Returns `*FieldError` or nil.
- Error files under the field's **camelCase** key.
- Used for single-field validation logic.

**Example:**
```go
type LoginForm struct {
    Email    *string `layrz:"email,required"`
    Password *string `layrz:"char,required,min_length=8"`
}

// Convention A: validate Email
func (f *LoginForm) CleanEmail(value *string) *FieldError {
    if value == nil {
        return nil  // required already caught this
    }
    if strings.HasSuffix(*value, "@blocked.com") {
        return &FieldError{Code: "blockedDomain"}
    }
    return nil
}

// Convention A: validate Password
func (f *LoginForm) CleanPassword(value *string) *FieldError {
    if value == nil {
        return nil
    }
    if strings.Contains(*value, " ") {
        return &FieldError{Code: "noSpacesAllowed"}
    }
    return nil
}
```

Parameter type must be exact:
```go
type Form struct {
    Age *int `layrz:"number,required"`
}

// CORRECT: parameter is *int
func (f *Form) CleanAge(value *int) *FieldError { ... }

// WRONG: parameter is int (does not match *int)
// → _config error: "parameter type mismatch"
func (f *Form) CleanAge(value int) *FieldError { ... }
```

### Convention B: cross-field

**Signature:** `func (f *MyForm) Clean<Suffix>() Errors`

- Method name is `Clean` + any suffix not matching a struct field name (e.g., `CleanPasswords`, `CleanConsistency`).
- No parameters (receiver only).
- Returns `Errors` (the full map), or nil.
- Can report under arbitrary keys.
- Used for cross-field logic.

**Example:**
```go
type PasswordResetForm struct {
    NewPassword     *string `layrz:"char,required,min_length=8"`
    ConfirmPassword *string `layrz:"char,required,min_length=8"`
}

// Convention B: check passwords match
func (f *PasswordResetForm) CleanPasswords() layrz.Errors {
    if f.NewPassword == nil || f.ConfirmPassword == nil {
        return nil  // field validation already caught missing
    }
    if *f.NewPassword != *f.ConfirmPassword {
        return layrz.Errors{
            "passwordMismatch": {{Code: "mismatch"}},
        }
    }
    return nil
}
```

### Execution order and visibility

1. **Phase 1:** Built-in field rules (tag validators) run first across the entire form.
2. **Phase 2:** Convention A hooks (per-field) run alphabetically by method name. They see the field value; custom logic can reject it.
3. **Phase 3:** Convention B hooks (cross-field) run alphabetically by method name. They see the full form state and can report under any key.

Custom hooks only execute if the pointer receiver exists and the method signature is correct. A signature mismatch or panic in a hook produces a `_config` error instead of crashing.

## Nested structures: subform and subform_list

### Subform (pointer to struct)

Use `subform` kind on a pointer-to-struct field. Recursively validates the nested struct and prefixes error keys with the subform's camelCase field name.

```go
type Address struct {
    Street *string `layrz:"char,required,min_length=5"`
    City   *string `layrz:"char,required"`
}

type CreateUserForm struct {
    Name    *string  `layrz:"char,required"`
    Address *Address `layrz:"subform"`
}
```

Validating:
```go
form := &CreateUserForm{
    Name: Ptr("Alice"),
    Address: &Address{
        Street: Ptr("123"),  // too short (min 5)
        City:   nil,         // required
    },
}
errs := Validate(&form)
// errs["address.street"] = [minLength error]
// errs["address.city"] = [required error]
```

**Nil subforms are skipped** — no errors for a nil `Address` field. This is a deliberate divergence from Python, which would report it as missing. In Go, a nil subform is assumed to be intentionally absent.

### Subform list (slice of struct or pointer-to-struct)

Use `subform_list` kind on a slice field. Each element is validated at index keys, and nil elements are skipped without shifting sibling indices.

```go
type Item struct {
    Name  *string  `layrz:"char,required"`
    Price *float64 `layrz:"number,required,min_value=0"`
}

type OrderForm struct {
    Items []*Item `layrz:"subform_list"`
}
```

Validating:
```go
form := &OrderForm{
    Items: []*Item{
        {Name: nil, Price: Ptr(10.0)},           // index 0: missing name
        {Name: Ptr("Widget"), Price: Ptr(-5.0)}, // index 1: negative price
        nil,                                      // index 2: skipped
    },
}
errs := Validate(&form)
// errs["items.0.name"] = [required error]
// errs["items.1.price"] = [minValue error]
// (no error for items.2; nil elements are skipped)
```

**Elements can be values or pointers:**
```go
type OrderForm struct {
    Items []Item   `layrz:"subform_list"`  // slice of values
    // OR
    Items []*Item  `layrz:"subform_list"`  // slice of pointers
}
```

**Recursion depth guard:** The engine tracks recursion depth (maximum 32 levels). Exceeding this limit adds an `internalError` to the `_config` key.

## Error structure and the _config key

**FieldError fields:**
- `Code` (string) — always set. Examples: `required`, `invalid`, `empty`, `minLength`.
- `Expected` (any) — when meaningful (e.g., min bound, allowed choices). Omitted from JSON if unset.
- `Received` (any) — the actual value. Omitted from JSON if unset.
- `Extra` (map[string]any) — arbitrary context. Omitted from JSON if unset.

**Errors type methods:**
- `Add(key, errs...)` — append to the key's error slice.
- `Merge(other Errors)` — merge another error map, converting keys to camelCase.
- `IsEmpty() bool` — true if no errors.
- `Keys() []string` — sorted field names.
- `json.Marshal(errs)` — produces `{key: [{code, expected?, received?, extra?}], ...}`.

**Reserved key: `_config`**

Configuration and internal errors (bad tag, type mismatch, hook signature error, recovered panic, recursion limit) are filed under `_config`. This key should never appear in production; its presence indicates a programming bug, not a validation failure:

```go
form := &struct{
    X *string `layrz:"invalid_kind"`  // bad tag
}{}
errs := Validate(&form)
if errs["_config"] != nil {
    // Programming error; fix the tag
}
```

## Casing: snake_case to camelCase gotcha

Field names are converted to camelCase via `ToCamelCase()`, which lowercases only the first character. Go acronyms mangle predictably:

```
ID → iD
URL → uRL
HTTPCode → hTTPCode
```

This behavior is intentional and matches Python's algorithm exactly, so cross-language forms produce identical keys. If a specific key is needed, name the Go field accordingly:

```go
type Form struct {
    IDValue *string  // → idValue (camelCase as intended)
    ID      *string  // → iD (acronym mangles, but this is correct)
}
```

## Validation error codes and rules

Each kind produces specific codes under specific conditions:

| Kind | Rule | Code | Expected | Received | Notes |
|------|------|------|----------|----------|-------|
| `id` | required | `required` | — | — | nil (pointer) |
| `id` | type/value | `invalid` | — | — | non-int/string, ≤ 0, wrong type |
| `email` | required | `required` | — | — | nil |
| `email` | empty | `empty` | — | — | "" and Empty==false |
| `email` | regex | `invalid` | — | — | doesn't match pattern |
| `uuid` | required | `required` | — | — | nil |
| `uuid` | format | `invalid` | — | — | not a valid UUID format |
| `char` | required | `required` | — | — | nil |
| `char` | empty | `empty` | — | — | "" and Empty==false |
| `char` | min_length | `minLength` | min (int) | length (int) | rune count < min |
| `char` | max_length | `maxLength` | max (int) | length (int) | rune count > max |
| `char` | choices | `invalidChoice` | []string | value (string) | not in list |
| `char` | regex | `invalidFormat` | pattern (string) | value (string) | doesn't match |
| `number` | required | `required` | — | — | nil |
| `number` | type/datatype | `invalid` | — | — | wrong numeric type, wrong datatype |
| `number` | min_value | `minValue` | bound (float or int64) | value (float or int64) | < min |
| `number` | max_value | `maxValue` | bound (float or int64) | value (float or int64) | > max |
| `bool` | required | `required` | — | — | nil |
| `json` | required | `required` | — | — | nil |
| `json` | type | `invalid` | — | — | wrong container type (list vs dict) or empty when Empty==false |

## Common mistakes

1. **Colon instead of equals:** `layrz:"char:required"` → fails with parse error. Use `char,required`.
2. **Passing a value instead of pointer:** `Validate(form)` → _config error. Use `Validate(&form)`.
3. **Value receiver on Clean method:** `func (f Form) CleanName(...) ...` → method not found. Use pointer receiver `func (f *Form) CleanName(...) ...`.
4. **Convention A parameter type mismatch:** Field is `*int`, but hook takes `int` → _config error. Match the type exactly.
5. **Untagged field:** `Name string` (no tag) is silently skipped. Add a tag to validate it.
6. **Expecting `required` on value field:** `Count int` with `required` → no error on zero. Use `empty:false` to reject zero-length containers.
7. **Regex not last:** `layrz:"char,regex=[0-9]+,required"` → regex eats the remaining tokens. Put regex at the end.
8. **Nil subform with subform kind:** `Address *Address` with `subform` tag — a nil subform is skipped, not errored. Either make it required (pointer, custom Convention A hook) or accept the absence.

## Verifying behavior

The Go implementation is pinned by cross-language test vectors in `vectors/fields/*.json` (one per kind, ~96 total cases). When unsure whether an input produces a specific error, check the vectors:

```bash
cd /home/mochi/Projects/layrz-forms
grep -E '"name"|"value_absent"|"expected_errors"' vectors/fields/NumberField.json | head -20
go test ./... -v  # run the full suite
make test         # enforces coverage threshold (90%)
```

After writing a form, run `go vet ./...` and `go test ./...` to confirm no compilation or logic errors.

### Example: comprehensive form with all features

```go
package main

import (
    "github.com/goldenm-software/layrz-forms/go/v3"
)

type Address struct {
    Street *string `layrz:"char,required,min_length=5"`
    City   *string `layrz:"char,required"`
    Zip    *string `layrz:"char,required,regex=^[0-9]{5}$"`
}

type Item struct {
    Name     *string  `layrz:"char,required"`
    Quantity *int     `layrz:"number,required,min_value=1"`
    Price    *float64 `layrz:"number,required,min_value=0"`
}

type CreateOrderForm struct {
    OrderID     *int       `layrz:"id,required"`
    Email       *string    `layrz:"email,required"`
    Description *string    `layrz:"char,empty"`
    Address     *Address   `layrz:"subform"`
    Items       []*Item    `layrz:"subform_list"`
}

// Convention A: email must not be a known spam domain
func (f *CreateOrderForm) CleanEmail(value *string) *layrz.FieldError {
    if value != nil && len(*value) > 0 {
        if strings.HasSuffix(*value, "@spam.com") {
            return &layrz.FieldError{Code: "spamDomain"}
        }
    }
    return nil
}

// Convention B: at least one item must be in the order
func (f *CreateOrderForm) CleanItems() layrz.Errors {
    if len(f.Items) == 0 {
        return layrz.Errors{
            "items": {{Code: "atLeastOne"}},
        }
    }
    return nil
}

func main() {
    // Build form
    form := &CreateOrderForm{
        OrderID: layrz.Ptr(123),
        Email:   layrz.Ptr("user@example.com"),
        Address: &Address{
            Street: layrz.Ptr("Main"),  // too short
            City:   layrz.Ptr("NYC"),
            Zip:    layrz.Ptr("10001"),
        },
        Items: []*Item{
            {Name: layrz.Ptr("Widget"), Quantity: layrz.Ptr(2), Price: layrz.Ptr(10.0)},
        },
    }

    // Validate
    errs := layrz.Validate(form)
    if !layrz.IsValid(form) {
        for _, key := range errs.Keys() {
            for _, err := range errs[key] {
                println(key, ":", err.Code)
            }
        }
    }
}
```

Output on invalid form:
```
address.street : minLength
```
