package layrz

// ValidateNumber validates a Number field. Rejects bool explicitly and
// enforces strict type matching for datatype (int does not accept float, and vice versa).
func ValidateNumber(value any, r NumberRules) []*FieldError {
	_, present := deref(value)
	if !present {
		if r.Required {
			return []*FieldError{{Code: "required"}}
		}
		return nil
	}

	// Explicitly reject bool (Python wrongly accepted True as int)
	if _, isBool := value.(bool); isBool {
		return []*FieldError{{Code: "invalid"}}
	}

	var numVal float64
	var isInt bool
	var isFloat bool

	// Check type and extract numeric value
	switch v := value.(type) {
	case int:
		numVal = float64(v)
		isInt = true
	case int8:
		numVal = float64(v)
		isInt = true
	case int16:
		numVal = float64(v)
		isInt = true
	case int32:
		numVal = float64(v)
		isInt = true
	case int64:
		numVal = float64(v)
		isInt = true
	case float32:
		numVal = float64(v)
		isFloat = true
	case float64:
		numVal = v
		isFloat = true
	default:
		return []*FieldError{{Code: "invalid"}}
	}

	// For datatype="float", reject int values (Python behavior)
	if r.Datatype == "float" && isInt {
		return []*FieldError{{Code: "invalid"}}
	}

	// For datatype="int", reject float values
	if r.Datatype == "int" && isFloat {
		return []*FieldError{{Code: "invalid"}}
	}

	var errs []*FieldError

	// MinValue check
	if r.MinValue != nil {
		if numVal < *r.MinValue {
			var expectedVal any
			var receivedVal any

			if r.Datatype == "int" {
				expectedVal = int64(*r.MinValue)
				if isInt {
					receivedVal = int64(numVal)
				} else {
					receivedVal = numVal
				}
			} else {
				expectedVal = *r.MinValue
				receivedVal = numVal
			}

			errs = append(errs, &FieldError{
				Code:     "minValue",
				Expected: expectedVal,
				Received: receivedVal,
			})
		}
	}

	// MaxValue check
	if r.MaxValue != nil {
		if numVal > *r.MaxValue {
			var expectedVal any
			var receivedVal any

			if r.Datatype == "int" {
				expectedVal = int64(*r.MaxValue)
				if isInt {
					receivedVal = int64(numVal)
				} else {
					receivedVal = numVal
				}
			} else {
				expectedVal = *r.MaxValue
				receivedVal = numVal
			}

			errs = append(errs, &FieldError{
				Code:     "maxValue",
				Expected: expectedVal,
				Received: receivedVal,
			})
		}
	}

	return errs
}
