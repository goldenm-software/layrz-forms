"""Form class"""

import inspect
from typing import Any, Callable, Dict, List, Optional, Self, Tuple, TypeAlias, cast

from layrz_forms.fields import Field
from layrz_forms.types import ErrorType

DataObjType: TypeAlias = Dict[str, Any]


class Form:
  """Form class"""

  _obj: DataObjType = {}
  _errors: ErrorType = {}
  _clean_functions: List[str] = []
  _attributes: Dict[str, Any] = {}
  _nested_attrs: Dict[str, List[Self | Field]] = {}
  _sub_forms_attrs: Dict[str, Self] = {}

  def __init__(self: Self, obj: Optional[DataObjType] = None) -> None:
    """
    Form constructor

    :param obj: Object to validate
    :type obj: Optional[DataObjType]
    """

    if obj is None:
      obj = {}
    self._obj = obj

    self.calculate_members()

  @property
  def cleaned_data(self: Self) -> DataObjType:
    """
    Returns the cleaned data

    :return: Cleaned data
    :rtype: DataObjType
    """
    return self._obj

  def calculate_members(self: Self) -> None:
    """Calculate members"""
    self._errors = {}
    self._clean_functions = []
    self._attributes = {}
    self._nested_attrs = {}
    self._sub_forms_attrs = {}

    for item in inspect.getmembers(self):
      if item[0] in self._reserved_words:
        continue
      if item[0].startswith('_'):
        continue

      if item[0].startswith('clean'):
        self._clean_functions.append(item[0])
        continue

      if isinstance(item[1], Field):
        self._attributes[item[0]] = item[1]
        continue

      if isinstance(item[1], list):
        self._nested_attrs[item[0]] = item[1]
        continue

      if isinstance(item[1], Form):
        self._sub_forms_attrs[item[0]] = cast(Self, item[1])
        continue

  @property
  def obj(self: Self) -> DataObjType:
    """
    Returns the object

    :return: Object
    :rtype: DataObjType
    """
    return self._obj

  @obj.setter
  def obj(self: Self, obj: DataObjType) -> None:
    """
    Set the object

    :param obj: Object to validate
    :type obj: DataObjType
    """
    self._obj = obj

  def is_valid(self: Self) -> bool:
    """
    Returns if the form is valid

    :return: True if the form is valid, False otherwise
    :rtype: bool
    """
    self._errors = {}

    for field in self._attributes.items():
      self._validate_field(field=field)

    for attr, form in self._sub_forms_attrs.items():
      self._validate_sub_form(
        field=attr,
        form=form,
        data=self._obj.get(attr, {}),
      )

    for nattr, nform in self._nested_attrs.items():
      if isinstance(nform[0], Field):
        self._validate_sub_form(
          field=nattr,
          form=nform[0],  # type: ignore[arg-type]
          data=self._obj.get(nattr, {}),
        )
      else:
        self._validate_sub_form_as_list(field=nattr, form=nform[0])

    for func in self._clean_functions:
      self._clean(clean_func=func)

    return len(self._errors) == 0

  def errors(self: Self) -> ErrorType:
    """Returns the list of errors"""
    return self._errors

  def add_errors(
    self: Self,
    *,
    key: str = '',
    code: str = '',
    extra_args: Optional[Dict[str, Any] | Callable[[Any], Any]] = None,
  ) -> None:
    """Add custom errors
    This function is designed to be used in a clean function

    :param key: Key of the field
    :type key: str
    :param code: Error code
    :type code: str
    :param extra_args: Extra arguments to add to the error
    :type extra_args: Optional[Dict[str, Any]]
    """
    if extra_args is None:
      extra_args = {}

    if key == '' or code == '':
      raise Exception('key and code are required')  # pylint: disable=W0719
    camel_key = self._convert_to_camel(key=key)

    if camel_key not in self._errors:
      self._errors[camel_key] = []

    new_error = {'code': code}
    if extra_args and isinstance(extra_args, dict):
      if callable(extra_args):
        new_error.update(extra_args())
      else:
        new_error.update(extra_args)

    self._errors[camel_key].append(new_error)

  def _validate_field(self: Self, *, field: Tuple[str, Field], new_key: Optional[str] = None) -> None:
    """
    Validate field

    :param field: Field to validate
    :type field: Tuple[str, ...]
    :param new_key: New key to use for the field
    :type new_key: Optional[str]

    :return: None
    :rtype: None
    """
    if isinstance(field[1], Field):
      func = field[1].validate
      if callable(func):
        # Validate if the validate function has the correct parameters
        params = [p for p, _ in inspect.signature(func).parameters.items()]
        valid_params = ['key', 'value', 'errors']

        if len(params) != len(valid_params):
          raise Exception(f'{type(field[1])} validate method has no the correct parameters')  # pylint: disable=W0719

        is_valid = False
        for param in params:
          if param in valid_params:
            is_valid = True
            continue
          is_valid = False
          break

        if not is_valid:
          raise Exception(  # pylint: disable=W0719
            f'{field[0]} of type {type(field[1]).__name__} validate method has no the correct '
            + f'parameters. Expected parameters: {", ".join(valid_params)}. '
            + f'Actual parameters: {", ".join(params)}'
          )

        field[1].validate(
          key=field[0] if new_key is None else new_key,
          value=self._obj.get(field[0], None),
          errors=self._errors,
        )
      else:
        raise Exception(f'{type(field[1])} has no validate method')  # pylint: disable=W0719

  def _clean(self: Self, clean_func: str) -> None:
    """Clean function"""
    func = getattr(self, clean_func)
    if callable(func):
      func()

  def _convert_to_camel(self: Self, *, key: str) -> str:
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

  def _validate_sub_form(self: Self, *, field: str, form: Self | Field, data: DataObjType) -> None:
    """Validate sub form"""
    if not isinstance(form, Form):
      return

    if not isinstance(data, dict):
      self.add_errors(
        key=field,
        code='invalid',
        extra_args={'message': 'Invalid data type'},
      )
      return

    form.obj = data

    form.calculate_members()
    if not form.is_valid():
      for key, errors in form.errors().items():
        for error in errors:
          code = error['code']
          del error['code']
          self.add_errors(key=f'{field}.{key}', code=code, extra_args=error)

  def _validate_sub_form_as_list(self: Self, *, field: str, form: Self | Field) -> None:
    """
    Validate sub form for list

    :param field: Field name
    :type field: str
    :param form: Form to validate
    :type form: Any

    :return: None
    :rtype: None
    """
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
            form=form,  # type: ignore[arg-type]
            data=obj,
          )

  @property
  def _reserved_words(self: Self) -> Tuple[str, ...]:
    """Reserved words"""
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
