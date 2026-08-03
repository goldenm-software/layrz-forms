"""Test the example from README.md."""

import layrz_forms as forms


class ExampleForm(forms.Form):
  """Example form from README."""

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

  def clean_func1(self) -> None:
    """Clean function 1."""
    self.add_errors(key='clean1', code='error1')
    self.add_errors(key='clean1', code='error2')

  def clean_func2(self) -> None:
    """Clean function 2."""
    self.add_errors(key='clean2', code='error1')


def test_readme_example() -> None:
  """Test the README example produces documented output."""
  obj = {
    'id_test': 1,
    'email_text': 'example@goldenmcorp.com',
    'json_dict_test': {'hola': 'mundo'},
    'json_list_test': ['hola mundo'],
    'int_test': 5,
    'float_test': 4.5,
    'bool_test': True,
    'plain_text_test': 'hola mundo',
    'empty_text_test': 'hola',
    'range_text_test': 'hola',
  }

  form = ExampleForm(obj)

  # README claims is_valid() returns None but it actually returns bool
  result = form.is_valid()
  assert isinstance(result, bool)
  assert result is False  # Actual behavior

  # README documents the errors (clean functions intentionally add errors)
  errors = form.errors()
  assert 'rangeTextTest' in errors
  assert errors['rangeTextTest'] == [{'code': 'minLength', 'expected': 5, 'received': 4}]
  assert 'clean1' in errors
  assert errors['clean1'] == [{'code': 'error1'}, {'code': 'error2'}]
  assert 'clean2' in errors
  assert errors['clean2'] == [{'code': 'error1'}]


def test_readme_example_fixed_range_text() -> None:
  """Test README example with range_text_test corrected."""
  obj = {
    'id_test': 1,
    'email_text': 'example@goldenmcorp.com',
    'json_dict_test': {'hola': 'mundo'},
    'json_list_test': ['hola mundo'],
    'int_test': 5,
    'float_test': 4.5,
    'bool_test': True,
    'plain_text_test': 'hola mundo',
    'empty_text_test': 'hola',
    'range_text_test': 'hola mundo',  # Fixed: now min_length=5
  }

  form = ExampleForm(obj)
  result = form.is_valid()
  # Still False due to clean functions adding errors
  assert result is False

  errors = form.errors()
  # range_text_test should be valid now
  assert 'rangeTextTest' not in errors
  # But clean functions still add errors
  assert 'clean1' in errors
  assert 'clean2' in errors
