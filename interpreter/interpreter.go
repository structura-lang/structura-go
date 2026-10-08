package interpreter

func Execute(program map[string]any, inputs map[string]any, rd *RuntimeData) (any, error) {
	if rd == nil {
		rd = NewRuntime()
	}

	return evalAnyFunction(program, rd, inputs)
}
