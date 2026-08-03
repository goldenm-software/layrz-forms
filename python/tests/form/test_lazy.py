"""Test lazy validation and property behavior of form.errors."""

import pytest

from layrz_forms import CharField, Form


class TestLazyValidation:
  """Test that .errors property performs lazy validation."""

  def test_errors_validates_without_explicit_is_valid(self) -> None:
    """Test .errors validates form without explicit is_valid() call."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True, min_length=5)

    form = TestForm({'name': 'ab'})
    # Access errors without calling is_valid()
    errors = form.errors
    assert 'name' in errors
    assert errors['name'][0].code == 'minLength'

  def test_errors_validates_fresh_form(self) -> None:
    """Test fresh form with no validation runs validation on .errors access."""

    class TestForm(Form):
      """Test form."""

      field = CharField(required=False)

    form = TestForm({})
    assert form._validated is False
    _ = form.errors
    # After accessing errors, _validated should be True
    assert form._validated is True


class TestErrorsNoDuplication:
  """Test that reading .errors twice doesn't re-validate or duplicate errors."""

  def test_errors_read_twice_no_duplication(self) -> None:
    """Test reading .errors twice doesn't accumulate errors."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True)

    form = TestForm({})
    errors1 = form.errors
    errors2 = form.errors
    assert 'name' in errors1
    assert len(errors1['name']) == 1
    assert len(errors2['name']) == 1
    assert errors1 is errors2  # Same object

  def test_errors_clean_function_runs_once(self) -> None:
    """Test clean function runs only once, not on each .errors read."""
    clean_call_count = 0

    class TestForm(Form):
      """Test form with counter."""

      field = CharField(required=False)

      def clean_track(self) -> None:
        """Track clean function calls."""
        nonlocal clean_call_count
        clean_call_count += 1

    form = TestForm({})
    _ = form.errors
    count_after_first = clean_call_count
    _ = form.errors
    count_after_second = clean_call_count
    # Clean function should run exactly once
    assert count_after_first == 1
    assert count_after_second == 1


class TestErrorsInvalidationOnObjChange:
  """Test that changing form.obj invalidates validation."""

  def test_obj_change_invalidates_errors(self) -> None:
    """Test setting form.obj invalidates _validated flag."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True)

    form = TestForm({})
    _ = form.errors  # Validate with empty dict
    assert form._validated is True

    # Change obj
    form.obj = {'name': 'John'}
    assert form._validated is False

  def test_obj_change_reflects_new_data(self) -> None:
    """Test changing form.obj then accessing .errors validates new data."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True)

    form = TestForm({})
    errors1 = form.errors
    assert 'name' in errors1

    form.obj = {'name': 'John'}
    errors2 = form.errors
    assert 'name' not in errors2


class TestCalculateMembersInvalidates:
  """Test that calculate_members() invalidates validation."""

  def test_calculate_members_clears_validated(self) -> None:
    """Test calculate_members() clears _validated flag."""

    class TestForm(Form):
      """Test form."""

      field = CharField(required=False)

    form = TestForm({})
    _ = form.errors
    assert form._validated is True

    form.calculate_members()
    assert form._validated is False


class TestAsyncCleanWithProperty:
  """Test that async clean functions require ais_valid()."""

  def test_errors_with_async_clean_raises_error(self) -> None:
    """Test .errors raises RuntimeError if form has async clean and hasn't been awaited."""

    class TestForm(Form):
      """Test form with async clean."""

      field = CharField(required=False)

      async def clean_validate(self) -> None:
        """Async clean function."""
        pass

    form = TestForm({})
    with pytest.raises(RuntimeError) as exc_info:
      _ = form.errors
    error_msg = str(exc_info.value)
    assert 'ais_valid' in error_msg
    assert 'async clean function' in error_msg.lower()

  @pytest.mark.asyncio
  async def test_errors_after_ais_valid_works(self) -> None:
    """Test .errors works after awaiting ais_valid()."""

    class TestForm(Form):
      """Test form with async clean."""

      name = CharField(required=False)

      async def clean_check(self) -> None:
        """Async clean function."""
        self.add_errors(key='name', code='async_error')

    form = TestForm({})
    await form.ais_valid()
    # Now accessing errors should work and have the async error
    errors = form.errors
    assert 'name' in errors
    assert errors['name'][0].code == 'async_error'


class TestExplicitIsValidThenErrors:
  """Test explicit is_valid() followed by .errors access."""

  def test_is_valid_then_errors_no_revalidation(self) -> None:
    """Test after is_valid(), accessing .errors doesn't re-validate."""
    clean_count = 0

    class TestForm(Form):
      """Test form."""

      field = CharField(required=False)

      def clean_track(self) -> None:
        """Track calls."""
        nonlocal clean_count
        clean_count += 1

    form = TestForm({})
    form.is_valid()
    assert clean_count == 1
    # Access errors multiple times
    _ = form.errors
    _ = form.errors
    # Still just 1 call
    assert clean_count == 1


class TestFreshFormState:
  """Test fresh form validation state."""

  def test_fresh_form_not_validated(self) -> None:
    """Test newly created form has _validated=False."""

    class TestForm(Form):
      """Test form."""

      field = CharField(required=False)

    form = TestForm({})
    # Member discovery doesn't validate
    assert form._validated is False

  def test_fresh_form_empty_errors_before_access(self) -> None:
    """Test _errors dict is empty before first validation."""

    class TestForm(Form):
      """Test form."""

      field = CharField(required=False)

    form = TestForm({})
    # _errors should be empty before access
    assert form._errors == {}
