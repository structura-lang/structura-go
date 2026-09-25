package interpreter

type RuntimeFunction func(inputs map[string]any) (any, error)

type runtimeData struct {
	values    map[string]any
	functions map[string]function
}

func NewRuntime() *runtimeData {
	return &runtimeData{
		values:    map[string]any{},
		functions: map[string]function{},
	}
}

func (r *runtimeData) RegisterValue(name string, value any) {
	r.values[name] = value
}

func (r *runtimeData) RegisterFunction(name string, fn RuntimeFunction, inputs []string) {
	r.functions[name] = function{
		FRuntime: fn,
		FInputs:  inputs,
	}
}
