""" JSON Field """
from .base import Field


class JsonField(Field):
  """
  JsonField class for validation
  ---
  Attributes
    required: bool
      Indicates if the field is required or not
    empty: bool
      Indicates if the field can be empty or not
    datatype: (dict, list)
      Indicates the datatype of the field
  """

  def __init__(
    self,
    required=False,
    empty=False,
    datatype=dict,
  ):
    super(JsonField, self).__init__(required=required)
    self.empty = empty
    self.datatype = datatype

  def validate(self, key, value, errors):
    """
    Validate the field with the following rules:
    - Should be a dict or list (Depending of the datatype)
    - If `empty` is False, the field should not be empty
      * For `dict`, should have at least 1 key
      * For `list`, should have at least 1 item
    ---
    Arguments
      key: str
        Key of the field
      value: any
        Value to validate
      errors: dict
        Dict of errors
    """

    super(JsonField, self).validate(key=key, value=value, errors=errors)

    if not isinstance(value, self.datatype):
      self._append_error(
        key=key,
        errors=errors,
        to_add={'code': 'invalid'},
      )
    elif not self.empty:
      length = 0

      if isinstance(self.datatype(), dict):
        length = len(value.keys())
      elif isinstance(self.datatype(), list):
        length = len(value)

      if length == 0:
        self._append_error(
          key=key,
          errors=errors,
          to_add={'code': 'invalid'},
        )
