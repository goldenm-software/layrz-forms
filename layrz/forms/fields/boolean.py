""" Boolean field """
from .base import Field


class BooleanField(Field):
  """
  IdField class for validation
  ---
  Attributes
    required: bool
      Indicates if the field is required or not
  """

  def __init__(self, required=False):
    super(BooleanField, self).__init__(required=required)

  def validate(self, key, value, errors, cleaned_data):
    """
    Validate the field with the following rules:
    - Should be a bool
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

    super(BooleanField, self).validate(key=key, value=value, errors=errors, cleaned_data=cleaned_data)

    if not isinstance(value, bool):
      self._append_error(
        key=key,
        errors=errors,
        to_add={'code': 'invalid'},
      )
