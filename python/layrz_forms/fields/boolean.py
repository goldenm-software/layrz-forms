"""Boolean field"""

from typing import Any, Self

from layrz_forms.errors import LayrzError
from layrz_forms.types import ErrorsType

from .base import Field


class BooleanField(Field):
  """Boolean Field"""

  def __init__(self: Self, required: bool = False) -> None:
    """
    BooleanField constructor

    :param required: Indicates if the field is required or not
    :type required: bool
    """
    super().__init__(required=required)

  def validate(self: Self, key: str, value: Any, errors: ErrorsType) -> None:
    """
    Validate the field with the following rules:
    - Should be a bool

    :param key: Key of the field
    :type key: str
    :param value: Value of the field
    :type value: Any
    :param errors: Errors mapping
    :type errors: ErrorsType
    """

    super().validate(key=key, value=value, errors=errors)

    if value is not None and not isinstance(value, bool):
      self._append_error(
        key=key,
        errors=errors,
        to_add=LayrzError(code='invalid'),
      )
