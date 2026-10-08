package builtins

import (
	d "structura/builtins/data"
	m "structura/builtins/math"
	"structura/interpreter"
)

func RegisterAll(rd *interpreter.RuntimeData) {
	m.RegisterMath(rd)
	d.RegisterData(rd)
}
