package util

import "fmt"

func CastSlice[T any](input []any) ([]T, error) {
	result := make([]T, len(input))

	for i, v := range input {
		value, ok := v.(T)
		if !ok {
			return nil, fmt.Errorf("Element %d has type %T, expected %T", i, v, *new(T))
		}

		result[i] = value
	}

	return result, nil
}
