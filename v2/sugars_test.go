package v2_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v2 "github.com/ysmood/gisp/v2"
)

func TestValueStack(t *testing.T) {
	fn := func(ctx *v2.Context) (interface{}, *v2.Error) {
		_, err := ctx.ArgNum(0)
		if err != nil {
			return nil, err
		}
		_, err = ctx.ArgNum(1)
		if err != nil {
			return nil, err
		}
		return nil, ctx.Error("err")
	}

	_, err := v2.Run(&v2.Context{
		AST: []interface{}{
			"bar", []interface{}{
				"foo", []interface{}{
					"bar", float64(1), 1,
				},
				1,
			},
			1,
		},
		Sandbox: v2.Sandbox{
			"bar": fn,
			"foo": fn,
		},
	})
	assert.Equal(t, "arg[1] not a number\nbar:0\nfoo:0\nbar:0\n", err.String())
}

func TestTypeErr(t *testing.T) {
	_, err := v2.Run(&v2.Context{
		AST: []interface{}{"foo", 1},
		Sandbox: v2.Sandbox{
			"foo": func(ctx *v2.Context) (interface{}, *v2.Error) {
				return ctx.ArgStr(0)
			},
		},
	})

	assert.Equal(t, "arg[0] not a string", err.Details)
}
