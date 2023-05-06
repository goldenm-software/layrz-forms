""" Init file """
from .fields.base import Field
from .fields.boolean import BooleanField
from .fields.char import CharField
from .fields.email import EmailField
from .fields.id import IdField
from .fields.json import JsonField
from .fields.number import NumberField
import inspect
from multiprocessing import Pool, Manager
from functools import partial


class Form:
  """
  Form class
  ---
  Notes:
    - Any member that starts with `_` will be ignored.
  """

  _obj = {}
  _errors = {}
  _cleaned_data = {}
  _clean_functions = []
  _attributes = {}

  def __init__(self, obj=dict):
    """ Constructor """
    self._errors = {}
    self._obj = obj
    self._clean_functions = []
    self._attributes = {}

    for item in inspect.getmembers(self):
      if item[0] in self._reserverd_words:
        continue
      if item[0].startswith('_'):
        continue
      elif item[0].startswith('clean'):
        self._clean_functions.append(item[0])
      elif isinstance(item[1], Field):
        self._attributes[item[0]] = item[1]

  @property
  def _reserverd_words(self):
    """ Reserved words """
    return ('add_errors', 'change_obj', 'clean', 'errors', 'is_valid')

  @property
  def cleaned_data(self):
    """ Returns the cleaned data """
    return self._cleaned_data

  def is_valid(self):
    """ Returns if the form is valid """
    manager = Manager()
    errors = manager.dict()
    cleaned_data = manager.dict()

    func = partial(self._validate_field, errors=errors, cleaned_data=cleaned_data)
    with Pool() as pool:
      pool.map(func, self._attributes.items())

    self._cleaned_data = dict(cleaned_data)
    self._errors = dict(errors)

    for func in self._clean_functions:
      self._clean(clean_func=func)

  def errors(self):
    """ Returns the list of errors """
    return self._errors

  def _validate_field(self, field, errors, cleaned_data):
    """ Validate field """
    if isinstance(field[1], Field):
      func = getattr(field[1], 'validate')
      if callable(func):
        # Validate if the validate function has the correct parameters
        params = [p for p, _ in inspect.signature(func).parameters.items()]

        if len(params) != 4:
          raise Exception(f'{type(field[1])} validate method has no the correct parameters')

        valid_params = ['key', 'value', 'errors', 'cleaned_data']
        is_valid = False
        for param in params:
          if param in valid_params:
            is_valid = True
            continue
          is_valid = False
          break

        if not is_valid:
          raise Exception(
            f"{field[0]} of type {type(field[1]).__name__} validate method has no the correct " +\
            f"parameters. Expected parameters: {', '.join(valid_params)}. " +\
            f"Actual parameters: {', '.join(params)}"
          )

        field[1].validate(
          key=field[0],
          value=self._obj.get(field[0], None),
          errors=errors,
          cleaned_data=cleaned_data,
        )
      else:
        raise Exception(f'{type(field[1])} has no validate method')

  def _clean(self, clean_func):
    """ Clean function """
    func = getattr(self, clean_func)
    if callable(func):
      func()

  def add_errors(self, key='', code='', extra_args=dict):
    """ Add custom errors
    This function is designed to be used in a clean function
    ---
    Arguments
      key: str
        Key of the field
      code: str
        Code of the error
      extra_args: dict
        Extra arguments to add to the error
    """
    if key == '' or code == '':
      raise
    camel_key = self._convert_to_camel(key=key)

    if camel_key not in self._errors:
      self._errors[camel_key] = []

    new_error = {'code': code}
    if extra_args and isinstance(extra_args, dict):
      if callable(extra_args):
        extra_args = extra_args()

      new_error.update(extra_args)
    self._errors[camel_key].append(new_error)

  def _convert_to_camel(self, key):
    """
    Convert the key to camel case
    """
    init, *temp = key.split('_')

    field = ''.join([init, *map(str.title, temp)])
    field_items = field.split(".")

    field_final = []
    for item in field_items:
      field_final.append(''.join([item[0].lower(), item[1:]]))

    return '.'.join(field_final)
