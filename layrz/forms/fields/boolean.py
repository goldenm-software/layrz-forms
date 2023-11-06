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

  def validate(self, key, value, errors):
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
    """

    super(BooleanField, self).validate(key=key, value=value, errors=errors)

    if not isinstance(value, bool) and (self.required and value is not None):
      self._append_error(
        key=key,
        errors=errors,
        to_add={'code': 'invalid'},
      )
