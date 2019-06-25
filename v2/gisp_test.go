package v2_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v2 "github.com/ysmood/gisp/v2"
)

func TestBasic(t *testing.T) {
	out, _ := v2.Run(&v2.Context{
		AST: []interface{}{"mul", 3, []interface{}{"add", 1, 2}},
		Sandbox: v2.Sandbox{
			"add": func(ctx *v2.Context) (interface{}, *v2.Error) {
				a, _ := ctx.Arg(0)
				b, _ := ctx.Arg(1)

				return a.(int) + b.(int), nil
			},
			"mul": func(ctx *v2.Context) (interface{}, *v2.Error) {
				a, _ := ctx.Arg(0)
				b, _ := ctx.Arg(1)

				return a.(int) * b.(int), nil
			},
		},
	})

	assert.Equal(t, 9, out)
}

func TestEmpty(t *testing.T) {
	out, _ := v2.Run(&v2.Context{
		AST: []interface{}{},
	})

	assert.Nil(t, out)
}

func TestErrNameNotString(t *testing.T) {
	_, err := v2.Run(&v2.Context{
		AST: []interface{}{1},
	})

	assert.EqualError(t, err, "function name is not string")
}

func TestErrNotArray(t *testing.T) {
	ctx := &v2.Context{AST: "foo"}
	_, err := ctx.Arg(0)

	assert.EqualError(t, err, "ast is not an array")
}

func TestErrArgNotDefined(t *testing.T) {
	ctx := &v2.Context{AST: []interface{}{"test"}}

	_, err := ctx.Arg(2)

	assert.EqualError(t, err, "argument not defined")
}

func TestErr(t *testing.T) {
	_, err := v2.Run(&v2.Context{
		AST:     []interface{}{"test"},
		Sandbox: v2.Sandbox{},
	})

	assert.EqualError(t, err, "function not defined")
}

func TestFnErr(t *testing.T) {
	_, err := v2.Run(&v2.Context{
		AST: []interface{}{"add", 3, []interface{}{"err"}},
		Sandbox: v2.Sandbox{
			"add": func(ctx *v2.Context) (interface{}, *v2.Error) {
				a, _ := ctx.Arg(0)
				b, err := ctx.Arg(1)
				if err != nil {
					return nil, err
				}

				return a.(int) + b.(int), nil
			},
			"err": func(ctx *v2.Context) (interface{}, *v2.Error) {
				return nil, ctx.Error("err")
			},
		},
	})

	assert.EqualError(t, err, "function error")
	assert.Equal(t, "err", err.Details)
	assert.Equal(t, 1, err.Context.Location)
}
