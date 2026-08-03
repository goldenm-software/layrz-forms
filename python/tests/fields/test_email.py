"""Test EmailField validation."""

from layrz_forms import EmailField, Form
from tests.helpers import dump_errors


class TestEmailFieldAbsent:
  """Test EmailField when field is absent."""

  def test_email_absent_required_true(self) -> None:
    """Test required EmailField absent from input."""

    class TestForm(Form):
      """Test form."""

      email = EmailField(required=True)

    form = TestForm({})
    assert not form.is_valid()
    # KNOWN BUG: EmailField emits both 'required' and 'invalid' when absent
    assert dump_errors(form.errors) == {'email': [{'code': 'required'}, {'code': 'invalid'}]}

  def test_email_absent_required_false(self) -> None:
    """Test optional EmailField absent from input."""

    class TestForm(Form):
      """Test form."""

      email = EmailField(required=False)

    form = TestForm({})
    assert not form.is_valid()
    # KNOWN BUG: optional EmailField emits 'invalid' when absent
    assert dump_errors(form.errors) == {'email': [{'code': 'invalid'}]}


class TestEmailFieldNone:
  """Test EmailField with None value."""

  def test_email_none_required_true(self) -> None:
    """Test required EmailField with None."""

    class TestForm(Form):
      """Test form."""

      email = EmailField(required=True)

    form = TestForm({'email': None})
    assert not form.is_valid()
    # KNOWN BUG: EmailField emits both 'required' and 'invalid' for None
    assert dump_errors(form.errors) == {'email': [{'code': 'required'}, {'code': 'invalid'}]}

  def test_email_none_required_false(self) -> None:
    """Test optional EmailField with None."""

    class TestForm(Form):
      """Test form."""

      email = EmailField(required=False)

    form = TestForm({'email': None})
    assert not form.is_valid()
    # KNOWN BUG: optional EmailField emits 'invalid' when None
    assert dump_errors(form.errors) == {'email': [{'code': 'invalid'}]}


class TestEmailFieldValid:
  """Test EmailField with valid emails."""

  def test_email_valid_basic(self) -> None:
    """Test EmailField with valid email."""

    class TestForm(Form):
      """Test form."""

      email = EmailField(required=True)

    form = TestForm({'email': 'test@example.com'})
    assert form.is_valid()
    assert dump_errors(form.errors) == {}
    assert form.cleaned_data == {'email': 'test@example.com'}

  def test_email_valid_with_plus(self) -> None:
    """Test EmailField with plus addressing."""

    class TestForm(Form):
      """Test form."""

      email = EmailField(required=True)

    form = TestForm({'email': 'test+tag@example.com'})
    assert form.is_valid()
    assert dump_errors(form.errors) == {}

  def test_email_valid_with_dots(self) -> None:
    """Test EmailField with dots in local part."""

    class TestForm(Form):
      """Test form."""

      email = EmailField(required=True)

    form = TestForm({'email': 'first.last@example.com'})
    assert form.is_valid()
    assert dump_errors(form.errors) == {}


class TestEmailFieldInvalid:
  """Test EmailField with invalid emails."""

  def test_email_invalid_no_at(self) -> None:
    """Test EmailField without @ symbol."""

    class TestForm(Form):
      """Test form."""

      email = EmailField(required=True)

    form = TestForm({'email': 'notanemail.com'})
    assert not form.is_valid()
    assert dump_errors(form.errors) == {'email': [{'code': 'invalid'}]}

  def test_email_invalid_no_domain(self) -> None:
    """Test EmailField without domain."""

    class TestForm(Form):
      """Test form."""

      email = EmailField(required=True)

    form = TestForm({'email': 'test@'})
    assert not form.is_valid()
    assert dump_errors(form.errors) == {'email': [{'code': 'invalid'}]}

  def test_email_invalid_no_local(self) -> None:
    """Test EmailField without local part."""

    class TestForm(Form):
      """Test form."""

      email = EmailField(required=True)

    form = TestForm({'email': '@example.com'})
    assert not form.is_valid()
    assert dump_errors(form.errors) == {'email': [{'code': 'invalid'}]}

  def test_email_invalid_no_tld(self) -> None:
    """Test EmailField without TLD."""

    class TestForm(Form):
      """Test form."""

      email = EmailField(required=True)

    form = TestForm({'email': 'test@example'})
    assert not form.is_valid()
    assert dump_errors(form.errors) == {'email': [{'code': 'invalid'}]}

  def test_email_invalid_spaces(self) -> None:
    """Test EmailField with spaces."""

    class TestForm(Form):
      """Test form."""

      email = EmailField(required=True)

    form = TestForm({'email': 'test @example.com'})
    assert not form.is_valid()
    assert dump_errors(form.errors) == {'email': [{'code': 'invalid'}]}


class TestEmailFieldEmpty:
  """Test EmailField with empty string handling."""

  def test_email_empty_string_empty_false_required_false(self) -> None:
    """Test optional EmailField with empty string when empty=False."""

    class TestForm(Form):
      """Test form."""

      email = EmailField(required=False, empty=False)

    form = TestForm({'email': ''})
    assert not form.is_valid()
    # KNOWN BUG: EmailField emits 'required' for '' where CharField emits 'empty'
    assert dump_errors(form.errors) == {'email': [{'code': 'required'}]}

  def test_email_empty_string_empty_true_required_false(self) -> None:
    """Test optional EmailField with empty string when empty=True."""

    class TestForm(Form):
      """Test form."""

      email = EmailField(required=False, empty=True)

    form = TestForm({'email': ''})
    assert form.is_valid()
    assert dump_errors(form.errors) == {}

  def test_email_empty_string_empty_false_required_true(self) -> None:
    """Test required EmailField with empty string when empty=False."""

    class TestForm(Form):
      """Test form."""

      email = EmailField(required=True, empty=False)

    form = TestForm({'email': ''})
    assert not form.is_valid()
    # KNOWN BUG: EmailField emits 'required' for ''
    assert dump_errors(form.errors) == {'email': [{'code': 'required'}]}


class TestEmailFieldNonString:
  """Test EmailField with non-string values."""

  def test_email_non_string_int(self) -> None:
    """Test EmailField with integer value."""

    class TestForm(Form):
      """Test form."""

      email = EmailField(required=True)

    form = TestForm({'email': 123})
    assert not form.is_valid()
    assert dump_errors(form.errors) == {'email': [{'code': 'invalid'}]}

  def test_email_non_string_list(self) -> None:
    """Test EmailField with list value."""

    class TestForm(Form):
      """Test form."""

      email = EmailField(required=True)

    form = TestForm({'email': []})
    assert not form.is_valid()
    assert dump_errors(form.errors) == {'email': [{'code': 'invalid'}]}

  def test_email_non_string_dict(self) -> None:
    """Test EmailField with dict value."""

    class TestForm(Form):
      """Test form."""

      email = EmailField(required=True)

    form = TestForm({'email': {}})
    assert not form.is_valid()
    assert dump_errors(form.errors) == {'email': [{'code': 'invalid'}]}

  def test_email_non_string_optional(self) -> None:
    """Test optional EmailField with non-string value."""

    class TestForm(Form):
      """Test form."""

      email = EmailField(required=False)

    form = TestForm({'email': 123})
    assert not form.is_valid()
    assert dump_errors(form.errors) == {'email': [{'code': 'invalid'}]}
