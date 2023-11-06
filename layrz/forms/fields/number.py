""" Boolean field """
from .base import Field


class NumberField(Field):
  """
  NumberField class for validation
  ---
  Attributes
    required: bool
      Indicates if the field is required or not
    datatype: (int, float)
      Indicates the datatype of the field
    min_value: int
      Indicates the minimum value of the field
    max_value: int
      Indicates the maximum value of the field
  """

  def __init__(
    self,
    required=False,
    datatype=float,
    min_value=None,
    max_value=None,
  ):
    super(NumberField, self).__init__(required=required)
    self.datatype = datatype
    self.min_value = min_value
    self.max_value = max_value

  def validate(self, key, value, errors):
    """
    Validate the field with the following rules:
    - Should be a int or float (Depending of the datatype)
    ---
    Arguments
      key: str
        Key of the field
      value: any
        Value to validate
      errors: dict
        Dict of errors
    """

    super(NumberField, self).validate(key=key, value=value, errors=errors)

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
