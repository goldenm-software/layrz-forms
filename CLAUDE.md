# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Install dev dependencies
uv sync --only-group dev

# Lint
uv run ruff check

# Type check
uv run mypy .

# Build distribution
uv run python -m build

# Publish to PyPI (requires PYPI_TOKEN env var)
uv run python manual_deploy.py
```

## Architecture

**layrz-forms** is a Python form validation library — a simpler alternative to Django Forms. It validates data from dicts, plain objects, or Strawberry GraphQL objects.

### Core abstractions

**`Form` (layrz_forms/form.py)** — The main orchestrator. Accepts an `obj` (dict, Strawberry type, or None) at init. On `is_valid()` / `is_valid_async()`, it:
1. Discovers `Field` members, nested `Form` instances, nested lists (`_attrs`), and methods prefixed with `clean` via introspection.
2. Runs each field's `validate()` method.
3. Recursively validates sub-forms and nested field/form lists.
4. Calls all `clean*` methods (async if using `is_valid_async()`).
5. Exposes `errors()` (camelCase keys) and `cleaned_data` (the validated object).

**`Field` subclasses (layrz_forms/fields/)** — Each field handles one data type. The base class (`base.py`) provides `_append_error()` and `_convert_to_camel()`. Available fields: `BooleanField`, `CharField`, `EmailField`, `IdField`, `JsonField`, `NumberField`, `UuidField`.

### Key conventions

- **snake_case → camelCase**: All field names and error keys are auto-converted (e.g., `id_test` → `idTest`) in error output.
- **Error structure**: `{'fieldName': [{'code': 'errorCode', 'expected': ..., 'received': ...}]}`
- **Clean functions**: Any method starting with `clean` is auto-discovered and called after field validation. Use `self.add_errors(key, code, extra_args)` inside them to add custom errors. Async clean functions are supported only via `is_valid_async()`.
- **Nested forms**: Assign a `Form` instance as a class attribute; it's auto-detected and validated recursively. Lists of fields/forms are declared via `_attrs`.
- **Strawberry support**: Pass a Strawberry input object directly to `Form(obj=...)`; it's converted internally via `strawberry_to_dict()`.

### Code style

- Line length: 120 chars, 2-space indents, single quotes (enforced by ruff).
- mypy strict mode — all functions must be fully typed.
- Ruff rule sets: isort, pycodestyle (E/W), pyflakes (F), bugbear (B), django (DJ), datetimez (DTZ). `F401` and `E701` are ignored.
