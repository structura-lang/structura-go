package interpreter

type runtimeData struct {
	values map[string]any
}

func NewRuntime() *runtimeData {
	return &runtimeData{
		values: map[string]any{},
	}
}

func (r runtimeData) RegisterValue(name string, value any) {
	r.values[name] = value
}
