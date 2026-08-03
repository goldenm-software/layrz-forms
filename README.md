# Layrz Forms

[![PyPI](https://img.shields.io/pypi/v/layrz_forms.svg)](https://pypi.org/project/layrz-forms/)
[![GitHub license](https://img.shields.io/github/license/goldenm-software/layrz-forms?logo=github)](https://github.com/goldenm-software/layrz-forms)

Form validation for Python and Go — a simpler alternative to Django Forms, implemented twice against
one shared spec.

Tired of complex form validations? Layrz Forms validates dicts, plain objects, Strawberry GraphQL
inputs and Go structs with a small API, and gives you the same error output in both languages.

## Pick your language

This is a monorepo. Each implementation documents itself:

| | Package | Documentation |
|---|---|---|
| **Python** | `pip install layrz-forms` | **[python/README.md](python/README.md)** |
| **Go** | `go get github.com/goldenm-software/layrz-forms/go/v3` | **[go/README.md](go/README.md)** |

Start there — the two APIs are shaped differently, and each README is the full reference for its
language.

## Repository layout

```
python/    Python library (layrz_forms) — the reference implementation
go/        Go library (github.com/goldenm-software/layrz-forms/go/v3)
vectors/   shared cross-language test vectors — the spec both sides satisfy
.claude/   Claude Code plugin (see below)
```

## One spec, two implementations

Python declares fields as class attributes; Go uses `layrz:` struct tags parsed by reflection. The
APIs differ by necessity — Go has no descriptors or metaclasses — but the observable behaviour is
identical:

- the same error codes (`required`, `invalid`, `empty`, `minLength`, `invalidChoice`, …)
- camelCase keys, with dotted paths for nested values (`address.streetName`, `items.0.name`)
- errors accumulate; nothing fails fast

That equivalence is enforced, not aspirational: `vectors/fields/*.json` holds 96 shared cases and
**both** test suites run them. A change that breaks a vector in one language is a change that has
diverged from the other.

## Claude Code skill

This repository includes a *Claude Code plugin* as part of our initiative to provide AI-assisted
development tools. The plugin contains skills that guide developers in writing and debugging forms
in either language: a `form-builder` router that detects which language you are working in, plus
`py-form-builder` and `go-form-builder`, each documenting one implementation in full — every field
argument and error code, custom validation hooks, nested structures, and the mistakes that actually
come up. The Go skill also covers interoperating with `graph-gophers/graphql-go`.

### Installation

Add this repository as a Claude Code plugin marketplace, then install the plugin:

```bash
/plugin marketplace add goldenm-software/layrz-forms
```

Once the marketplace is added, install the plugin from the **Discover** tab in `/plugin`, or run:

```bash
/plugin install layrz-forms@layrz-forms
```

Then reload your plugins:

```bash
/reload-plugins
```

## Development

From the repository root, `make` runs both languages:

```bash
make checks          # lint, typecheck, build and test, Python and Go
make test            # tests with coverage thresholds enforced
make format          # ruff format + gofmt -w
make install-hooks   # enable the pre-commit hook, which runs `make checks`
```

For a single language, use `make -C python <target>` or `make -C go <target>`.

## FAQ

### Why is this package called `layrz-forms`?

All packages developed by [Layrz](https://layrz.com) are prefixed with `layrz`, check out our other packages on [PyPi](https://pypi.org/user/layrz-software/) and [GitHub](https://github.com/goldenm-software).

### Why this library exists?

We validate a lot of structured input across our services — API payloads, GraphQL inputs, message bodies — and Django Forms is heavier than we need for that. So we built `layrz-forms` as a smaller alternative, and then ported it to Go so both halves of our stack validate identically against the same shared spec. We think it could be useful for other developers, so we decided to share it with the community.

### Do you have other libraries?

Of course! We have multiple libraries (for Layrz or general purpose) that you can use in your projects, you can find us on [PyPi of Golden M](https://pypi.org/user/goldenm/) or [PyPi of Layrz](https://pypi.org/user/layrz-software/) for Python libraries, [RubyGems](https://rubygems.org/profiles/goldenm) for Ruby gems, [NPM of Golden M](https://www.npmjs.com/~goldenm) or [NPM of Layrz](https://www.npmjs.com/~layrz-software) for NodeJS libraries or here in [Pub.dev](https://pub.dev/publishers/goldenm.com/packages) for Dart/Flutter libraries.

### I need to pay to use this package?

**No!** This library is free and open source, you can use it in your projects without any cost, but if you want to support us, give us an star on our [Repository](https://github.com/goldenm-software/layrz-forms)!

### Can I contribute to this package?

**Yes!** We are open to contributions, feel free to open a pull request or an issue on the [Repository](https://github.com/goldenm-software/layrz-forms)!

### I have a question, how can I contact you?

If you need more assistance, you open an issue on the [Repository](https://github.com/goldenm-software/layrz-forms) and we're happy to help you :)

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/goldenm-software/layrz-forms/blob/main/LICENSE) file for details.

This project is maintained by [Golden M](https://goldenm.com) with authorization of [Layrz LTD](https://layrz.com).

## Who are you? / Want to work with us?

**Golden M** is a software and hardware development company what is working on a new, innovative and disruptive technologies. For more information, contact us at [sales@goldenm.com](mailto:sales@goldenm.com) or via WhatsApp at [+(507)-6979-3073](https://wa.me/50769793073?text="From%20layrz-forms%20library.%20Hello").
