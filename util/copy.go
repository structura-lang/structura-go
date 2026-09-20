package util

import "reflect"

func MakeCopy[T any](v T) T {
	copied := deepCopy(reflect.ValueOf(v))

	if !copied.IsValid() {
		var zero T
		return zero
	}

	return copied.Interface().(T)
}

func deepCopy(v reflect.Value) reflect.Value {
	if !v.IsValid() {
		return v
	}

	switch v.Kind() {
	case reflect.Interface:
		if v.IsNil() {
			return reflect.Zero(v.Type())
		}

		elem := deepCopy(v.Elem())
		result := reflect.New(v.Type()).Elem()
		result.Set(elem)
		return result

	case reflect.Map:
		if v.IsNil() {
			return reflect.Zero(v.Type())
		}

		result := reflect.MakeMapWithSize(v.Type(), v.Len())

		iter := v.MapRange()
		for iter.Next() {
			key := deepCopy(iter.Key())
			value := deepCopy(iter.Value())
			result.SetMapIndex(key, value)
		}

		return result

	case reflect.Slice:
		if v.IsNil() {
			return reflect.Zero(v.Type())
		}

		result := reflect.MakeSlice(v.Type(), v.Len(), v.Len())
		for i := 0; i < v.Len(); i++ {
			result.Index(i).Set(deepCopy(v.Index(i)))
		}

		return result

	case reflect.Array:
		result := reflect.New(v.Type()).Elem()
		for i := 0; i < v.Len(); i++ {
			result.Index(i).Set(deepCopy(v.Index(i)))
		}
		return result

	case reflect.Struct:
		result := reflect.New(v.Type()).Elem()

		for i := 0; i < v.NumField(); i++ {
			// Unexported fields cannot safely be set through reflection.
			if result.Field(i).CanSet() {
				result.Field(i).Set(deepCopy(v.Field(i)))
			}
		}

		return result

	default:
		// Strings, numbers, bools, channels, funcs, etc. are treated as values.
		return v
	}
}
