# Form tools - Layrz
Managed by Golden M, Inc.

## Motivation
This project is a collection of tools that we use to make django developers life easier. I hope you find them useful too.

## Usage
The idea is simple, replace the django forms to a more easier way to use them. Also provide the ability to return the errors key to support i18n.

```python
import layrz.forms as forms


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
  #> form.is_valid(): None
  print('form.errors():', form.errors())
  #> form.errors(): {'rangeTextTest': [{'code': 'minLength', 'expected': 5, 'received': 4}], 'clean1': [{'code': 'error1'}, {'code': 'error2'}], 'clean2': [{'code': 'error1'}]}
```

## Work with us
Golden M is a software/hardware development company what is working on
a new, innovative and disruptive technologies.

For more information, contact us at [sales@goldenmcorp.com](mailto:sales@goldenmcorp.com)

## License
This project is under MIT License, for more information, check out the `LICENCE`
