"""Test base Field class."""

from typing import Any

import pytest

from layrz_forms import Form
from layrz_forms.fields import Field
from layrz_forms.types import ErrorType


class TestFieldRequiredSemantics:
  """Test base Field required semantics."""

  def test_required_field_with_none(self) -> None:
    """Test required field raises error when value is None."""
    field = Field(required=True)
    errors: ErrorType = {}
    field.validate(key='test_field', value=None, errors=errors)
    assert 'testField' in errors
    assert errors['testField'] == [{'code': 'required'}]

  def test_required_field_with_value(self) -> None:
    """Test required field accepts non-None value."""
    field = Field(required=True)
    errors: ErrorType = {}
    field.validate(key='test_field', value='value', errors=errors)
    assert 'testField' not in errors

  def test_optional_field_with_none(self) -> None:
    """Test optional field accepts None."""
    field = Field(required=False)
    errors: ErrorType = {}
    field.validate(key='test_field', value=None, errors=errors)
    assert 'testField' not in errors

  def test_optional_field_with_value(self) -> None:
    """Test optional field accepts value."""
    field = Field(required=False)
    errors: ErrorType = {}
    field.validate(key='test_field', value='value', errors=errors)
    assert 'testField' not in errors


class TestCustomFieldSubclass:
  """Test custom Field subclass with correct signature."""

  def test_custom_field_correct_signature(self) -> None:
    """Test custom field with correct validate signature works."""

    class CustomField(Field):
      """Custom field with correct signature."""

      def validate(self, key: str, value: Any, errors: ErrorType) -> None:
        """Validate method with correct parameters."""
        super().validate(key=key, value=value, errors=errors)
        if value == 'bad':
          self._append_error(key=key, errors=errors, to_add={'code': 'custom_error'})

    class TestForm(Form):
      """Test form with custom field."""

      custom = CustomField(required=True)

    form = TestForm({'custom': 'bad'})
    assert not form.is_valid()
    assert form.errors() == {'custom': [{'code': 'custom_error'}]}

  def test_custom_field_wrong_signature_missing_param(self) -> None:
    """Test custom field with wrong signature (missing param) raises RuntimeError."""

    class BadField(Field):
      """Custom field with wrong signature."""

      def validate(self, key: str) -> None:  # ty: ignore[invalid-method-override]
        """Validate method with missing parameter."""
        pass

    class TestForm(Form):
      """Test form with bad field."""

      bad = BadField(required=True)

    form = TestForm({'bad': 'value'})
    with pytest.raises(RuntimeError) as exc_info:
      form.is_valid()
    assert 'validate method has no the correct parameters' in str(exc_info.value)

  def test_custom_field_wrong_signature_extra_param(self) -> None:
    """Test custom field with wrong signature (extra param) raises RuntimeError."""

    class BadField(Field):
      """Custom field with wrong signature."""

      def validate(
        self,
        key: str,
        value: Any,
        errors: ErrorType,
        extra: str,
      ) -> None:  # ty: ignore[invalid-method-override]
        """Validate method with extra parameter."""
        pass

    class TestForm(Form):
      """Test form with bad field."""

      bad = BadField(required=True)

    form = TestForm({'bad': 'value'})
    with pytest.raises(RuntimeError) as exc_info:
      form.is_valid()
    assert 'validate method has no the correct parameters' in str(exc_info.value)

  def test_custom_field_wrong_param_names(self) -> None:
    """Test custom field with wrong parameter names raises RuntimeError."""

    class BadField(Field):
      """Custom field with wrong parameter names."""

      def validate(
        self,
        name: str,
        val: Any,
        err: ErrorType,
      ) -> None:  # ty: ignore[invalid-method-override]
        """Validate method with wrong parameter names."""
        pass

    class TestForm(Form):
      """Test form with bad field."""

      bad = BadField(required=True)

    form = TestForm({'bad': 'value'})
    with pytest.raises(RuntimeError) as exc_info:
      form.is_valid()
    assert 'validate method has no the correct parameters' in str(exc_info.value)
