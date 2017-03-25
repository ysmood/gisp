package gisp

import "encoding/json"

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

// ArgNum Get argument as number
func (ctx *Context) ArgNum(index int) float64 {
	return ctx.Arg(index).(float64)
}

// ArgStr Get argument as string
func (ctx *Context) ArgStr(index int) string {
	return ctx.Arg(index).(string)
}

// ArgBool Get argument as bool
func (ctx *Context) ArgBool(index int) bool {
	return ctx.Arg(index).(bool)
}

// ArgObj Get argument as object
func (ctx *Context) ArgObj(index int) map[string]interface{} {
	return ctx.Arg(index).(map[string]interface{})
}

// ArgArr Get argument as array
func (ctx *Context) ArgArr(index int) []interface{} {
	return ctx.Arg(index).([]interface{})
}

// Len Arg sugar
func (ctx *Context) Len() int {
	return len(ctx.AST.([]interface{}))
}
