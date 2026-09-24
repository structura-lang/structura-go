package interpreter

import (
	"fmt"
	"structura/interpreter/util"
)

type function struct {
	FInputs     []string              `structura:"inputs"`
	FVariables  map[string]expression `structura:"variables"`
	FOperations []operation           `structura:"operations"`
	FOutput     expression            `structura:"output"`
}

type parentVariables struct {
	Inputs    map[string]any
	Variables map[string]any
	Other     map[string]any
}

func (f function) evaluate(inputs map[string]any) (any, error) {
	var iInputs map[string]any = map[string]any{}

	// validate inputs
	for _, inputName := range f.FInputs {
		input, hasValue := inputs[inputName]
		if !hasValue {
			return nil, fmt.Errorf("Input %s missing!", inputName)
		}

		iInputs[inputName] = input
	}

	inputsPV := parentVariables{
		Inputs:    util.MakeCopy(iInputs),
		Variables: map[string]any{},
	}

	var iVariables map[string]any = map[string]any{}

	for varName, varExp := range f.FVariables {
		eval, err := varExp.evaluate(inputsPV)
		if err != nil {
			return nil, fmt.Errorf("EXP[var:%s]: %w", varName, err)
		}

		iVariables[varName] = eval
	}

	pv := parentVariables{
		Inputs:    util.MakeCopy(iInputs),
		Variables: iVariables,
	}

	for i, operation := range f.FOperations {
		reason, err := operation.evaluate(pv)
		if err != nil {
			return nil, fmt.Errorf("OP[%d]: %w", i, err)
		}

		if reason == rReturn {
			break
		}

		if reason != rDone {
			return nil, fmt.Errorf("OP[%d]: Invalid return reason!", i)
		}
	}

	eval, err := f.FOutput.evaluate(pv)
	if err != nil {
		return nil, fmt.Errorf("EXP[out]: %w", err)
	}

	return eval, nil
}

func evalAnyFunction(fn any, inputs map[string]any) (any, error) {
	if f, ok := fn.(function); ok {
		return f.evaluate(inputs)
	} else {
		mapFn, ok := fn.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("Unsupported type %T, expected function or map[string]any", fn)
		}

		f, err := util.MapToStruct[function](mapFn, "structura")
		if err != nil {
			return nil, fmt.Errorf("MapToStruct failed: %w", err)
		}

		return f.evaluate(inputs)
	}
}
