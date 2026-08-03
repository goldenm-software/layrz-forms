package layrz

// deref returns the concrete value from an interface{}, handling nil interface,
// nil pointer, and non-pointer values. Returns (concrete value, present bool).
// If v is nil or a nil pointer, returns (nil, false). Otherwise returns (v or *v, true).
func deref(v any) (any, bool) {
	if v == nil {
		return nil, false
	}

	// Check if it's a pointer
	type ptrVal interface {
		isPtr()
	}

	// For pointers, we need to check if the pointer itself is nil
	// Use a type assertion approach to safely handle this
	switch pv := v.(type) {
	case *int:
		if pv == nil {
			return nil, false
		}
		return *pv, true
	case *int8:
		if pv == nil {
			return nil, false
		}
		return *pv, true
	case *int16:
		if pv == nil {
			return nil, false
		}
		return *pv, true
	case *int32:
		if pv == nil {
			return nil, false
		}
		return *pv, true
	case *int64:
		if pv == nil {
			return nil, false
		}
		return *pv, true
	case *float32:
		if pv == nil {
			return nil, false
		}
		return *pv, true
	case *float64:
		if pv == nil {
			return nil, false
		}
		return *pv, true
	case *string:
		if pv == nil {
			return nil, false
		}
		return *pv, true
	case *bool:
		if pv == nil {
			return nil, false
		}
		return *pv, true
	case *[]any:
		if pv == nil {
			return nil, false
		}
		return *pv, true
	case *map[string]any:
		if pv == nil {
			return nil, false
		}
		return *pv, true
	default:
		// Non-pointer value is always present
		return v, true
	}
}
