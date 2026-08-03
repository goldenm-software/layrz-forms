# Layrz Forms

[![PyPI](https://img.shields.io/pypi/v/layrz_forms.svg)](https://pypi.org/project/layrz-forms/)
[![GitHub license](https://img.shields.io/github/license/goldenm-software/layrz-forms?logo=github)](https://github.com/goldenm-software/layrz-forms)

A collection of tools that we use to make django developers life easier. I hope you find them useful too.

Tired of complex forms validations? The default Django schema is too complex? Layrz Forms is for you! Works like Django Forms but with a simpler API and more features. Highly personalizable and extensible, you can use it in your projects without any problem.

```python
import layrz_forms as forms


class ExampleForm(forms.Form):
  """ Example form """
  id_test = forms.IdField(required=True)
  email_text = forms.EmailField(required=True)
  json_list_test = forms.JsonField(required=True, datatype=list)
  json_dict_test = forms.JsonField(required=True, datatype=dict)
  int_test = forms.NumberField(required=True, datatype=int, min_value=0, max_value=5)
  float_test = forms.NumberField(required=True, datatype=float, min_value=0, max_value=5)
  bool_test = forms.BooleanField(required=True)
  plain_text_test = forms.CharField(required=True, empty=False)
  empty_text_test = forms.CharField(required=True, empty=True)
  range_text_test = forms.CharField(required=True, empty=False, min_length=5, max_length=10)

  def clean_func1(self):
    """ Print clean """
    self.add_errors(key='clean1', code='error1')
    self.add_errors(key='clean1', code='error2')

  def clean_func2(self):
    self.add_errors(key='clean2', code='error1')


if __name__ == '__main__':
  obj = {
    'id_test': 1,
    'email_text': 'example@goldenmcorp.com',
    'json_dict_test': {
      'hola': 'mundo'
    },
    'json_list_test': ['hola mundo'],
    'int_test': 5,
    'float_test': 4.5,
    'bool_test': True,
    'plain_text_test': 'hola mundo',
    'empty_text_test': 'hola',
    'range_text_test': 'hola'
  }

  form = ExampleForm(obj)

  print('form.is_valid():', form.is_valid())
  #> form.is_valid(): False
  print('form.errors:', form.errors)
  #> form.errors: {'rangeTextTest': [LayrzError(code='minLength', expected=5, received=4, extra=None)], 'clean1': [LayrzError(code='error1', expected=None, received=None, extra=None), LayrzError(code='error2', expected=None, received=None, extra=None)], 'clean2': [LayrzError(code='error1', expected=None, received=None, extra=None)]}

  # Convert to plain dicts for JSON serialization:
  errors_as_dicts = {k: [e.model_dump() for e in v] for k, v in form.errors.items()}
  print('errors_as_dicts:', errors_as_dicts)
  #> errors_as_dicts: {'rangeTextTest': [{'code': 'minLength', 'expected': 5, 'received': 4}], 'clean1': [{'code': 'error1'}, {'code': 'error2'}], 'clean2': [{'code': 'error1'}]}
```

## Errors

Form errors are represented as `LayrzError` Pydantic models with four fields:

| Field | Type | Description |
|-------|------|-------------|
| `code` | `str` | Error code (e.g. `'required'`, `'minLength'`, `'invalidChoice'`) |
| `expected` | `Any` | The constraint that was violated (e.g. `5`, `['a','b']`, a regex); `None` if not applicable |
| `received` | `Any` | The offending value; `None` if not applicable |
| `extra` | `dict\|None` | Additional context keys from custom validators; `None` if empty |

Access errors via the `form.errors` property (note: not a method). The property **lazily validates** on first access — if `is_valid()` hasn't been called, reading `.errors` runs synchronous validation automatically.

Serialize errors to plain dicts using `.model_dump()`, which excludes unset fields by default:

```python
{k: [e.model_dump() for e in v] for k, v in form.errors.items()}
```

When adding custom errors via `self.add_errors(key, code, extra_args={...})`, keys `expected` and `received` are lifted into their own fields; all other keys nest under `extra`:

```python
self.add_errors('password', 'weak_password', extra_args={'min_length': 8})
# Produces: LayrzError(code='weak_password', extra={'min_length': 8})
```

**Async clean functions caveat:** If your form has `async def clean_*` methods, you must await `form.ais_valid()` before reading `form.errors`. Reading `.errors` without awaiting `ais_valid()` raises `RuntimeError`.

## Clean Methods

All methods named with the `clean_*` prefix are automatically discovered and executed after field validation. They are executed in **alphabetical order by method name** (e.g., `clean_apple`, then `clean_banana`, then `clean_zebra`). This behavior is stable and intentional to match the Go implementation, where declaration order is not available via reflection.

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
