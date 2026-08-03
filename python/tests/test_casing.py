"""Test snake_case to camelCase conversion."""

from layrz_forms import CharField, Form


class TestFormCasing:
  """Test Form._convert_to_camel method."""

  def test_simple_snake_case(self) -> None:
    """Test simple snake_case conversion."""
    form = Form()
    result = form._convert_to_camel(key='user_name')
    assert result == 'userName'

  def test_single_word(self) -> None:
    """Test single word (no underscore)."""
    form = Form()
    result = form._convert_to_camel(key='name')
    assert result == 'name'

  def test_multiple_underscores(self) -> None:
    """Test multiple consecutive underscores."""
    form = Form()
    result = form._convert_to_camel(key='a__b')
    assert result == 'aB'

  def test_dotted_keys(self) -> None:
    """Test dotted notation with snake_case segments."""
    form = Form()
    result = form._convert_to_camel(key='a_b.c_d')
    assert result == 'aB.cD'

  def test_already_camel_case(self) -> None:
    """Test already camelCase input."""
    form = Form()
    result = form._convert_to_camel(key='userName')
    assert result == 'userName'

  def test_all_caps(self) -> None:
    """Test all caps conversion."""
    form = Form()
    result = form._convert_to_camel(key='ID_TEST')
    assert result == 'iDTest'

  def test_empty_string(self) -> None:
    """Test empty string conversion."""
    form = Form()
    result = form._convert_to_camel(key='')
    assert result == ''

  def test_trailing_underscore(self) -> None:
    """Test key with trailing underscore."""
    form = Form()
    result = form._convert_to_camel(key='test_')
    assert result == 'test'

  def test_leading_underscore(self) -> None:
    """Test key with leading underscore."""
    form = Form()
    result = form._convert_to_camel(key='_test')
    assert result == 'test'


class TestFieldCasing:
  """Test Field._convert_to_camel method."""

  def test_simple_snake_case(self) -> None:
    """Test simple snake_case conversion via field."""
    field = CharField()
    result = field._convert_to_camel('user_name')
    assert result == 'userName'

  def test_single_word(self) -> None:
    """Test single word via field."""
    field = CharField()
    result = field._convert_to_camel('name')
    assert result == 'name'

  def test_dotted_keys(self) -> None:
    """Test dotted notation via field."""
    field = CharField()
    result = field._convert_to_camel('a_b.c_d')
    assert result == 'aB.cD'

  def test_empty_string(self) -> None:
    """Test empty string via field."""
    field = CharField()
    result = field._convert_to_camel('')
    assert result == ''
