package gisp_test

import (
	"encoding/json"
	"testing"

	"strings"

	"github.com/a8m/djson"
	"github.com/ysmood/gisp"
	"github.com/ysmood/gisp/lib"
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

func BenchmarkComplexGisp(b *testing.B) {
	code := []byte(`
		[
			"do",
			[
				"def",
				"profession",
				[
					"|",
					"debuglog",
					"holmes",
					"hotfix_metrics",
					"hotfix_log",
					"meituan-update",
					"catchexception",
					"anr",
					"env",
					"large_picture",
					"QRCodeImg",
					"hydra",
					"flexbox",
					"update-downloadmanager",
					"aid",
					"multidex",
					"config_monitor",
					"httpdns",
					"timeout"
				]
			],
			["def", "ret", [":"]],
			[
				"for",
				"index",
				"item",
				[
					"profession"
				],
				[
					"set",
					[
						"ret"
					],
					[
						"item"
					],
					[
						":",
						"commons",
						"ok"
					]
				]
			]
		]
	`)
	ast, err := djson.Decode(code)

	if err != nil {
		panic(err)
	}

	sandbox := gisp.New(gisp.Box{
		"do":  lib.Do,
		"for": lib.For,
		"def": lib.Def,
		"set": lib.Set,
		"get": lib.Get,
		":":   lib.Dict,
		"|":   lib.Arr,
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gisp.Run(&gisp.Context{
			AST:     ast,
			Sandbox: sandbox,
		})
	}
}

// We can see that Gisp is only 4 times slower than the simplified Go code.
func BenchmarkComplexGo(b *testing.B) {
	program := func() interface{} {
		profression := []string{
			"debuglog",
			"holmes",
			"hotfix_metrics",
			"hotfix_log",
			"meituan-update",
			"catchexception",
			"anr",
			"env",
			"large_picture",
			"QRCodeImg",
			"hydra",
			"flexbox",
			"update-downloadmanager",
			"aid",
			"multidex",
			"config_monitor",
			"httpdns",
			"timeout",
		}

		ret := map[string]interface{}{}

		for _, el := range profression {
			// simulate the json path
			paths := strings.Split(el, ".")

			for _, p := range paths {
				ret[p] = map[string]interface{}{
					"commons": "ok",
				}
				break
			}
		}

		return ret
	}

	for i := 0; i < b.N; i++ {
		program()
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
