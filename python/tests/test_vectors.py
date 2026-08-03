"""Test field validation against JSON vectors."""

import json
from pathlib import Path
from typing import Any

import pytest

import layrz_forms


def load_vector_files() -> list[tuple[str, list[dict[str, Any]]]]:
  """Load all vector JSON files from the vectors directory."""
  vector_dir = Path(__file__).parent.parent.parent / 'vectors' / 'fields'
  files = sorted(vector_dir.glob('*.json'))

  results = []
  for file_path in files:
    with open(file_path) as f:
      cases = json.load(f)
    field_name = file_path.stem
    results.append((field_name, cases))
  return results


@pytest.mark.parametrize(
  'field_name, cases',
  load_vector_files(),
  ids=lambda x: x[0] if isinstance(x, tuple) else '',
)
def test_vectors(field_name: str, cases: list[dict[str, Any]]) -> None:
  """Test field validation against vector cases."""
  for case in cases:
    field_type = getattr(layrz_forms, field_name)
    case_id = case['name']

    # Get kwargs, converting datatype strings to actual types
    kwargs = case['kwargs'].copy()
    if 'datatype' in kwargs:
      datatype_str = kwargs['datatype']
      if datatype_str == 'dict':
        kwargs['datatype'] = dict
      elif datatype_str == 'list':
        kwargs['datatype'] = list
      elif datatype_str == 'int':
        kwargs['datatype'] = int
      elif datatype_str == 'float':
        kwargs['datatype'] = float

    # Create form class with the field
    class TestForm(layrz_forms.Form):
      """Test form."""

      test_field = field_type(**kwargs)

    # Prepare input dict
    input_dict: dict[str, Any] = {}
    if 'value' in case:
      input_dict['test_field'] = case['value']
    elif 'value_absent' not in case or not case['value_absent']:
      # If neither value nor value_absent is specified, treat as absent
      pass

    # Create and validate form
    form = TestForm(input_dict)
    form.is_valid()
    errors = form.errors()

    # Convert snake_case field name to camelCase for assertion
    camel_field = 'testField'
    expected_errors = case['expected_errors']

    if expected_errors:
      # Errors expected
      assert camel_field in errors, f'{case_id}: Expected field {camel_field} in errors, got {errors}'
      actual_errors = errors[camel_field]
      assert actual_errors == expected_errors, (
        f'{case_id}: Error mismatch.\nExpected: {expected_errors}\nGot: {actual_errors}'
      )
    else:
      # No errors expected
      if camel_field in errors:
        pytest.fail(f'{case_id}: Expected no errors, but got: {errors[camel_field]}')
