"""Test NumberField validation."""

from layrz_forms import Form, NumberField


class TestNumberFieldIntAbsent:
  """Test NumberField (int) when field is absent."""

  def test_number_int_absent_required_true(self) -> None:
    """Test required NumberField(int) absent from input."""

    class TestForm(Form):
      """Test form."""

      count = NumberField(datatype=int, required=True)

    form = TestForm({})
    assert not form.is_valid()
    assert form.errors() == {'count': [{'code': 'required'}]}

  def test_number_int_absent_required_false(self) -> None:
    """Test optional NumberField(int) absent from input."""

    class TestForm(Form):
      """Test form."""

      count = NumberField(datatype=int, required=False)

    form = TestForm({})
    assert form.is_valid()
    assert form.errors() == {}


class TestNumberFieldFloatAbsent:
  """Test NumberField (float) when field is absent."""

  def test_number_float_absent_required_true(self) -> None:
    """Test required NumberField(float) absent from input."""

    class TestForm(Form):
      """Test form."""

      price = NumberField(datatype=float, required=True)

    form = TestForm({})
    assert not form.is_valid()
    assert form.errors() == {'price': [{'code': 'required'}]}

  def test_number_float_absent_required_false(self) -> None:
    """Test optional NumberField(float) absent from input."""

    class TestForm(Form):
      """Test form."""

      price = NumberField(datatype=float, required=False)

    form = TestForm({})
    assert form.is_valid()
    assert form.errors() == {}


class TestNumberFieldIntNone:
  """Test NumberField (int) with None value."""

  def test_number_int_none_required_true(self) -> None:
    """Test required NumberField(int) with None."""

    class TestForm(Form):
      """Test form."""

      count = NumberField(datatype=int, required=True)

    form = TestForm({'count': None})
    assert not form.is_valid()
    assert form.errors() == {'count': [{'code': 'required'}]}

  def test_number_int_none_required_false(self) -> None:
    """Test optional NumberField(int) with None."""

    class TestForm(Form):
      """Test form."""

      count = NumberField(datatype=int, required=False)

    form = TestForm({'count': None})
    assert form.is_valid()
    assert form.errors() == {}


class TestNumberFieldFloatNone:
  """Test NumberField (float) with None value."""

  def test_number_float_none_required_true(self) -> None:
    """Test required NumberField(float) with None."""

    class TestForm(Form):
      """Test form."""

      price = NumberField(datatype=float, required=True)

    form = TestForm({'price': None})
    assert not form.is_valid()
    assert form.errors() == {'price': [{'code': 'required'}]}

  def test_number_float_none_required_false(self) -> None:
    """Test optional NumberField(float) with None."""

    class TestForm(Form):
      """Test form."""

      price = NumberField(datatype=float, required=False)

    form = TestForm({'price': None})
    assert form.is_valid()
    assert form.errors() == {}


class TestNumberFieldIntValid:
  """Test NumberField (int) with valid values."""

  def test_number_int_valid_positive(self) -> None:
    """Test NumberField(int) with positive integer."""

    class TestForm(Form):
      """Test form."""

      count = NumberField(datatype=int, required=True)

    form = TestForm({'count': 42})
    assert form.is_valid()
    assert form.errors() == {}
    assert form.cleaned_data == {'count': 42}

  def test_number_int_valid_zero(self) -> None:
    """Test NumberField(int) with zero."""

    class TestForm(Form):
      """Test form."""

      count = NumberField(datatype=int, required=True)

    form = TestForm({'count': 0})
    assert form.is_valid()
    assert form.errors() == {}

  def test_number_int_valid_negative(self) -> None:
    """Test NumberField(int) with negative integer."""

    class TestForm(Form):
      """Test form."""

      count = NumberField(datatype=int, required=True)

    form = TestForm({'count': -10})
    assert form.is_valid()
    assert form.errors() == {}

  def test_number_int_valid_bool_true(self) -> None:
    """Test NumberField(int) with True (isinstance(True, int) is True)."""

    class TestForm(Form):
      """Test form."""

      count = NumberField(datatype=int, required=True)

    form = TestForm({'count': True})
    # KNOWN BUG: NumberField(datatype=int) accepts True since isinstance(True, int)
    assert form.is_valid()
    assert form.errors() == {}


class TestNumberFieldFloatValid:
  """Test NumberField (float) with valid values."""

  def test_number_float_valid_positive(self) -> None:
    """Test NumberField(float) with positive float."""

    class TestForm(Form):
      """Test form."""

      price = NumberField(datatype=float, required=True)

    form = TestForm({'price': 3.14})
    assert form.is_valid()
    assert form.errors() == {}
    assert form.cleaned_data == {'price': 3.14}

  def test_number_float_valid_integer_as_float(self) -> None:
    """Test NumberField(float) with integer value."""

    class TestForm(Form):
      """Test form."""

      price = NumberField(datatype=float, required=True)

    form = TestForm({'price': 42})
    # KNOWN BUG: NumberField(float) rejects int values (isinstance(42, float) is False)
    assert not form.is_valid()
    assert form.errors() == {'price': [{'code': 'invalid'}]}

  def test_number_float_valid_zero(self) -> None:
    """Test NumberField(float) with zero."""

    class TestForm(Form):
      """Test form."""

      price = NumberField(datatype=float, required=True)

    form = TestForm({'price': 0.0})
    assert form.is_valid()
    assert form.errors() == {}


class TestNumberFieldIntMinValue:
  """Test NumberField (int) with min_value constraint."""

  def test_number_int_min_value_exact(self) -> None:
    """Test NumberField(int) at exact min_value."""

    class TestForm(Form):
      """Test form."""

      count = NumberField(datatype=int, required=True, min_value=5)

    form = TestForm({'count': 5})
    assert form.is_valid()
    assert form.errors() == {}

  def test_number_int_min_value_under(self) -> None:
    """Test NumberField(int) under min_value."""

    class TestForm(Form):
      """Test form."""

      count = NumberField(datatype=int, required=True, min_value=5)

    form = TestForm({'count': 4})
    assert not form.is_valid()
    assert form.errors() == {'count': [{'code': 'minValue', 'expected': 5, 'received': 4}]}

  def test_number_int_min_value_over(self) -> None:
    """Test NumberField(int) over min_value."""

    class TestForm(Form):
      """Test form."""

      count = NumberField(datatype=int, required=True, min_value=5)

    form = TestForm({'count': 10})
    assert form.is_valid()
    assert form.errors() == {}


class TestNumberFieldIntMaxValue:
  """Test NumberField (int) with max_value constraint."""

  def test_number_int_max_value_exact(self) -> None:
    """Test NumberField(int) at exact max_value."""

    class TestForm(Form):
      """Test form."""

      count = NumberField(datatype=int, required=True, max_value=10)

    form = TestForm({'count': 10})
    assert form.is_valid()
    assert form.errors() == {}

  def test_number_int_max_value_under(self) -> None:
    """Test NumberField(int) under max_value."""

    class TestForm(Form):
      """Test form."""

      count = NumberField(datatype=int, required=True, max_value=10)

    form = TestForm({'count': 5})
    assert form.is_valid()
    assert form.errors() == {}

  def test_number_int_max_value_over(self) -> None:
    """Test NumberField(int) over max_value."""

    class TestForm(Form):
      """Test form."""

      count = NumberField(datatype=int, required=True, max_value=10)

    form = TestForm({'count': 15})
    assert not form.is_valid()
    assert form.errors() == {'count': [{'code': 'maxValue', 'expected': 10, 'received': 15}]}


class TestNumberFieldFloatMinMaxValue:
  """Test NumberField (float) with min/max constraints."""

  def test_number_float_min_value(self) -> None:
    """Test NumberField(float) with min_value."""

    class TestForm(Form):
      """Test form."""

      price = NumberField(datatype=float, required=True, min_value=1.5)

    form = TestForm({'price': 1.5})
    assert form.is_valid()
    assert form.errors() == {}

  def test_number_float_min_value_under(self) -> None:
    """Test NumberField(float) under min_value."""

    class TestForm(Form):
      """Test form."""

      price = NumberField(datatype=float, required=True, min_value=1.5)

    form = TestForm({'price': 1.4})
    assert not form.is_valid()
    assert form.errors() == {'price': [{'code': 'minValue', 'expected': 1.5, 'received': 1.4}]}

  def test_number_float_max_value(self) -> None:
    """Test NumberField(float) with max_value."""

    class TestForm(Form):
      """Test form."""

      price = NumberField(datatype=float, required=True, max_value=99.99)

    form = TestForm({'price': 99.99})
    assert form.is_valid()
    assert form.errors() == {}

  def test_number_float_max_value_over(self) -> None:
    """Test NumberField(float) over max_value."""

    class TestForm(Form):
      """Test form."""

      price = NumberField(datatype=float, required=True, max_value=99.99)

    form = TestForm({'price': 100.0})
    assert not form.is_valid()
    assert form.errors() == {'price': [{'code': 'maxValue', 'expected': 99.99, 'received': 100.0}]}


class TestNumberFieldIntInvalid:
  """Test NumberField (int) with invalid values."""

  def test_number_int_wrong_type_string_required_true(self) -> None:
    """Test required NumberField(int) with string value."""

    class TestForm(Form):
      """Test form."""

      count = NumberField(datatype=int, required=True)

    form = TestForm({'count': 'not a number'})
    assert not form.is_valid()
    assert form.errors() == {'count': [{'code': 'invalid'}]}

  def test_number_int_wrong_type_string_required_false(self) -> None:
    """Test optional NumberField(int) with string value."""

    class TestForm(Form):
      """Test form."""

      count = NumberField(datatype=int, required=False)

    form = TestForm({'count': 'not a number'})
    # KNOWN BUG: optional fields silently accept wrong types
    assert form.is_valid()
    assert form.errors() == {}

  def test_number_int_string_numeric_required_true(self) -> None:
    """Test required NumberField(int) with numeric string."""

    class TestForm(Form):
      """Test form."""

      count = NumberField(datatype=int, required=True)

    form = TestForm({'count': '42'})
    assert not form.is_valid()
    assert form.errors() == {'count': [{'code': 'invalid'}]}

  def test_number_float_string_numeric_required_true(self) -> None:
    """Test required NumberField(float) with numeric string."""

    class TestForm(Form):
      """Test form."""

      price = NumberField(datatype=float, required=True)

    form = TestForm({'price': '3.14'})
    assert not form.is_valid()
    assert form.errors() == {'price': [{'code': 'invalid'}]}

  def test_number_int_float_required_true(self) -> None:
    """Test required NumberField(int) with float value."""

    class TestForm(Form):
      """Test form."""

      count = NumberField(datatype=int, required=True)

    form = TestForm({'count': 3.14})
    assert not form.is_valid()
    assert form.errors() == {'count': [{'code': 'invalid'}]}

  def test_number_float_list_required_true(self) -> None:
    """Test required NumberField(float) with list value."""

    class TestForm(Form):
      """Test form."""

      price = NumberField(datatype=float, required=True)

    form = TestForm({'price': []})
    assert not form.is_valid()
    assert form.errors() == {'price': [{'code': 'invalid'}]}
