"""Test Strawberry GraphQL input object integration."""

import strawberry

from layrz_forms import CharField, Form, NumberField, UuidField
from tests.helpers import dump_errors


@strawberry.input
class AddressInput:
  """Test Strawberry input for address."""

  street: str
  city: str


@strawberry.input
class PersonInput:
  """Test Strawberry input type."""

  name: str
  email: str = ''
  age: int = 0
  custom_id: int = strawberry.field(name='customId')


class TestStrawberryToDict:
  """Test strawberry_to_dict static method."""

  def test_strawberry_to_dict_basic(self) -> None:
    """Test converting Strawberry input to dict."""
    person_input = PersonInput(name='John', email='john@example.com', age=30, custom_id=1)
    result = Form.strawberry_to_dict(person_input)
    assert result == {'name': 'John', 'email': 'john@example.com', 'age': 30, 'customId': 1}

  def test_strawberry_to_dict_graphql_name(self) -> None:
    """Test strawberry_to_dict respects name (graphql_name)."""
    person_input = PersonInput(name='Jane', custom_id=2)
    result = Form.strawberry_to_dict(person_input)
    # custom_id field has name='customId'
    assert 'customId' in result
    assert result['customId'] == 2
    # Original field name should NOT be present
    assert 'custom_id' not in result

  def test_strawberry_to_dict_nested(self) -> None:
    """Test strawberry_to_dict with nested Strawberry objects."""
    address = AddressInput(street='123 Main St', city='NYC')
    result = Form.strawberry_to_dict(address)
    assert result == {'street': '123 Main St', 'city': 'NYC'}


class TestFormWithStrawberryInput:
  """Test Form validation with Strawberry input objects."""

  def test_form_with_strawberry_input_valid(self) -> None:
    """Test Form with valid Strawberry input."""

    class PersonForm(Form):
      """Form for person data."""

      name = CharField(required=True)
      email = CharField(required=False)

    person_input = PersonInput(name='John', email='john@example.com', custom_id=42)
    form = PersonForm(person_input)
    assert form.is_valid() is True
    assert dump_errors(form.errors) == {}
    # custom_id is mapped to customId via strawberry field name
    assert form.cleaned_data == {'name': 'John', 'email': 'john@example.com', 'age': 0, 'customId': 42}

  def test_form_with_strawberry_input_invalid(self) -> None:
    """Test Form validation with invalid Strawberry input."""

    class PersonForm(Form):
      """Form for person data."""

      name = CharField(required=True, min_length=3)
      email = CharField(required=True)

    person_input = PersonInput(name='Jo', custom_id=1)  # name too short
    form = PersonForm(person_input)
    assert form.is_valid() is False
    assert 'name' in form.errors
    assert form.errors['name'][0].code == 'minLength'

  def test_form_missing_required_field_in_strawberry(self) -> None:
    """Test Form with Strawberry input missing required field."""

    @strawberry.input
    class MinimalInput:
      """Minimal input."""

      name: str

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True)
      email = CharField(required=True)

    input_obj = MinimalInput(name='John')
    form = TestForm(input_obj)
    assert form.is_valid() is False
    assert 'email' in form.errors
    assert form.errors['email'][0].code == 'required'


class TestStrawberryInputWithComplexForm:
  """Test Strawberry input with complex forms."""

  def test_strawberry_input_with_custom_clean(self) -> None:
    """Test Strawberry input with custom clean method."""

    class PersonForm(Form):
      """Form for person data."""

      name = CharField(required=True)
      email = CharField(required=False)

      def clean_email_check(self) -> None:
        """Validate email."""
        if self._obj.get('email') == 'blocked@example.com':
          self.add_errors(key='email', code='blocked_email')

    person_input = PersonInput(name='John', email='blocked@example.com', custom_id=1)
    form = PersonForm(person_input)
    assert form.is_valid() is False
    assert 'email' in form.errors
    assert form.errors['email'][0].code == 'blocked_email'

  def test_strawberry_input_camel_case_in_errors(self) -> None:
    """Test Strawberry input error keys are in camelCase."""

    @strawberry.input
    class InputWithSnake:
      """Input with snake_case field."""

      first_name: str = ''
      last_name: str = ''

    class TestForm(Form):
      """Test form."""

      first_name = CharField(required=True)
      last_name = CharField(required=True)

    input_obj = InputWithSnake()
    form = TestForm(input_obj)
    assert form.is_valid() is False
    assert 'firstName' in form.errors
    assert 'lastName' in form.errors
