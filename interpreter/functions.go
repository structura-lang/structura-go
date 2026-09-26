package interpreter

import (
	"fmt"
	"structura/interpreter/util"
)

type function struct {
	FRuntime    RuntimeFunction
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

func (f function) evaluate(rd *runtimeData, inputs map[string]any) (any, error) {
	var iInputs map[string]any = map[string]any{}

	// validate inputs
	for _, inputName := range f.FInputs {
		input, hasValue := inputs[inputName]
		if !hasValue {
			return nil, fmt.Errorf("Input %s missing!", inputName)
		}

		iInputs[inputName] = input
	}

	if f.FRuntime != nil {
		if (f.FVariables != nil || f.FOperations != nil || f.FOutput != expression{}) {
			return nil, fmt.Errorf("Runtime functions cannot have `variables`, `operations`, or `output` set!")
		}

		return f.FRuntime(util.MakeCopy(iInputs)) // copy input, in case function mutates it
	} else {
		pv := parentVariables{
			Inputs:    iInputs,
			Variables: map[string]any{},
		}

		var iVariables map[string]any = map[string]any{}

		for varName, varExp := range f.FVariables {
			eval, err := varExp.evaluate(rd, pv)
			if err != nil {
				return nil, fmt.Errorf("EXP[var:%s]: %w", varName, err)
			}

			iVariables[varName] = util.MakeCopy(eval)
		}

		pv.Variables = iVariables

		for i, operation := range f.FOperations {
			reason, err := operation.evaluate(rd, pv)
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

		eval, err := f.FOutput.evaluate(rd, pv)
		if err != nil {
			return nil, fmt.Errorf("EXP[out]: %w", err)
		}

		return eval, nil
	}
}

func evalAnyFunction(fn any, rd *runtimeData, inputs map[string]any) (any, error) {
	if f, ok := fn.(function); ok {
		return f.evaluate(rd, inputs)
	} else {
		mapFn, ok := fn.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("Unsupported type %T, expected function or map[string]any", fn)
		}

		f, err := util.MapToStruct[function](mapFn, "structura")
		if err != nil {
			return nil, fmt.Errorf("MapToStruct failed: %w", err)
		}

		return f.evaluate(rd, inputs)
	}
}
