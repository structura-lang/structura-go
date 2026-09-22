package interpreter

import (
	"fmt"
	"math"
	"strings"
	"structura/interpreter/util"
)

type operation struct {
	OType      string         `structura:"type"`
	OArguments map[string]any `structura:"arguments"`
}

func (o operation) evaluate(pv parentVariables) error {
	switch o.OType {
	case "store":
		arguments := o.OArguments

		if arguments == nil {
			return fmt.Errorf("`arguments` field missing!")
		}

		valueExp, exists := arguments["value"]
		if !exists {
			return fmt.Errorf("`arguments/value` expression missing!")
		}

		value, err := evalAnyExpression[any](valueExp, pv)
		if err != nil {
			return fmt.Errorf("EXP[arguments/value]: %w", err)
		}

		targetExp, exists := arguments["target_var"]
		if !exists {
			return fmt.Errorf("`arguments/target_var` expression missing!")
		}

		target, err := evalAnyExpression[string](targetExp, pv)
		if err != nil {
			return fmt.Errorf("EXP[arguments/target_var]: %w", err)
		}

		_, targetExists := pv.Variables[target]
		if !targetExists {
			return fmt.Errorf("Target variable does not exist!")
		}

		pv.Variables[target] = value

		return nil

	case "set":
		arguments := o.OArguments

		if arguments == nil {
			return fmt.Errorf("`arguments` field missing!")
		}

		raw_value, vExists := arguments["value"]
		if !vExists {
			return fmt.Errorf("`arguments/value` expression missing!")
		}

		value, err := evalAnyExpression[any](raw_value, pv)
		if err != nil {
			return fmt.Errorf("EXP[value]: %w", err)
		}

		raw_path, pExists := arguments["path"]
		raw_index, iExists := arguments["index"]

		if pExists && iExists {
			return fmt.Errorf("`path` and `index` are mutually exclusive!")
		}

		if !pExists && !iExists {
			return fmt.Errorf("`path` or `index` expressions missing!")
		}

		targetExp, exists := arguments["target_var"]
		if !exists {
			return fmt.Errorf("`arguments/target_var` expression missing!")
		}

		target, err := evalAnyExpression[string](targetExp, pv)
		if err != nil {
			return fmt.Errorf("EXP[arguments/target_var]: %w", err)
		}

		// object is map
		if pExists {
			targetMap, ok := pv.Variables[target].(map[string]any)
			if !ok {
				return fmt.Errorf("Variable `target` is not a map!")
			}

			path, err := evalAnyExpression[[]any](raw_path, pv)
			if err != nil {
				return fmt.Errorf("EXP[path]: %w", err)
			}

			castPath, err := util.CastSlice[string](path)
			if err != nil {
				return fmt.Errorf("EXP[path]: %w", err)
			}

			ok = util.SetPath(targetMap, castPath, value)
			if !ok {
				return fmt.Errorf("Error setting value of variable %s at %s!", target, strings.Join(castPath, "/"))
			}

			return nil
		}

		// object is array
		if iExists {
			targetArray, ok := pv.Variables[target].([]any)
			if !ok {
				return fmt.Errorf("Variable `target` is not an array!")
			}

			fIndex, err := evalAnyExpression[float64](raw_index, pv)
			if err != nil {
				return fmt.Errorf("EXP[index]: %w", err)
			}

			if !(fIndex == math.Trunc(fIndex)) {
				return fmt.Errorf("EXP[index]: value is not an integer!")
			}

			index := int(fIndex)

			if index < 0 || index >= len(targetArray) {
				return fmt.Errorf("EXP[index]: index out of range!")
			}

			targetArray[index] = value

			return nil
		}

		return fmt.Errorf("This error should be unreachable. How did you get here?")

	case "append":
		arguments := o.OArguments

		if arguments == nil {
			return fmt.Errorf("`arguments` field missing!")
		}

		valueExp, exists := arguments["value"]
		if !exists {
			return fmt.Errorf("`arguments/value` expression missing!")
		}

		value, err := evalAnyExpression[any](valueExp, pv)
		if err != nil {
			return fmt.Errorf("EXP[arguments/value]: %w", err)
		}

		targetExp, exists := arguments["target_var"]
		if !exists {
			return fmt.Errorf("`arguments/target_var` expression missing!")
		}

		target, err := evalAnyExpression[string](targetExp, pv)
		if err != nil {
			return fmt.Errorf("EXP[arguments/target_var]: %w", err)
		}

		targetArray, ok := pv.Variables[target].([]any)
		if !ok {
			return fmt.Errorf("Variable `target` is not an array!")
		}

		targetArray = append(targetArray, value)

		pv.Variables[target] = targetArray

		return nil

	default:
		return fmt.Errorf("Unknown operation type: %s", o.OType)
	}
}
