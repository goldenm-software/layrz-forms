"""Test async Form validation."""

import pytest

from layrz_forms import CharField, Form
from tests.helpers import dump_errors


class TestAsyncCleanFunction:
  """Test async clean functions."""

  @pytest.mark.asyncio
  async def test_async_clean_function(self) -> None:
    """Test async clean function via ais_valid."""

    class TestForm(Form):
      """Test form with async clean."""

      name = CharField(required=False)

      async def clean_validate(self) -> None:
        """Async clean function."""
        if self._obj.get('name') == 'banned':
          self.add_errors(key='name', code='banned_word')

    form = TestForm({'name': 'banned'})
    result = await form.ais_valid()
    assert result is False
    assert dump_errors(form.errors) == {'name': [{'code': 'banned_word'}]}

  @pytest.mark.asyncio
  async def test_async_clean_function_valid(self) -> None:
    """Test async clean function passes."""

    class TestForm(Form):
      """Test form with async clean."""

      name = CharField(required=False)

      async def clean_validate(self) -> None:
        """Async clean function."""
        if self._obj.get('name') == 'banned':
          self.add_errors(key='name', code='banned_word')

    form = TestForm({'name': 'allowed'})
    result = await form.ais_valid()
    assert result is True
    assert dump_errors(form.errors) == {}


class TestSyncCleanViaAsync:
  """Test sync clean function via is_valid_async."""

  @pytest.mark.asyncio
  async def test_sync_clean_via_ais_valid(self) -> None:
    """Test sync clean function works via ais_valid."""

    class TestForm(Form):
      """Test form with sync clean."""

      name = CharField(required=False)

      def clean_validate(self) -> None:
        """Sync clean function."""
        if self._obj.get('name') == 'banned':
          self.add_errors(key='name', code='banned_word')

    form = TestForm({'name': 'banned'})
    result = await form.ais_valid()
    assert result is False
    assert dump_errors(form.errors) == {'name': [{'code': 'banned_word'}]}


class TestAsyncCleanViaSyncRaises:
  """Test async clean function via sync is_valid raises error."""

  def test_async_clean_via_sync_is_valid_raises(self) -> None:
    """Test calling sync is_valid with async clean raises RuntimeError."""

    class TestForm(Form):
      """Test form with async clean."""

      field = CharField(required=False)

      async def clean_validate(self) -> None:
        """Async clean function."""
        pass

    form = TestForm({'field': 'value'})
    with pytest.raises(RuntimeError) as exc_info:
      form.is_valid()
    assert 'Cannot call async clean function in sync context' in str(exc_info.value)
    assert 'please use ais_valid method' in str(exc_info.value)


class TestMixedSyncAsyncClean:
  """Test mix of sync and async clean functions."""

  @pytest.mark.asyncio
  async def test_mixed_sync_async_clean(self) -> None:
    """Test form with both sync and async clean functions."""

    class TestForm(Form):
      """Test form with mixed clean."""

      name = CharField(required=False)

      def clean_apple(self) -> None:
        """Sync clean function."""
        self.add_errors(key='name', code='sync_error')

      async def clean_zebra(self) -> None:
        """Async clean function."""
        self.add_errors(key='name', code='async_error')

    form = TestForm({'name': 'test'})
    result = await form.ais_valid()
    assert result is False
    errors = form.errors
    assert len(errors['name']) == 2
    # Alphabetical order
    assert errors['name'][0].code == 'sync_error'
    assert errors['name'][1].code == 'async_error'
