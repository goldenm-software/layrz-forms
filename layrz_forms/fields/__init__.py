""" Layrz Forms fields
Available classes
  - Field => Base class for the fields, you can define your own using this base
  - BooleanField => Field for boolean values
  - CharField => Field for string values and choices validation
  - EmailField => Field for email values
  - IdField => Field for id values
  - JsonField => Field for json values
  - NumberField => Field for number values (int, float)
"""
from .base import Field
from .boolean import BooleanField
from .char import CharField
from .email import EmailField
from .id import IdField
from .json import JsonField
from .number import NumberField

__all__ = [
  'Field',
  'BooleanField',
  'CharField',
  'EmailField',
  'IdField',
  'JsonField',
  'NumberField',
]
