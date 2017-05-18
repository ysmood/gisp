package gisp

import (
	"encoding/json"
)

// Box which contains the raw userspace values
type Box map[string]interface{}

// Sandbox sandbox
// It implements prototype design pattern
type Sandbox struct {
	dict   Box
	parent *Sandbox
}

// New create a new sandbox
func New(dict Box) *Sandbox {
	return &Sandbox{
		dict: dict,
	}
}

// Create create a new sandbox which dirives from current sandbox
func (sandbox *Sandbox) Create() *Sandbox {
	return &Sandbox{
		dict:   Box{},
		parent: sandbox,
	}
}

// Get get property from the prototype chain
func (sandbox *Sandbox) Get(name string) (interface{}, bool) {
	for sandbox != nil {
		val, has := sandbox.dict[name]

		if has {
			return val, true
		}

		sandbox = sandbox.parent
	}

	return nil, false
}

// Set set property
func (sandbox *Sandbox) Set(name string, val interface{}) {
	sandbox.dict[name] = val
}

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
		case func(*Context):
			val.(func(*Context))(ctx)
			return nil
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

		// if val is string
		default:
			name, isStr := val.(string)
			var has bool

			if isStr {
				val, has = ctx.Sandbox.Get(name)
			}

			if has {
				switch val.(type) {
				case func(*Context):
					val.(func(*Context))(ctx)
					return nil
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
