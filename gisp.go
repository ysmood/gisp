package gisp

import (
	"encoding/json"
)

// Sandbox sandbox
type Sandbox map[string]interface{}

// Context context
type Context struct {
	// AST is simply a json data
	AST interface{}

	// The functions exposed to the vm
	Sandbox Sandbox

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

		switch val.(type) {
		case string:
			var has bool
			name := val.(string)
			val, has = ctx.Sandbox[name]

			if !has {
				ctx.Error("\"" + name + "\" is undefined")
			}
		default:
			msg, _ := json.Marshal(ctx.AST.([]interface{})[0])
			ctx.Error(string(msg) + " should return function")

		}

		switch val.(type) {
		case func(*Context) float64:
			return val.(func(*Context) float64)(ctx)
		case func(*Context) string:
			return val.(func(*Context) string)(ctx)
		case func(*Context) bool:
			return val.(func(*Context) bool)(ctx)
		case func(*Context) map[string]interface{}:
			return val.(func(*Context) map[string]interface{})(ctx)
		case func(*Context) []interface{}:
			return val.(func(*Context) []interface{})(ctx)
		case func(*Context) interface{}:
			return val.(func(*Context) interface{})(ctx)
		}

		return val
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
