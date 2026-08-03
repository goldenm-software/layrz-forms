"""Test Form member discovery."""

import pytest

from layrz_forms import CharField, Form, NumberField


class TestFieldDiscovery:
  """Test discovery of Field instances."""

  def test_single_field_discovery(self) -> None:
    """Test discovery of a single field."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True)

    form = TestForm({})
    assert 'name' in form._attributes
    assert isinstance(form._attributes['name'], CharField)

  def test_multiple_fields_discovery(self) -> None:
    """Test discovery of multiple fields."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True)
      age = NumberField(required=False)

    form = TestForm({})
    assert 'name' in form._attributes
    assert 'age' in form._attributes
    assert len(form._attributes) == 2


class TestFormDiscovery:
  """Test discovery of nested Form instances."""

  def test_single_nested_form_discovery(self) -> None:
    """Test discovery of a single nested form."""

    class NestedForm(Form):
      """Nested form."""

      pass

    class TestForm(Form):
      """Test form."""

      nested = NestedForm()

    form = TestForm({})
    assert 'nested' in form._sub_forms_attrs
    assert isinstance(form._sub_forms_attrs['nested'], Form)


class TestListDiscovery:
  """Test discovery of nested lists."""

  def test_list_of_fields_discovery(self) -> None:
    """Test discovery of list of fields."""

    class TestForm(Form):
      """Test form."""

      tags = [CharField(required=True)]

    form = TestForm({})
    assert 'tags' in form._nested_attrs
    assert len(form._nested_attrs['tags']) == 1
    assert isinstance(form._nested_attrs['tags'][0], CharField)

  def test_list_of_forms_discovery(self) -> None:
    """Test discovery of list of forms."""

    class NestedForm(Form):
      """Nested form."""

      pass

    class TestForm(Form):
      """Test form."""

      items = [NestedForm()]

    form = TestForm({})
    assert 'items' in form._nested_attrs
    assert len(form._nested_attrs['items']) == 1
    assert isinstance(form._nested_attrs['items'][0], Form)

  def test_plain_list_attribute(self) -> None:
    """Test plain list attribute (not a form declaration)."""

    class TestForm(Form):
      """Test form with plain list attribute."""

      config = [1, 2, 3]

    form = TestForm({})
    # KNOWN BUG: any list class attribute is claimed as a nested-form declaration
    assert 'config' in form._nested_attrs
    assert form._nested_attrs['config'] == [1, 2, 3]

  def test_empty_list_attribute_skipped(self) -> None:
    """Test empty list attribute is skipped during validation."""

    class TestForm(Form):
      """Test form with empty list."""

      items = []

    form = TestForm({})
    assert 'items' in form._nested_attrs
    assert form.is_valid()
    assert form.errors() == {}


class TestCleanFunctionDiscovery:
  """Test discovery of clean methods."""

  def test_single_clean_function_discovery(self) -> None:
    """Test discovery of a single clean method."""

    class TestForm(Form):
      """Test form."""

      def clean_name(self) -> None:
        """Clean method."""
        pass

    form = TestForm({})
    assert 'clean_name' in form._clean_functions

  def test_multiple_clean_functions_discovery(self) -> None:
    """Test discovery of multiple clean methods."""

    class TestForm(Form):
      """Test form."""

      def clean_name(self) -> None:
        """Clean method."""
        pass

      def clean_age(self) -> None:
        """Clean method."""
        pass

    form = TestForm({})
    assert 'clean_name' in form._clean_functions
    assert 'clean_age' in form._clean_functions
    assert len(form._clean_functions) == 2

  def test_clean_functions_alphabetical_order(self) -> None:
    """Test clean functions are discovered in alphabetical order."""

    class TestForm(Form):
      """Test form."""

      def clean_zebra(self) -> None:
        """Clean function Z."""
        self.add_errors(key='order', code='z')

      def clean_apple(self) -> None:
        """Clean function A."""
        self.add_errors(key='order', code='a')

    form = TestForm({})
    # inspect.getmembers returns in alphabetical order
    assert form._clean_functions == ['clean_apple', 'clean_zebra']

    # Verify the order of execution
    form.is_valid()
    errors = form.errors()
    assert 'order' in errors
    # Both errors should be present; 'a' should be first (apple executes first)
    assert len(errors['order']) == 2
    assert errors['order'][0]['code'] == 'a'
    assert errors['order'][1]['code'] == 'z'


class TestReservedWordsSkipped:
  """Test that reserved words are skipped during discovery."""

  def test_reserved_words_not_discovered(self) -> None:
    """Test reserved words like 'errors', 'is_valid' are not discovered."""

    class TestForm(Form):
      """Test form."""

      pass

    form = TestForm({})
    # Reserved words should not appear in _attributes, _clean_functions, etc.
    assert 'errors' not in form._attributes
    assert 'is_valid' not in form._clean_functions
    assert 'is_valid_async' not in form._clean_functions
    assert 'cleaned_data' not in form._attributes


class TestUnderscorePrefixedSkipped:
  """Test that underscore-prefixed attributes are skipped."""

  def test_private_field_not_discovered(self) -> None:
    """Test private (underscore-prefixed) fields are not discovered."""

    class TestForm(Form):
      """Test form."""

      _private_field = CharField(required=True)

    form = TestForm({})
    assert '_private_field' not in form._attributes

  def test_private_method_not_discovered(self) -> None:
    """Test private methods are not discovered."""

    class TestForm(Form):
      """Test form."""

      def _private_method(self) -> None:
        """Private method."""
        pass

    form = TestForm({})
    assert '_private_method' not in form._clean_functions


class TestComplexDiscovery:
  """Test discovery with mixed attributes."""

  def test_mixed_discovery(self) -> None:
    """Test discovery with fields, nested forms, and clean methods."""

    class NestedForm(Form):
      """Nested form."""

      pass

    class TestForm(Form):
      """Test form with mixed attributes."""

      name = CharField(required=True)
      nested = NestedForm()
      tags = [CharField()]
      _private = CharField()

      def clean_validate(self) -> None:
        """Clean method."""
        pass

    form = TestForm({})
    assert 'name' in form._attributes
    assert 'nested' in form._sub_forms_attrs
    assert 'tags' in form._nested_attrs
    assert '_private' not in form._attributes
    assert 'clean_validate' in form._clean_functions
