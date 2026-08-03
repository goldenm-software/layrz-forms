"""Email field"""

import re
from typing import Any, Self

from layrz_forms.errors import LayrzError
from layrz_forms.types import ErrorsType

from .base import Field


class EmailField(Field):
  """Email Field"""

  def __init__(
    self: Self,
    required: bool = False,
    empty: bool = False,
    regex: str = r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-z]{2,63}$',
  ) -> None:
    """
    EmailField constructor

    :param required: Indicates if the field is required or not
    :type required: bool
    :param empty: Indicates if the field can be empty
    :type empty: bool
    :param regex: Regex to validate the email
    :type regex: str
    """
    super().__init__(required=required)
    self.empty = empty
    self.regex = regex

  def validate(self: Self, key: str, value: Any, errors: ErrorsType) -> None:
    """
    Validate the field with the following rules:
    - Should be a valid email, the validation will compile the regex

    :param key: Key of the field
    :type key: str
    :param value: Value of the field
    :type value: Any
    :param errors: Errors mapping
    :type errors: ErrorsType
    """

    super().validate(key=key, value=value, errors=errors)

    if isinstance(value, str):
      if not self.empty:
        if value == '':
          self._append_error(
            key=key,
            errors=errors,
            to_add=LayrzError(code='required'),
          )
        else:
          if not re.match(self.regex, value):
            self._append_error(
              key=key,
              errors=errors,
              to_add=LayrzError(code='invalid'),
            )
    else:
      self._append_error(
        key=key,
        errors=errors,
        to_add=LayrzError(code='invalid'),
      )
