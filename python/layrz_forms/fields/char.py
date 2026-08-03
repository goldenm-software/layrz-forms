"""Char field"""

import re
from enum import Enum, StrEnum
from typing import Any, Self

from layrz_forms.errors import LayrzError
from layrz_forms.types import ErrorsType

from .base import Field


class CharField(Field):
  """Char Field"""

  def __init__(
    self: Self,
    required: bool = False,
    max_length: int | None = None,
    min_length: int | None = None,
    empty: bool = False,
    regex: str | None = None,
    choices: tuple[tuple[str, str], ...] | None = None,
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
    :type choices: Optional[tuple[tuple[str, str], ...]]
    """
    super().__init__(required=required)
    self.max_length = max_length
    self.min_length = min_length
    self.empty = empty
    self.choices = choices
    self.regex = regex

  def validate(self: Self, key: str, value: Any, errors: ErrorsType) -> None:
    """
    Validate the field with the following rules:
    - Value must be a string (or Enum/StrEnum, which are converted to strings)
    - Should not be empty if required
    - Should be one of the choices indicated if choices is not None
    - Should be less than max_length if max_length is not None
    - Should be greater than min_length if min_length is not None
    - Should match the regex if regex is not None

    :param key: Key of the field
    :type key: str
    :param value: Value of the field
    :type value: Any
    :param errors: Errors mapping
    :type errors: ErrorsType
    """

    super().validate(key=key, value=value, errors=errors)

    if value is not None:
      if isinstance(value, Enum):
        value = value.name
      elif isinstance(value, StrEnum):
        value = value.value

      if not isinstance(value, str):
        self._append_error(
          key=key,
          errors=errors,
          to_add=LayrzError(code='invalid'),
        )
        return

      if not self.empty:
        if len(value) == 0:
          self._append_error(
            key=key,
            errors=errors,
            to_add=LayrzError(code='empty'),
          )

      if self.max_length is not None:
        if len(value) > self.max_length:
          self._append_error(
            key=key,
            errors=errors,
            to_add=LayrzError(
              code='maxLength',
              expected=self.max_length,
              received=len(value),
            ),
          )

      if self.min_length is not None:
        if len(value) < self.min_length:
          self._append_error(
            key=key,
            errors=errors,
            to_add=LayrzError(
              code='minLength',
              expected=self.min_length,
              received=len(value),
            ),
          )

      if self.choices is not None:
        mapped_choices = [choice[0] for choice in self.choices]
        if value not in mapped_choices:
          self._append_error(
            key=key,
            errors=errors,
            to_add=LayrzError(
              code='invalidChoice',
              expected=mapped_choices,
              received=value,
            ),
          )

      if self.regex is not None:
        if not re.match(self.regex, value):
          self._append_error(
            key=key,
            errors=errors,
            to_add=LayrzError(
              code='invalidFormat',
              expected=self.regex,
              received=value,
            ),
          )
