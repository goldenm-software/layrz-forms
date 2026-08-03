"""Form member discovery via introspection."""

import inspect
from dataclasses import dataclass
from typing import TYPE_CHECKING, Any, cast

if TYPE_CHECKING:
  from .form import Form

from .fields import Field


@dataclass(frozen=True)
class MemberDiscovery:
  """Discovered form members."""

  fields: dict[str, Any]
  sub_forms: dict[str, Any]
  nested: dict[str, Any]
  clean_functions: list[str]
  reserved_words: tuple[str, ...]


def discover_members(form_instance: 'Form') -> MemberDiscovery:
  """
  Discover Form members via introspection.

  Uses inspect.getmembers to discover Field instances, nested Form instances,
  nested lists, and clean methods. Discovery order is alphabetical (as returned
  by inspect.getmembers), with filtering rules applied in order:
    1. Skip reserved words
    2. Skip names starting with '_'
    3. Collect clean* methods
    4. Collect Field instances
    5. Collect list instances
    6. Collect Form instances

  :param form_instance: Form instance to introspect
  :type form_instance: Form

  :return: Discovered members
  :rtype: MemberDiscovery
  """
  fields: dict[str, Any] = {}
  sub_forms: dict[str, Any] = {}
  nested: dict[str, Any] = {}
  clean_functions: list[str] = []

  reserved = form_instance._reserved_words

  # Import Form here to avoid circular dependency at module load time
  from .form import Form

  for item in inspect.getmembers(form_instance):
    if item[0] in reserved:
      continue
    if item[0].startswith('_'):
      continue

    if item[0].startswith('clean'):
      clean_functions.append(item[0])
      continue

    if isinstance(item[1], Field):
      fields[item[0]] = item[1]
      continue

    if isinstance(item[1], list):
      nested[item[0]] = item[1]
      continue

    if isinstance(item[1], Form):
      sub_forms[item[0]] = cast('Form', item[1])
      continue

  return MemberDiscovery(
    fields=fields,
    sub_forms=sub_forms,
    nested=nested,
    clean_functions=clean_functions,
    reserved_words=reserved,
  )
