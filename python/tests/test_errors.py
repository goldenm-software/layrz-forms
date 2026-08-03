"""Test LayrzError model and error handling."""

import pytest
from pydantic import ValidationError

from layrz_forms import CharField, Form, LayrzError
from layrz_forms.types import ErrorsType


class TestLayrzErrorDefaults:
  """Test LayrzError field defaults."""

  def test_error_code_required(self) -> None:
    """Test LayrzError requires code field."""
    with pytest.raises(ValidationError):
      LayrzError()  # ty: ignore[missing-argument]

  def test_error_code_only(self) -> None:
    """Test LayrzError with only code field."""
    error = LayrzError(code='required')
    assert error.code == 'required'
    assert error.expected is None
    assert error.received is None
    assert error.extra is None

  def test_error_all_fields(self) -> None:
    """Test LayrzError with all fields."""
    error = LayrzError(
      code='minLength',
      expected=5,
      received=3,
      extra={'note': 'too short'},
    )
    assert error.code == 'minLength'
    assert error.expected == 5
    assert error.received == 3
    assert error.extra == {'note': 'too short'}


class TestLayrzErrorModelDump:
  """Test LayrzError.model_dump() serialization."""

  def test_model_dump_excludes_none_by_default(self) -> None:
    """Test model_dump excludes None values by default."""
    error = LayrzError(code='required')
    dumped = error.model_dump()
    assert dumped == {'code': 'required'}
    assert 'expected' not in dumped
    assert 'received' not in dumped
    assert 'extra' not in dumped

  def test_model_dump_with_expected_received(self) -> None:
    """Test model_dump includes expected and received when set."""
    error = LayrzError(code='minLength', expected=5, received=3)
    dumped = error.model_dump()
    assert dumped == {'code': 'minLength', 'expected': 5, 'received': 3}
    assert 'extra' not in dumped

  def test_model_dump_with_extra(self) -> None:
    """Test model_dump includes extra when set."""
    error = LayrzError(code='custom', extra={'key': 'value'})
    dumped = error.model_dump()
    assert dumped == {'code': 'custom', 'extra': {'key': 'value'}}

  def test_model_dump_exclude_none_false(self) -> None:
    """Test model_dump with exclude_none=False includes None fields."""
    error = LayrzError(code='required')
    dumped = error.model_dump(exclude_none=False)
    assert dumped == {
      'code': 'required',
      'expected': None,
      'received': None,
      'extra': None,
    }

  def test_model_dump_partial_values(self) -> None:
    """Test model_dump with some None values."""
    error = LayrzError(code='test', expected=5, received=None)
    dumped = error.model_dump()
    assert dumped == {'code': 'test', 'expected': 5}
    assert 'received' not in dumped
    assert 'extra' not in dumped


class TestLayrzErrorModelDumpJson:
  """Test LayrzError.model_dump_json() serialization."""

  def test_model_dump_json_excludes_none_by_default(self) -> None:
    """Test model_dump_json excludes None by default."""
    error = LayrzError(code='required')
    json_str = error.model_dump_json()
    assert '{"code":"required"}' in json_str or '{"code": "required"}' in json_str
    assert 'expected' not in json_str

  def test_model_dump_json_exclude_none_false(self) -> None:
    """Test model_dump_json with exclude_none=False."""
    error = LayrzError(code='required')
    json_str = error.model_dump_json(exclude_none=False)
    assert 'expected' in json_str
    assert 'null' in json_str


class TestLayrzErrorExtraForbid:
  """Test LayrzError rejects unknown fields."""

  def test_extra_forbid_rejects_unknown_field(self) -> None:
    """Test LayrzError with extra='forbid' rejects unknown kwargs."""
    with pytest.raises(ValidationError) as exc_info:
      LayrzError(code='test', unknown_field='value')  # ty: ignore[unknown-argument]
    assert 'unknown_field' in str(exc_info.value)


class TestAddErrorsRouting:
  """Test Form.add_errors routes expected/received correctly."""

  def test_add_errors_lifts_expected_received(self) -> None:
    """Test add_errors lifts expected/received from extra_args."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=False)

      def clean_check(self) -> None:
        """Add error in clean."""
        self.add_errors(
          key='name',
          code='weak_password',
          extra_args={'expected': 8, 'received': 4},
        )

    form = TestForm({})
    form.is_valid()
    errors = form.errors
    assert 'name' in errors
    error = errors['name'][0]
    assert error.code == 'weak_password'
    assert error.expected == 8
    assert error.received == 4
    assert error.extra is None

  def test_add_errors_extra_dict_separate(self) -> None:
    """Test add_errors keeps other keys in extra."""

    class TestForm(Form):
      """Test form."""

      password = CharField(required=False)

      def clean_check(self) -> None:
        """Add error in clean."""
        self.add_errors(
          key='password',
          code='weak',
          extra_args={'expected': 8, 'min_length': 8, 'note': 'too short'},
        )

    form = TestForm({})
    form.is_valid()
    errors = form.errors
    error = errors['password'][0]
    assert error.code == 'weak'
    assert error.expected == 8
    assert error.received is None
    assert error.extra == {'min_length': 8, 'note': 'too short'}

  def test_add_errors_no_extra_args(self) -> None:
    """Test add_errors with no extra_args."""

    class TestForm(Form):
      """Test form."""

      field = CharField(required=False)

      def clean_check(self) -> None:
        """Add error in clean."""
        self.add_errors(key='field', code='error')

    form = TestForm({})
    form.is_valid()
    errors = form.errors
    error = errors['field'][0]
    assert error.code == 'error'
    assert error.expected is None
    assert error.received is None
    assert error.extra is None

  def test_add_errors_only_extra_keys(self) -> None:
    """Test add_errors with only extra keys (no expected/received)."""

    class TestForm(Form):
      """Test form."""

      field = CharField(required=False)

      def clean_check(self) -> None:
        """Add error in clean."""
        self.add_errors(
          key='field',
          code='custom',
          extra_args={'min_length': 5, 'max_length': 10},
        )

    form = TestForm({})
    form.is_valid()
    errors = form.errors
    error = errors['field'][0]
    assert error.code == 'custom'
    assert error.expected is None
    assert error.received is None
    assert error.extra == {'min_length': 5, 'max_length': 10}


class TestAddErrorsValidation:
  """Test Form.add_errors validation."""

  def test_add_errors_empty_key_raises(self) -> None:
    """Test add_errors with empty key raises RuntimeError."""

    class TestForm(Form):
      """Test form."""

      field = CharField(required=False)

    form = TestForm({})
    with pytest.raises(RuntimeError) as exc_info:
      form.add_errors(key='', code='error')
    assert 'key and code are required' in str(exc_info.value)

  def test_add_errors_empty_code_raises(self) -> None:
    """Test add_errors with empty code raises RuntimeError."""

    class TestForm(Form):
      """Test form."""

      field = CharField(required=False)

    form = TestForm({})
    with pytest.raises(RuntimeError) as exc_info:
      form.add_errors(key='field', code='')
    assert 'key and code are required' in str(exc_info.value)


class TestNestedFormErrorPreservation:
  """Test that nested form errors preserve expected/received."""

  def test_nested_form_error_expected_received_preserved(self) -> None:
    """Test nested form field error preserves expected/received values."""

    class InnerForm(Form):
      """Inner form."""

      name = CharField(required=True, min_length=5)

    class OuterForm(Form):
      """Outer form."""

      child = InnerForm()

    form = OuterForm({'child': {'name': 'ab'}})
    form.is_valid()
    errors = form.errors

    # Nested field error should preserve expected/received
    assert 'child.name' in errors
    error = errors['child.name'][0]
    assert error.code == 'minLength'
    assert error.expected == 5
    assert error.received == 2


class TestFormErrorsPropertyCallable:
  """Test that form.errors property raises TypeError if called."""

  def test_errors_not_callable(self) -> None:
    """Test form.errors returns a dict, which cannot be called as method."""

    class TestForm(Form):
      """Test form."""

      field = CharField(required=False)

    form = TestForm({})
    errors = form.errors
    # errors is a dict, not callable
    with pytest.raises(TypeError):
      errors()  # ty: ignore[call-non-callable]
