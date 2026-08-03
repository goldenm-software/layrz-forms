"""Test CharField validation."""

from enum import Enum, StrEnum

from layrz_forms import CharField, Form


class StringChoice(Enum):
  """String choice enum."""

  OPTION_A = 'optionA'
  OPTION_B = 'optionB'


class StringChoiceStr(StrEnum):
  """String choice enum (StrEnum)."""

  OPTION_X = 'optionX'
  OPTION_Y = 'optionY'


class TestCharFieldAbsent:
  """Test CharField when field is absent."""

  def test_char_absent_required_true(self) -> None:
    """Test required CharField absent from input."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True)

    form = TestForm({})
    assert not form.is_valid()
    assert form.errors() == {'name': [{'code': 'required'}]}

  def test_char_absent_required_false(self) -> None:
    """Test optional CharField absent from input."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=False)

    form = TestForm({})
    assert form.is_valid()
    assert form.errors() == {}


class TestCharFieldNone:
  """Test CharField with None value."""

  def test_char_none_required_true(self) -> None:
    """Test required CharField with None."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True)

    form = TestForm({'name': None})
    assert not form.is_valid()
    assert form.errors() == {'name': [{'code': 'required'}]}

  def test_char_none_required_false(self) -> None:
    """Test optional CharField with None."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=False)

    form = TestForm({'name': None})
    assert form.is_valid()
    assert form.errors() == {}


class TestCharFieldValid:
  """Test CharField with valid values."""

  def test_char_valid_string(self) -> None:
    """Test CharField with valid string."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True)

    form = TestForm({'name': 'John'})
    assert form.is_valid()
    assert form.errors() == {}
    assert form.cleaned_data == {'name': 'John'}

  def test_char_valid_empty_string_empty_true(self) -> None:
    """Test CharField with empty string when empty=True."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=False, empty=True)

    form = TestForm({'name': ''})
    assert form.is_valid()
    assert form.errors() == {}
    assert form.cleaned_data == {'name': ''}

  def test_char_valid_empty_string_empty_false(self) -> None:
    """Test CharField with empty string when empty=False."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True, empty=False)

    form = TestForm({'name': ''})
    assert not form.is_valid()
    assert form.errors() == {'name': [{'code': 'empty'}]}


class TestCharFieldMinLength:
  """Test CharField min_length constraint."""

  def test_char_min_length_exact(self) -> None:
    """Test CharField at exact min_length."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True, min_length=5)

    form = TestForm({'name': 'hello'})
    assert form.is_valid()
    assert form.errors() == {}

  def test_char_min_length_under(self) -> None:
    """Test CharField under min_length."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True, min_length=5)

    form = TestForm({'name': 'hola'})
    assert not form.is_valid()
    assert form.errors() == {'name': [{'code': 'minLength', 'expected': 5, 'received': 4}]}

  def test_char_min_length_over(self) -> None:
    """Test CharField over min_length."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True, min_length=5)

    form = TestForm({'name': 'hello world'})
    assert form.is_valid()
    assert form.errors() == {}


class TestCharFieldMaxLength:
  """Test CharField max_length constraint."""

  def test_char_max_length_exact(self) -> None:
    """Test CharField at exact max_length."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True, max_length=5)

    form = TestForm({'name': 'hello'})
    assert form.is_valid()
    assert form.errors() == {}

  def test_char_max_length_under(self) -> None:
    """Test CharField under max_length."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True, max_length=5)

    form = TestForm({'name': 'hi'})
    assert form.is_valid()
    assert form.errors() == {}

  def test_char_max_length_over(self) -> None:
    """Test CharField over max_length."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True, max_length=5)

    form = TestForm({'name': 'hello world'})
    assert not form.is_valid()
    assert form.errors() == {'name': [{'code': 'maxLength', 'expected': 5, 'received': 11}]}


class TestCharFieldChoices:
  """Test CharField with choices constraint."""

  def test_char_choices_valid(self) -> None:
    """Test CharField with valid choice."""

    class TestForm(Form):
      """Test form."""

      status = CharField(required=True, choices=(('active', 'Active'), ('inactive', 'Inactive')))

    form = TestForm({'status': 'active'})
    assert form.is_valid()
    assert form.errors() == {}

  def test_char_choices_invalid(self) -> None:
    """Test CharField with invalid choice."""

    class TestForm(Form):
      """Test form."""

      status = CharField(required=True, choices=(('active', 'Active'), ('inactive', 'Inactive')))

    form = TestForm({'status': 'pending'})
    assert not form.is_valid()
    assert form.errors() == {
      'status': [{'code': 'invalidChoice', 'expected': ['active', 'inactive'], 'received': 'pending'}]
    }


class TestCharFieldRegex:
  """Test CharField with regex constraint."""

  def test_char_regex_match(self) -> None:
    """Test CharField matching regex."""

    class TestForm(Form):
      """Test form."""

      code = CharField(required=True, regex=r'^[A-Z]{3}$')

    form = TestForm({'code': 'ABC'})
    assert form.is_valid()
    assert form.errors() == {}

  def test_char_regex_no_match(self) -> None:
    """Test CharField not matching regex."""

    class TestForm(Form):
      """Test form."""

      code = CharField(required=True, regex=r'^[A-Z]{3}$')

    form = TestForm({'code': 'abc'})
    assert not form.is_valid()
    assert form.errors() == {'code': [{'code': 'invalidFormat', 'expected': r'^[A-Z]{3}$', 'received': 'abc'}]}

  def test_char_regex_empty_string_empty_true(self) -> None:
    """Test CharField regex validation with empty=True and empty string."""

    class TestForm(Form):
      """Test form."""

      code = CharField(required=False, empty=True, regex=r'^[A-Z]{3}$')

    form = TestForm({'code': ''})
    # When empty=True and string is empty, regex is still validated
    assert not form.is_valid()
    assert form.errors() == {'code': [{'code': 'invalidFormat', 'expected': r'^[A-Z]{3}$', 'received': ''}]}


class TestCharFieldEnumConversion:
  """Test CharField with Enum values."""

  def test_char_enum_string_choice(self) -> None:
    """Test CharField with regular Enum (uses .name)."""

    class TestForm(Form):
      """Test form."""

      choice = CharField(required=True)

    form = TestForm({'choice': StringChoice.OPTION_A})
    assert form.is_valid()
    # Enum instance is stored as-is in cleaned_data (by reference);
    # .name conversion is only during validation
    assert form.cleaned_data == {'choice': StringChoice.OPTION_A}

  def test_char_enum_str_enum(self) -> None:
    """Test CharField with StrEnum (uses .value)."""

    class TestForm(Form):
      """Test form."""

      choice = CharField(required=True)

    form = TestForm({'choice': StringChoiceStr.OPTION_X})
    assert form.is_valid()
    # StrEnum.value is 'optionX'
    assert form.cleaned_data == {'choice': 'optionX'}


class TestCharFieldInvalidType:
  """Test CharField with non-string types."""

  def test_char_invalid_type_int_required_true(self) -> None:
    """Test CharField with int when required=True."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True)

    form = TestForm({'name': 42})
    assert not form.is_valid()
    assert form.errors() == {'name': [{'code': 'invalid'}]}

  def test_char_invalid_type_int_required_false(self) -> None:
    """Test CharField with int when required=False."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=False)

    form = TestForm({'name': 42})
    assert not form.is_valid()
    assert form.errors() == {'name': [{'code': 'invalid'}]}

  def test_char_invalid_type_float_required_true(self) -> None:
    """Test CharField with float when required=True."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True)

    form = TestForm({'name': 1.5})
    assert not form.is_valid()
    assert form.errors() == {'name': [{'code': 'invalid'}]}

  def test_char_invalid_type_float_required_false(self) -> None:
    """Test CharField with float when required=False."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=False)

    form = TestForm({'name': 1.5})
    assert not form.is_valid()
    assert form.errors() == {'name': [{'code': 'invalid'}]}

  def test_char_invalid_type_bool_required_true(self) -> None:
    """Test CharField with bool when required=True."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True)

    form = TestForm({'name': True})
    assert not form.is_valid()
    assert form.errors() == {'name': [{'code': 'invalid'}]}

  def test_char_invalid_type_bool_required_false(self) -> None:
    """Test CharField with bool when required=False."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=False)

    form = TestForm({'name': False})
    assert not form.is_valid()
    assert form.errors() == {'name': [{'code': 'invalid'}]}

  def test_char_invalid_type_list_required_true(self) -> None:
    """Test CharField with list when required=True."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True)

    form = TestForm({'name': [1, 2]})
    assert not form.is_valid()
    assert form.errors() == {'name': [{'code': 'invalid'}]}

  def test_char_invalid_type_list_required_false(self) -> None:
    """Test CharField with list when required=False."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=False)

    form = TestForm({'name': [1, 2]})
    assert not form.is_valid()
    assert form.errors() == {'name': [{'code': 'invalid'}]}

  def test_char_invalid_type_dict_required_true(self) -> None:
    """Test CharField with dict when required=True."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True)

    form = TestForm({'name': {'a': 1}})
    assert not form.is_valid()
    assert form.errors() == {'name': [{'code': 'invalid'}]}

  def test_char_invalid_type_dict_required_false(self) -> None:
    """Test CharField with dict when required=False."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=False)

    form = TestForm({'name': {'a': 1}})
    assert not form.is_valid()
    assert form.errors() == {'name': [{'code': 'invalid'}]}

  def test_char_invalid_type_list_with_min_length(self) -> None:
    """Test that invalid type is caught before min_length check."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True, min_length=3)

    form = TestForm({'name': [1, 2]})
    assert not form.is_valid()
    # Only 'invalid' error, not 'minLength'
    assert form.errors() == {'name': [{'code': 'invalid'}]}

  def test_char_invalid_type_dict_with_min_length(self) -> None:
    """Test that invalid type is caught before min_length check."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True, min_length=3)

    form = TestForm({'name': {'a': 1}})
    assert not form.is_valid()
    # Only 'invalid' error, not 'minLength'
    assert form.errors() == {'name': [{'code': 'invalid'}]}

  def test_char_invalid_type_int_with_min_length(self) -> None:
    """Test that invalid type is caught before min_length check."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True, min_length=3)

    form = TestForm({'name': 42})
    assert not form.is_valid()
    # Only 'invalid' error, not 'minLength'
    assert form.errors() == {'name': [{'code': 'invalid'}]}


class TestCharFieldCombined:
  """Test CharField with multiple constraints."""

  def test_char_min_and_max_length_valid(self) -> None:
    """Test CharField with both min and max length constraints."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True, min_length=3, max_length=10)

    form = TestForm({'name': 'Alice'})
    assert form.is_valid()
    assert form.errors() == {}

  def test_char_min_and_max_length_under_min(self) -> None:
    """Test CharField under min when both constraints applied."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True, min_length=3, max_length=10)

    form = TestForm({'name': 'Al'})
    assert not form.is_valid()
    assert form.errors() == {'name': [{'code': 'minLength', 'expected': 3, 'received': 2}]}

  def test_char_min_and_max_length_over_max(self) -> None:
    """Test CharField over max when both constraints applied."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True, min_length=3, max_length=10)

    form = TestForm({'name': 'Alice Wonderland'})
    assert not form.is_valid()
    assert form.errors() == {'name': [{'code': 'maxLength', 'expected': 10, 'received': 16}]}

  def test_char_empty_false_takes_precedence(self) -> None:
    """Test that empty and min_length both validate on empty string."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=False, empty=False, min_length=3)

    form = TestForm({'name': ''})
    assert not form.is_valid()
    # Both empty and min_length errors are emitted
    assert form.errors() == {'name': [{'code': 'empty'}, {'code': 'minLength', 'expected': 3, 'received': 0}]}
