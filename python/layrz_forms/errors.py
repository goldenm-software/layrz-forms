"""Structured error model for form validation."""

from typing import Any, Literal, Self

from pydantic import BaseModel, ConfigDict


class LayrzError(BaseModel):
  """A single validation error attached to a form field."""

  model_config = ConfigDict(extra='forbid')

  code: str
  """The error code identifying the type of validation failure."""
  expected: Any = None
  """The expected value or constraint that was violated."""
  received: Any = None
  """The actual value received that failed validation."""
  extra: dict[str, Any] | None = None
  """Additional contextual data about the error."""

  def model_dump(
    self: Self,
    *,
    mode: str = 'python',
    include: Any = None,
    exclude: Any = None,
    context: Any = None,
    by_alias: bool | None = None,
    exclude_unset: bool = False,
    exclude_defaults: bool = False,
    exclude_none: bool | None = None,
    exclude_computed_fields: bool = False,
    round_trip: bool = False,
    warnings: bool | Literal['none', 'warn', 'error'] = True,
    fallback: Any | None = None,
    serialize_as_any: bool = False,
    polymorphic_serialization: bool | None = None,
  ) -> dict[str, Any]:
    """
    Dump the error, omitting unset (``None``) slots by default.

    :param mode: Mode for serialization (default 'python')
    :type mode: str
    :param include: Fields to include
    :type include: Any
    :param exclude: Fields to exclude
    :type exclude: Any
    :param context: Serialization context
    :type context: Any
    :param by_alias: Use field aliases
    :type by_alias: bool | None
    :param exclude_unset: Exclude unset fields
    :type exclude_unset: bool
    :param exclude_defaults: Exclude default values
    :type exclude_defaults: bool
    :param exclude_none: Exclude None values (default True)
    :type exclude_none: bool | None
    :param exclude_computed_fields: Exclude computed fields
    :type exclude_computed_fields: bool
    :param round_trip: Enable round-trip serialization
    :type round_trip: bool
    :param warnings: Warnings mode
    :type warnings: bool | str
    :param fallback: Fallback function
    :type fallback: Any | None
    :param serialize_as_any: Serialize as any type
    :type serialize_as_any: bool
    :param polymorphic_serialization: Polymorphic serialization mode
    :type polymorphic_serialization: bool | None
    :return: Dictionary representation of the error
    :rtype: dict[str, Any]
    """
    if exclude_none is None:
      exclude_none = True
    return super().model_dump(
      mode=mode,
      include=include,
      exclude=exclude,
      context=context,
      by_alias=by_alias,
      exclude_unset=exclude_unset,
      exclude_defaults=exclude_defaults,
      exclude_none=exclude_none,
      exclude_computed_fields=exclude_computed_fields,
      round_trip=round_trip,
      warnings=warnings,
      fallback=fallback,
      serialize_as_any=serialize_as_any,
      polymorphic_serialization=polymorphic_serialization,
    )

  def model_dump_json(
    self: Self,
    *,
    indent: int | None = None,
    include: Any = None,
    exclude: Any = None,
    context: Any = None,
    by_alias: bool | None = None,
    exclude_unset: bool = False,
    exclude_defaults: bool = False,
    exclude_none: bool | None = None,
    exclude_computed_fields: bool = False,
    round_trip: bool = False,
    warnings: bool | Literal['none', 'warn', 'error'] = True,
    fallback: Any | None = None,
    serialize_as_any: bool = False,
    polymorphic_serialization: bool | None = None,
    ensure_ascii: bool = False,
  ) -> str:
    """
    Dump the error as JSON, omitting unset (``None``) slots by default.

    :param indent: Indentation level
    :type indent: int | None
    :param include: Fields to include
    :type include: Any
    :param exclude: Fields to exclude
    :type exclude: Any
    :param context: Serialization context
    :type context: Any
    :param by_alias: Use field aliases
    :type by_alias: bool | None
    :param exclude_unset: Exclude unset fields
    :type exclude_unset: bool
    :param exclude_defaults: Exclude default values
    :type exclude_defaults: bool
    :param exclude_none: Exclude None values (default True)
    :type exclude_none: bool | None
    :param exclude_computed_fields: Exclude computed fields
    :type exclude_computed_fields: bool
    :param round_trip: Enable round-trip serialization
    :type round_trip: bool
    :param warnings: Warnings mode
    :type warnings: bool | str
    :param fallback: Fallback function
    :type fallback: Any | None
    :param serialize_as_any: Serialize as any type
    :type serialize_as_any: bool
    :param polymorphic_serialization: Polymorphic serialization mode
    :type polymorphic_serialization: bool | None
    :param ensure_ascii: Ensure ASCII characters only
    :type ensure_ascii: bool
    :return: JSON representation of the error
    :rtype: str
    """
    if exclude_none is None:
      exclude_none = True
    return super().model_dump_json(
      indent=indent,
      include=include,
      exclude=exclude,
      context=context,
      by_alias=by_alias,
      exclude_unset=exclude_unset,
      exclude_defaults=exclude_defaults,
      exclude_none=exclude_none,
      exclude_computed_fields=exclude_computed_fields,
      round_trip=round_trip,
      warnings=warnings,
      fallback=fallback,
      serialize_as_any=serialize_as_any,
      polymorphic_serialization=polymorphic_serialization,
      ensure_ascii=ensure_ascii,
    )
