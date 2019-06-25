package v2

import (
	"fmt"

	"github.com/a8m/djson"
)

// Stack ...
func (e *Error) Stack() []*Context {
	stack := []*Context{}
	node := e.Context

	for node != nil {
		switch node.AST.(type) {
		case []interface{}:
			stack = append(stack, node)
			node = node.Parent
		default:
			stack = append(stack, node)
			node = node.Parent
		}
	}

	return stack
}

// StackString ...
func (e *Error) String() string {
	out := ""
	if e.Details == nil {
		out += fmt.Sprintf("%s:%d\n", e.Error(), e.Location)
	} else {
		out += fmt.Sprintf("%v\n", e.Details)
	}
	for _, s := range e.Stack() {
		ast, ok := s.AST.([]interface{})
		name := s.AST
		if ok {
			name = ast[0]
		}
		out += fmt.Sprintf("%v:%d\n", name, s.Location)
	}

	return out
}

// RunJSON json entrance
func RunJSON(code string, ctx *Context) (interface{}, error) {
	return RunJSONRaw([]byte(code), ctx)
}

// RunJSONRaw json entrance
func RunJSONRaw(code []byte, ctx *Context) (ret interface{}, err error) {
	ast, err := djson.Decode(code)
	ctx.AST = ast
	return Run(ctx)
}

// ArgNum Get argument as number
func (ctx *Context) ArgNum(index int) (float64, *Error) {
	arg, err := ctx.Arg(index)
	if err != nil {
		return 0, err
	}
	val, ok := arg.(float64)
	if !ok {
		return 0, ctx.Error(fmt.Sprintf("arg[%d] not a number", index))
	}
	return val, nil
}

// ArgStr Get argument as string
func (ctx *Context) ArgStr(index int) (string, *Error) {
	arg, err := ctx.Arg(index)
	if err != nil {
		return "", err
	}
	val, ok := arg.(string)
	if !ok {
		return "", ctx.Error(fmt.Sprintf("arg[%d] not a string", index))
	}
	return val, nil
}

// ArgBool Get argument as bool
func (ctx *Context) ArgBool(index int) (bool, *Error) {
	arg, err := ctx.Arg(index)
	if err != nil {
		return false, err
	}
	val, ok := arg.(bool)
	if !ok {
		return false, ctx.Error(fmt.Sprintf("arg[%d] not a boolean", index))
	}
	return val, nil

}

// ArgObj Get argument as object
func (ctx *Context) ArgObj(index int) (map[string]interface{}, *Error) {
	arg, err := ctx.Arg(index)
	if err != nil {
		return nil, err
	}
	val, ok := arg.(map[string]interface{})
	if !ok {
		return nil, ctx.Error(fmt.Sprintf("arg[%d] not a object", index))
	}
	return val, nil
}

// ArgArr Get argument as array
func (ctx *Context) ArgArr(index int) ([]interface{}, *Error) {
	arg, err := ctx.Arg(index)
	if err != nil {
		return nil, err
	}
	val, ok := arg.([]interface{})
	if !ok {
		return nil, ctx.Error(fmt.Sprintf("arg[%d] not a array", index))
	}
	return val, nil
}

// Len Arg sugar
func (ctx *Context) Len() int {
	return len(ctx.AST.([]interface{}))
}
