package gisp

import "encoding/json"

// Sandbox sandbox
type Sandbox map[string]func(*Context) interface{}

// Context context
type Context struct {
	AST     interface{}
	Sandbox Sandbox
	ENV     interface{}
	Index   int
	Parent  *Context
}

type GispError struct {
	Message string
	Stack   []interface{}
}

func (e GispError) Error() string {
	return e.Message
}

// Run entrance
func Run(ctx *Context) interface{} {
	switch ctx.AST.(type) {
	case []interface{}:
		action := ctx.Arg(0)

		switch action.(type) {
		case string:
			name := action.(string)
			fn := ctx.Sandbox[name]
			if fn == nil {
				ctx.Error("function is undefined: " + name)
			}
			return fn(ctx)
		}

		return action.(func(*Context) interface{})(ctx)
	}

	return ctx.AST
}

// Error used to throw error
func (ctx *Context) Error(msg string) {
	stack := []interface{}{}
	node := ctx

	for node != nil {
		name := node.AST.([]interface{})[0]
		stack = append(stack, name, node.Index)
		node = node.Parent
	}

	panic(GispError{
		Message: msg,
		Stack:   stack,
	})
}

// RunJSON json entrance
func RunJSON(code []byte, ctx *Context) (ret interface{}, err error) {
	var ast interface{}
	err = json.Unmarshal(code, &ast)

	ret = Run(&Context{
		AST:     ast,
		Sandbox: ctx.Sandbox,
		ENV:     ctx.ENV,
	})

	return
}

// Arg sugar
func (ctx *Context) Arg(index int) interface{} {
	return Run(&Context{
		AST:     ctx.AST.([]interface{})[index],
		Sandbox: ctx.Sandbox,
		ENV:     ctx.ENV,
		Index:   index,
		Parent:  ctx,
	})
}

// Arg sugar
func (ctx *Context) Len() int {
	return len(ctx.AST.([]interface{}))
}
