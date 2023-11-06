""" ID field """
from .base import Field


class IdField(Field):
  """
  IdField class for validation
  ---
  Attributes
    required: bool
      Indicates if the field is required or not
  """

  def __init__(self, required=False):
    super(IdField, self).__init__(required=required)

  def validate(self, key, value, errors):
    """
    Validate the field with the following rules:
    - Should be a number or a string that can be converted to a number
    - The number should be greater than 0
    ---
    Arguments
      key: str
        Key of the field
      value: any
        Value to validate
      errors: dict
        Dict of errors
    """

    super(IdField, self).validate(key=key, value=value, errors=errors)

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
