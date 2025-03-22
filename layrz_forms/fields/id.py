"""ID field"""

from typing import Any, Self

from layrz_forms.types import ErrorType

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

  def validate(self: Self, key: str, value: Any, errors: ErrorType) -> None:
    """
    Validate the field with the following rules:
    - Should be a number or a string that can be converted to a number
    - The number should be greater than 0

    :param key: Key of the field
    :type key: str
    :param value: Value of the field
    :type value: Any
    :param errors: Errors dict
    :type errors: ErrorType
    """

    super().validate(key=key, value=value, errors=errors)

    if not isinstance(value, (int, str)) and (self.required and value is not None):
      self._append_error(
        key=key,
        errors=errors,
        to_add={'code': 'invalid'},
      )
    else:
      if value is None:
        return
      if isinstance(value, str):
        value = int(value)
      if value <= 0:
        self._append_error(
          key=key,
          errors=errors,
          to_add={'code': 'invalid'},
        )
