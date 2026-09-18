package interpreter

import (
	"fmt"
	"structura/util"
)

func Execute(program map[string]any, inputs map[string]any) (any, error) {
	function, err := util.MapToStruct[function](program, "structura")
	if err != nil {
		return nil, err
	}

	fmt.Printf("%+v\n", function)

	return function.evaluate(inputs)
}
