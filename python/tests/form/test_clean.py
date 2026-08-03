"""Test Form clean methods."""

import pytest

from layrz_forms import CharField, Form
from tests.helpers import dump_errors


class TestSingleCleanFunction:
  """Test a single clean function."""

  def test_clean_function_adds_error(self) -> None:
    """Test clean function can add an error."""

    class TestForm(Form):
      """Test form with clean function."""

      name = CharField(required=False)

      def clean_validate(self) -> None:
        """Custom validation."""
        if self._obj.get('name') == 'banned':
          self.add_errors(key='name', code='banned_word')

    form = TestForm({'name': 'banned'})
    assert form.is_valid() is False
    assert dump_errors(form.errors) == {'name': [{'code': 'banned_word'}]}

  def test_clean_function_field_valid(self) -> None:
    """Test clean function when validation passes."""

    class TestForm(Form):
      """Test form with clean function."""

      name = CharField(required=False)

      def clean_validate(self) -> None:
        """Custom validation."""
        if self._obj.get('name') == 'banned':
          self.add_errors(key='name', code='banned_word')

    form = TestForm({'name': 'allowed'})
    assert form.is_valid() is True
    assert dump_errors(form.errors) == {}


class TestMultipleCleanFunctions:
  """Test multiple clean functions."""

  def test_two_clean_functions_both_add_errors(self) -> None:
    """Test two clean functions adding errors to same key."""

    class TestForm(Form):
      """Test form with multiple clean functions."""

      name = CharField(required=False)

      def clean_apple(self) -> None:
        """First clean function."""
        self.add_errors(key='name', code='error_a')

      def clean_zebra(self) -> None:
        """Second clean function."""
        self.add_errors(key='name', code='error_z')

    form = TestForm({'name': 'test'})
    assert form.is_valid() is False
    errors = form.errors
    assert 'name' in errors
    assert len(errors['name']) == 2
    # Executed alphabetically
    assert errors['name'][0].code == 'error_a'
    assert errors['name'][1].code == 'error_z'


class TestCleanFunctionExtraArgs:
  """Test clean function with extra_args."""

  def test_clean_with_extra_args(self) -> None:
    """Test add_errors with extra_args."""

    class TestForm(Form):
      """Test form."""

      password = CharField(required=False)

      def clean_strength(self) -> None:
        """Validate password strength."""
        if len(self._obj.get('password', '')) < 8:
          self.add_errors(
            key='password',
            code='weak_password',
            extra_args={'min_length': 8, 'received_length': len(self._obj.get('password', ''))},
          )

    form = TestForm({'password': 'short'})
    assert form.is_valid() is False
    errors = form.errors
    assert dump_errors(errors)['password'][0] == {
      'code': 'weak_password',
      'extra': {'min_length': 8, 'received_length': 5},
    }


class TestCleanEmptyKeyOrCode:
  """Test add_errors with empty key or code."""

  def test_add_errors_empty_key_raises_runtime_error(self) -> None:
    """Test add_errors with empty key raises RuntimeError."""

    class TestForm(Form):
      """Test form."""

      def clean_validate(self) -> None:
        """Clean function."""
        self.add_errors(key='', code='some_error')

    form = TestForm({})
    with pytest.raises(RuntimeError) as exc_info:
      form.is_valid()
    assert 'key and code are required' in str(exc_info.value)

  def test_add_errors_empty_code_raises_runtime_error(self) -> None:
    """Test add_errors with empty code raises RuntimeError."""

    class TestForm(Form):
      """Test form."""

      def clean_validate(self) -> None:
        """Clean function."""
        self.add_errors(key='field', code='')

    form = TestForm({})
    with pytest.raises(RuntimeError) as exc_info:
      form.is_valid()
    assert 'key and code are required' in str(exc_info.value)

  def test_add_errors_both_empty_raises_runtime_error(self) -> None:
    """Test add_errors with both key and code empty raises RuntimeError."""

    class TestForm(Form):
      """Test form."""

      def clean_validate(self) -> None:
        """Clean function."""
        self.add_errors(key='', code='')

    form = TestForm({})
    with pytest.raises(RuntimeError) as exc_info:
      form.is_valid()
    assert 'key and code are required' in str(exc_info.value)


class TestCleanMultipleDifferentKeys:
  """Test clean function adding errors to different keys."""

  def test_clean_multiple_keys(self) -> None:
    """Test clean function adding errors to multiple keys."""

    class TestForm(Form):
      """Test form."""

      password = CharField(required=False)
      confirm_password = CharField(required=False)

      def clean_match(self) -> None:
        """Validate passwords match."""
        if self._obj.get('password') != self._obj.get('confirm_password'):
          self.add_errors(key='confirm_password', code='passwords_do_not_match')

    form = TestForm({'password': 'secret123', 'confirm_password': 'wrong'})
    assert form.is_valid() is False
    errors = form.errors
    assert 'confirmPassword' in errors
    assert errors['confirmPassword'][0].code == 'passwords_do_not_match'


class TestCleanWithFieldErrors:
  """Test clean function runs after field validation."""

  def test_clean_runs_after_field_validation(self) -> None:
    """Test clean runs after field validation."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True, min_length=5)

      def clean_unique(self) -> None:
        """Check uniqueness."""
        if self._obj.get('name') == 'admin':
          self.add_errors(key='name', code='reserved')

    # Field validation fails first
    form = TestForm({'name': 'bob'})
    assert form.is_valid() is False
    errors = form.errors
    assert 'name' in errors
    # Only field error, clean didn't run
    assert len(errors['name']) == 1
    assert errors['name'][0].code == 'minLength'

  def test_clean_runs_when_field_valid(self) -> None:
    """Test clean runs when field is valid."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True, min_length=5)

      def clean_unique(self) -> None:
        """Check uniqueness."""
        if self._obj.get('name') == 'admin':
          self.add_errors(key='name', code='reserved')

    # Field is valid, clean runs
    form = TestForm({'name': 'admin'})
    assert form.is_valid() is False
    errors = form.errors
    assert 'name' in errors
    assert errors['name'][0].code == 'reserved'
