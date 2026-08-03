# Layrz Forms (Python)

A form validation library for Python with a simple, flexible API. Validates data from dicts, plain objects, and Strawberry GraphQL inputs. Simpler than Django Forms, yet powerful enough for complex nested forms and custom validation.

## Install

```bash
pip install layrz-forms
```

Or with `uv`:

```bash
uv add layrz-forms
```

**Requires Python 3.14+.**

## Quick Start

```python
import layrz_forms as forms


class ExampleForm(forms.Form):
  """Example form with basic and nested validation."""
  id_test = forms.IdField(required=True)
  email_text = forms.EmailField(required=True)
  json_dict_test = forms.JsonField(required=True, datatype=dict)
  json_list_test = forms.JsonField(required=True, datatype=list)
  int_test = forms.NumberField(required=True, datatype=int, min_value=0, max_value=5)
  float_test = forms.NumberField(required=True, datatype=float, min_value=0, max_value=5)
  bool_test = forms.BooleanField(required=True)
  plain_text = forms.CharField(required=True, empty=False)
  empty_text = forms.CharField(required=True, empty=True)
  range_text = forms.CharField(required=True, empty=False, min_length=5, max_length=10)

  def clean_func1(self) -> None:
    """Cross-field validation (alphabetically first)."""
    self.add_errors(key='clean1', code='error1')
    self.add_errors(key='clean1', code='error2')

  def clean_func2(self) -> None:
    """Cross-field validation (alphabetically second)."""
    self.add_errors(key='clean2', code='error1')


if __name__ == '__main__':
  obj = {
    'id_test': 1,
    'email_text': 'example@goldenmcorp.com',
    'json_dict_test': {'hola': 'mundo'},
    'json_list_test': ['hola mundo'],
    'int_test': 5,
    'float_test': 4.5,
    'bool_test': True,
    'plain_text': 'hola mundo',
    'empty_text': 'hola',
    'range_text': 'hola',  # 4 chars, min_length=5 — error!
  }

  form = ExampleForm(obj)

  print('form.is_valid():', form.is_valid())
  # Output: form.is_valid(): False

  print('form.errors:')
  for key, errors in form.errors.items():
    print(f'  {key}:')
    for e in errors:
      print(f'    code={e.code!r}, expected={e.expected}, received={e.received}')
  # Output:
  #   rangeTextTest:
  #     code='minLength', expected=5, received=4
  #   clean1:
  #     code='error1', expected=None, received=None
  #     code='error2', expected=None, received=None
  #   clean2:
  #     code='error1', expected=None, received=None

  # Serialize errors to JSON
  errors_as_dicts = {k: [e.model_dump() for e in v] for k, v in form.errors.items()}
  print('errors_as_dicts:', errors_as_dicts)
  # Output: errors_as_dicts: {'rangeTextTest': [{'code': 'minLength', 'expected': 5, 'received': 4}], 'clean1': [{'code': 'error1'}, {'code': 'error2'}], 'clean2': [{'code': 'error1'}]}
```

## Field Reference

### BooleanField

Validates boolean values.

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `required` | `bool` | `False` | Whether the field must be present (non-None) |

**Error Codes:**
- `required` — Value is `None` and `required=True`
- `invalid` — Value is present but not a `bool`

### CharField

Validates strings with optional length, choice, and regex constraints.

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `required` | `bool` | `False` | Whether the field must be present (non-None) |
| `empty` | `bool` | `False` | Whether an empty string (`''`) is allowed |
| `min_length` | `int \| None` | `None` | Minimum number of characters |
| `max_length` | `int \| None` | `None` | Maximum number of characters |
| `choices` | `tuple[tuple[str, str], ...] \| None` | `None` | Allowed values as `(('value', 'Label'), ...)` |
| `regex` | `str \| None` | `None` | PCRE regex pattern to match (checked on non-empty strings) |

**Error Codes:**
- `required` — Value is `None` and `required=True`
- `invalid` — Value is present but not a `str` (or `Enum`/`StrEnum`)
- `empty` — Value is `''` and `empty=False`
- `minLength` — String length < `min_length`; `expected` and `received` are length values
- `maxLength` — String length > `max_length`; `expected` and `received` are length values
- `invalidChoice` — Value not in `choices`; `expected` is the list of allowed values, `received` is the value
- `invalidFormat` — Value does not match `regex`; `expected` is the regex pattern, `received` is the value

### EmailField

Validates email addresses.

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `required` | `bool` | `False` | Whether the field must be present (non-None) |
| `empty` | `bool` | `False` | Whether an empty string (`''`) is allowed |
| `regex` | `str` | `r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-z]{2,63}$'` | Email regex pattern |

**Error Codes:**
- `required` — Value is `None` and `required=True`
- `invalid` — Value is present but not a `str`, or non-empty string fails regex match
- `empty` — Value is `''` and `empty=False`

### IdField

Validates positive integer IDs (can be `int` or numeric `str`).

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `required` | `bool` | `False` | Whether the field must be present (non-None) |

**Error Codes:**
- `required` — Value is `None` and `required=True`
- `invalid` — Value is not an `int` or numeric `str`, is `bool`, or is ≤ 0

### JsonField

Validates JSON objects (dicts) or arrays (lists).

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `required` | `bool` | `False` | Whether the field must be present (non-None) |
| `empty` | `bool` | `False` | Whether an empty container is allowed |
| `datatype` | `type[list] \| type[dict]` | `dict` | The container type to validate |

**Error Codes:**
- `required` — Value is `None` and `required=True`
- `invalid` — Value is not an instance of `datatype`, or is empty and `empty=False`

### NumberField

Validates numeric values (int or float, configurable).

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `required` | `bool` | `False` | Whether the field must be present (non-None) |
| `datatype` | `type[int] \| type[float]` | `float` | The numeric type to validate |
| `min_value` | `float \| None` | `None` | Minimum allowed value |
| `max_value` | `float \| None` | `None` | Maximum allowed value |

**Error Codes:**
- `required` — Value is `None` and `required=True`
- `invalid` — Value is not an instance of `datatype`, or is `bool`
- `minValue` — Value < `min_value`; `expected` and `received` are converted to `datatype`
- `maxValue` — Value > `max_value`; `expected` and `received` are converted to `datatype`

### UuidField

Validates UUID strings or `uuid.UUID` instances in any standard format (hyphenated, unhyphenated, braced, URN).

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `required` | `bool` | `False` | Whether the field must be present (non-None) |

**Error Codes:**
- `required` — Value is `None` and `required=True`
- `invalid` — Value is not a `str` or `uuid.UUID`, or string is not a valid UUID

## Errors

Form errors are accessed via the `form.errors` **property** (not a method), returning `dict[str, list[LayrzError]]`.

```python
form = ExampleForm({'email_text': None})
# Reading .errors triggers lazy validation if not yet done:
print(form.errors)
# Output: {'emailText': [LayrzError(code='required')]}
```

### LayrzError Model

Each error in the list is a Pydantic `LayrzError` with four fields:

| Field | Type | Description |
|-------|------|-------------|
| `code` | `str` | Error code (always set) |
| `expected` | `Any` | The constraint that was violated (e.g., `5` for min_length); `None` if not applicable |
| `received` | `Any` | The offending value; `None` if not applicable |
| `extra` | `dict \| None` | Additional contextual keys from custom validators (e.g., `{'message': '...'}`); `None` if empty |

### Serializing to JSON

Use `.model_dump()` to convert errors to plain dicts. It excludes `None` fields by default:

```python
errors_as_dicts = {k: [e.model_dump() for e in v] for k, v in form.errors.items()}
# {'emailText': [{'code': 'required'}], 'rangeText': [{'code': 'minLength', 'expected': 5, 'received': 4}]}
```

### Custom Errors via add_errors()

Inside a `clean_*` method, use `self.add_errors(key, code, extra_args={...})` to add custom errors. The `extra_args` dict is processed as follows:

- Keys `expected` and `received` are lifted into their own `LayrzError` fields
- All other keys nest under the `extra` field

```python
def clean_password(self) -> None:
  # Emit a custom error with expected/received:
  self.add_errors(
    key='password',
    code='weak_password',
    extra_args={
      'expected': 'at least 12 chars',
      'received': len(self._obj.get('password', '')),
      'min_entropy': 25,  # Nests under extra
    }
  )
  # Result: LayrzError(code='weak_password', expected='at least 12 chars', received=6, extra={'min_entropy': 25})
```

## snake_case to camelCase Conversion

All field names in error keys are automatically converted from snake_case to camelCase:

| snake_case | camelCase |
|-----------|-----------|
| `id_test` | `idTest` |
| `email_text` | `emailText` |
| `range_text_test` | `rangeTextTest` |

For nested forms, each segment is converted independently:

| snake_case (nested) | camelCase (nested) |
|-----|-----|
| `address.street_name` | `address.streetName` |
| `items.0.name_display` | `items.0.nameDisplay` |

## Clean Methods

Any method starting with `clean_` is auto-discovered and invoked after field validation completes. Clean methods run in **alphabetical order by method name** (e.g., `clean_apple`, then `clean_banana`, then `clean_zebra`). This ordering is stable and intentional—it matches the Go implementation, where declaration order is not available via reflection.

### Sync Clean Methods

```python
def clean_password_match(self) -> None:
  """Sync clean method."""
  password = self._obj.get('password')
  password_confirm = self._obj.get('password_confirm')
  if password != password_confirm:
    self.add_errors('password_confirm', 'mismatch')
```

### Async Clean Methods

Async clean methods are supported but **only when using `await form.ais_valid()`**:

```python
async def clean_username_available(self) -> None:
  """Async clean method."""
  username = self._obj.get('username')
  if username:
    # Check database, API, etc.
    is_taken = await check_username_db(username)
    if is_taken:
      self.add_errors('username', 'already_taken')
```

**Important:** If your form has `async def clean_*` methods, reading `.errors` without awaiting `ais_valid()` raises `RuntimeError`:

```python
form = MyAsyncForm(obj)
# ❌ This raises RuntimeError:
print(form.errors)

# ✅ This is correct:
await form.ais_valid()
print(form.errors)
```

## Nested Forms

Assign a `Form` instance as a class attribute to nest it. Nested form fields are validated recursively, and error keys are prefixed with the form's name using dot notation:

```python
class Address(forms.Form):
  street_name = forms.CharField(required=True, min_length=5)
  zip_code = forms.CharField(required=True)

class User(forms.Form):
  name = forms.CharField(required=True)
  address = Address()  # Nested form

form = User({'name': 'Alice', 'address': {'street_name': 'oak', 'zip_code': '12345'}})
form.is_valid()
# Errors keyed as: 'address.streetName' (not 'address'), 'address.zipCode'
print(form.errors)
# Output: {'address.streetName': [LayrzError(code='minLength', expected=5, received=3)]}
```

## Lists of Fields or Forms

Declare a list of fields or forms using the special `_attrs` naming convention. The first element defines the type:

```python
class ShoppingCart(forms.Form):
  items = [ItemForm()]  # List of nested forms
  tags = [forms.CharField()]  # List of char fields

form = ShoppingCart({
  'items': [
    {'name': 'Widget', 'price': 9.99},
    {'name': 'Gadget', 'price': -5},  # Error: negative price
  ],
  'tags': ['electronics', 'sale'],
})
form.is_valid()
# Errors keyed as: 'items.0.name', 'items.1.price', 'tags.2.myError', etc.
print(form.errors)
```

### Current Behavior & Quirks

**Only lists whose first element is a `Field` or `Form` are treated as nested declarations.** Plain lists (e.g., `colors = ['red', 'green']`) are silently ignored during validation.

If a non-list value is supplied for a declared nested list field, a validation error is emitted:

```python
form = ShoppingCart({'items': 'not-a-list'})
form.is_valid()
print(form.errors)
# Output: {'items': [LayrzError(code='invalid', extra={'message': 'Invalid data type'})]}
```

Empty list class attributes (e.g., `items = []`) are silently skipped and produce no errors.

## Strawberry GraphQL

Pass a Strawberry input object directly to `Form(obj=...)`. It is automatically converted to a dict via `Form.strawberry_to_dict()`:

```python
import strawberry

@strawberry.input
class UserInput:
  name: str
  email: str

class UserForm(forms.Form):
  name = forms.CharField(required=True)
  email = forms.EmailField(required=True)

# From a Strawberry resolver:
strawberry_obj = UserInput(name='Alice', email='alice@example.com')
form = UserForm(obj=strawberry_obj)
form.is_valid()
```

## cleaned_data

The `form.cleaned_data` property returns a **deep copy** of the validated object. Mutations to the returned dict (at any nesting depth) cannot affect the caller's original object:

```python
form = MyForm({'name': 'Alice', 'nested': {'value': 1}})
form.is_valid()

cleaned = form.cleaned_data
cleaned['nested']['value'] = 999  # Does not affect form._obj
```

Note: Deep copying can fail if the payload contains non-serializable objects (e.g., file handles, custom classes without `__deepcopy__`); such failures propagate naturally and are not caught.

## Differences from the Go Implementation

The Python and Go implementations share the same error codes, validation rules, and cross-language test vectors (96 test cases in `../vectors/fields/`). However, they diverge in several aspects:

| Aspect | Python | Go |
|--------|--------|-----|
| **Field Declaration** | Class attributes (e.g., `name = CharField()`) | Struct tags (e.g., `` `layrz:"char,required"` ``) |
| **Async Support** | `async def clean_*` + `await form.ais_valid()` | No async equivalent (all sync) |
| **Nil/None Handling** | Python uses `None` for absence; value fields always present | Go uses pointers for absence; pointer to pointer is an error |
| **Value Field Semantics** | N/A; Python doesn't distinguish pointer vs. value | Go value fields always present, so `required` never fires; pointer fields model absence as nil |
| **Boolean in Number Fields** | Rejected explicitly | Rejected explicitly |
| **Clean Method Ordering** | Alphabetical (intentional, stable) | Alphabetical (intentional, stable) |
| **Nested Form Absence** | Validates against empty dict | Skipped entirely (nil pointers produce no errors) |

Both implementations are correct and pass their respective test suites. Bugs in one will be fixed in future major releases (Python 4.0.0, Go 1.0.0 when it exists) to bring them into full alignment.

## Development

All commands run from the `python/` directory.

### Install dev dependencies

```bash
uv sync --only-group dev
```

### Lint

```bash
uv run ruff check
```

### Type Check

Uses **`ty`** (not mypy), so type suppressions are `# ty: ignore[rule-code]`:

```bash
uv run ty check
```

### Run Tests

```bash
uv run pytest -q
```

### Coverage

The CI threshold is 90%; current coverage is 97%:

```bash
uv run pytest --cov=layrz_forms --cov-report=term-missing
```

### Build Distribution

```bash
uv run python -m build
```

## Deployment

Releases are **tag-triggered** via GitHub Actions. To release a new version:

1. Update `version` in `python/pyproject.toml` (e.g., `3.1.0`)
2. Update `CHANGELOG.md`
3. Push a tag matching `v[0-9]+.[0-9]+.[0-9]+` (e.g., `git tag v3.1.0`) to the `main` branch
4. The workflow in `.github/workflows/deploy.yaml` builds and publishes to PyPI automatically

## License

MIT License. See the repository for details.

---

Maintained by [Golden M](https://goldenm.com) with authorization of [Layrz LTD](https://layrz.com).
