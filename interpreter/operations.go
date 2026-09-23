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

type returnReason int

const (
	rDone returnReason = iota
	rError
)

func (o operation) evaluate(pv parentVariables) (returnReason, error) {
	arguments := o.OArguments

	if arguments == nil {
		return rError, fmt.Errorf("`arguments` field missing!")
	}

	switch o.OType {
	case "store":
		valueExp, exists := arguments["value"]
		if !exists {
			return rError, fmt.Errorf("`arguments/value` expression missing!")
		}

		value, err := evalAnyExpression[any](valueExp, pv)
		if err != nil {
			return rError, fmt.Errorf("EXP[arguments/value]: %w", err)
		}

		targetExp, exists := arguments["target_var"]
		if !exists {
			return rError, fmt.Errorf("`arguments/target_var` expression missing!")
		}

		target, err := evalAnyExpression[string](targetExp, pv)
		if err != nil {
			return rError, fmt.Errorf("EXP[arguments/target_var]: %w", err)
		}

		_, targetExists := pv.Variables[target]
		if !targetExists {
			return rError, fmt.Errorf("Target variable does not exist!")
		}

		pv.Variables[target] = value

		return rDone, nil

	case "set":
		raw_value, vExists := arguments["value"]
		if !vExists {
			return rError, fmt.Errorf("`arguments/value` expression missing!")
		}

		value, err := evalAnyExpression[any](raw_value, pv)
		if err != nil {
			return rError, fmt.Errorf("EXP[value]: %w", err)
		}

		raw_path, pExists := arguments["path"]
		raw_index, iExists := arguments["index"]

		if pExists && iExists {
			return rError, fmt.Errorf("`path` and `index` are mutually exclusive!")
		}

		if !pExists && !iExists {
			return rError, fmt.Errorf("`path` or `index` expressions missing!")
		}

		targetExp, exists := arguments["target_var"]
		if !exists {
			return rError, fmt.Errorf("`arguments/target_var` expression missing!")
		}

		target, err := evalAnyExpression[string](targetExp, pv)
		if err != nil {
			return rError, fmt.Errorf("EXP[arguments/target_var]: %w", err)
		}

		// object is map
		if pExists {
			targetMap, ok := pv.Variables[target].(map[string]any)
			if !ok {
				return rError, fmt.Errorf("Variable `target` is not a map!")
			}

			path, err := evalAnyExpression[[]any](raw_path, pv)
			if err != nil {
				return rError, fmt.Errorf("EXP[path]: %w", err)
			}

			castPath, err := util.CastSlice[string](path)
			if err != nil {
				return rError, fmt.Errorf("EXP[path]: %w", err)
			}

			ok = util.SetPath(targetMap, castPath, value)
			if !ok {
				return rError, fmt.Errorf("Error setting value of variable %s at %s!", target, strings.Join(castPath, "/"))
			}

			return rDone, nil
		}

		// object is array
		if iExists {
			targetArray, ok := pv.Variables[target].([]any)
			if !ok {
				return rError, fmt.Errorf("Variable `target` is not an array!")
			}

			fIndex, err := evalAnyExpression[float64](raw_index, pv)
			if err != nil {
				return rError, fmt.Errorf("EXP[index]: %w", err)
			}

			if !(fIndex == math.Trunc(fIndex)) {
				return rError, fmt.Errorf("EXP[index]: value is not an integer!")
			}

			index := int(fIndex)

			if index < 0 || index >= len(targetArray) {
				return rError, fmt.Errorf("EXP[index]: index out of range!")
			}

			targetArray[index] = value

			return rDone, nil
		}

		return rError, fmt.Errorf("This error should be unreachable. How did you get here?")

	case "append":
		valueExp, exists := arguments["value"]
		if !exists {
			return rError, fmt.Errorf("`arguments/value` expression missing!")
		}

		value, err := evalAnyExpression[any](valueExp, pv)
		if err != nil {
			return rError, fmt.Errorf("EXP[arguments/value]: %w", err)
		}

		targetExp, exists := arguments["target_var"]
		if !exists {
			return rError, fmt.Errorf("`arguments/target_var` expression missing!")
		}

		target, err := evalAnyExpression[string](targetExp, pv)
		if err != nil {
			return rError, fmt.Errorf("EXP[arguments/target_var]: %w", err)
		}

		targetArray, ok := pv.Variables[target].([]any)
		if !ok {
			return rError, fmt.Errorf("Variable `target` is not an array!")
		}

		targetArray = append(targetArray, value)

		pv.Variables[target] = targetArray

		return rDone, nil

	case "if":
		conditionExp, exists := arguments["condition"]
		if !exists {
			return rError, fmt.Errorf("`arguments/condition` expression missing!")
		}

		condition, err := evalAnyExpression[bool](conditionExp, pv)
		if err != nil {
			return rError, fmt.Errorf("EXP[arguments/condition]: %w", err)
		}

		rawThen, exists := arguments["then"]
		if !exists {
			return rError, fmt.Errorf("`arguments/then` expression missing!")
		}

		thenOperations, ok := rawThen.([]any)
		if !ok {
			return rError, fmt.Errorf("`arguments/then` is not a list of expressions!")
		}

		rawElse, exists := arguments["else"]
		if !exists {
			return rError, fmt.Errorf("`arguments/else` expression missing!")
		}

		elseOperations, ok := rawElse.([]any)
		if !ok {
			return rError, fmt.Errorf("`arguments/else` is not a list of expressions!")
		}

		if condition {
			for _, op := range thenOperations {
				_, err := evalAnyOperation(op, pv)
				if err != nil {
					return rError, err
				}
			}
		} else {
			for _, op := range elseOperations {
				_, err := evalAnyOperation(op, pv)
				if err != nil {
					return rError, err
				}
			}
		}

		return rDone, nil

	case "while":
		conditionExp, exists := arguments["condition"]
		if !exists {
			return rError, fmt.Errorf("`arguments/condition` expression missing!")
		}

		rawOperations, exists := arguments["operations"]
		if !exists {
			return rError, fmt.Errorf("`arguments/operations` expression missing!")
		}

		operations, ok := rawOperations.([]any)
		if !ok {
			return rError, fmt.Errorf("`arguments/operations` is not a list of expressions!")
		}

		for {
			condition, err := evalAnyExpression[bool](conditionExp, pv)
			if err != nil {
				return rError, fmt.Errorf("EXP[arguments/condition]: %w", err)
			}

			if !condition {
				break
			}

			for _, op := range operations {
				_, err := evalAnyOperation(op, pv)
				if err != nil {
					return rError, err
				}
			}
		}

		return rDone, nil

	case "for_each":
		inExp, exists := arguments["in"]
		if !exists {
			return rError, fmt.Errorf("`arguments/in` expression missing!")
		}

		in, err := evalAnyExpression[[]any](inExp, pv)
		if err != nil {
			return rError, fmt.Errorf("EXP[arguments/in]: %w", err)
		}

		rawOperations, exists := arguments["operations"]
		if !exists {
			return rError, fmt.Errorf("`arguments/operations` expression missing!")
		}

		operations, ok := rawOperations.([]any)
		if !ok {
			return rError, fmt.Errorf("`arguments/operations` is not a list of expressions!")
		}

		for index, item := range in {
			newPv := parentVariables{
				Inputs:    pv.Inputs,
				Variables: pv.Variables,
				Other: map[string]any{
					"loop": map[string]any{
						"index": index,
						"item":  item,
					},
				},
			}

			for _, op := range operations {
				_, err := evalAnyOperation(op, newPv)
				if err != nil {
					return rError, err
				}
			}
		}

		return rDone, nil

	default:
		return rError, fmt.Errorf("Unknown operation type: %s", o.OType)
	}
}

func evalAnyOperation(rawOp any, pv parentVariables) (returnReason, error) {
	mapOp, ok := rawOp.(map[string]any)
	if !ok {
		return rError, fmt.Errorf("Operation has the wrong type!")
	}

	op, err := util.MapToStruct[operation](mapOp, "structura")
	if err != nil {
		return rError, err
	}

	return op.evaluate(pv)
}
