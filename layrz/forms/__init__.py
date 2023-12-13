""" Init file """
import inspect

from .fields.base import Field
from .fields.boolean import BooleanField
from .fields.char import CharField
from .fields.email import EmailField
from .fields.id import IdField
from .fields.json import JsonField
from .fields.number import NumberField


class Form:
  """
  Form class
  ---
  Notes:
    - Any member that starts with `_` will be ignored.
  """

  _obj = {}
  _errors = {}
  _clean_functions = []
  _attributes = {}
  _nested_attrs = {}
  _sub_forms_attrs = {}

  def __init__(self, obj=dict):
    """ Constructor """
    self._obj = obj
    self.calculate_members()

  @property
  def cleaned_data(self):
    """ Returns the cleaned data """
    return self._obj

  def calculate_members(self):
    """ Calculate members """
    self._errors = {}
    self._clean_functions = []
    self._attributes = {}
    self._nested_attrs = {}
    self._sub_forms_attrs = {}

    for item in inspect.getmembers(self):
      if item[0] in self._reserverd_words:
        continue
      if item[0].startswith('_'):
        continue
      elif item[0].startswith('clean'):
        self._clean_functions.append(item[0])
      elif isinstance(item[1], Field):
        self._attributes[item[0]] = item[1]
      elif isinstance(item[1], list):
        self._nested_attrs[item[0]] = item[1]
      elif isinstance(item[1], Form):
        self._sub_forms_attrs[item[0]] = item[1]
      else:
        print('Unknown field', item[0])

  def set_obj(self, obj):
    """ Set the object """
    self._obj = obj

  def is_valid(self):
    """ Returns if the form is valid """
    self._errors = {}

    for field in self._attributes.items():
      self._validate_field(field=field)

    for field, form in self._sub_forms_attrs.items():
      self._validate_sub_form(field=field, form=form, data=self._obj.get(field, {}))

    for field, form in self._nested_attrs.items():
      self._validate_sub_form_as_list(field=field, form=form[0])

    for func in self._clean_functions:
      self._clean(clean_func=func)

    return len(self._errors) == 0

  def errors(self):
    """ Returns the list of errors """
    return self._errors

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
      raise Exception('key and code are required')  #pylint: disable=W0719
    camel_key = self._convert_to_camel(key=key)

    if camel_key not in self._errors:
      self._errors[camel_key] = []

    new_error = {'code': code}
    if extra_args and isinstance(extra_args, dict):
      if callable(extra_args):
        extra_args = extra_args()

      new_error.update(extra_args)
    self._errors[camel_key].append(new_error)

  def _validate_field(self, field, new_key=None):
    """ Validate field """
    if isinstance(field[1], Field):
      func = getattr(field[1], 'validate')
      if callable(func):
        # Validate if the validate function has the correct parameters
        params = [p for p, _ in inspect.signature(func).parameters.items()]
        valid_params = ['key', 'value', 'errors']

        if len(params) != len(valid_params):
          raise Exception(f'{type(field[1])} validate method has no the correct parameters')  #pylint: disable=W0719

        is_valid = False
        for param in params:
          if param in valid_params:
            is_valid = True
            continue
          is_valid = False
          break

        if not is_valid:
          raise Exception(   #pylint: disable=W0719
            f'{field[0]} of type {type(field[1]).__name__} validate method has no the correct ' +\
            f'parameters. Expected parameters: {", ".join(valid_params)}. ' +\
            f'Actual parameters: {", ".join(params)}'
          )

        field[1].validate(
          key=field[0] if new_key is None else new_key,
          value=self._obj.get(field[0], None),
          errors=self._errors,
        )
      else:
        raise Exception(f'{type(field[1])} has no validate method')  #pylint: disable=W0719

  def _clean(self, clean_func):
    """ Clean function """
    func = getattr(self, clean_func)
    if callable(func):
      func()

  def _convert_to_camel(self, key):
    """
    Convert the key to camel case
    """
    init, *temp = key.split('_')

    field = ''.join([init, *map(str.title, temp)])
    field_items = field.split('.')

    field_final = []
    for item in field_items:
      field_final.append(''.join([item[0].lower(), item[1:]]))

    return '.'.join(field_final)

  def _validate_sub_form(self, field, form, data):
    """ Validate sub form """
    if not isinstance(form, Form):
      return
    form.set_obj(data)
    form.calculate_members()
    if not form.is_valid():
      for key, errors in form.errors().items():
        for error in errors:
          code = error['code']
          del error['code']
          self.add_errors(key=f'{field}.{key}', code=code, extra_args=error)

  def _validate_sub_form_as_list(self, field, form):
    """ Validate sub form for list """
    list_obj = self._obj.get(field, [])

    if isinstance(list_obj, (list, tuple)):
      for i, obj in enumerate(list_obj):
        if isinstance(form, Field):
          self._validate_field(
            field=obj,
            new_key=f'{field}.{i}',
          )
        elif isinstance(form, Form):
          self._validate_sub_form(
            field=f'{field}.{i}',
            form=form,
            data=obj,
          )

  @property
  def _reserverd_words(self):
    """ Reserved words """
    return (
      'add_errors',
      'change_obj',
      'clean',
      'errors',
      'is_valid',
      'set_obj',
      'calculate_members',
      'cleaned_data',
    )
