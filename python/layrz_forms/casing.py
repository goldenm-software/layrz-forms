"""Convert snake_case keys to camelCase."""


def to_camel_case(key: str) -> str:
  """
  Convert snake_case keys to camelCase, handling dotted notation.

  For example:
    - 'user_name' → 'userName'
    - 'a_b.c_d' → 'aB.cD'
    - '' → ''

  :param key: Key to convert
  :type key: str

  :return: Key in camelCase
  :rtype: str
  """
  result = []
  for segment in key.split('.'):
    init, *temp = segment.split('_')
    camel = ''.join([init, *map(str.title, temp)])
    result.append(''.join([camel[0].lower(), camel[1:]]) if camel else camel)
  return '.'.join(result)
