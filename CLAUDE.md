# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

From the repository root, `make` targets run **both** languages — that is the pre-push gate:

```bash
make checks          # lint + typecheck + build + test, Python and Go
make test            # tests with coverage thresholds enforced (Python 90%, Go 90%)
make format          # ruff format + gofmt -w
make install-hooks   # enable .githooks/pre-commit, which runs `make checks`
```

Per language, run `make -C python <target>` or `make -C go <target>`. The root Makefile is the
union of both; it deliberately has no `py-*`/`go-*` targets.

Python, from `python/`:

```bash
uv sync --only-group dev                                    # install dev dependencies
uv run ruff check                                           # lint
uv run ty check                                             # type check — ty, NOT mypy
uv run pytest -q                                            # tests
uv run pytest --cov=layrz_forms --cov-report=term-missing    # tests with coverage
uv run python -m build                                      # build distribution
```

Go, from `go/`:

```bash
gofmt -l .           # must print nothing
go vet ./...
go test ./...
make test            # enforces the coverage threshold in go/Makefile
```

Type checking uses **`ty`**, not mypy. Suppressions must be `# ty: ignore[rule-code]` — `ty`
silently ignores mypy's `# type: ignore[...]` and the check still fails.

Deployment to PyPI is tag-triggered via `.github/workflows/deploy.yaml` — push a `v[0-9]+.[0-9]+.[0-9]+` tag to main to build and publish.

## Architecture

**layrz-forms** is a form validation library — a simpler alternative to Django Forms — implemented
**twice**, in Python and in Go, against one shared behavioural spec.

```
python/    Python library (layrz_forms) — the REFERENCE implementation
go/        Go library (module github.com/goldenm-software/layrz-forms/go/v3)
vectors/   shared cross-language test vectors — 96 cases, the spec both sides must satisfy
.claude/   Claude Code plugin: form-builder router + py/go-form-builder skills
```

Python validates dicts, plain objects, and Strawberry GraphQL inputs. Go validates structs via
`layrz:` tags and reflection. The APIs differ by necessity — Go has no descriptors or metaclasses —
but the error output is identical: same codes, same camelCase keys, same dotted nested keys.

**`vectors/fields/*.json` is the spec, not a convenience.** Both test suites consume it. When
changing validation behaviour in either language, the vectors decide what is correct; if a change
makes a vector fail, the change is wrong unless the vector is deliberately being corrected — and
then BOTH implementations must be updated together.

### Core abstractions

**`Form` (python/layrz_forms/form.py)** — The main orchestrator. Accepts an `obj` (dict, Strawberry type, or None) at init. On `is_valid()` or `ais_valid()`, it:
1. Discovers `Field` members, nested `Form` instances, nested lists (`_attrs`), and methods prefixed with `clean` via introspection (delegated to `introspection.discover_members()`).
2. Runs each field's `validate()` method.
3. Recursively validates sub-forms and nested field/form lists.
4. Calls all `clean*` methods (async if using `ais_valid()`).
5. Exposes `errors` property (a `dict[str, list[LayrzError]]` with camelCase keys, lazily validated) and `cleaned_data` (the validated object).

**`Field` subclasses (python/layrz_forms/fields/)** — Each field handles one data type. The base class (`base.py`) provides `_append_error()` method; camelCase conversion is now delegated to `casing.to_camel_case()`. Available fields: `BooleanField`, `CharField`, `EmailField`, `IdField`, `JsonField`, `NumberField`, `UuidField`.

**`errors` (python/layrz_forms/errors.py)** — Defines `LayrzError`, a Pydantic model representing a single validation error with fields `code` (str), `expected` (Any), `received` (Any), and `extra` (dict|None). Exported from `layrz_forms.__init__`.

**`casing` (python/layrz_forms/casing.py)** — Centralized `to_camel_case(key: str) -> str` function for snake_case → camelCase conversion. Both `Form` and `Field` delegate to it; legacy `_convert_to_camel()` methods remain for backwards compatibility.

**`introspection` (python/layrz_forms/introspection.py)** — Member discovery via `discover_members(form)` returning a frozen `MemberDiscovery` dataclass containing field, nested form, and clean method lists. Uses `inspect.getmembers_static()` to avoid invoking the `errors` property as a side effect. `Form.calculate_members()` is the public entry point.

**`validators` (python/layrz_forms/validators.py)** — Module-level validation functions `_validate_field`, `_validate_sub_form`, `_validate_sub_form_as_list` extracted for clarity; `Form` retains thin delegating wrappers for backwards compatibility.

**`types` (python/layrz_forms/types.py)** — Type aliases: `ErrorType = LayrzError` (a single error) and `ErrorsType = dict[str, list[LayrzError]]` (the full error mapping). Both exported from `layrz_forms.__init__`.

### Key conventions

- **snake_case → camelCase**: All field names and error keys are auto-converted (e.g., `id_test` → `idTest`) in error output.
- **Error structure**: `{'fieldName': [LayrzError(code='errorCode', expected=..., received=..., extra=None)]}`. Each error is a Pydantic model with `code` (always set), `expected`, `received`, and `extra` fields. Use `model_dump()` to convert to plain dicts, which excludes unset fields by default.
- **Clean functions**: Any method starting with `clean` is auto-discovered and called after field validation. Use `self.add_errors(key, code, extra_args)` inside them to add custom errors; `expected` and `received` are lifted into their own fields, other keys nest under `extra`. Async clean functions are supported only via `ais_valid()`; reading `.errors` without awaiting `ais_valid()` raises `RuntimeError` if async clean methods are present.
- **Nested forms**: Assign a `Form` instance as a class attribute; it's auto-detected and validated recursively. Lists of fields/forms are declared via `_attrs`.
- **Strawberry support**: Pass a Strawberry input object directly to `Form(obj=...)`; it's converted internally via `strawberry_to_dict()`.

### Code style

- Line length: 120 chars, 2-space indents, single quotes (enforced by ruff).
- `ty` type checking (config: `[tool.ty.environment]` / `[tool.ty.src]` in `python/pyproject.toml`) — all functions must be fully typed.
- Ruff rule sets: `I` (isort), `E`/`W` (pycodestyle), `F` (pyflakes), `B` (bugbear), `TD` (flake8-todo), `DJ` (django), `DTZ` (datetimez), `T20` (flake8-print), `PYI` (flake8-pyi), `ANN` (flake8-annotations). Ignored: `F401` (unused import), `E701` (multiple statements), `TD003` (missing issue link), `ANN401` (Any type).
- Python minimum version: `>=3.14`
- Dependencies: `strawberry-graphql>=0.287.0` (for `Form.strawberry_to_dict()`), `pydantic>=2.13.4`

## Go implementation

`go/` is a redesign, not a transliteration. Rules come from `layrz:` struct tags parsed by
reflection; custom validation comes from method-name conventions:

- `Clean<FieldName>(value <FieldType>) *FieldError` — per field
- `Clean<Suffix>() Errors` — cross-field, returns the key→errors map so it can use arbitrary keys

Both need pointer receivers. Order is tag rules → per-field hooks → cross-field hooks, accumulating.

**Absence is modelled by pointers.** A pointer field is absent when nil, and `required` reports on
it. A **value** field is always present, so `required` never fires on it and a zero value is a real
value — `Name string` = `""` reports `empty`, not `required`. Value fields exist because
`graph-gophers/graphql-go` rejects a pointer for a non-null (`!`) GraphQL input, so a `String!`
must be declared as a value.

Configuration problems — a malformed tag, a tag/type mismatch, a bad hook signature, a recovered
panic — surface under the reserved `_config` key rather than panicking. A `_config` error means a
programming bug, not invalid user input.

`ToCamelCase` lowercases only the first character, so Go acronyms mangle: `ID` → `iD`, `URL` →
`uRL`. This is intentional and matches Python exactly. Do not "fix" it.

## Testing

Python: tests in `python/tests/`, currently 295 tests at ~97% coverage. Go: tests alongside the
sources in `go/`, currently ~90.7% coverage. Both enforce a 90% floor.

Behavior is pinned by tests. A test failing after a source change means the change altered
observable behavior — do not silence it without intent.

`vectors/fields/*.json` holds 96 shared cases (one file per field type) consumed by BOTH suites.
The schema uses an explicit `value_absent` flag because "absent" and "present but None" are
distinct cases. Python drives them via `python/tests/test_vectors.py`, which is parametrized per
FILE — so pytest reports 7 tests, each looping over its file's cases. Seven passing tests means all
96 cases passed; it does not mean cases were skipped.

Numbers in this file and in CHANGELOG.md drift. Measure before quoting:

```bash
python3 -c "import json,glob; print(sum(len(json.load(open(f))) for f in glob.glob('vectors/fields/*.json')))"
```
