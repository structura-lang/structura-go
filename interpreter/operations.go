package interpreter

import "fmt"

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

		rawOutputTo, exists := arguments["output_to"]
		if !exists {
			return fmt.Errorf("`arguments/output_to` string missing!")
		}

		outputTo, ok := rawOutputTo.(string)
		if !ok {
			return fmt.Errorf("`arguments/output_to` field is not a string!")
		}

		pv.variables[outputTo] = value

		return nil

	default:
		return fmt.Errorf("Unknown operation type: %s", o.OType)
	}
}
