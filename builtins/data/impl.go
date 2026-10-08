package data

import "fmt"

func length(inputs map[string]any) (any, error) {
	arr, arrOk := inputs["arr"].([]any)
	if !arrOk {
		return nil, fmt.Errorf("Input `arr` missing")
	}

	return len(arr), nil
}
