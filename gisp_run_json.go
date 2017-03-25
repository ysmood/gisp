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
