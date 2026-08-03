---
name: py-form-builder
description: Build, edit, or debug Python forms using layrz-forms Form subclasses. Use this for any Python validation work — when defining or modifying Form subclasses, adding BooleanField/CharField/EmailField/IdField/JsonField/NumberField/UuidField fields, implementing clean_* methods, calling is_valid() or ais_valid(), reading .errors or .cleaned_data, integrating with Strawberry GraphQL, or diagnosing error codes. Do NOT use for Go forms — this is Python-specific.
---

# Python form builder

## Minimal form

A form is a class that inherits `Form` and declares `Field` instances as attributes. Each field
defines validation rules. Call `is_valid()` to run validation, then inspect `.errors` (a
property that auto-validates if not yet run):

```python
from layrz_forms import Form, CharField, NumberField

class UserForm(Form):
  name = CharField(required=True)
  age = NumberField(datatype=int, required=False)

form = UserForm({'name': 'Alice', 'age': 30})
if form.is_valid():
  print(form.cleaned_data)  # {'name': 'Alice', 'age': 30}
else:
  print(form.errors)  # {} (no errors in this case)
```

Missing required field:

```python
form = UserForm({'name': 'Bob'})
form.is_valid()
print(form.errors)
# {'age': [LayrzError(code='required')]} (lazy property: .errors auto-validates)
# Note: age is not in the dict, so it's required but absent (error code 'required')
```

Present but invalid:

```python
form = UserForm({'name': 'Charlie', 'age': 'not a number'})
form.is_valid()
print(form.errors)
# {'age': [LayrzError(code='invalid')]}
```

To serialize errors to JSON, use:

```python
{k: [e.model_dump() for e in v] for k, v in form.errors.items()}
# {'age': [{'code': 'invalid'}]}
```

Field names are auto-converted to camelCase: `user_age` → `userAge` in error keys.

## Field reference

| Field | Constructor Kwargs | Default | Error Codes |
| --- | --- | --- | --- |
| `BooleanField` | `required` | `False` | `required`, `invalid` |
| `CharField` | `required`, `max_length`, `min_length`, `empty`, `regex`, `choices` | `required=False`, all others `None`/`False` | `required`, `invalid`, `empty`, `maxLength`, `minLength`, `invalidChoice`, `invalidFormat` |
| `EmailField` | `required`, `empty`, `regex` | `required=False`, `empty=False`, `regex=r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-z]{2,63}$'` | `required`, `invalid`, `empty` |
| `IdField` | `required` | `False` | `required`, `invalid` |
| `JsonField` | `required`, `empty`, `datatype` | `required=False`, `empty=False`, `datatype=dict` | `required`, `invalid` |
| `NumberField` | `required`, `datatype`, `min_value`, `max_value` | `required=False`, `datatype=float`, min/max=`None` | `required`, `invalid`, `minValue`, `maxValue` |
| `UuidField` | `required` | `False` | `required`, `invalid` |

### BooleanField

Accepts only `bool` values (or `None` if not required).

```python
field = BooleanField(required=True)
# {'active': [LayrzError(code='invalid')]}  — if value is 1 or 'true' (not a bool)
# {'active': [LayrzError(code='required')]}  — if value is None and required=True
```

### CharField

Accepts `str`, `Enum`, or `StrEnum` (auto-converted to string). Returns early on non-string type,
then accumulates all length/choice/regex errors in one call.

String lengths are counted by code point (Python `len()`), not bytes, so non-ASCII characters
count as single characters. If `choices` is defined, value must be in the first element of each
tuple:

```python
field = CharField(
  required=True,
  min_length=2,
  max_length=10,
  choices=(('admin', 'Administrator'), ('user', 'User')),
  regex=r'^[a-z]+$',
)

# value='admin' → valid
# value='a' → [minLength error (expected=2, received=1)]
# value='toolong' → [maxLength error (expected=10, received=7)]
# value='Admin' → [invalidChoice (expected=['admin','user']), invalidFormat]
# value=123 → [invalid] (early return; no length/choice/regex checks)
```

With `empty=True`, empty string `''` is accepted *without* regex validation. Non-empty strings
are always regex-checked if `regex` is set.

```python
field = CharField(empty=True, regex=r'^\d+$')
# value='' → valid (empty string allowed, no regex)
# value='123' → valid (matches regex)
# value='abc' → [invalidFormat] (non-empty, must match regex)
```

Error codes:
- `invalid`: not a string
- `empty`: empty string when `empty=False`
- `minLength`: `expected` is the minimum, `received` is the actual length
- `maxLength`: `expected` is the maximum, `received` is the actual length
- `invalidChoice`: `expected` is list of allowed choices, `received` is the value
- `invalidFormat`: `expected` is the regex pattern, `received` is the value

### EmailField

Validates against a regex (default: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-z]{2,63}$`).
With `empty=True`, empty string is allowed without regex check. Non-empty strings always
validate against the regex.

```python
field = EmailField(required=True)
# value='user@example.com' → valid
# value='invalid' → [invalid]
# value=None → [required]
# value='' → [empty]

field = EmailField(empty=True)
# value='' → valid (no regex check for empty string)
# value='test@domain.co' → valid
# value='bad@domain' → [invalid] (non-empty, must match regex)
```

Error codes:
- `required`: value is `None` and `required=True`
- `invalid`: not a string or doesn't match regex (both non-empty and custom regex)
- `empty`: empty string when `empty=False`

### IdField

Accepts `int` or numeric `str` (e.g., `"123"`). Rejects `bool` explicitly (since
`isinstance(True, int)` is true in Python). Must be > 0.

```python
field = IdField(required=True)
# value=123 → valid
# value='456' → valid (string converted to int)
# value=0 → [invalid] (not > 0)
# value=-5 → [invalid] (not > 0)
# value=True → [invalid] (bool rejected)
# value='abc' → [invalid] (cannot convert to int)
# value=None → [required]
```

Error code:
- `required`: value is `None` and `required=True`
- `invalid`: non-int/str, bool, ≤ 0, or unparseable string

### JsonField

Accepts `dict` or `list` (type specified by `datatype` parameter). With `empty=False`, dict
must have at least one key and list must have at least one element. An absent optional field
(not in the dict at all) produces NO errors.

```python
field = JsonField(datatype=dict, required=False, empty=False)
# value={} → [invalid] (empty dict when empty=False)
# value={'key': 'val'} → valid
# value=[] → [invalid] (not a dict)
# value=None (and field absent from input) → valid (optional, no error)

field = JsonField(datatype=list, empty=True)
# value=[] → valid (empty allowed)
# value=[1, 2, 3] → valid
# value={} → [invalid] (not a list)
```

Error codes:
- `required`: value is `None` and `required=True`
- `invalid`: wrong datatype, or empty dict/list when `empty=False`

### NumberField

Validates using `isinstance(value, self.datatype)`, so `int` does NOT pass `datatype=float`.
Rejects `bool` explicitly. Checks `min_value` and `max_value` constraints, emitting both if
violated.

```python
field = NumberField(datatype=float, min_value=0.0, max_value=100.0)
# value=50.5 → valid
# value=50 → [invalid] (int, not float)
# value=True → [invalid] (bool rejected)
# value=-10.0 → [minValue (expected=0.0, received=-10.0)]
# value=150.0 → [maxValue (expected=100.0, received=150.0)]
# value='50' → [invalid]

field = NumberField(datatype=int)
# value=42 → valid
# value=3.14 → [invalid] (float, not int)
```

Error codes:
- `required`: value is `None` and `required=True`
- `invalid`: wrong datatype, bool, or constraint parse error
- `minValue`: `expected` is min_value, `received` is the value
- `maxValue`: `expected` is max_value, `received` is the value

### UuidField

Accepts `str` or `uuid.UUID` instance. Strings are validated against UUID format using
Python's `uuid.UUID()` constructor.

```python
field = UuidField(required=True)
# value='550e8400-e29b-41d4-a716-446655440000' → valid
# value=uuid.UUID('550e8400-e29b-41d4-a716-446655440000') → valid
# value='not-a-uuid' → [invalid]
# value=None → [required]
```

Error codes:
- `required`: value is `None` and `required=True`
- `invalid`: not a string/UUID or invalid UUID format

## Clean methods

After all field validation, any method starting with `clean` is auto-discovered and called in
**alphabetical order** (not declaration order). This allows cross-field validation.

```python
class RegistrationForm(Form):
  password = CharField(required=True)
  password_confirm = CharField(required=True)

  def clean_password_match(self) -> None:
    if self._obj.get('password') != self._obj.get('password_confirm'):
      self.add_errors(
        key='passwordConfirm',
        code='mismatch',
        extra_args={'message': 'Passwords do not match'},
      )
```

When errors are added in `add_errors()`, `expected` and `received` are lifted into their own
fields; other keys nest under `extra`. The resulting error:

```python
LayrzError(
  code='mismatch',
  expected=None,
  received=None,
  extra={'message': 'Passwords do not match'},
)
```

Use `expected` and `received` when they matter:

```python
def clean_age_range(self) -> None:
  min_age = 18
  actual = self._obj.get('age')
  if actual and actual < min_age:
    self.add_errors(
      key='age',
      code='tooYoung',
      extra_args={
        'expected': min_age,
        'received': actual,
        'note': 'Must be 18+',
      },
    )

# Output:
# LayrzError(
#   code='tooYoung',
#   expected=18,
#   received=16,
#   extra={'note': 'Must be 18+'},
# )
```

Clean methods always see the result of field validation — if a field failed, it is still in
`self._obj` (the original input), not yet in `cleaned_data`.

### Async clean methods

Use `async def clean_*` for async validation (e.g., database checks). Must use `await
form.ais_valid()` to wait for them.

```python
async def clean_email_unique(self) -> None:
  email = self._obj.get('email')
  if email and await db.email_exists(email):
    self.add_errors(key='email', code='duplicate')
```

Reading `.errors` without awaiting `ais_valid()` when async clean methods exist raises
`RuntimeError`:

```python
form = MyForm(data)
# form.errors  # RuntimeError: Cannot call async clean function in sync context
await form.ais_valid()
form.errors  # Now safe to read
```

If mixing sync and async clean methods, all run within `ais_valid()` (sync methods are wrapped
in a no-op `await asyncio.sleep(0)`).

## Nested forms and lists

A `Form` instance as a class attribute is validated recursively; errors are keyed by dot-notation:

```python
class AddressForm(Form):
  street = CharField(required=True)
  city = CharField(required=True)

class UserForm(Form):
  name = CharField(required=True)
  address = AddressForm()

form = UserForm({'name': 'Alice', 'address': {'street': '', 'city': 'NYC'}})
form.is_valid()
print(form.errors)
# {'address.street': [LayrzError(code='empty')]}
```

A list attribute whose first element is a `Field` or `Form` declares a nested list:

```python
class PhoneListForm(Form):
  phones = [EmailField()]  # List of emails

form = PhoneListForm({'phones': ['alice@ex.com', 'invalid', 'bob@ex.com']})
form.is_valid()
print(form.errors)
# {'phones.1': [LayrzError(code='invalid')]}
```

Similarly for nested forms:

```python
class TeamForm(Form):
  members = [AddressForm()]

form = TeamForm({
  'members': [
    {'street': '123 Main', 'city': 'NYC'},
    {'street': '', 'city': 'LA'},
  ]
})
form.is_valid()
# {'members.1.street': [LayrzError(code='empty')]}
```

A non-list value supplied where a list belongs emits `{'code': 'invalid', 'extra': {'message':
'Invalid data type'}}`:

```python
form = PhoneListForm({'phones': 'single@email.com'})
form.is_valid()
print(form.errors)
# {'phones': [LayrzError(code='invalid', extra={'message': 'Invalid data type'})]}
```

An empty list attribute is skipped (no validation performed on it).

## Errors and serialisation

`.errors` is a lazy property: the first access triggers sync validation if it has not run yet.
Each error is a `LayrzError` Pydantic model; calling `.model_dump()` excludes `None` by default:

```python
error = LayrzError(code='minLength', expected=5, received=3)
error.model_dump()  # {'code': 'minLength', 'expected': 5, 'received': 3}

error = LayrzError(code='required')
error.model_dump()  # {'code': 'required'}
```

To serialize the entire error dict to a plain dict (no Pydantic models):

```python
serialized = {k: [e.model_dump() for e in v] for k, v in form.errors.items()}
```

To serialize to JSON:

```python
import json
serialized = {k: [e.model_dump() for e in v] for k, v in form.errors.items()}
json_str = json.dumps(serialized)
```

## cleaned_data

Returns a **deep copy** of the original input dict if validation passed. Mutations to the
returned object at any depth cannot affect the caller's original dict. Deep copying can fail on
non-copyable values; such errors propagate naturally.

```python
form = UserForm({'name': 'Alice', 'address': {'city': 'NYC'}})
form.is_valid()
clean = form.cleaned_data
clean['address']['city'] = 'LA'
# Original input is unaffected; clean is a separate deep copy
```

If validation failed, `cleaned_data` still returns the deep copy (garbage-in, garbage-out —
validation failure does not prevent access to the copy).

## Strawberry GraphQL

Pass a Strawberry input object directly to `Form(obj=...)`. It is converted internally via
`strawberry_to_dict()`:

```python
import strawberry
from layrz_forms import Form, CharField

@strawberry.input
class UserInput:
  name: str
  email: str

class UserForm(Form):
  name = CharField(required=True)
  email = CharField(required=True)

user_input = UserInput(name='Alice', email='alice@ex.com')
form = UserForm(user_input)
form.is_valid()  # Conversion is automatic
```

## Reserved names

These names cannot be used as field or nested-form attributes (they collide with Form methods):

`add_errors`, `change_obj`, `clean`, `errors`, `is_valid`, `ais_valid`, `set_obj`,
`calculate_members`, `cleaned_data`.

Declaring one raises an error during discovery.

## Common mistakes

**Calling `.errors()` as a method:** `.errors` is a property, not a method. Write `form.errors`,
not `form.errors()`.

**Confusing `required` and `empty`:** `required` means the field must be present in the input
dict. `empty` is only for `CharField` and `EmailField` and means an empty string is acceptable.
A missing optional field produces no errors; an empty string in a non-empty field produces
`empty`, never `required`.

**Expecting an int to satisfy `datatype=float`:** `NumberField` uses `isinstance()` strictly.
`NumberField(datatype=float)` rejects integers; use `NumberField(datatype=(int, float))` is not
supported — use `datatype=float` and handle both in clean methods, or split the form.

**Declaring a plain list and expecting it to be a nested form:**

```python
class BadForm(Form):
  colors = ['red', 'green', 'blue']  # Plain list, IGNORED

# This does not create a nested list field; 'colors' is simply ignored.
```

Do this instead:

```python
class GoodForm(Form):
  colors = [CharField()]  # First element is a Field, so it's a nested list
```

**Forgetting `await ais_valid()`:** Reading `.errors` on a form with async clean methods
without awaiting `ais_valid()` raises `RuntimeError`. Always `await` if any `async def
clean_*` methods are present.

**Expecting field validation to short-circuit:** All fields are validated, and errors
accumulate. A field that is too short AND fails a regex reports both errors.

## Verify

Run tests from `python/`:

```bash
uv run pytest -q                                    # Run all tests
uv run pytest -q -k CharField                       # Run CharField tests
uv run ruff check                                   # Lint
uv run ty check                                     # Type check
uv run pytest --cov=layrz_forms --cov-report=term-missing  # Coverage
```

The shared vectors at `vectors/fields/*.json` pin behaviour for all fields. When unsure whether
a value produces an error, grep the vectors:

```bash
grep -l 'required\|empty' vectors/fields/*.json
python3 -c "import json;[print(c['name'],c.get('value','ABSENT'),c['expected_errors']) for c in json.load(open('vectors/fields/CharField.json'))]"
```

Type checking uses **`ty`**, not mypy. Suppress with `# ty: ignore[rule-code]`.

## See also

See the form-builder router for cross-language contract details and Go form guidance.
