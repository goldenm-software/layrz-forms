"""Test Form instance isolation."""

from layrz_forms import CharField, Form, NumberField
from tests.helpers import dump_errors


class TestInstanceIsolation:
  """Test that different instances don't share state."""

  def test_two_instances_different_errors(self) -> None:
    """Test two instances of same form class have separate errors."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True)

    form1 = TestForm({})
    form2 = TestForm({'name': 'John'})

    form1.is_valid()
    form2.is_valid()

    assert dump_errors(form1.errors) == {'name': [{'code': 'required'}]}
    assert dump_errors(form2.errors) == {}

  def test_two_instances_different_cleaned_data(self) -> None:
    """Test two instances have separate cleaned data."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=False)

    data1 = {'name': 'Alice'}
    data2 = {'name': 'Bob'}

    form1 = TestForm(data1)
    form2 = TestForm(data2)

    form1.is_valid()
    form2.is_valid()

    assert form1.cleaned_data == {'name': 'Alice'}
    assert form2.cleaned_data == {'name': 'Bob'}
    assert form1.cleaned_data is not form2.cleaned_data

  def test_validation_does_not_affect_other_instance(self) -> None:
    """Test validating one instance doesn't affect another."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True)
      age = NumberField(required=False)

    form1 = TestForm({'name': 'John'})
    form2 = TestForm({})

    form1.is_valid()
    # form2 hasn't been validated yet
    assert form2._errors == {}

    form2.is_valid()
    assert dump_errors(form1.errors) == {}
    assert dump_errors(form2.errors) == {'name': [{'code': 'required'}]}


class TestSubclassIsolation:
  """Test that subclasses with inherited fields behave correctly."""

  def test_subclass_inherits_parent_fields(self) -> None:
    """Test subclass inherits parent fields."""

    class ParentForm(Form):
      """Parent form."""

      name = CharField(required=True)

    class ChildForm(ParentForm):
      """Child form."""

      age = NumberField(required=False)

    child = ChildForm({'name': 'John', 'age': 30})
    assert child.is_valid() is True
    assert 'name' in child._attributes
    assert 'age' in child._attributes

  def test_subclass_does_not_modify_parent_fields(self) -> None:
    """Test subclass doesn't modify parent class fields."""

    class ParentForm(Form):
      """Parent form."""

      name = CharField(required=True)

    class ChildForm(ParentForm):
      """Child form."""

      age = NumberField(required=False)

    parent = ParentForm({'name': 'John'})
    child = ChildForm({'name': 'Jane', 'age': 25})

    parent.is_valid()
    child.is_valid()

    # Parent should only have its own fields
    parent_attrs = set(parent._attributes.keys())
    child_attrs = set(child._attributes.keys())

    assert 'name' in parent_attrs
    assert 'age' not in parent_attrs
    assert 'name' in child_attrs
    assert 'age' in child_attrs


class TestRepeatedValidation:
  """Test repeated validation on same instance."""

  def test_repeated_validation_resets_errors(self) -> None:
    """Test repeated validation resets old errors."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=True)

    form = TestForm({})
    form.is_valid()
    assert dump_errors(form.errors) == {'name': [{'code': 'required'}]}

    # Update and validate again
    form.obj = {'name': 'John'}
    form.is_valid()
    assert dump_errors(form.errors) == {}

  def test_multiple_clean_runs_dont_accumulate(self) -> None:
    """Test calling is_valid multiple times doesn't accumulate clean errors."""

    class TestForm(Form):
      """Test form."""

      name = CharField(required=False)

      def clean_check(self) -> None:
        """Custom check."""
        self.add_errors(key='name', code='always_error')

    form = TestForm({})

    form.is_valid()
    first_errors = form.errors
    assert len(first_errors['name']) == 1

    form.is_valid()
    second_errors = form.errors
    assert len(second_errors['name']) == 1
