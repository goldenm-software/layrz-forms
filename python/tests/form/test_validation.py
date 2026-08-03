"""Test Form validation."""

from layrz_forms import BooleanField, CharField, Form, NumberField


class TestFormValidationAllValid:
  """Test Form validation with all valid data."""

  def test_all_valid_single_field(self) -> None:
    """Test form with valid single field."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True)

    form = TestForm({'name': 'John'})
    assert form.is_valid() is True
    assert form.errors() == {}
    assert form.cleaned_data == {'name': 'John'}

  def test_all_valid_multiple_fields(self) -> None:
    """Test form with valid multiple fields."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True)
      age = NumberField(datatype=int, required=True)
      active = BooleanField(required=False)

    form = TestForm({'name': 'Jane', 'age': 30, 'active': True})
    assert form.is_valid() is True
    assert form.errors() == {}


class TestFormValidationMultiInvalid:
  """Test Form validation with multiple invalid fields."""

  def test_multiple_errors(self) -> None:
    """Test form reports errors for multiple invalid fields."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True, min_length=5)
      age = NumberField(datatype=int, required=True, min_value=0)

    form = TestForm({'name': 'Bob', 'age': -5})
    assert form.is_valid() is False
    errors = form.errors()
    assert 'name' in errors
    assert 'age' in errors
    assert errors['name'][0]['code'] == 'minLength'
    assert errors['age'][0]['code'] == 'minValue'


class TestFormValidationNone:
  """Test Form validation with None as input."""

  def test_form_none_input(self) -> None:
    """Test form initialized with None."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=False)

    form = TestForm(None)
    assert form.is_valid() is True
    assert form.errors() == {}
    assert form.cleaned_data == {}


class TestFormValidationEmpty:
  """Test Form validation with empty dict."""

  def test_form_empty_dict_input(self) -> None:
    """Test form initialized with empty dict."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=False)

    form = TestForm({})
    assert form.is_valid() is True
    assert form.errors() == {}
    assert form.cleaned_data == {}

  def test_form_empty_dict_required_field(self) -> None:
    """Test form with empty dict and required field."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True)

    form = TestForm({})
    assert form.is_valid() is False
    assert form.errors() == {'name': [{'code': 'required'}]}


class TestFormValidationTwice:
  """Test calling is_valid() twice."""

  def test_is_valid_called_twice_errors_reset(self) -> None:
    """Test calling is_valid() twice doesn't accumulate errors."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True)

    form = TestForm({})
    assert form.is_valid() is False
    assert form.errors() == {'name': [{'code': 'required'}]}

    # Call is_valid again with same data
    assert form.is_valid() is False
    # Errors should not accumulate
    assert form.errors() == {'name': [{'code': 'required'}]}
    assert len(form.errors()['name']) == 1

  def test_is_valid_called_twice_with_new_data(self) -> None:
    """Test calling is_valid() twice with different data."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True)

    form = TestForm({'name': 'John'})
    assert form.is_valid() is True
    assert form.errors() == {}

    # Update data and validate again
    form.obj = {}
    assert form.is_valid() is False
    assert form.errors() == {'name': [{'code': 'required'}]}


class TestFormCleanedDataReference:
  """Test that cleaned_data is returned by reference."""

  def test_cleaned_data_is_input_dict_reference(self) -> None:
    """Test cleaned_data is the input dict by reference."""
    input_dict = {'name': 'John', 'age': 30}

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True)
      age = NumberField(datatype=int, required=False)

    form = TestForm(input_dict)
    form.is_valid()
    assert form.cleaned_data is input_dict
    assert id(form.cleaned_data) == id(input_dict)

  def test_cleaned_data_mutation_leaks_to_caller(self) -> None:
    """Test mutations to cleaned_data leak to caller's dict."""
    input_dict = {'name': 'John'}

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True)

    form = TestForm(input_dict)
    form.is_valid()

    # KNOWN BUG: cleaned_data is returned by reference; mutations leak
    form.cleaned_data['name'] = 'Jane'
    assert input_dict['name'] == 'Jane'
    assert input_dict is form.cleaned_data

  def test_cleaned_data_add_key_leaks_to_caller(self) -> None:
    """Test adding keys to cleaned_data leaks to caller's dict."""
    input_dict = {'name': 'John'}

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True)

    form = TestForm(input_dict)
    form.is_valid()

    form.cleaned_data['extra'] = 'value'
    assert input_dict['extra'] == 'value'


class TestFormReturnsBoolean:
  """Test that is_valid() returns a boolean."""

  def test_is_valid_returns_bool_true(self) -> None:
    """Test is_valid() returns True (not None or other truthy)."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=False)

    form = TestForm({'name': 'John'})
    result = form.is_valid()
    assert result is True
    assert type(result) is bool

  def test_is_valid_returns_bool_false(self) -> None:
    """Test is_valid() returns False (not None or other falsy)."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True)

    form = TestForm({})
    result = form.is_valid()
    assert result is False
    assert type(result) is bool
