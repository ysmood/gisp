package lib_test

import (
	"testing"

	"github.com/a8m/djson"
	"github.com/stretchr/testify/assert"
	"github.com/ysmood/gisp"
	"github.com/ysmood/gisp/lib"
)

func TestThrow(t *testing.T) {
	defer func() {
		err := recover()

		if err != nil {
			assert.Equal(t, "err", err.(gisp.Error).Message)
		} else {
			panic("should throw error")
		}
	}()

	gisp.RunJSON(`["throw", "err"]`, &gisp.Context{
		Sandbox: gisp.New(gisp.Box{
			"throw": lib.Throw,
		}),
	})
}

func TestGet(t *testing.T) {
	out, _ := gisp.RunJSON(`["get", { "a": 1.1 }, "a"]`, &gisp.Context{
		Sandbox: gisp.New(gisp.Box{
			"get": lib.Get,
		}),
	})
	assert.Equal(t, float64(1.1), out)
}

func TestGetPath(t *testing.T) {
	out, _ := gisp.RunJSON(`["get", { "a": {"b": [1,2,3]} }, "a.b.1"]`, &gisp.Context{
		Sandbox: gisp.New(gisp.Box{
			"get": lib.Get,
		}),
	})
	assert.Equal(t, float64(2), out)
}

func TestGetDefault(t *testing.T) {
	out, _ := gisp.RunJSON(`["get", {}, 1]`, &gisp.Context{
		Sandbox: gisp.New(gisp.Box{
			"get": lib.Get,
		}),
	})
	assert.Equal(t, nil, out)
}

func TestGetArrDefault(t *testing.T) {
	out, _ := gisp.RunJSON(`["get", [], "10"]`, &gisp.Context{
		Sandbox: gisp.New(gisp.Box{
			"get": lib.Get,
		}),
	})
	assert.Equal(t, nil, out)
}

func TestGetDefaultValFromObj(t *testing.T) {
	out, _ := gisp.RunJSON(`["get", { "a": 1.1 }, "b", false]`, &gisp.Context{
		Sandbox: gisp.New(gisp.Box{
			"get": lib.Get,
		}),
	})
	assert.Equal(t, false, out)
}

func TestGetDefaultValFromArr(t *testing.T) {
	out, _ := gisp.RunJSON(`["get", ["$", []], "10", false]`, &gisp.Context{
		Sandbox: gisp.New(gisp.Box{
			"$":   lib.Raw,
			"get": lib.Get,
		}),
	})
	assert.Equal(t, false, out)
}

func TestSet(t *testing.T) {
	out, _ := gisp.RunJSON(`
		["set", ["$", {}], "a.b", "ok"]
	`, &gisp.Context{
		Sandbox: gisp.New(gisp.Box{
			"$":   lib.Raw,
			"set": lib.Set,
		}),
	})
	exp, _ := djson.Decode([]byte(`
		{"a": {"b": "ok"}}
	`))
	assert.Equal(t, exp, out)
}

func TestSetArr(t *testing.T) {
	out, _ := gisp.RunJSON(`
		["set", ["|"], "1.2", "ok"]
	`, &gisp.Context{
		Sandbox: gisp.New(gisp.Box{
			"|":   lib.Arr,
			"set": lib.Set,
		}),
	})
	exp, _ := djson.Decode([]byte(`
		[null, [null, null, "ok"]]
	`))
	assert.Equal(t, exp, out)
}

func TestSetObj(t *testing.T) {
	out, _ := gisp.RunJSON(`
		["set", [":", "a", 10], "a.2", "ok"]
	`, &gisp.Context{
		Sandbox: gisp.New(gisp.Box{
			":":   lib.Dict,
			"set": lib.Set,
		}),
	})
	exp, _ := djson.Decode([]byte(`
		{"a": [null, null, "ok"]}
	`))
	assert.Equal(t, exp, out)
}

func TestSetCircular(t *testing.T) {
	out, _ := gisp.RunJSON(`["do",
		["def", "a", [":"]],
		["set", ["a"], "a", ["a"]]
	]`, &gisp.Context{
		Sandbox: gisp.New(gisp.Box{
			":":   lib.Dict,
			"do":  lib.Do,
			"def": lib.Def,
			"set": lib.Set,
		}),
	})
	exp, _ := djson.Decode([]byte(`
		{"a": {}}
	`))
	assert.Equal(t, exp, out)
}

func TestSwitchHasExpr(t *testing.T) {
	out, _ := gisp.RunJSON(`["do",
		["def", "id", 2],
		["switch",
			["id"],
			["case", 1, 1],
			["case", 2, 2],
			["default", 3]
		]
	]`, &gisp.Context{
		Sandbox: gisp.New(gisp.Box{
			"switch": lib.Switch,
			"do":     lib.Do,
			"def":    lib.Def,
		}),
	})
	exp, _ := djson.Decode([]byte(`
		2
	`))
	assert.Equal(t, exp, out)
}

func TestSwitchNoExpr(t *testing.T) {
	out, _ := gisp.RunJSON(`["do",
		["def", "id", false],
		["def", "id2", true],
		["switch",
			["case", ["id"], 1],
			["case", ["id2"], 2],
			["default", 3]
		]
	]`, &gisp.Context{
		Sandbox: gisp.New(gisp.Box{
			"switch": lib.Switch,
			"do":     lib.Do,
			"def":    lib.Def,
		}),
	})
	exp, _ := djson.Decode([]byte(`
		2
	`))
	assert.Equal(t, exp, out)
}

func TestSwitchDefault(t *testing.T) {
	out, _ := gisp.RunJSON(`["do",
		["def", "id", 1000],
		["switch",
			["id"],
			["case", 1, 1],
			["case", 2, 2],
			["default", 3]
		]
	]`, &gisp.Context{
		Sandbox: gisp.New(gisp.Box{
			"switch": lib.Switch,
			"do":     lib.Do,
			"def":    lib.Def,
		}),
	})
	exp, _ := djson.Decode([]byte(`
		3
	`))
	assert.Equal(t, exp, out)
}

func TestFn(t *testing.T) {
	out, _ := gisp.RunJSON(`
		["do",
			["def",
				"foo",
				["fn", ["a"],
					["+", ["a"], 1]
				]
			],

			["foo", 1]
		]
		
	`, &gisp.Context{
		Sandbox: gisp.New(gisp.Box{
			"do":  lib.Do,
			"fn":  lib.Fn,
			"def": lib.Def,
			"+":   lib.Add,
		}),
	})

	assert.Equal(t, float64(2), out)
}

func TestFor(t *testing.T) {
	out, err := gisp.RunJSON(`
		["do",
			["def", "sum", 0],

			["for", "i", "el", ["arr"],
				["redef", "sum", ["+", ["sum"], ["el"]]]
			],

			["sum"]
		]
		
	`, &gisp.Context{
		Sandbox: gisp.New(gisp.Box{
			"do":    lib.Do,
			"|":     lib.Arr,
			"def":   lib.Def,
			"redef": lib.Redef,
			"+":     lib.Add,
			"for":   lib.For,

			"arr": []interface{}{float64(1), float64(2), float64(3)},
		}),
	})

	if err != nil {
		panic(err)
	}

	assert.Equal(t, float64(6), out)
}

func TestConcat(t *testing.T) {
	out, _ := gisp.RunJSON(`
		["concat", ["arr1"], ["arr2"], ["item"]]
	`, &gisp.Context{
		Sandbox: gisp.New(gisp.Box{
			"concat": lib.Concat,
			"arr1":   []interface{}{float64(1), float64(2)},
			"arr2":   []interface{}{float64(3), "4"},
			"item":   "ok",
		}),
	})
	exp, _ := djson.Decode([]byte(`
		[1, 2, 3, "4", "ok"]
	`))
	assert.Equal(t, exp, out)
}

func TestAppend(t *testing.T) {
	out, _ := gisp.RunJSON(`
		["append", ["arr1"], ["item"]]
	`, &gisp.Context{
		Sandbox: gisp.New(gisp.Box{
			"append": lib.Append,
			"arr1":   []interface{}{float64(1), float64(2)},
			"item":   "ok",
		}),
	})
	exp, _ := djson.Decode([]byte(`
		[1, 2, "ok"]
	`))
	assert.Equal(t, exp, out)
}
