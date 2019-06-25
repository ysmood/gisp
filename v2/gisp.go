package v2

// Func ...
type Func func(*Context) (interface{}, *Error)

// Sandbox ...
type Sandbox map[string]Func

// Context ...
type Context struct {
	// AST is simply a json data
	AST interface{}

	Sandbox Sandbox

	Location int
	Parent   *Context
}

// Run ...
func Run(ctx *Context) (interface{}, *Error) {
	ast, ok := ctx.AST.([]interface{})
	if !ok {
		return ctx.AST, nil
	}

	if len(ast) == 0 {
		return nil, nil
	}

	name, ok := ast[0].(string)
	if !ok {
		return nil, ctx.error(ErrNameNotString, 0)
	}

	fn, has := ctx.Sandbox[name]
	if has {
		return fn(ctx)
	}
	return nil, ctx.error(ErrNotDefined, 0)
}

// Arg ...
func (ctx *Context) Arg(index int) (interface{}, *Error) {
	ast, ok := ctx.AST.([]interface{})
	if !ok {
		return nil, ctx.error(ErrNotArray, index)
	}

	i := index + 1
	if i >= len(ast) || i < 0 {
		return nil, ctx.error(ErrArgNotDefined, index)
	}

	return Run(&Context{
		AST:      ast[i],
		Sandbox:  ctx.Sandbox,
		Location: index,
		Parent:   ctx,
	})
}

func (ctx *Context) Error(details interface{}) *Error {
	return &Error{
		Code:    ErrFromFunction,
		Context: ctx,
		Details: details,
	}
}

func (ctx *Context) error(code ErrorCode, loc int) *Error {
	return &Error{
		Code:     code,
		Context:  ctx,
		Location: loc,
	}
}

// ErrorCode ...
type ErrorCode int

const (
	// ErrNotDefined ...
	ErrNotDefined ErrorCode = iota
	// ErrNameNotString ...
	ErrNameNotString
	// ErrNotArray ...
	ErrNotArray
	// ErrArgNotDefined ...
	ErrArgNotDefined
	// ErrFromFunction ...
	ErrFromFunction
)

func (c ErrorCode) String() string {
	switch c {
	case ErrNotDefined:
		return "function not defined"
	case ErrNameNotString:
		return "function name is not string"
	case ErrNotArray:
		return "ast is not an array"
	case ErrArgNotDefined:
		return "argument not defined"
	default:
		return "function error"
	}
}

// Error ...
type Error struct {
	Context  *Context
	Code     ErrorCode
	Location int
	Details  interface{}
}

func (e *Error) Error() string {
	return e.Code.String()
}
