package gisp_test

import (
	"encoding/json"
	"testing"

	"github.com/a8m/djson"
	"github.com/ysmood/gisp"
	lua "github.com/yuin/gopher-lua"
)

func noop(v interface{}) {}

func BenchmarkTypeAssertion(b *testing.B) {
	var s interface{} = "test"

	for n := 0; n < b.N; n++ {
		switch s.(type) {
		case string:
			noop(s.(string))
		case int32:
			noop(s.(int32))
		case int64:
			noop(s.(int64))
		}
	}
}

func BenchmarkTypeNonAssertion(b *testing.B) {
	var s string = "test"

	for n := 0; n < b.N; n++ {
		noop(s)
	}
}

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
		L.DoString("add(1.1,1.2)")
	}
}

func BenchmarkGisp(b *testing.B) {
	code := []byte(`["+", 1.1, 1.2]`)
	ast, _ := djson.Decode(code)

	sandbox := gisp.New(gisp.Box{
		"+": func(ctx *gisp.Context) float64 {
			return ctx.ArgNum(1) + ctx.ArgNum(2)
		},
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gisp.Run(&gisp.Context{
			AST:     ast,
			Sandbox: sandbox,
		})
	}
}

func BenchmarkLiftPanic(b *testing.B) {
	code := []byte(`["+", 1, 1]`)
	ast, _ := djson.Decode(code)

	sandbox := gisp.New(gisp.Box{
		"+": func(ctx *gisp.Context) float64 {
			return ctx.ArgNum(1) + ctx.ArgNum(2)
		},
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gisp.Run(&gisp.Context{
			AST:         ast,
			Sandbox:     sandbox,
			IsLiftPanic: true,
		})
	}
}

func BenchmarkJSON(b *testing.B) {
	sandbox := gisp.New(gisp.Box{
		"+": func(ctx *gisp.Context) float64 {
			return ctx.ArgNum(1) + ctx.ArgNum(2)
		},
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gisp.RunJSON(`["+", ["+", 1, 1], ["+", 1, 1]]`, &gisp.Context{
			Sandbox: sandbox,
		})
	}
}

func BenchmarkDJSONBase(b *testing.B) {
	for i := 0; i < b.N; i++ {
		code := []byte(`["+", ["+", 1, 1], ["+", 1, 1]]`)
		djson.Decode(code)
	}
}

func BenchmarkJSONBase(b *testing.B) {
	for i := 0; i < b.N; i++ {
		code := []byte(`["+", ["+", 1, 1], ["+", 1, 1]]`)
		var ast interface{}
		json.Unmarshal(code, &ast)
	}
}
