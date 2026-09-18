package interpreter

import (
	"fmt"
)

type Function struct {
	FInputs     []string              `json:"inputs"`
	FVariables  map[string]Expression `json:"variables"`
	FOperations []Operation           `json:"operations"`
	FOutput     Expression            `json:"output"`
}

func (f Function) evaluate(inputs map[string]any) (any, error) {
	// validate inputs
	for _, inputName := range f.FInputs {
		_, hasValue := inputs[inputName]
		if !hasValue {
			return nil, fmt.Errorf("Input %s missing!", inputName)
		}
	}

	var iVariables map[string]any = map[string]any{}

	pvVar := parentVariables{
		inputs:    inputs,
		variables: iVariables,
	}

	for varName, varExp := range f.FVariables {
		eval, err := varExp.evaluate(pvVar)
		if err != nil {
			return nil, fmt.Errorf("EXP[var:%s]: %w", varName, err)
		}

		iVariables[varName] = eval
	}

	//	for _, operation := range f.FOperations {
	//		// TODO
	//	}

	pvOut := parentVariables{
		inputs:    inputs,
		variables: iVariables,
	}

	eval, err := f.FOutput.evaluate(pvOut)
	if err != nil {
		return nil, fmt.Errorf("EXP[out]: %w", err)
	}

	return eval, nil
}
