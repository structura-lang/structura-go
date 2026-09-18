package interpreter

import (
	"fmt"
	"structura/util"
)

type expression struct {
	EType string `structura:"type"`
	EData any    `structura:"data"`
}

type parentVariables struct {
	inputs    map[string]any
	variables map[string]any
}

func (e expression) evaluate(pv parentVariables) (any, error) {
	switch e.EType {
	case "literal":
		return e.EData, nil

	case "variable":
		varName, ok := e.EData.(string)
		if !ok {
			return nil, fmt.Errorf("`data` field is not a string!")
		}

		variable, exists := pv.variables[varName]
		if !exists {
			return nil, fmt.Errorf("Variable %s is not defined!", varName)
		}

		return variable, nil

	case "input":
		inputName, ok := e.EData.(string)
		if !ok {
			return nil, fmt.Errorf("`data` field is not a string!")
		}

		input, exists := pv.inputs[inputName]
		if !exists {
			return nil, fmt.Errorf("Input %s is not defined!", inputName)
		}

		return input, nil

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
			return nil, err
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
				return nil, err
			}

			inputs[input] = res
		}

		res, err := evalMapFunction(fn, inputs)
		if err != nil {
			return nil, fmt.Errorf("FUNC[function]: %w", err)
		}

		return res, nil

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
