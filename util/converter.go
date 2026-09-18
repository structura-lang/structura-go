package util

import (
	"fmt"
	"reflect"
	"strings"
)

func MapToStruct[T any](data map[string]any, tagName string) (T, error) {
	var result T

	dst := reflect.ValueOf(&result).Elem()

	if dst.Kind() != reflect.Struct {
		return result, fmt.Errorf("mapper: target must be a struct, got %s", dst.Kind())
	}

	if err := mapIntoStruct(dst, data, tagName); err != nil {
		return result, err
	}

	return result, nil
}

func mapIntoStruct(dst reflect.Value, data map[string]any, tagName string) error {
	dstType := dst.Type()

	for i := 0; i < dst.NumField(); i++ {
		fieldType := dstType.Field(i)
		fieldValue := dst.Field(i)

		if fieldType.PkgPath != "" {
			continue
		}

		tag, ok := fieldType.Tag.Lookup(tagName)
		if !ok {
			continue
		}

		name, _, _ := strings.Cut(tag, ",")

		// `tag:"-"`
		if name == "-" {
			continue
		}

		raw, exists := data[name]
		if !exists {
			continue
		}

		if err := setValue(fieldValue, raw, tagName); err != nil {
			return fmt.Errorf("mapper: field %q: %w", name, err)
		}
	}

	return nil
}

func setValue(dst reflect.Value, src any, tagName string) error {
	if !dst.CanSet() {
		return fmt.Errorf("cannot set value")
	}

	if src == nil {
		switch dst.Kind() {
		case reflect.Interface,
			reflect.Pointer,
			reflect.Map,
			reflect.Slice,
			reflect.Func,
			reflect.Chan:
			dst.SetZero()
			return nil

		default:
			return fmt.Errorf("cannot assign nil to %s", dst.Type())
		}
	}

	srcValue := reflect.ValueOf(src)
	srcType := srcValue.Type()

	if dst.Kind() == reflect.Interface {
		if srcType.Implements(dst.Type()) {
			dst.Set(srcValue)
			return nil
		}

		return fmt.Errorf("type %s does not implement %s", srcType, dst.Type())
	}

	if srcType == dst.Type() {
		dst.Set(srcValue)
		return nil
	}

	switch dst.Kind() {
	case reflect.Struct:
		return setStruct(dst, src, tagName)

	case reflect.Slice:
		return setSlice(dst, src, tagName)

	case reflect.Array:
		return setArray(dst, src, tagName)

	case reflect.Map:
		return setMap(dst, src, tagName)

	default:
		return fmt.Errorf("cannot assign %s to %s", srcType, dst.Type())
	}
}

func setStruct(dst reflect.Value, src any, tagName string) error {
	srcMap, ok := src.(map[string]any)
	if !ok {
		return fmt.Errorf("expected map[string]any for %s, got %T", dst.Type(), src)
	}

	return mapIntoStruct(dst, srcMap, tagName)
}

func setSlice(dst reflect.Value, src any, tagName string) error {
	srcValue := reflect.ValueOf(src)

	if srcValue.Kind() != reflect.Slice && srcValue.Kind() != reflect.Array {
		return fmt.Errorf("expected slice or array for %s, got %T", dst.Type(), src)
	}

	result := reflect.MakeSlice(dst.Type(), srcValue.Len(), srcValue.Len())

	for i := 0; i < srcValue.Len(); i++ {
		if err := setValue(result.Index(i), srcValue.Index(i).Interface(), tagName); err != nil {
			return fmt.Errorf("index %d: %w", i, err)
		}
	}

	dst.Set(result)
	return nil
}

func setArray(dst reflect.Value, src any, tagName string) error {
	srcValue := reflect.ValueOf(src)

	if srcValue.Kind() != reflect.Array && srcValue.Kind() != reflect.Slice {
		return fmt.Errorf("expected array or slice for %s, got %T", dst.Type(), src)
	}

	if srcValue.Len() != dst.Len() {
		return fmt.Errorf("expected %d elements, got %d", dst.Len(), srcValue.Len())
	}

	for i := 0; i < dst.Len(); i++ {
		if err := setValue(dst.Index(i), srcValue.Index(i).Interface(), tagName); err != nil {
			return fmt.Errorf("index %d: %w", i, err)
		}
	}

	return nil
}

func setMap(dst reflect.Value, src any, tagName string) error {
	srcValue := reflect.ValueOf(src)

	if srcValue.Kind() != reflect.Map {
		return fmt.Errorf("expected map for %s, got %T", dst.Type(), src)
	}

	result := reflect.MakeMapWithSize(dst.Type(), srcValue.Len())

	iter := srcValue.MapRange()

	for iter.Next() {
		newKey := reflect.New(dst.Type().Key()).Elem()
		newValue := reflect.New(dst.Type().Elem()).Elem()

		if err := setValue(newKey, iter.Key().Interface(), tagName); err != nil {
			return fmt.Errorf("map key: %w", err)
		}

		if err := setValue(newValue, iter.Value().Interface(), tagName); err != nil {
			return fmt.Errorf("map value: %w", err)
		}

		result.SetMapIndex(newKey, newValue)
	}

	dst.Set(result)
	return nil
}
