"""JSON Field"""

from typing import Any, Self, Type

from layrz_forms.errors import LayrzError
from layrz_forms.types import ErrorsType

from .base import Field


class JsonField(Field):
  """JSON Field"""

  def __init__(
    self: Self,
    required: bool = False,
    empty: bool = False,
    datatype: Type[list[Any] | dict[Any, Any]] = dict,
  ) -> None:
    """
    JsonField constructor

    :param required: Indicates if the field is required or not
    :type required: bool
    :param empty: Indicates if the field can be empty
    :type empty: bool
    :param datatype: Type of the field
    :type datatype: Type[list[Any] | dict[Any, Any]]
    """
    super().__init__(required=required)
    self.empty = empty
    self.datatype = datatype

  def validate(self: Self, key: str, value: Any, errors: ErrorsType) -> None:
    """
    Validate the field with the following rules:
    - Should be a dict or list (Depending of the datatype)
    - If `empty` is False, the field should not be empty
      * For `dict`, should have at least 1 key
      * For `list`, should have at least 1 item

    :param key: Key of the field
    :type key: str
    :param value: Value of the field
    :type value: Any
    :param errors: Errors mapping
    :type errors: ErrorsType
    """

    super().validate(key=key, value=value, errors=errors)

    if not isinstance(value, self.datatype):
      self._append_error(
        key=key,
        errors=errors,
        to_add=LayrzError(code='invalid'),
      )

    elif not self.empty:
      length = 0

      if issubclass(self.datatype, dict):
        length = len(value.keys())
      elif issubclass(self.datatype, list):
        length = len(value)

      if length == 0:
        self._append_error(
          key=key,
          errors=errors,
          to_add=LayrzError(code='invalid'),
        )
