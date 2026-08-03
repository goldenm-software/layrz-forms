"""Layrz Forms Types"""

from typing import TypeAlias

from .errors import LayrzError

ErrorType: TypeAlias = LayrzError
ErrorsType: TypeAlias = dict[str, list[LayrzError]]
