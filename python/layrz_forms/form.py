import asyncio
import inspect
from collections.abc import Callable
from typing import Any, Self, cast

from strawberry.types import get_object_definition, has_object_definition

from .casing import to_camel_case
from .errors import LayrzError
from .fields import Field
from .introspection import discover_members
from .types import ErrorsType
from .validators import validate_field, validate_sub_form, validate_sub_form_as_list


class Form:
  """
  Base class for the forms, this class contains the logic to validate the fields, sub-forms and nested forms.

  To use it, you need to inherit from this class and define the fields as class attributes.
  """

  _obj: dict[str, Any]
  _errors: ErrorsType
  _validated: bool
  _clean_functions: list[str]
  _attributes: dict[str, Any]
  _nested_attrs: dict[str, list[Self | Field]]
  _sub_forms_attrs: dict[str, Self]

  @staticmethod
  def strawberry_to_dict(obj: object) -> dict[str, Any]:
    """
    Convert a strawberry object to a dictionary

    :param obj: Strawberry object
    :return: Dictionary representation of the object
    """
    data = {}
    definitions = get_object_definition(obj)
    if definitions:
      for f in definitions.fields:
        name = f.graphql_name or f.name
        data[name] = getattr(obj, f.name)
    return data

  def __init__(self: Self, obj: object | None = None) -> None:
    """
    Form constructor

    :param obj: Object to validate
    """
    self._obj = {}
    self._errors = {}
    self._validated = False
    self._clean_functions = []
    self._attributes = {}
    self._nested_attrs = {}
    self._sub_forms_attrs = {}

    if isinstance(obj, dict):
      self.obj = cast(dict[str, Any], obj)
    elif obj is not None and has_object_definition(obj):
      self.obj = self.strawberry_to_dict(obj=obj)
    else:
      self.obj = {}

    self.calculate_members()

  @property
  def cleaned_data(self: Self) -> dict[str, Any]:
    """
    Returns the cleaned data

    :return: Cleaned data
    """
    return self._obj

  def calculate_members(self: Self) -> None:
    """Calculate members"""
    self._errors = {}
    self._validated = False
    discovery = discover_members(self)
    self._clean_functions = discovery.clean_functions
    self._attributes = discovery.fields
    self._nested_attrs = discovery.nested
    self._sub_forms_attrs = discovery.sub_forms

  @property
  def obj(self: Self) -> dict[str, Any]:
    """
    Returns the object

    :return: Object
    """
    return self._obj

  @obj.setter
  def obj(self: Self, obj: dict[str, Any]) -> None:
    """
    Set the object

    :param obj: Object to validate
    """
    self._obj = obj
    self._validated = False
    self._errors = {}

  def _run_field_validations(self: Self) -> None:
    """Run field, sub-form, and nested validations (shared by sync/async)."""
    for field in self._attributes.items():
      self._validate_field(field=field)

    for attr, form in self._sub_forms_attrs.items():
      self._validate_sub_form(
        field=attr,
        form=form,
        data=self._obj.get(attr, {}),
      )

    for nattr, nform in self._nested_attrs.items():
      if len(nform) == 0:
        # Skip empty list attributes
        continue
      if isinstance(nform[0], Field):
        self._validate_sub_form(
          field=nattr,
          form=nform[0],
          data=self._obj.get(nattr, {}),
        )
      else:
        self._validate_sub_form_as_list(field=nattr, form=nform[0])

  async def ais_valid(self: Self) -> bool:
    """
    Returns if the form is valid asynchronously

    :return: True if the form is valid, False otherwise
    """
    self._errors = {}

    self._run_field_validations()

    for func in self._clean_functions:
      await self._clean_async(clean_func=func)

    self._validated = True
    return len(self._errors) == 0

  def is_valid(self: Self) -> bool:
    """
    Returns if the form is valid

    :return: True if the form is valid, False otherwise
    """
    self._errors = {}

    self._run_field_validations()

    for func in self._clean_functions:
      self._clean_sync(clean_func=func)

    self._validated = True
    return len(self._errors) == 0

  @property
  def errors(self: Self) -> ErrorsType:
    """
    The validation errors, keyed by camelCase field name.

    Runs synchronous validation on first access if it has not run yet.

    :return: Mapping of camelCase field name to its list of errors
    :rtype: ErrorsType
    """
    if not self._validated:
      self.is_valid()
    return self._errors

  def add_errors(
    self: Self,
    key: str = '',
    code: str = '',
    extra_args: dict[str, Any] | Callable[[Any], Any] | None = None,
  ) -> None:
    """
    Add custom errors
    This function is designed to be used in a clean function

    :param key: Key of the field
    :param code: Error code
    :param extra_args: Extra arguments to add to the error
    """
    if key == '' or code == '':
      raise RuntimeError('key and code are required')
    camel_key = self._convert_to_camel(key=key)

    args: dict[str, Any] = dict(cast(dict[str, Any], extra_args)) if isinstance(extra_args, dict) else {}
    error = LayrzError(
      code=code,
      expected=args.pop('expected', None),
      received=args.pop('received', None),
      extra=args or None,
    )
    self._errors.setdefault(camel_key, []).append(error)

  def _validate_field(self: Self, *, field: tuple[str, Field], new_key: str | None = None) -> None:
    """
    Validate field

    :param field: Field to validate
    :param new_key: New key to use for the field
    """
    validate_field(self, field=field, new_key=new_key)

  def _clean_sync(self: Self, clean_func: str) -> None:
    """Clean function"""
    func = getattr(self, clean_func)
    if callable(func):
      if inspect.iscoroutinefunction(func):
        raise RuntimeError(
          'Cannot call async clean function in sync context',
          'please use ais_valid method',
        )
      # It is sync call it
      func()

  async def _clean_async(self: Self, clean_func: str) -> None:
    """Clean function async"""
    func = getattr(self, clean_func)
    if callable(func):
      if inspect.iscoroutinefunction(func):
        await func()
      else:
        await asyncio.sleep(0)  # This is to ensure the function is awaitable
        func()

  def _convert_to_camel(self: Self, *, key: str) -> str:
    """
    Convert the key to camel case

    :param key: Key to convert
    :type key: str

    :return: Key in camelCase
    :rtype: str
    """
    return to_camel_case(key)

  def _validate_sub_form(self: Self, *, field: str, form: Self | Field, data: dict[str, Any]) -> None:
    """Validate sub form"""
    validate_sub_form(self, field=field, sub_form=form, data=data)

  def _validate_sub_form_as_list(self: Self, *, field: str, form: Self | Field) -> None:
    """
    Validate sub form for list

    :param field: Field name
    :param form: Form to validate
    """
    validate_sub_form_as_list(self, field=field, sub_form=form)

  @property
  def _reserved_words(self: Self) -> tuple[str, ...]:
    """Reserved words"""
    return (
      'add_errors',
      'change_obj',
      'clean',
      'errors',
      'is_valid',
      'ais_valid',
      'set_obj',
      'calculate_members',
      'cleaned_data',
    )
