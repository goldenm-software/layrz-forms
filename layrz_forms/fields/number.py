"""Number field"""

from typing import Any, Optional, Self

from layrz_forms.types import ErrorType

from .base import Field


class NumberField(Field):
  """Number Field"""

  def __init__(
    self: Self,
    required: bool = False,
    datatype: type[int | float] = float,
    min_value: Optional[float] = None,
    max_value: Optional[float] = None,
  ) -> None:
    """
    NumberField constructor

    :param required: Indicates if the field is required or not
    :type required: bool
    :param datatype: Type of the field
    :type datatype: Type[int | float]
    :param min_value: Minimum value of the field
    :type min_value: Optional[float]
    :param max_value: Maximum value of the field
    :type max_value: Optional[float]
    """
    super().__init__(required=required)
    self.datatype = datatype
    self.min_value = min_value
    self.max_value = max_value

  def validate(self: Self, key: str, value: Any, errors: ErrorType) -> None:
    """
    Validate the field with the following rules:
    - Should be a int or float (Depending of the datatype)

    :param key: Key of the field
    :type key: str
    :param value: Value of the field
    :type value: Any
    :param errors: Errors dict
    :type errors: ErrorType
    """

    super().validate(key=key, value=value, errors=errors)

    if not isinstance(value, self.datatype) and (self.required and value is not None):
      self._append_error(key=key, errors=errors, to_add={'code': 'invalid'})
    else:
      try:
        if self.min_value is not None:
          if self.datatype(value) < self.datatype(self.min_value):
            self._append_error(
              key=key,
              errors=errors,
              to_add={
                'code': 'minValue',
                'expected': self.datatype(self.min_value),
                'received': self.datatype(value),
              },
            )
        if self.max_value is not None:
          if self.datatype(value) > self.datatype(self.max_value):
            self._append_error(
              key=key,
              errors=errors,
              to_add={
                'code': 'maxValue',
                'expected': self.datatype(self.max_value),
                'received': self.datatype(value),
              },
            )
      except ValueError:
        if self.required:
          self._append_error(
            key=key,
            errors=errors,
            to_add={'code': 'invalid'},
          )
