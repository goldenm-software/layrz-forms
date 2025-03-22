"""Char field"""

from typing import Any, List, Optional, Self

from layrz_forms.types import ErrorType

from .base import Field


class CharField(Field):
  """Char Field"""

  def __init__(
    self: Self,
    required: bool = False,
    max_length: Optional[int] = None,
    min_length: Optional[int] = None,
    empty: bool = False,
    choices: Optional[List[str]] = None,
  ) -> None:
    """
    CharField constructor

    :param required: Indicates if the field is required or not
    :type required: bool
    :param max_length: Maximum length of the field
    :type max_length: Optional[int]
    :param min_length: Minimum length of the field
    :type min_length: Optional[int]
    :param empty: Indicates if the field can be empty
    :type empty: bool
    :param choices: List of choices for the field
    :type choices: Optional[List[str]]
    """
    super().__init__(required=required)
    self.max_length = max_length
    self.min_length = min_length
    self.empty = empty
    self.choices = choices

  def validate(self: Self, key: str, value: Any, errors: ErrorType) -> None:
    """
    Validate the field with the following rules:
    - Should not be empty if required
    - Should be one of the choices indicated if choices is not None
    - Should be less than max_length if max_length is not None
    - Should be greater than min_length if min_length is not None

    :param key: Key of the field
    :type key: str
    :param value: Value of the field
    :type value: Any
    :param errors: Errors dict
    :type errors: ErrorType
    """

    super().validate(key=key, value=value, errors=errors)

    if value is not None:
      if not self.empty:
        if len(value) == 0:
          self._append_error(
            key=key,
            errors=errors,
            to_add={'code': 'empty'},
          )

      if self.max_length is not None:
        if len(value) > self.max_length:
          self._append_error(
            key=key,
            errors=errors,
            to_add={
              'code': 'maxLength',
              'expected': self.max_length,
              'received': len(value),
            },
          )

      if self.min_length is not None:
        if len(value) < self.min_length:
          self._append_error(
            key=key,
            errors=errors,
            to_add={
              'code': 'minLength',
              'expected': self.min_length,
              'received': len(value),
            },
          )

      if self.choices is not None:
        mapped_choices = [choice[0] for choice in self.choices]
        if value not in mapped_choices:
          self._append_error(
            key=key,
            errors=errors,
            to_add={
              'code': 'invalidChoice',
              'expected': mapped_choices,
              'received': value,
            },
          )
