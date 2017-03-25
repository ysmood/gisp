package gisp

import "encoding/json"

// Sandbox sandbox
type Sandbox map[string]func(*Context) interface{}

// Context context
type Context struct {
	AST         interface{}
	Sandbox     Sandbox
	ENV         interface{}
	Index       int
	Parent      *Context
	IsLiftPanic bool
}

// Error ...
type Error struct {
	Message string
	Stack   []interface{}
}

func (e Error) Error() string {
	return e.Message
}

func (ctx *Context) liftPanic() {
	if r := recover(); r != nil {
		err, ok := r.(Error)
		if ok {
			panic(err)
		} else {
			ctx.Error(r.(error).Error())
		}
	}
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
				ctx.Error("function \"" + name + "\" is undefined")
			}
			if ctx.IsLiftPanic {
				defer ctx.liftPanic()
			}
			return fn(ctx)
		}

		if ctx.IsLiftPanic {
			defer ctx.liftPanic()
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

	panic(Error{
		Message: msg,
		Stack:   stack,
	})
}

// RunJSON json entrance
func RunJSON(code []byte, ctx *Context) (ret interface{}, err error) {
	var ast interface{}
	err = json.Unmarshal(code, &ast)
	ctx.AST = ast
	ret = Run(ctx)
	return
}

// Arg sugar
func (ctx *Context) Arg(index int) interface{} {
	ast := ctx.AST.([]interface{})

	if index >= len(ast) {
		return nil
	}

	return Run(&Context{
		AST:         ast[index],
		Sandbox:     ctx.Sandbox,
		ENV:         ctx.ENV,
		Index:       index,
		Parent:      ctx,
		IsLiftPanic: ctx.IsLiftPanic,
	})
}

// Len Arg sugar
func (ctx *Context) Len() int {
	return len(ctx.AST.([]interface{}))
}
