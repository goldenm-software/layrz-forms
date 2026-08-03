"""Form field and sub-form validators."""

import inspect
from typing import TYPE_CHECKING, Any

if TYPE_CHECKING:
  from .form import Form

from .fields import Field


def validate_field(form: 'Form', *, field: tuple[str, Field], new_key: str | None = None) -> None:
  """
  Validate a field.

  :param form: Form instance
  :type form: Form
  :param field: Field to validate (name, Field instance)
  :type field: tuple[str, Field]
  :param new_key: New key to use for the field (for nested validation)
  :type new_key: str | None
  """
  if isinstance(field[1], Field):
    func = field[1].validate
    if callable(func):
      # Validate if the validate function has the correct parameters
      params = [p for p, _ in inspect.signature(func).parameters.items()]
      valid_params = ['key', 'value', 'errors']

      if len(params) != len(valid_params):
        raise RuntimeError(f'{type(field[1])} validate method has no the correct parameters')

      is_valid = False
      for param in params:
        if param in valid_params:
          is_valid = True
          continue
        is_valid = False
        break

      if not is_valid:
        raise RuntimeError(
          f'{field[0]} of type {type(field[1]).__name__} validate method has no the correct '
          + f'parameters. Expected parameters: {", ".join(valid_params)}. '
          + f'Actual parameters: {", ".join(params)}'
        )

      field[1].validate(
        key=field[0] if new_key is None else new_key,
        value=form._obj.get(field[0], None),
        errors=form._errors,
      )
    else:
      raise RuntimeError(f'{type(field[1])} has no validate method')


def validate_sub_form(form: 'Form', *, field: str, sub_form: 'Form | Field', data: dict[str, Any]) -> None:
  """
  Validate a sub-form.

  :param form: Parent form instance
  :type form: Form
  :param field: Field name
  :type field: str
  :param sub_form: Sub-form or field to validate
  :type sub_form: Form | Field
  :param data: Data to validate
  :type data: dict[str, Any]
  """
  # Import here to avoid circular import at module load
  from .form import Form

  if not isinstance(sub_form, Form):
    return

  if not isinstance(data, dict):
    form.add_errors(
      key=field,
      code='invalid',
      extra_args={'message': 'Invalid data type'},
    )
    return

  sub_form.obj = data

  sub_form.calculate_members()
  if not sub_form.is_valid():
    for key, errors in sub_form.errors().items():
      for error in errors:
        code = error['code']
        del error['code']
        form.add_errors(key=f'{field}.{key}', code=code, extra_args=error)


def validate_sub_form_as_list(form: 'Form', *, field: str, sub_form: 'Form | Field') -> None:
  """
  Validate a sub-form as a list.

  :param form: Form instance
  :type form: Form
  :param field: Field name
  :type field: str
  :param sub_form: Form or field to validate as list
  :type sub_form: Form | Field
  """
  # Import here to avoid circular import at module load
  from .form import Form

  list_obj = form._obj.get(field, [])

  if isinstance(list_obj, (list, tuple)):
    for i, obj in enumerate(list_obj):
      if isinstance(sub_form, Field):
        validate_field(
          form,
          field=obj,
          new_key=f'{field}.{i}',
        )
      elif isinstance(sub_form, Form):
        validate_sub_form(
          form,
          field=f'{field}.{i}',
          sub_form=sub_form,
          data=obj,
        )
