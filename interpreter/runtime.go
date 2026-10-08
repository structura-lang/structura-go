package interpreter

type RuntimeFunction func(inputs map[string]any) (any, error)

type RuntimeData struct {
	values    map[string]any
	functions map[string]function
}

func NewRuntime() *RuntimeData {
	return &RuntimeData{
		values:    map[string]any{},
		functions: map[string]function{},
	}
}

func (r *RuntimeData) RegisterValue(name string, value any) {
	r.values[name] = value
}

func (r *RuntimeData) RegisterFunction(name string, fn RuntimeFunction, inputs []string) {
	r.functions[name] = function{
		FRuntime: fn,
		FInputs:  inputs,
	}
}
