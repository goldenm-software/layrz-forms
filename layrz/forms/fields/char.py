""" Email field """
from .base import Field


class CharField(Field):
  """
  CharField class for validation
  ---
  Attributes
    required: bool
      Indicates if the field is required or not
    max_length: int
      Indicates the maximum length of the field
    min_length: int
      Indicates the minimum length of the field
    empty: bool
      Indicates if the field can be empty or not
    choices: tuple
      Indicates the choices of the field
  """

  def __init__(
    self,
    required=False,
    max_length=None,
    min_length=None,
    empty=False,
    choices=None,
  ):
    super(CharField, self).__init__(required=required)
    self.max_length = max_length
    self.min_length = min_length
    self.empty = empty
    self.choices = choices

  def validate(self, key, value, errors):
    """
    Validate the field with the following rules:
    - Should not be empty if required
    - Should be one of the choices indicated if choices is not None
    - Should be less than max_length if max_length is not None
    - Should be greater than min_length if min_length is not None
    ---
    Arguments
      key: str
        Key of the field
      value: any
        Value to validate
      errors: dict
        Dict of errors
    """

    super(CharField, self).validate(key=key, value=value, errors=errors)

    if value is not None:
      if not self.empty:
        if len(value) == 0:
          self._append_error(
            key=key,
            errors=errors,
            to_add={'code': 'empty'},
          )

      if self.max_length is not None:
        if len(value) > self.max_length:
          self._append_error(
            key=key,
            errors=errors,
            to_add={
              'code': 'maxLength',
              'expected': self.max_length,
              'received': len(value),
            },
          )

      if self.min_length is not None:
        if len(value) < self.min_length:
          self._append_error(
            key=key,
            errors=errors,
            to_add={
              'code': 'minLength',
              'expected': self.min_length,
              'received': len(value),
            },
          )

      if self.choices is not None:
        mapped_choices = [choice[0] for choice in self.choices]
        if value not in mapped_choices:
          self._append_error(
            key=key,
            errors=errors,
            to_add={
              'code': 'invalidChoice',
              'expected': mapped_choices,
              'received': value,
            },
          )
