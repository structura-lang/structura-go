package math

import "structura/interpreter"

func RegisterMath(rd *interpreter.RuntimeData) {
	rd.RegisterFunction("add", add, []string{"a", "b"})
	rd.RegisterFunction("subract", subtract, []string{"a", "b"})
	rd.RegisterFunction("multiply", multiply, []string{"a", "b"})
	rd.RegisterFunction("divide", divide, []string{"a", "b"})
	rd.RegisterFunction("lessThan", lessThan, []string{"a", "b"})
	rd.RegisterFunction("greaterThan", greaterThan, []string{"a", "b"})
	rd.RegisterFunction("equal", equal, []string{"a", "b"})
}
