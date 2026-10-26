package data

import "structura/interpreter"

func RegisterData(rd *interpreter.RuntimeData) {
	rd.RegisterFunction("length", length, []string{"arr"})
	rd.RegisterFunction("equal", equal, []string{"a", "b"})
}
