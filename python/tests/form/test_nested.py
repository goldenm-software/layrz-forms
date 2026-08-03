"""Test nested Form and list validation."""

from layrz_forms import CharField, Form, NumberField
from tests.helpers import dump_errors


class TestNestedFormValid:
  """Test nested form validation with valid data."""

  def test_nested_form_valid(self) -> None:
    """Test nested form with valid data."""

    class NestedForm(Form):
      """Nested form."""

      city = CharField(required=True)

    class ParentForm(Form):
      """Parent form."""

      name = CharField(required=True)
      address = NestedForm()

    form = ParentForm({'name': 'John', 'address': {'city': 'NYC'}})
    assert form.is_valid() is True
    assert dump_errors(form.errors) == {}


class TestNestedFormInvalid:
  """Test nested form validation with invalid data."""

  def test_nested_form_invalid_child_error(self) -> None:
    """Test nested form reports child errors with dotted keys."""

    class NestedForm(Form):
      """Nested form."""

      city = CharField(required=True)

    class ParentForm(Form):
      """Parent form."""

      name = CharField(required=True)
      address = NestedForm()

    form = ParentForm({'name': 'John', 'address': {}})
    assert form.is_valid() is False
    errors = form.errors
    assert 'address.city' in errors
    assert dump_errors(errors)['address.city'] == [{'code': 'required'}]

  def test_nested_form_multiple_child_errors(self) -> None:
    """Test nested form with multiple child errors."""

    class NestedForm(Form):
      """Nested form."""

      city = CharField(required=True, min_length=3)
      zip_code = NumberField(datatype=int, required=True)

    class ParentForm(Form):
      """Parent form."""

      address = NestedForm()

    form = ParentForm({'address': {'city': 'LA', 'zip_code': 'not_int'}})
    assert form.is_valid() is False
    errors = form.errors
    assert 'address.city' in errors
    assert 'address.zipCode' in errors
    assert errors['address.city'][0].code == 'minLength'
    assert errors['address.zipCode'][0].code == 'invalid'

  def test_nested_form_camel_case_keys(self) -> None:
    """Test nested form error keys are camelCase."""

    class NestedForm(Form):
      """Nested form."""

      first_name = CharField(required=True)

    class ParentForm(Form):
      """Parent form."""

      address_info = NestedForm()

    form = ParentForm({'address_info': {}})
    assert form.is_valid() is False
    errors = form.errors
    assert 'addressInfo.firstName' in errors


class TestNestedFormInvalidDataType:
  """Test nested form with non-dict data."""

  def test_nested_form_non_dict_required_true(self) -> None:
    """Test required nested form with non-dict value."""

    class NestedForm(Form):
      """Nested form."""

      city = CharField(required=True)

    class ParentForm(Form):
      """Parent form."""

      address = NestedForm()

    form = ParentForm({'address': 'not a dict'})
    assert form.is_valid() is False
    errors = form.errors
    assert 'address' in errors
    assert dump_errors(errors)['address'] == [{'code': 'invalid', 'extra': {'message': 'Invalid data type'}}]

  def test_nested_form_non_dict_list(self) -> None:
    """Test nested form with list value."""

    class NestedForm(Form):
      """Nested form."""

      city = CharField(required=True)

    class ParentForm(Form):
      """Parent form."""

      address = NestedForm()

    form = ParentForm({'address': ['not', 'a', 'dict']})
    assert form.is_valid() is False
    errors = form.errors
    assert 'address' in errors
    assert errors['address'][0].code == 'invalid'
    extra = errors['address'][0].extra
    assert extra is not None
    assert extra['message'] == 'Invalid data type'


class TestNestedListOfForms:
  """Test nested list of forms."""

  def test_nested_list_of_forms_valid(self) -> None:
    """Test nested list of forms with valid data."""

    class ItemForm(Form):
      """Item form."""

      name = CharField(required=True)

    class ParentForm(Form):
      """Parent form."""

      items = [ItemForm()]

    form = ParentForm({'items': [{'name': 'Item1'}, {'name': 'Item2'}]})
    assert form.is_valid() is True
    assert dump_errors(form.errors) == {}

  def test_nested_list_of_forms_invalid(self) -> None:
    """Test nested list of forms with invalid data."""

    class ItemForm(Form):
      """Item form."""

      name = CharField(required=True)

    class ParentForm(Form):
      """Parent form."""

      items = [ItemForm()]

    form = ParentForm({'items': [{'name': 'Item1'}, {}]})
    assert form.is_valid() is False
    errors = form.errors
    assert 'items.1.name' in errors
    assert dump_errors(errors)['items.1.name'] == [{'code': 'required'}]

  def test_nested_list_of_forms_empty_list(self) -> None:
    """Test nested list of forms with empty list."""

    class ItemForm(Form):
      """Item form."""

      name = CharField(required=True)

    class ParentForm(Form):
      """Parent form."""

      items = [ItemForm()]

    form = ParentForm({'items': []})
    assert form.is_valid() is True
    assert dump_errors(form.errors) == {}


class TestNestedListOfFields:
  """Test nested list of fields."""

  def test_nested_list_of_fields_valid(self) -> None:
    """Test nested list of fields with valid data."""

    class ParentForm(Form):
      """Parent form."""

      tags = [CharField(required=True)]

    form = ParentForm({'tags': ['tag1', 'tag2', 'tag3']})
    assert form.is_valid() is True
    assert dump_errors(form.errors) == {}

  def test_nested_list_of_fields_invalid_item(self) -> None:
    """Test nested list of fields with invalid item."""

    class ParentForm(Form):
      """Parent form."""

      tags = [CharField(required=True, min_length=2)]

    form = ParentForm({'tags': ['a', 'bb', 'cc']})
    # First item 'a' fails min_length=2 validation
    assert form.is_valid() is False
    errors = form.errors
    assert 'tags.0' in errors
    assert errors['tags.0'][0].code == 'minLength'

  def test_nested_list_of_fields_empty_list(self) -> None:
    """Test nested list of fields with empty list."""

    class ParentForm(Form):
      """Parent form."""

      tags = [CharField(required=False)]

    form = ParentForm({'tags': []})
    assert form.is_valid() is True
    assert dump_errors(form.errors) == {}


class TestNestedNonListValue:
  """Test nested list attribute with non-list value."""

  def test_nested_list_non_list_value_reports_error(self) -> None:
    """Test non-list value for nested list reports an error."""

    class ParentForm(Form):
      """Parent form."""

      tags = [CharField(required=True)]

    form = ParentForm({'tags': 'not_a_list'})
    assert form.is_valid() is False
    errors = form.errors
    assert 'tags' in errors
    assert dump_errors(errors)['tags'] == [{'code': 'invalid', 'extra': {'message': 'Invalid data type'}}]

  def test_nested_list_absent_value_silently_skipped(self) -> None:
    """Test absent value for nested list is silently skipped."""

    class ParentForm(Form):
      """Parent form."""

      tags = [CharField(required=False)]

    form = ParentForm({})
    # Absent field should not produce an error
    assert form.is_valid() is True
    assert dump_errors(form.errors) == {}

  def test_nested_list_dict_value_reports_error(self) -> None:
    """Test dict value for nested list reports an error."""

    class ParentForm(Form):
      """Parent form."""

      tags = [CharField(required=True)]

    form = ParentForm({'tags': {'not': 'a_list'}})
    assert form.is_valid() is False
    errors = form.errors
    assert 'tags' in errors
    assert errors['tags'][0].code == 'invalid'


class TestNestedFormErrorMutation:
  """Test that parent and child errors are independent."""

  def test_parent_child_errors_are_independent_copies(self) -> None:
    """Test parent and child form errors are independent copies."""

    class NestedForm(Form):
      """Nested form."""

      name = CharField(required=True)

    class ParentForm(Form):
      """Parent form."""

      data = NestedForm()

    # Create parent form and validate
    form = ParentForm({'data': {}})
    assert form.is_valid() is False

    parent_errors = form.errors
    assert 'data.name' in parent_errors

    # Get the nested form and its errors
    nested_form = form._sub_forms_attrs['data']
    nested_errors = nested_form.errors

    # The errors should be the same code but different objects
    assert parent_errors['data.name'][0].code == 'required'
    assert nested_errors['name'][0].code == 'required'

    # Verify they are different objects (the parent error is a copy)
    assert parent_errors['data.name'][0] is not nested_errors['name'][0]
