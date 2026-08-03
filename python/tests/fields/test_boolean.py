"""Test BooleanField validation."""

from layrz_forms import BooleanField, Form


class TestBooleanFieldAbsent:
  """Test BooleanField when field is absent."""

  def test_boolean_absent_required_true(self) -> None:
    """Test required BooleanField absent from input."""

    class TestForm(Form):
      """Test form."""

      flag = BooleanField(required=True)

    form = TestForm({})
    assert not form.is_valid()
    assert form.errors() == {'flag': [{'code': 'required'}]}
    assert form.cleaned_data == {}

  def test_boolean_absent_required_false(self) -> None:
    """Test optional BooleanField absent from input."""

    class TestForm(Form):
      """Test form."""

      flag = BooleanField(required=False)

    form = TestForm({})
    assert form.is_valid()
    assert form.errors() == {}
    assert form.cleaned_data == {}


class TestBooleanFieldNone:
  """Test BooleanField when value is None."""

  def test_boolean_none_required_true(self) -> None:
    """Test required BooleanField with None value."""

    class TestForm(Form):
      """Test form."""

      flag = BooleanField(required=True)

    form = TestForm({'flag': None})
    assert not form.is_valid()
    assert form.errors() == {'flag': [{'code': 'required'}]}

  def test_boolean_none_required_false(self) -> None:
    """Test optional BooleanField with None value."""

    class TestForm(Form):
      """Test form."""

      flag = BooleanField(required=False)

    form = TestForm({'flag': None})
    assert form.is_valid()
    assert form.errors() == {}


class TestBooleanFieldValid:
  """Test BooleanField with valid boolean values."""

  def test_boolean_true(self) -> None:
    """Test BooleanField with True value."""

    class TestForm(Form):
      """Test form."""

      flag = BooleanField(required=True)

    form = TestForm({'flag': True})
    assert form.is_valid()
    assert form.errors() == {}
    assert form.cleaned_data == {'flag': True}

  def test_boolean_false(self) -> None:
    """Test BooleanField with False value."""

    class TestForm(Form):
      """Test form."""

      flag = BooleanField(required=True)

    form = TestForm({'flag': False})
    assert form.is_valid()
    assert form.errors() == {}
    assert form.cleaned_data == {'flag': False}


class TestBooleanFieldInvalid:
  """Test BooleanField with invalid (non-bool) values."""

  def test_boolean_non_bool_string_required_true(self) -> None:
    """Test required BooleanField with string value."""

    class TestForm(Form):
      """Test form."""

      flag = BooleanField(required=True)

    form = TestForm({'flag': 'not_bool'})
    assert not form.is_valid()
    assert form.errors() == {'flag': [{'code': 'invalid'}]}

  def test_boolean_non_bool_string_required_false(self) -> None:
    """Test optional BooleanField with string value."""

    class TestForm(Form):
      """Test form."""

      flag = BooleanField(required=False)

    form = TestForm({'flag': 'not_bool'})
    # KNOWN BUG: optional fields silently accept wrong types
    assert form.is_valid()
    assert form.errors() == {}

  def test_boolean_non_bool_int_required_true(self) -> None:
    """Test required BooleanField with integer value."""

    class TestForm(Form):
      """Test form."""

      flag = BooleanField(required=True)

    form = TestForm({'flag': 1})
    assert not form.is_valid()
    assert form.errors() == {'flag': [{'code': 'invalid'}]}

  def test_boolean_non_bool_int_required_false(self) -> None:
    """Test optional BooleanField with integer value."""

    class TestForm(Form):
      """Test form."""

      flag = BooleanField(required=False)

    form = TestForm({'flag': 1})
    # KNOWN BUG: optional fields silently accept wrong types
    assert form.is_valid()
    assert form.errors() == {}

  def test_boolean_non_bool_list_required_true(self) -> None:
    """Test required BooleanField with list value."""

    class TestForm(Form):
      """Test form."""

      flag = BooleanField(required=True)

    form = TestForm({'flag': []})
    assert not form.is_valid()
    assert form.errors() == {'flag': [{'code': 'invalid'}]}

  def test_boolean_non_bool_dict_required_false(self) -> None:
    """Test optional BooleanField with dict value."""

    class TestForm(Form):
      """Test form."""

      flag = BooleanField(required=False)

    form = TestForm({'flag': {}})
    # KNOWN BUG: optional fields silently accept wrong types
    assert form.is_valid()
    assert form.errors() == {}
