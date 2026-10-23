package builtins

import (
	d "structura/builtins/data"
	l "structura/builtins/logic"
	m "structura/builtins/math"
	"structura/interpreter"
)

func RegisterAll(rd *interpreter.RuntimeData) {
	m.RegisterMath(rd)
	d.RegisterData(rd)
	l.RegisterLogic(rd)
}
