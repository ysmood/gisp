package gisp

import (
	"encoding/json"
)

// Context context
type Context struct {
	// AST is simply a json data
	AST interface{}

	// The functions exposed to the vm
	Sandbox *Sandbox

	// The state exposed to the functions in the vm
	// It's not directly visible to user.
	ENV interface{}

	// The index of parent context
	Index int

	// Parent AST
	Parent *Context

	// Whether auto lift sandbox panic with informal stack info or not
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
		if ctx.IsLiftPanic {
			defer ctx.liftPanic()
		}

		if ctx.Len() == 0 {
			return nil
		}

		val := ctx.Arg(0)

		// if val is function
		switch val.(type) {
		case func(*Context) interface{}:
			return val.(func(*Context) interface{})(ctx)

		// if val is string
		default:
			name, isStr := val.(string)
			var has bool

			if isStr {
				val, has = ctx.Sandbox.Get(name)
			}

			if has {
				switch val.(type) {
				case func(*Context) interface{}:
					return val.(func(*Context) interface{})(ctx)
				default:
					return val
				}
			}

			msg, _ := json.Marshal(ctx.AST.([]interface{})[0])
			ctx.Error("function " + string(msg) + " is undefined")
		}
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
