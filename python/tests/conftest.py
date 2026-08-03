"""Pytest configuration and shared fixtures."""

from typing import Any

import pytest

from layrz_forms import Form


@pytest.fixture
def simple_form() -> type[Form]:
  """Return a simple form class for testing."""

  class SimpleForm(Form):
    """Simple test form."""

    pass

  return SimpleForm


@pytest.fixture
def form_with_fields() -> type[Form]:
  """Return a form with various fields for testing."""
  from layrz_forms import BooleanField, CharField, NumberField

  class FormWithFields(Form):
    """Test form with multiple fields."""

    name = CharField(required=True)
    age = NumberField(datatype=int, required=False)
    active = BooleanField(required=False)

  return FormWithFields


@pytest.fixture
def empty_dict() -> dict[str, Any]:
  """Return an empty dictionary."""
  return {}


@pytest.fixture
def sample_dict() -> dict[str, Any]:
  """Return a sample dictionary for testing."""
  return {'name': 'John', 'age': 30, 'active': True}
