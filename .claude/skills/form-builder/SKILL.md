---
name: form-builder
description: Router for building, reviewing, or debugging a validation form with the layrz-forms library. Use this skill whenever the user wants to validate structured input — request payloads, GraphQL inputs, dicts, structs, API bodies — or mentions layrz-forms, layrz_forms, a "form class", a `layrz:` struct tag, `is_valid()`, `ais_valid()`, `Validate()`, `LayrzError`, `FieldError`, a `clean_*` method, a `Clean<Field>` hook, or error codes like `required`, `invalid`, `empty`, `minLength`, `invalidChoice`. Also use it when the user asks "why is my form not catching X", "how do I validate a nested list", or wants to add a field to an existing form — even if they never say the library's name. This skill only picks the language and hands off; it does not itself explain how to write forms.
---

# Form builder — router

Detect the language, then invoke that sub-skill. Do not answer the question from this file — the
Python and Go APIs are shaped completely differently, and every detail that matters lives in the
sub-skill.

## Detect

Work through these in order and stop at the first that resolves:

1. **The user names a language.** "in Go", "for my Python service", "the Django side" → done.
2. **The user pasted code.** `class MyForm(Form)`, `CharField(`, `def clean_`, `self.add_errors` →
   Python. A struct with backtick tags, `layrz:"`, `*string`, `func (f *Form)` → Go.
3. **The files under discussion.** A path ending `.py` → Python; `.go` → Go.
4. **The working directory.** Check what is actually present:

```bash
ls go/go.mod python/pyproject.toml 2>/dev/null
```

In this monorepo both exist, so fall through to the next step rather than guessing from this alone.
In a consuming project, whichever is present is the answer.

5. **Still ambiguous.** Ask. One short question — "Python or Go?" — beats writing the wrong API and
   having it thrown away. Do not guess when a wrong guess wastes the whole turn.

## Hand off

| Language | Invoke |
| --- | --- |
| Python | `py-form-builder` |
| Go | `go-form-builder` |

Use the Skill tool with that name, then follow what it says.

**Both languages.** When the task genuinely spans both — porting a form, keeping two
implementations aligned, comparing behaviour — load `py-form-builder` first, because Python is the
reference implementation the Go port was validated against, then `go-form-builder`.

**Already routed.** If a sub-skill is loaded in this turn, keep using it. Do not come back here.

## One thing worth knowing before you route

Both implementations are pinned by the same test vectors at `vectors/fields/*.json` (96 cases), so
their error codes and output shape are identical by construction. If a question is purely about what
output some input produces — and not about how to write the form — a vector answers it directly and
authoritatively:

```bash
python3 -c "import json;[print(c['name'], c.get('value','ABSENT'), c['expected_errors']) for c in json.load(open('vectors/fields/NumberField.json'))]"
```

A vector case beats reasoning it out in either language. Everything else — syntax, field
declaration, hooks, nesting, absence semantics — comes from the sub-skill.
