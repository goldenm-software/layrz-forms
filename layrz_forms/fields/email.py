""" Email field """
import re
from typing import Any

from .base import Field


class EmailField(Field):
  """
  EmailField class for validation
  ---
  Attributes
    required: bool
      Indicates if the field is required or not
    empty: bool
      Indicates if the field can be empty or not
    regex: str
      Regex to validate the email
  """

  def __init__(
    self,
    required: bool = False,
    empty: bool = False,
    regex: str = r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-z]{2,63}$',
  ) -> None:
    super().__init__(required=required)
    self.empty = empty
    self.regex = regex

  def validate(self, key: str, value: Any, errors: dict) -> None:
    """
    Validate the field with the following rules:
    - Should be a valid email, the validation will compile the regex
    ---
    Arguments
      key: str
        Key of the field
      value: any
        Value to validate
      errors: dict
        Dict of errors
    """

    super().validate(key=key, value=value, errors=errors)

    if isinstance(value, str):
      if not self.empty:
        if value == '' or value is None:
          self._append_error(
            key=key,
            errors=errors,
            to_add={'code': 'required'},
          )
        else:
          if not re.match(self.regex, value):
            self._append_error(
              key=key,
              errors=errors,
              to_add={'code': 'invalid'},
            )
    else:
      self._append_error(
        key=key,
        errors=errors,
        to_add={'code': 'invalid'},
      )
