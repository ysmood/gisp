package main

import (
	"fmt"

	"github.com/ysmood/gisp"
	"github.com/ysmood/gisp/lib"
)

func main() {
	//  arithmetic expression: 1 + (2 * 3)
	//    function expression: add(1, multiply(2, 3))
	// gisp (json) expression: ["+", 1, ["*", 2, 3]]
	code := `["+", 1, ["*", 2, 3]]`

	out, _ := gisp.RunJSON(code, &gisp.Context{
		Sandbox: gisp.New(gisp.Box{
			"+": lib.Add,

			// define custom function
			"*": func(ctx *gisp.Context) interface{} {
				arg1 := ctx.ArgNum(1) // get the first argument as number
				arg2 := ctx.ArgNum(2) // get the second argument as number
				return arg1 * arg2    // 2 * 3
			},
		}),
	})

	fmt.Println(out) // print 7
}
