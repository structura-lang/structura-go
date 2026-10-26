package data

import "fmt"

func length(inputs map[string]any) (any, error) {
	arr, arrOk := inputs["arr"].([]any)
	if !arrOk {
		return nil, fmt.Errorf("Input `arr` missing")
	}

	return len(arr), nil
}

func equal(inputs map[string]any) (any, error) {
	a, ok := inputs["a"]
	if !ok {
		return nil, fmt.Errorf("Input `a` missing")
	}

	b, ok := inputs["b"]
	if !ok {
		return nil, fmt.Errorf("Input `b` missing")
	}

	return a == b, nil
}
