"""Test UuidField validation."""

import uuid

from layrz_forms import Form, UuidField


class TestUuidFieldAbsent:
  """Test UuidField when field is absent."""

  def test_uuid_absent_required_true(self) -> None:
    """Test required UuidField absent from input."""

    class TestForm(Form):
      """Test form."""

      uid = UuidField(required=True)

    form = TestForm({})
    assert not form.is_valid()
    # KNOWN BUG: UuidField emits both 'required' and 'invalid' when absent
    assert form.errors() == {'uid': [{'code': 'required'}, {'code': 'invalid'}]}

  def test_uuid_absent_required_false(self) -> None:
    """Test optional UuidField absent from input."""

    class TestForm(Form):
      """Test form."""

      uid = UuidField(required=False)

    form = TestForm({})
    assert form.is_valid()
    assert form.errors() == {}


class TestUuidFieldNone:
  """Test UuidField with None value."""

  def test_uuid_none_required_true(self) -> None:
    """Test required UuidField with None."""

    class TestForm(Form):
      """Test form."""

      uid = UuidField(required=True)

    form = TestForm({'uid': None})
    assert not form.is_valid()
    # KNOWN BUG: UuidField emits both 'required' and 'invalid' for None
    assert form.errors() == {'uid': [{'code': 'required'}, {'code': 'invalid'}]}

  def test_uuid_none_required_false(self) -> None:
    """Test optional UuidField with None."""

    class TestForm(Form):
      """Test form."""

      uid = UuidField(required=False)

    form = TestForm({'uid': None})
    assert form.is_valid()
    assert form.errors() == {}


class TestUuidFieldValidString:
  """Test UuidField with valid UUID strings."""

  def test_uuid_valid_string_v4(self) -> None:
    """Test UuidField with valid UUID v4 string."""
    valid_uuid = '550e8400-e29b-41d4-a716-446655440000'

    class TestForm(Form):
      """Test form."""

      uid = UuidField(required=True)

    form = TestForm({'uid': valid_uuid})
    assert form.is_valid()
    assert form.errors() == {}
    assert form.cleaned_data == {'uid': valid_uuid}

  def test_uuid_valid_string_uppercase(self) -> None:
    """Test UuidField with uppercase UUID string."""
    valid_uuid = '550E8400-E29B-41D4-A716-446655440000'

    class TestForm(Form):
      """Test form."""

      uid = UuidField(required=True)

    form = TestForm({'uid': valid_uuid})
    assert form.is_valid()
    assert form.errors() == {}

  def test_uuid_valid_string_no_hyphens(self) -> None:
    """Test UuidField with UUID string without hyphens."""
    valid_uuid = '550e8400e29b41d4a716446655440000'

    class TestForm(Form):
      """Test form."""

      uid = UuidField(required=True)

    form = TestForm({'uid': valid_uuid})
    assert form.is_valid()
    assert form.errors() == {}


class TestUuidFieldValidInstance:
  """Test UuidField with uuid.UUID instances."""

  def test_uuid_valid_uuid_instance(self) -> None:
    """Test UuidField with uuid.UUID instance."""
    valid_uuid = uuid.UUID('550e8400-e29b-41d4-a716-446655440000')

    class TestForm(Form):
      """Test form."""

      uid = UuidField(required=True)

    form = TestForm({'uid': valid_uuid})
    assert form.is_valid()
    assert form.errors() == {}
    assert form.cleaned_data == {'uid': valid_uuid}

  def test_uuid_valid_uuid_instance_optional(self) -> None:
    """Test optional UuidField with uuid.UUID instance."""
    valid_uuid = uuid.UUID('550e8400-e29b-41d4-a716-446655440000')

    class TestForm(Form):
      """Test form."""

      uid = UuidField(required=False)

    form = TestForm({'uid': valid_uuid})
    assert form.is_valid()
    assert form.errors() == {}


class TestUuidFieldInvalidString:
  """Test UuidField with invalid UUID strings."""

  def test_uuid_invalid_malformed_required_true(self) -> None:
    """Test required UuidField with malformed UUID string."""

    class TestForm(Form):
      """Test form."""

      uid = UuidField(required=True)

    form = TestForm({'uid': 'not-a-uuid'})
    assert not form.is_valid()
    assert form.errors() == {'uid': [{'code': 'invalid'}]}

  def test_uuid_invalid_malformed_required_false(self) -> None:
    """Test optional UuidField with malformed UUID string."""

    class TestForm(Form):
      """Test form."""

      uid = UuidField(required=False)

    form = TestForm({'uid': 'not-a-uuid'})
    assert not form.is_valid()
    assert form.errors() == {'uid': [{'code': 'invalid'}]}

  def test_uuid_invalid_partial_required_true(self) -> None:
    """Test required UuidField with partial UUID."""

    class TestForm(Form):
      """Test form."""

      uid = UuidField(required=True)

    form = TestForm({'uid': '550e8400-e29b-41d4'})
    assert not form.is_valid()
    assert form.errors() == {'uid': [{'code': 'invalid'}]}

  def test_uuid_invalid_empty_string_required_true(self) -> None:
    """Test required UuidField with empty string."""

    class TestForm(Form):
      """Test form."""

      uid = UuidField(required=True)

    form = TestForm({'uid': ''})
    assert not form.is_valid()
    assert form.errors() == {'uid': [{'code': 'invalid'}]}

  def test_uuid_invalid_empty_string_required_false(self) -> None:
    """Test optional UuidField with empty string."""

    class TestForm(Form):
      """Test form."""

      uid = UuidField(required=False)

    form = TestForm({'uid': ''})
    assert not form.is_valid()
    assert form.errors() == {'uid': [{'code': 'invalid'}]}


class TestUuidFieldNonStringNonUuid:
  """Test UuidField with non-string, non-UUID values."""

  def test_uuid_invalid_int_required_true(self) -> None:
    """Test required UuidField with integer."""

    class TestForm(Form):
      """Test form."""

      uid = UuidField(required=True)

    form = TestForm({'uid': 123})
    assert not form.is_valid()
    assert form.errors() == {'uid': [{'code': 'invalid'}]}

  def test_uuid_invalid_int_required_false(self) -> None:
    """Test optional UuidField with integer."""

    class TestForm(Form):
      """Test form."""

      uid = UuidField(required=False)

    form = TestForm({'uid': 123})
    assert not form.is_valid()
    assert form.errors() == {'uid': [{'code': 'invalid'}]}

  def test_uuid_invalid_list_required_true(self) -> None:
    """Test required UuidField with list."""

    class TestForm(Form):
      """Test form."""

      uid = UuidField(required=True)

    form = TestForm({'uid': []})
    assert not form.is_valid()
    assert form.errors() == {'uid': [{'code': 'invalid'}]}

  def test_uuid_invalid_dict_required_true(self) -> None:
    """Test required UuidField with dict."""

    class TestForm(Form):
      """Test form."""

      uid = UuidField(required=True)

    form = TestForm({'uid': {}})
    assert not form.is_valid()
    assert form.errors() == {'uid': [{'code': 'invalid'}]}

  def test_uuid_invalid_bool_required_true(self) -> None:
    """Test required UuidField with boolean."""

    class TestForm(Form):
      """Test form."""

      uid = UuidField(required=True)

    form = TestForm({'uid': True})
    assert not form.is_valid()
    assert form.errors() == {'uid': [{'code': 'invalid'}]}
