package interpreter

func Execute(program map[string]any, inputs map[string]any) (any, error) {
	return evalAnyFunction(program, inputs)
}
