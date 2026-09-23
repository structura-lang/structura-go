package interpreter

import (
	"fmt"
	"math"
	"strings"
	"structura/interpreter/util"
)

type expression struct {
	EType string `structura:"type"`
	EData any    `structura:"data"`
}

func (e expression) evaluate(rpv parentVariables) (any, error) {
	pv := util.MakeCopy(rpv)

	switch e.EType {
	case "literal":
		return e.EData, nil

	case "variable":
		varName, ok := e.EData.(string)
		if !ok {
			return nil, fmt.Errorf("`data` field is not a string!")
		}

		variable, exists := pv.Variables[varName]
		if !exists {
			return nil, fmt.Errorf("Variable %s is not defined!", varName)
		}

		return variable, nil

	case "input":
		inputName, ok := e.EData.(string)
		if !ok {
			return nil, fmt.Errorf("`data` field is not a string!")
		}

		input, exists := pv.Inputs[inputName]
		if !exists {
			return nil, fmt.Errorf("Input %s is not defined!", inputName)
		}

		return input, nil

	case "get":
		data, ok := e.EData.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("`data` field is not a map!")
		}

		raw_from, exists := data["from"]
		if !exists {
			return nil, fmt.Errorf("`from` expression missing!")
		}

		from, err := evalAnyExpression[any](raw_from, pv)
		if err != nil {
			return nil, fmt.Errorf("EXP[from]: %w", err)
		}

		raw_path, pExists := data["path"]
		raw_index, iExists := data["index"]

		if pExists && iExists {
			return nil, fmt.Errorf("`path` and `index` are mutually exclusive!")
		}

		if !pExists && !iExists {
			return nil, fmt.Errorf("`path` or `index` expressions missing!")
		}

		// object is map
		if pExists {
			fromMap, ok := from.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("`from` expression returned invalid type!")
			}

			path, err := evalAnyExpression[[]any](raw_path, pv)
			if err != nil {
				return nil, fmt.Errorf("EXP[path]: %w", err)
			}

			castPath, err := util.CastSlice[string](path)
			if err != nil {
				return nil, fmt.Errorf("EXP[path]: %w", err)
			}

			result, ok := util.GetPath(fromMap, castPath)
			if !ok {
				return nil, fmt.Errorf("Failed getting object at path `%s`", strings.Join(castPath, "/"))
			}

			return result, nil
		}

		// object is array
		if iExists {
			fromArray, ok := from.([]any)
			if !ok {
				return nil, fmt.Errorf("`from` expression returned invalid type!")
			}

			fIndex, err := evalAnyExpression[float64](raw_index, pv)
			if err != nil {
				return nil, fmt.Errorf("EXP[index]: %w", err)
			}

			if !(fIndex == math.Trunc(fIndex)) {
				return nil, fmt.Errorf("EXP[index]: value is not an integer!")
			}

			index := int(fIndex)

			if index < 0 || index >= len(fromArray) {
				return nil, fmt.Errorf("EXP[index]: value out of bounds!")
			}

			return fromArray[index], nil
		}

		return nil, fmt.Errorf("This error should be unreachable. How did you get here?")

	case "call":
		data, ok := e.EData.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("`data` field is not a map!")
		}

		rawFunction, exists := data["function"]
		if !exists {
			return nil, fmt.Errorf("`function` expression missing!")
		}

		fn, err := evalAnyExpression[map[string]any](rawFunction, pv)
		if err != nil {
			return nil, fmt.Errorf("EXP[function]: %w", err)
		}

		rawInputExpressions, exists := data["inputs_from"]
		if !exists {
			return nil, fmt.Errorf("`inputs_from` expression missing!")
		}

		inputExpressions, ok := rawInputExpressions.(map[string]any)

		var inputs map[string]any = map[string]any{}

		for input, exp := range inputExpressions {
			res, err := evalAnyExpression[any](exp, pv)
			if err != nil {
				return nil, fmt.Errorf("EXP[inputs_from]: %w", err)
			}

			inputs[input] = res
		}

		res, err := evalMapFunction(fn, inputs)
		if err != nil {
			return nil, fmt.Errorf("FUNC[function]: %w", err)
		}

		return res, nil

	case "loop":
		loopObj, ok := e.EData.(string)
		if !ok {
			return nil, fmt.Errorf("`data` field is not a string!")
		}

		rawLoopData, exists := pv.Other["loop"]
		if !exists {
			return nil, fmt.Errorf("Cannot evaluate `loop` expression outside loop!")
		}

		loopData, ok := rawLoopData.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("Invalid loop data!")
		}

		result, exists := loopData[loopObj]
		if !exists {
			return nil, fmt.Errorf("Loop object `%s` does not exist!", loopObj)
		}

		return result, nil

	default:
		return nil, fmt.Errorf("Unknown expression type: %s", e.EType)
	}
}

func evalAnyExpression[T any](rawExp any, pv parentVariables) (T, error) {
	var null T

	mapExp, ok := rawExp.(map[string]any)
	if !ok {
		return null, fmt.Errorf("Expression has the wrong type!")
	}

	exp, err := util.MapToStruct[expression](mapExp, "structura")
	if err != nil {
		return null, err
	}

	eval, err := exp.evaluate(pv)
	if err != nil {
		return null, err
	}

	result, ok := eval.(T)
	if !ok {
		return null, fmt.Errorf("Expression returned the wrong type!")
	}

	return result, nil
}
