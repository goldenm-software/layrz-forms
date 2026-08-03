"""Test IdField validation."""

import pytest

from layrz_forms import Form, IdField
from tests.helpers import dump_errors


class TestIdFieldAbsent:
  """Test IdField when field is absent."""

  def test_id_absent_required_true(self) -> None:
    """Test required IdField absent from input."""

    class TestForm(Form):
      """Test form."""

      item_id = IdField(required=True)

    form = TestForm({})
    assert not form.is_valid()
    assert dump_errors(form.errors) == {'itemId': [{'code': 'required'}]}

  def test_id_absent_required_false(self) -> None:
    """Test optional IdField absent from input."""

    class TestForm(Form):
      """Test form."""

      item_id = IdField(required=False)

    form = TestForm({})
    assert form.is_valid()
    assert dump_errors(form.errors) == {}


class TestIdFieldNone:
  """Test IdField with None value."""

  def test_id_none_required_true(self) -> None:
    """Test required IdField with None."""

    class TestForm(Form):
      """Test form."""

      item_id = IdField(required=True)

    form = TestForm({'item_id': None})
    assert not form.is_valid()
    assert dump_errors(form.errors) == {'itemId': [{'code': 'required'}]}

  def test_id_none_required_false(self) -> None:
    """Test optional IdField with None."""

    class TestForm(Form):
      """Test form."""

      item_id = IdField(required=False)

    form = TestForm({'item_id': None})
    assert form.is_valid()
    assert dump_errors(form.errors) == {}


class TestIdFieldValid:
  """Test IdField with valid values."""

  def test_id_valid_positive_int(self) -> None:
    """Test IdField with positive integer."""

    class TestForm(Form):
      """Test form."""

      item_id = IdField(required=True)

    form = TestForm({'item_id': 42})
    assert form.is_valid()
    assert dump_errors(form.errors) == {}
    assert form.cleaned_data == {'item_id': 42}

  def test_id_valid_positive_str(self) -> None:
    """Test IdField with positive numeric string."""

    class TestForm(Form):
      """Test form."""

      item_id = IdField(required=True)

    form = TestForm({'item_id': '123'})
    assert form.is_valid()
    assert dump_errors(form.errors) == {}
    assert form.cleaned_data == {'item_id': '123'}

  def test_id_valid_large_int(self) -> None:
    """Test IdField with large integer."""

    class TestForm(Form):
      """Test form."""

      item_id = IdField(required=True)

    form = TestForm({'item_id': 999999999})
    assert form.is_valid()
    assert dump_errors(form.errors) == {}


class TestIdFieldInvalid:
  """Test IdField with invalid values."""

  def test_id_invalid_zero(self) -> None:
    """Test IdField with zero (not > 0)."""

    class TestForm(Form):
      """Test form."""

      item_id = IdField(required=True)

    form = TestForm({'item_id': 0})
    assert not form.is_valid()
    assert dump_errors(form.errors) == {'itemId': [{'code': 'invalid'}]}

  def test_id_invalid_negative(self) -> None:
    """Test IdField with negative integer."""

    class TestForm(Form):
      """Test form."""

      item_id = IdField(required=True)

    form = TestForm({'item_id': -5})
    assert not form.is_valid()
    assert dump_errors(form.errors) == {'itemId': [{'code': 'invalid'}]}

  def test_id_invalid_negative_str(self) -> None:
    """Test IdField with negative numeric string."""

    class TestForm(Form):
      """Test form."""

      item_id = IdField(required=True)

    form = TestForm({'item_id': '-42'})
    assert not form.is_valid()
    assert dump_errors(form.errors) == {'itemId': [{'code': 'invalid'}]}

  def test_id_invalid_non_numeric_string_required_true(self) -> None:
    """Test required IdField with non-numeric string emits invalid error."""

    class TestForm(Form):
      """Test form."""

      item_id = IdField(required=True)

    form = TestForm({'item_id': 'abc'})
    assert not form.is_valid()
    assert dump_errors(form.errors) == {'itemId': [{'code': 'invalid'}]}

  def test_id_invalid_non_numeric_string_required_false(self) -> None:
    """Test optional IdField with non-numeric string emits invalid error."""

    class TestForm(Form):
      """Test form."""

      item_id = IdField(required=False)

    form = TestForm({'item_id': 'abc'})
    assert not form.is_valid()
    assert dump_errors(form.errors) == {'itemId': [{'code': 'invalid'}]}

  def test_id_invalid_float_required_true(self) -> None:
    """Test required IdField with float value."""

    class TestForm(Form):
      """Test form."""

      item_id = IdField(required=True)

    form = TestForm({'item_id': 3.14})
    assert not form.is_valid()
    assert dump_errors(form.errors) == {'itemId': [{'code': 'invalid'}]}

  def test_id_invalid_float_required_false(self) -> None:
    """Test optional IdField with float value."""

    class TestForm(Form):
      """Test form."""

      item_id = IdField(required=False)

    form = TestForm({'item_id': 3.14})
    assert not form.is_valid()
    assert dump_errors(form.errors) == {'itemId': [{'code': 'invalid'}]}

  def test_id_invalid_list_required_true(self) -> None:
    """Test required IdField with list value."""

    class TestForm(Form):
      """Test form."""

      item_id = IdField(required=True)

    form = TestForm({'item_id': []})
    assert not form.is_valid()
    assert dump_errors(form.errors) == {'itemId': [{'code': 'invalid'}]}

  def test_id_invalid_list_required_false(self) -> None:
    """Test optional IdField with list value emits invalid error."""

    class TestForm(Form):
      """Test form."""

      item_id = IdField(required=False)

    form = TestForm({'item_id': []})
    assert not form.is_valid()
    assert dump_errors(form.errors) == {'itemId': [{'code': 'invalid'}]}

  def test_id_invalid_dict_required_true(self) -> None:
    """Test required IdField with dict value."""

    class TestForm(Form):
      """Test form."""

      item_id = IdField(required=True)

    form = TestForm({'item_id': {}})
    assert not form.is_valid()
    assert dump_errors(form.errors) == {'itemId': [{'code': 'invalid'}]}

  def test_id_invalid_dict_required_false(self) -> None:
    """Test optional IdField with dict value emits invalid error."""

    class TestForm(Form):
      """Test form."""

      item_id = IdField(required=False)

    form = TestForm({'item_id': {'a': 1}})
    assert not form.is_valid()
    assert dump_errors(form.errors) == {'itemId': [{'code': 'invalid'}]}
