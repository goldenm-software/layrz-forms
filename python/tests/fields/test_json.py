"""Test JsonField validation."""

from layrz_forms import Form, JsonField
from tests.helpers import dump_errors


class TestJsonFieldDictAbsent:
  """Test JsonField (dict) when field is absent."""

  def test_json_dict_absent_required_true(self) -> None:
    """Test required JsonField(dict) absent from input."""

    class TestForm(Form):
      """Test form."""

      data = JsonField(datatype=dict, required=True)

    form = TestForm({})
    assert not form.is_valid()
    # KNOWN BUG: JsonField emits both 'required' and 'invalid' when absent
    assert dump_errors(form.errors) == {'data': [{'code': 'required'}, {'code': 'invalid'}]}

  def test_json_dict_absent_required_false(self) -> None:
    """Test optional JsonField(dict) absent from input."""

    class TestForm(Form):
      """Test form."""

      data = JsonField(datatype=dict, required=False)

    form = TestForm({})
    assert form.is_valid()
    assert dump_errors(form.errors) == {}


class TestJsonFieldListAbsent:
  """Test JsonField (list) when field is absent."""

  def test_json_list_absent_required_true(self) -> None:
    """Test required JsonField(list) absent from input."""

    class TestForm(Form):
      """Test form."""

      items = JsonField(datatype=list, required=True)

    form = TestForm({})
    assert not form.is_valid()
    # KNOWN BUG: JsonField emits both 'required' and 'invalid' when absent
    assert dump_errors(form.errors) == {'items': [{'code': 'required'}, {'code': 'invalid'}]}

  def test_json_list_absent_required_false(self) -> None:
    """Test optional JsonField(list) absent from input."""

    class TestForm(Form):
      """Test form."""

      items = JsonField(datatype=list, required=False)

    form = TestForm({})
    assert form.is_valid()
    assert dump_errors(form.errors) == {}


class TestJsonFieldDictNone:
  """Test JsonField (dict) with None value."""

  def test_json_dict_none_required_true(self) -> None:
    """Test required JsonField(dict) with None."""

    class TestForm(Form):
      """Test form."""

      data = JsonField(datatype=dict, required=True)

    form = TestForm({'data': None})
    assert not form.is_valid()
    # KNOWN BUG: JsonField emits both 'required' and 'invalid' for None
    assert dump_errors(form.errors) == {'data': [{'code': 'required'}, {'code': 'invalid'}]}

  def test_json_dict_none_required_false(self) -> None:
    """Test optional JsonField(dict) with None."""

    class TestForm(Form):
      """Test form."""

      data = JsonField(datatype=dict, required=False)

    form = TestForm({'data': None})
    assert form.is_valid()
    assert dump_errors(form.errors) == {}


class TestJsonFieldListNone:
  """Test JsonField (list) with None value."""

  def test_json_list_none_required_true(self) -> None:
    """Test required JsonField(list) with None."""

    class TestForm(Form):
      """Test form."""

      items = JsonField(datatype=list, required=True)

    form = TestForm({'items': None})
    assert not form.is_valid()
    # KNOWN BUG: JsonField emits both 'required' and 'invalid' for None
    assert dump_errors(form.errors) == {'items': [{'code': 'required'}, {'code': 'invalid'}]}

  def test_json_list_none_required_false(self) -> None:
    """Test optional JsonField(list) with None."""

    class TestForm(Form):
      """Test form."""

      items = JsonField(datatype=list, required=False)

    form = TestForm({'items': None})
    assert form.is_valid()
    assert dump_errors(form.errors) == {}


class TestJsonFieldDictValid:
  """Test JsonField (dict) with valid values."""

  def test_json_dict_populated(self) -> None:
    """Test JsonField(dict) with populated dict."""

    class TestForm(Form):
      """Test form."""

      data = JsonField(datatype=dict, required=True)

    form = TestForm({'data': {'key': 'value'}})
    assert form.is_valid()
    assert dump_errors(form.errors) == {}
    assert form.cleaned_data == {'data': {'key': 'value'}}

  def test_json_dict_empty_empty_true(self) -> None:
    """Test JsonField(dict) with empty dict when empty=True."""

    class TestForm(Form):
      """Test form."""

      data = JsonField(datatype=dict, required=False, empty=True)

    form = TestForm({'data': {}})
    assert form.is_valid()
    assert dump_errors(form.errors) == {}

  def test_json_dict_empty_empty_false(self) -> None:
    """Test JsonField(dict) with empty dict when empty=False."""

    class TestForm(Form):
      """Test form."""

      data = JsonField(datatype=dict, required=False, empty=False)

    form = TestForm({'data': {}})
    # KNOWN BUG: optional JsonField emits 'invalid' when the value is absent/empty
    assert not form.is_valid()
    assert dump_errors(form.errors) == {'data': [{'code': 'invalid'}]}


class TestJsonFieldListValid:
  """Test JsonField (list) with valid values."""

  def test_json_list_populated(self) -> None:
    """Test JsonField(list) with populated list."""

    class TestForm(Form):
      """Test form."""

      items = JsonField(datatype=list, required=True)

    form = TestForm({'items': [1, 2, 3]})
    assert form.is_valid()
    assert dump_errors(form.errors) == {}
    assert form.cleaned_data == {'items': [1, 2, 3]}

  def test_json_list_single_item(self) -> None:
    """Test JsonField(list) with single item."""

    class TestForm(Form):
      """Test form."""

      items = JsonField(datatype=list, required=True)

    form = TestForm({'items': ['one']})
    assert form.is_valid()
    assert dump_errors(form.errors) == {}

  def test_json_list_empty_empty_true(self) -> None:
    """Test JsonField(list) with empty list when empty=True."""

    class TestForm(Form):
      """Test form."""

      items = JsonField(datatype=list, required=False, empty=True)

    form = TestForm({'items': []})
    assert form.is_valid()
    assert dump_errors(form.errors) == {}

  def test_json_list_empty_empty_false(self) -> None:
    """Test JsonField(list) with empty list when empty=False."""

    class TestForm(Form):
      """Test form."""

      items = JsonField(datatype=list, required=False, empty=False)

    form = TestForm({'items': []})
    # KNOWN BUG: optional JsonField emits 'invalid' when the value is absent/empty
    assert not form.is_valid()
    assert dump_errors(form.errors) == {'items': [{'code': 'invalid'}]}


class TestJsonFieldDictInvalid:
  """Test JsonField (dict) with invalid values."""

  def test_json_dict_wrong_type_list(self) -> None:
    """Test JsonField(dict) with list value."""

    class TestForm(Form):
      """Test form."""

      data = JsonField(datatype=dict, required=True)

    form = TestForm({'data': []})
    assert not form.is_valid()
    assert dump_errors(form.errors) == {'data': [{'code': 'invalid'}]}

  def test_json_dict_wrong_type_string(self) -> None:
    """Test JsonField(dict) with string value."""

    class TestForm(Form):
      """Test form."""

      data = JsonField(datatype=dict, required=True)

    form = TestForm({'data': 'not a dict'})
    assert not form.is_valid()
    assert dump_errors(form.errors) == {'data': [{'code': 'invalid'}]}

  def test_json_dict_wrong_type_int(self) -> None:
    """Test JsonField(dict) with integer value."""

    class TestForm(Form):
      """Test form."""

      data = JsonField(datatype=dict, required=True)

    form = TestForm({'data': 42})
    assert not form.is_valid()
    assert dump_errors(form.errors) == {'data': [{'code': 'invalid'}]}


class TestJsonFieldListInvalid:
  """Test JsonField (list) with invalid values."""

  def test_json_list_wrong_type_dict(self) -> None:
    """Test JsonField(list) with dict value."""

    class TestForm(Form):
      """Test form."""

      items = JsonField(datatype=list, required=True)

    form = TestForm({'items': {}})
    assert not form.is_valid()
    assert dump_errors(form.errors) == {'items': [{'code': 'invalid'}]}

  def test_json_list_wrong_type_string(self) -> None:
    """Test JsonField(list) with string value."""

    class TestForm(Form):
      """Test form."""

      items = JsonField(datatype=list, required=True)

    form = TestForm({'items': 'not a list'})
    assert not form.is_valid()
    assert dump_errors(form.errors) == {'items': [{'code': 'invalid'}]}

  def test_json_list_wrong_type_int(self) -> None:
    """Test JsonField(list) with integer value."""

    class TestForm(Form):
      """Test form."""

      items = JsonField(datatype=list, required=True)

    form = TestForm({'items': 42})
    assert not form.is_valid()
    assert dump_errors(form.errors) == {'items': [{'code': 'invalid'}]}
