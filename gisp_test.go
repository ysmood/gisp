package gisp_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ysmood/gisp"
	"github.com/yuin/gopher-lua"
)

func Add(L *lua.LState) int {
	a := L.ToInt(1)            /* get argument */
	b := L.ToInt(2)            /* get argument */
	L.Push(lua.LNumber(a + b)) /* push result */
	return 1                   /* number of results */
}

func BenchmarkLua(b *testing.B) {
	L := lua.NewState()
	defer L.Close()
	L.SetGlobal("add", L.NewFunction(Add))

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		L.DoString("add(1,1)")
	}
}

func TestJSON(t *testing.T) {
	sandbox := gisp.Sandbox{
		"+": func(ctx *gisp.Context) interface{} {
			a := ctx.Arg(1).(float64)
			b := ctx.Arg(2).(float64)
			return a + b
		},
		"-": func(ctx *gisp.Context) interface{} {
			a := ctx.Arg(1).(float64)
			b := ctx.Arg(2).(float64)
			return a - b
		},
	}

	out, _ := gisp.RunJSON([]byte(`["-", ["+", 5, 1], ["+", 1, 1]]`), &gisp.Context{
		Sandbox: sandbox,
	})

	assert.Equal(t, float64(4), out)
}

func TestReturnFn(t *testing.T) {
	sandbox := gisp.Sandbox{
		"foo": func(ctx *gisp.Context) interface{} {
			return func(ctx *gisp.Context) interface{} {
				return ctx.Arg(1).(float64) + ctx.Arg(2).(float64)
			}
		},
	}

	out, _ := gisp.RunJSON([]byte(`[["foo"], 1, 2]`), &gisp.Context{
		Sandbox: sandbox,
	})

	assert.Equal(t, float64(3), out)
}

func TestStr(t *testing.T) {
	sandbox := gisp.Sandbox{}

	out, _ := gisp.RunJSON([]byte(`"foo"`), &gisp.Context{
		Sandbox: sandbox,
	})

	assert.Equal(t, "foo", out)
}

func TestAST(t *testing.T) {
	code := []byte(`["*", ["*", 2, 5], ["*", 9, 3]]`)
	var ast interface{}
	json.Unmarshal(code, &ast)

	sandbox := gisp.Sandbox{
		"*": func(ctx *gisp.Context) interface{} {
			a := ctx.Arg(1).(float64)
			b := ctx.Arg(2).(float64)
			return a * b
		},
	}

	out := gisp.Run(&gisp.Context{
		AST:     ast,
		Sandbox: sandbox,
	})

	assert.Equal(t, float64(270), out)
}

func TestMissName(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("should panic")
		} else {
			assert.Equal(
				t,
				"function \"foo\" is undefined",
				r.(error).Error(),
			)
		}
	}()

	gisp.RunJSON([]byte(`["@", ["@", 1, 1], ["@", ["foo"], 1]]`), &gisp.Context{
		Sandbox: gisp.Sandbox{
			"@": func(ctx *gisp.Context) interface{} {
				ctx.Arg(1)
				ctx.Arg(2)
				return nil
			},
		},
	})
}

func TestRuntimeErr(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("should panic")
		} else {
			assert.Equal(
				t,
				"[foo 1 @ 2 @ 0]",
				fmt.Sprint(r.(gisp.Error).Stack),
			)
		}
	}()

	gisp.RunJSON([]byte(`["@", ["@", 1, 1], ["@", ["foo"], 1]]`), &gisp.Context{
		Sandbox: gisp.Sandbox{
			"foo": func(ctx *gisp.Context) interface{} {
				a := []int{}
				a[100] = 1
				return nil
			},
			"@": func(ctx *gisp.Context) interface{} {
				ctx.Arg(1)
				ctx.Arg(2)
				return nil
			},
		},
	})
}

func BenchmarkAST(b *testing.B) {
	code := []byte(`["+", 1, 1]`)
	var ast interface{}
	json.Unmarshal(code, &ast)

	sandbox := gisp.Sandbox{
		"+": func(ctx *gisp.Context) interface{} {
			a := ctx.Arg(1).(float64)
			b := ctx.Arg(2).(float64)
			return a + b
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gisp.Run(&gisp.Context{
			AST:     ast,
			Sandbox: sandbox,
		})
	}
}

func BenchmarkJSON(b *testing.B) {
	sandbox := gisp.Sandbox{
		"+": func(ctx *gisp.Context) interface{} {
			a := ctx.Arg(1).(float64)
			b := ctx.Arg(2).(float64)
			return a + b
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gisp.RunJSON([]byte(`["+", ["+", 1, 1], ["+", 1, 1]]`), &gisp.Context{
			Sandbox: sandbox,
		})
	}
}

func BenchmarkJSONBase(b *testing.B) {
	for i := 0; i < b.N; i++ {
		code := []byte(`["+", ["+", 1, 1], ["+", 1, 1]]`)
		var ast interface{}
		json.Unmarshal(code, &ast)
	}
}
