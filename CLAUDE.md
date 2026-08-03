# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

All commands run from the `python/` directory.

```bash
# Install dev dependencies
uv sync --only-group dev

# Lint
uv run ruff check

# Type check
uv run ty check

# Run tests
uv run pytest -q

# Run tests with coverage (CI requires 90% threshold; current is 98%)
uv run pytest --cov=layrz_forms --cov-report=term-missing

# Build distribution
uv run python -m build
```

Deployment to PyPI is tag-triggered via `.github/workflows/deploy.yaml` — push a `v[0-9]+.[0-9]+.[0-9]+` tag to main to build and publish.

## Architecture

**layrz-forms** is a Python form validation library — a simpler alternative to Django Forms. It validates data from dicts, plain objects, or Strawberry GraphQL objects. The package is located at `python/layrz_forms/` (this is a monorepo; `vectors/` holds shared cross-language test vectors for a future Go implementation).

### Core abstractions

**`Form` (python/layrz_forms/form.py)** — The main orchestrator. Accepts an `obj` (dict, Strawberry type, or None) at init. On `is_valid()` / `is_valid_async()`, it:
1. Discovers `Field` members, nested `Form` instances, nested lists (`_attrs`), and methods prefixed with `clean` via introspection (delegated to `introspection.discover_members()`).
2. Runs each field's `validate()` method.
3. Recursively validates sub-forms and nested field/form lists.
4. Calls all `clean*` methods (async if using `is_valid_async()`).
5. Exposes `errors()` (camelCase keys) and `cleaned_data` (the validated object).

**`Field` subclasses (python/layrz_forms/fields/)** — Each field handles one data type. The base class (`base.py`) provides `_append_error()` method; camelCase conversion is now delegated to `casing.to_camel_case()`. Available fields: `BooleanField`, `CharField`, `EmailField`, `IdField`, `JsonField`, `NumberField`, `UuidField`.

**`casing` (python/layrz_forms/casing.py)** — Centralized `to_camel_case(key: str) -> str` function for snake_case → camelCase conversion. Both `Form` and `Field` delegate to it; legacy `_convert_to_camel()` methods remain for backwards compatibility.

**`introspection` (python/layrz_forms/introspection.py)** — Member discovery via `discover_members(form)` returning a frozen `MemberDiscovery` dataclass containing field, nested form, and clean method lists. `Form.calculate_members()` is the public entry point.

**`validators` (python/layrz_forms/validators.py)** — Module-level validation functions `_validate_field`, `_validate_sub_form`, `_validate_sub_form_as_list` extracted for clarity; `Form` retains thin delegating wrappers for backwards compatibility.

### Key conventions

- **snake_case → camelCase**: All field names and error keys are auto-converted (e.g., `id_test` → `idTest`) in error output.
- **Error structure**: `{'fieldName': [{'code': 'errorCode', 'expected': ..., 'received': ...}]}`
- **Clean functions**: Any method starting with `clean` is auto-discovered and called after field validation. Use `self.add_errors(key, code, extra_args)` inside them to add custom errors. Async clean functions are supported only via `is_valid_async()`.
- **Nested forms**: Assign a `Form` instance as a class attribute; it's auto-detected and validated recursively. Lists of fields/forms are declared via `_attrs`.
- **Strawberry support**: Pass a Strawberry input object directly to `Form(obj=...)`; it's converted internally via `strawberry_to_dict()`.

### Code style

- Line length: 120 chars, 2-space indents, single quotes (enforced by ruff).
- `ty` type checking (config: `[tool.ty.environment]` / `[tool.ty.src]` in `python/pyproject.toml`) — all functions must be fully typed.
- Ruff rule sets: `I` (isort), `E`/`W` (pycodestyle), `F` (pyflakes), `B` (bugbear), `TD` (flake8-todo), `DJ` (django), `DTZ` (datetimez), `T20` (flake8-print), `PYI` (flake8-pyi), `ANN` (flake8-annotations). Ignored: `F401` (unused import), `E701` (multiple statements), `TD003` (missing issue link), `ANN401` (Any type).
- Python minimum version: `>=3.14`
- Dependencies: `strawberry-graphql>=0.287.0` (for `Form.strawberry_to_dict()`), `pydantic>=2.13.4`

## Testing

Tests live in `python/tests/` with 258 test cases achieving 98% coverage. Run `uv run pytest -q` from `python/` to execute.

Behavior is pinned by tests, including known bugs (see CHANGELOG.md). A test failing after a source change indicates the change altered observable behavior — do not silence it without intent.

Cross-language test vectors are in `vectors/fields/*.json` (one file per field type) with ~82 cases each. This is the shared spec for the future Go implementation; the schema uses an explicit `value_absent` flag because "absent" and "present but None" are distinct cases.

CI enforces a 90% coverage threshold; current coverage is 98%.
