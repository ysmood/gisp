package gisp

// Sandbox sandbox
type Sandbox map[string]func(*Context) interface{}

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
