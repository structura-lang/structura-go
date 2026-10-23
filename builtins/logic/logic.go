package logic

import "structura/interpreter"

func RegisterLogic(rd *interpreter.RuntimeData) {
	rd.RegisterFunction("and", and, []string{"a", "b"})
	rd.RegisterFunction("or", or, []string{"a", "b"})
	rd.RegisterFunction("not", not, []string{"a"})
}
