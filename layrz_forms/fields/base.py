"""Base class for the fields"""

from typing import Any, Self

from layrz_forms.types import ErrorType


class Field:
  """Field abstract class"""

  def __init__(self: Self, required: bool = False) -> None:
    self.required = required

  def validate(self: Self, key: str, value: Any, errors: ErrorType) -> None:
    """
    Validate is the field is blank or None if is required

    :param key: Key of the field
    :type key: str
    :param value: Value of the field
    :type value: Any
    :param errors: Errors dict
    :type errors: ErrorType
    """

    if self.required:
      if value is None:
        self._append_error(key=key, errors=errors, to_add={'code': 'required'})

  def _convert_to_camel(self: Self, key: str) -> str:
    """
    Convert the key to camel case

    :param key: Key to convert
    :type key: str

    :return: Key in camel case
    :rtype: str
    """
    init, *temp = key.split('_')

    field = ''.join([init, *map(str.title, temp)])
    field_items = field.split('.')

    field_final = []
    for item in field_items:
      field_final.append(''.join([item[0].lower(), item[1:]]))

    return '.'.join(field_final)

  def _append_error(self: Self, key: str, errors: ErrorType, to_add: ErrorType) -> None:
    """
    Append an error to a dict of errors

    :param key: Key of the field
    :type key: str
    :param errors: Errors dict
    :type errors: ErrorType
    :param to_add: Error to add
    :type to_add: ErrorType
    """

    key = self._convert_to_camel(key=key)
    if key in errors:
      errors[key].append(to_add)
    else:
      errors[key] = [to_add]
