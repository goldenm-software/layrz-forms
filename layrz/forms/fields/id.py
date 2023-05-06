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

  def validate(self, key, value, errors, cleaned_data):
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
      cleaned_data: dict
        Dict of cleaned data
    """

    super(IdField, self).validate(key=key, value=value, errors=errors, cleaned_data=cleaned_data)

    if not isinstance(value, (int, str)):
      self._append_error(
        key=key,
        errors=errors,
        to_add={'code': 'invalid'},
      )
    else:
      if isinstance(value, str):
        value = int(value)
      if value <= 0:
        self._append_error(
          key=key,
          errors=errors,
          to_add={'code': 'invalid'},
        )
