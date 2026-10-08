package math

import "fmt"

func add(inputs map[string]any) (any, error) {
	rawA, ok := inputs["a"]
	if !ok {
		return nil, fmt.Errorf("Input `a` missing")
	}

	a, ok := toFloat64(rawA)
	if !ok {
		return nil, fmt.Errorf("Input `a` has the wrong type")
	}

	rawB, ok := inputs["b"].(float64)
	if !ok {
		return nil, fmt.Errorf("Input `b` missing")
	}

	b, ok := toFloat64(rawB)
	if !ok {
		return nil, fmt.Errorf("Input `b` has the wrong type")
	}

	return a + b, nil
}

func subtract(inputs map[string]any) (any, error) {
	rawA, ok := inputs["a"]
	if !ok {
		return nil, fmt.Errorf("Input `a` missing")
	}

	a, ok := toFloat64(rawA)
	if !ok {
		return nil, fmt.Errorf("Input `a` has the wrong type")
	}

	rawB, ok := inputs["b"].(float64)
	if !ok {
		return nil, fmt.Errorf("Input `b` missing")
	}

	b, ok := toFloat64(rawB)
	if !ok {
		return nil, fmt.Errorf("Input `b` has the wrong type")
	}

	return a - b, nil
}

func multiply(inputs map[string]any) (any, error) {
	rawA, ok := inputs["a"]
	if !ok {
		return nil, fmt.Errorf("Input `a` missing")
	}

	a, ok := toFloat64(rawA)
	if !ok {
		return nil, fmt.Errorf("Input `a` has the wrong type")
	}

	rawB, ok := inputs["b"].(float64)
	if !ok {
		return nil, fmt.Errorf("Input `b` missing")
	}

	b, ok := toFloat64(rawB)
	if !ok {
		return nil, fmt.Errorf("Input `b` has the wrong type")
	}

	return a * b, nil
}

func divide(inputs map[string]any) (any, error) {
	rawA, ok := inputs["a"]
	if !ok {
		return nil, fmt.Errorf("Input `a` missing")
	}

	a, ok := toFloat64(rawA)
	if !ok {
		return nil, fmt.Errorf("Input `a` has the wrong type")
	}

	rawB, ok := inputs["b"].(float64)
	if !ok {
		return nil, fmt.Errorf("Input `b` missing")
	}

	b, ok := toFloat64(rawB)
	if !ok {
		return nil, fmt.Errorf("Input `b` has the wrong type")
	}

	return a / b, nil
}

func lessThan(inputs map[string]any) (any, error) {
	rawA, ok := inputs["a"]
	if !ok {
		return nil, fmt.Errorf("Input `a` missing")
	}

	a, ok := toFloat64(rawA)
	if !ok {
		return nil, fmt.Errorf("Input `a` has the wrong type")
	}

	rawB, ok := inputs["b"].(float64)
	if !ok {
		return nil, fmt.Errorf("Input `b` missing")
	}

	b, ok := toFloat64(rawB)
	if !ok {
		return nil, fmt.Errorf("Input `b` has the wrong type")
	}

	return a < b, nil
}

func greaterThan(inputs map[string]any) (any, error) {
	rawA, ok := inputs["a"]
	if !ok {
		return nil, fmt.Errorf("Input `a` missing")
	}

	a, ok := toFloat64(rawA)
	if !ok {
		return nil, fmt.Errorf("Input `a` has the wrong type")
	}

	rawB, ok := inputs["b"].(float64)
	if !ok {
		return nil, fmt.Errorf("Input `b` missing")
	}

	b, ok := toFloat64(rawB)
	if !ok {
		return nil, fmt.Errorf("Input `b` has the wrong type")
	}

	return a > b, nil
}

func equal(inputs map[string]any) (any, error) {
	rawA, ok := inputs["a"]
	if !ok {
		return nil, fmt.Errorf("Input `a` missing")
	}

	a, ok := toFloat64(rawA)
	if !ok {
		return nil, fmt.Errorf("Input `a` has the wrong type")
	}

	rawB, ok := inputs["b"].(float64)
	if !ok {
		return nil, fmt.Errorf("Input `b` missing")
	}

	b, ok := toFloat64(rawB)
	if !ok {
		return nil, fmt.Errorf("Input `b` has the wrong type")
	}

	return a == b, nil
}
