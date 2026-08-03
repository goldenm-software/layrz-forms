"""ID field"""

from typing import Any, Self

from layrz_forms.errors import LayrzError
from layrz_forms.types import ErrorsType

from .base import Field


class IdField(Field):
  """ID Field"""

  def __init__(self: Self, required: bool = False) -> None:
    """
    IdField constructor

    :param required: Indicates if the field is required or not
    :type required: bool
    """
    super().__init__(required=required)

  def validate(self: Self, key: str, value: Any, errors: ErrorsType) -> None:
    """
    Validate the field with the following rules:
    - Should be a number or a string that can be converted to a number
    - The number should be greater than 0

    :param key: Key of the field
    :type key: str
    :param value: Value of the field
    :type value: Any
    :param errors: Errors mapping
    :type errors: ErrorsType
    """

    super().validate(key=key, value=value, errors=errors)

    if value is None:
      return

    # Reject booleans explicitly (isinstance(True, int) is True in Python)
    if isinstance(value, bool):
      self._append_error(
        key=key,
        errors=errors,
        to_add=LayrzError(code='invalid'),
      )
      return

    # Check if value is int or str
    if not isinstance(value, (int, str)):
      self._append_error(
        key=key,
        errors=errors,
        to_add=LayrzError(code='invalid'),
      )
      return

    # If it's a string, try to convert to int
    if isinstance(value, str):
      try:
        value = int(value)
      except ValueError:
        self._append_error(
          key=key,
          errors=errors,
          to_add=LayrzError(code='invalid'),
        )
        return

    # Check if value is positive
    try:
      if value <= 0:
        self._append_error(
          key=key,
          errors=errors,
          to_add=LayrzError(code='invalid'),
        )
    except TypeError:
      self._append_error(
        key=key,
        errors=errors,
        to_add=LayrzError(code='invalid'),
      )
