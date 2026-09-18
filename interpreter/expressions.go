package interpreter

import "fmt"

type Expression struct {
	EType string `json:"type"`
	EData any    `json:"data"`
}

type parentVariables struct {
	inputs    map[string]any
	variables map[string]any
}

func (e Expression) evaluate(pv parentVariables) (any, error) {
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
	default:
		return nil, fmt.Errorf("Unknown expression type: %s", e.EType)
	}
}
