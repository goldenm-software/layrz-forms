"""UUID field"""

import uuid
from typing import Any, Self

from layrz_forms.types import ErrorType

from .base import Field


class UuidField(Field):
  """UUID Field"""

  def __init__(self: Self, required: bool = False) -> None:
    """
    UuidField constructor

    :param required: Indicates if the field is required or not
    :type required: bool
    """
    super().__init__(required=required)

  def validate(self: Self, key: str, value: Any, errors: ErrorType) -> None:
    """
    Validate the field with the following rules:
    - Should be a string with valid UUID format or a uuid.UUID instance
    - Uses Python's built-in uuid.UUID() constructor for validation

    :param key: Key of the field
    :type key: str
    :param value: Value of the field
    :type value: Any
    :param errors: Errors dict
    :type errors: ErrorType
    """

    super().validate(key=key, value=value, errors=errors)

    if value is None and not self.required:
      return

    # Validate the value is a str or a UUID class
    if not isinstance(value, (str, uuid.UUID)):
      self._append_error(
        key=key,
        errors=errors,
        to_add={'code': 'invalid'},
      )
      return
    # If it's already a UUID instance, it's valid
    if isinstance(value, uuid.UUID):
      return
    # If it's a string, validate UUID format
    try:
      uuid.UUID(value)
    except ValueError:
      self._append_error(
        key=key,
        errors=errors,
        to_add={'code': 'invalid'},
      )
