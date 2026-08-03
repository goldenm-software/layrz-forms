# Changelog

## 3.0.0

### Breaking changes

- **`Form.errors` is now a property, not a method**: Access as `form.errors` instead of `form.errors()`. The property lazily validates on first access — if `is_valid()` hasn't been called, reading `.errors` runs synchronous validation automatically. Reassigning `form.obj` or calling `calculate_members()` invalidates cached errors.

- **Error entries are now `LayrzError` Pydantic models, not plain dicts**: `form.errors` returns `dict[str, list[LayrzError]]`. Each `LayrzError` has four fields: `code` (str, always set), `expected` (Any), `received` (Any), and `extra` (dict|None for custom keys). `model_dump()` defaults to `exclude_none=True`; use `exclude_none=False` to include unset fields. The model is `extra='forbid'` — unknown constructor kwargs raise.

- **`add_errors()` extra_args routing**: The `key`, `code` signature remains unchanged, but `extra_args` dict keys `expected` and `received` are now lifted into their own fields; all other keys nest under `extra`. Example: `self.add_errors('password', 'weak', extra_args={'min_length': 8})` produces `LayrzError(code='weak', extra={'min_length': 8})`.

- **`is_valid_async()` renamed to `ais_valid()`**: Use the `a`-prefix async convention (e.g. `await form.ais_valid()` instead of `await form.is_valid_async()`). The old name is removed entirely; no alias exists.

- **Async clean functions + lazy property caveat**: Forms with `async def clean_*` methods raise `RuntimeError` if `.errors` is read without awaiting `ais_valid()` first. Always use: `await form.ais_valid()` then `form.errors`.

- **`CharField` now validates type**: Previously accepted non-string types (lists, dicts, numbers) and silently length-checked them. Now emits `{'code': 'invalid'}` for any present non-string value (except `Enum`/`StrEnum` members, which still convert to strings). Impacts forms accepting arbitrary input without type validation.

- **`IdField` validation robustness**: Now emits `{'code': 'invalid'}` instead of raising `ValueError` for non-numeric strings like `'abc'`, and instead of raising `TypeError` when comparing non-comparable types on optional fields.

- **Empty list class attributes**: Forms declaring empty list class attributes (e.g. `empty_attr = []`) previously raised `IndexError`; now silently skipped during validation.

### Fixed

- 49 uncaught exceptions across 7 field types eliminated (9 new validation errors added in their place, all emitted as `{'code': 'invalid'}`). Forms accepting untrusted input are now crash-proof.

### Added

- **`LayrzError` Pydantic model** (module: `layrz_forms/errors.py`): The new error representation with four fields (`code`, `expected`, `received`, `extra`). Exported from `layrz_forms.__init__`.

- **Error type aliases** (module: `layrz_forms/types.py`):
  - `ErrorType`: A single `LayrzError`.
  - `ErrorsType`: The full mapping, `dict[str, list[LayrzError]]`.
  Both exported from `layrz_forms.__init__`.

- **New internal modules** for cleaner architecture:
  - `casing.py`: centralized `to_camel_case(key)` function (replaces duplicated `_convert_to_camel()` methods; originals retained for backwards compatibility).
  - `introspection.py`: member discovery via `discover_members()` returning a `MemberDiscovery` dataclass; uses `inspect.getmembers_static()` to avoid side effects when introspecting the new `errors` property. `Form.calculate_members()` remains the public API.
  - `validators.py`: module-level `_validate_field`, `_validate_sub_form`, `_validate_sub_form_as_list` functions extracted from `Form` for clarity.

- **Test suite**: 258 tests achieving 98% coverage (pytest + pytest-asyncio + pytest-cov).

- **Cross-language test vectors**: `vectors/fields/*.json` with 89 cases across the seven field types, to be shared with the future Go implementation.

### Known issues

The following issues are pinned by tests and slated for follow-up releases:

- `EmailField(empty=True)` skips regex validation entirely, so invalid emails pass.
- `EmailField` emits `'required'` for empty string where `CharField` emits `'empty'`.
- Optional `JsonField` emits `'invalid'` when the value is absent.
- Optional `BooleanField`, `NumberField`, `IdField` silently accept wrong types (error only when `required=True`).
- `NumberField(datatype=int)` accepts `True` (since `isinstance(True, int)` is `True`).
- Parent form mutates child form's error dicts (via `del error['code']`).
- Non-list value supplied for nested list is silently skipped with no error.
- Any `list` class attribute is claimed as a nested-form declaration.
- `cleaned_data` returned by reference (mutations leak to caller's original dict).
- `clean*` methods run in alphabetical order, not declaration order.

### Internal

- Consolidated 25 duplicated lines in `is_valid()` / `is_valid_async()` into shared helper.
- Six mutable class-level attributes promoted to instance attributes; annotations remain at class level.
- Removed unreachable code: `callable(extra_args)` branch (dict never callable), `value is None` inside `isinstance(value, str)` branch, `isinstance(self.datatype(), dict)` replaced with `issubclass`.

## 2.1.12

- Add `regex` to `CharField`

## 2.1.10

- Add `TypeError` to `NumberField`

## 2.1.9

- Exposed a method to convert strawberry objects to dicts

## 2.1.8

- Added support for strawberry objects in forms

## 2.1.7

- Fixed char field validation to accept Enum and StrEnum types

## 2.1.6

- Added support for async clean functions in forms

## 2.1.4

- Adjustments on all typings

## 2.1.3

- Fixed again the `choices` typing on `CharField`

## 2.1.2

- Removed forcing kwargs on `add_errors()` method

## 2.1.1

- Fixed typing of `choices` on `CharField`.

## 2.1.0

- Migrated to use `ruff` and `mypy` for linting and type checking
- Improvements on documentation to use PEP style

## 2.0.1

- Removed namespace scan in `pyproject.toml` to avoid issues with other libraries

## 2.0.0

- Due to a incompatibility with shared namespaces, we changed the usage of the library from `from layrz import forms` to `import layrz_forms as forms`

## 1.0.12

- Fixed issue with nested fields (Previous fix was not working properly)

## 1.0.11

- Fixed issue with nested fields

## 1.0.10

- Fixed issue with nested forms, now will be validated correctly

## v1.0.9

- Migrated to GitHub

## v1.0.7

- Fixes on package namespace

## v1.0.6

- Update logic for required in Number and Id fields

## v1.0.5

- Added support for nested fields along with forms
- Added extra condition at the moment of validate the type of the field (boolean, number, id) to avoid false invalid when is not required

## v1.0.4

- Hotfix

## v1.0.3

- Removed multiprocessing option

## v1.0.2

- hotfix

## v1.0.1

- Fixed long_description in setup

## v1.0.0

- Initial release
