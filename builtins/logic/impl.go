package logic

import "fmt"

func and(inputs map[string]any) (any, error) {
	rawA, ok := inputs["a"]
	if !ok {
		return nil, fmt.Errorf("Input `a` missing")
	}

	a, ok := rawA.(bool)
	if !ok {
		return nil, fmt.Errorf("Input `a` has the wrong type")
	}

	rawB, ok := inputs["b"]
	if !ok {
		return nil, fmt.Errorf("Input `b` missing")
	}

	b, ok := rawB.(bool)
	if !ok {
		return nil, fmt.Errorf("Input `b` has the wrong type")
	}

	return a && b, nil
}

func or(inputs map[string]any) (any, error) {
	rawA, ok := inputs["a"]
	if !ok {
		return nil, fmt.Errorf("Input `a` missing")
	}

	a, ok := rawA.(bool)
	if !ok {
		return nil, fmt.Errorf("Input `a` has the wrong type")
	}

	rawB, ok := inputs["b"]
	if !ok {
		return nil, fmt.Errorf("Input `b` missing")
	}

	b, ok := rawB.(bool)
	if !ok {
		return nil, fmt.Errorf("Input `b` has the wrong type")
	}

	return a || b, nil
}

func not(inputs map[string]any) (any, error) {
	rawA, ok := inputs["a"]
	if !ok {
		return nil, fmt.Errorf("Input `a` missing")
	}

	a, ok := rawA.(bool)
	if !ok {
		return nil, fmt.Errorf("Input `a` has the wrong type")
	}

	return !a, nil
}
