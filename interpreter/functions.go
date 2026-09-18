package interpreter

import (
	"fmt"
	"structura/util"
)

type function struct {
	FInputs     []string              `structura:"inputs"`
	FVariables  map[string]expression `structura:"variables"`
	FOperations []operation           `structura:"operations"`
	FOutput     expression            `structura:"output"`
}

func (f function) evaluate(inputs map[string]any) (any, error) {
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

func evalMapFunction(fn map[string]any, inputs map[string]any) (any, error) {
	function, err := util.MapToStruct[function](fn, "structura")
	if err != nil {
		return nil, err
	}

	return function.evaluate(inputs)
}
