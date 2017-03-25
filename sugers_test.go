package gisp_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ysmood/gisp"
)

func TestJSON(t *testing.T) {
	sandbox := gisp.Sandbox{
		"+": func(ctx *gisp.Context) float64 {
			a := ctx.ArgNum(1)
			b := ctx.ArgNum(2)
			return a + b
		},
		"-": func(ctx *gisp.Context) float64 {
			a := ctx.ArgNum(1)
			b := ctx.ArgNum(2)
			return a - b
		},
	}

	out, _ := gisp.RunJSON(`["-", ["+", 5, 1], ["+", 1, 1]]`, &gisp.Context{
		Sandbox: sandbox,
	})

	assert.Equal(t, float64(4), out)
}

func TestTypes(t *testing.T) {
	sandbox := gisp.Sandbox{
		"$": func(g *gisp.Context) interface{} {
			return g.AST.([]interface{})[1]
		},
		"echo": func(g *gisp.Context) []interface{} {
			return []interface{}{
				g.ArgNum(1),
				g.ArgBool(2),
				g.ArgArr(3),
				g.ArgObj(4),
				g.ArgStr(5),
			}
		},
	}

	out, _ := gisp.RunJSON(`["echo", 1.2, true, ["$", []], {}, "ok"]`, &gisp.Context{
		Sandbox: sandbox,
	})

	assert.Equal(t, []interface{}{
		1.2,
		true,
		[]interface{}{},
		map[string]interface{}{},
		"ok",
	}, out)
}
