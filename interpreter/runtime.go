package interpreter

import "fmt"

func Execute(program Function, inputs map[string]any) (any, error) {
	fmt.Printf("%+v\n", program)

	return program.evaluate(inputs)
}
