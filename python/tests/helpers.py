"""Helper functions for testing form errors."""

from typing import Any

from layrz_forms.types import ErrorsType


def dump_errors(errors: ErrorsType) -> dict[str, list[dict[str, Any]]]:
  """
  Dump a form's errors to plain dicts for comparison.

  :param errors: The errors mapping from form.errors
  :return: Dictionary of field names to lists of error dicts
  """
  return {key: [error.model_dump() for error in value] for key, value in errors.items()}
