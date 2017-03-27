package lib_test

import (
	"fmt"
	"testing"

	"github.com/a8m/djson"
	"github.com/stretchr/testify/assert"
	"github.com/ysmood/gisp"
	"github.com/ysmood/gisp/lib"
)

func TestGet(t *testing.T) {
	out, _ := gisp.RunJSON(`["get", { "a": 1.1 }, "a"]`, &gisp.Context{
		Sandbox: map[string]interface{}{
			"get": lib.Get,
		},
	})
	assert.Equal(t, float64(1.1), out)
}

func TestGetPath(t *testing.T) {
	out, _ := gisp.RunJSON(`["get", { "a": {"b": [1,2,3]} }, "a.b.1"]`, &gisp.Context{
		Sandbox: map[string]interface{}{
			"get": lib.Get,
		},
	})
	assert.Equal(t, float64(2), out)
}

func TestGetDefault(t *testing.T) {
	out, _ := gisp.RunJSON(`["get", {}, 1]`, &gisp.Context{
		Sandbox: map[string]interface{}{
			"get": lib.Get,
		},
	})
	fmt.Println(out)
	assert.Equal(t, nil, out)
}

func TestGetArrDefault(t *testing.T) {
	out, _ := gisp.RunJSON(`["get", [], "10"]`, &gisp.Context{
		Sandbox: map[string]interface{}{
			"get": lib.Get,
		},
	})
	assert.Equal(t, nil, out)
}

func TestGetDefaultValFromObj(t *testing.T) {
	out, _ := gisp.RunJSON(`["get", { "a": 1.1 }, "b", false]`, &gisp.Context{
		Sandbox: map[string]interface{}{
			"get": lib.Get,
		},
	})
	assert.Equal(t, false, out)
}

func TestGetDefaultValFromArr(t *testing.T) {
	out, _ := gisp.RunJSON(`["get", ["$", []], "10", false]`, &gisp.Context{
		Sandbox: map[string]interface{}{
			"$":   lib.Raw,
			"get": lib.Get,
		},
	})
	assert.Equal(t, false, out)
}

func TestSet(t *testing.T) {
	out, _ := gisp.RunJSON(`
		["set", ["$", {}], "a.b", "ok"]
	`, &gisp.Context{
		Sandbox: map[string]interface{}{
			"$":   lib.Raw,
			"set": lib.Set,
		},
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
		Sandbox: map[string]interface{}{
			"|":   lib.Arr,
			"set": lib.Set,
		},
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
		Sandbox: map[string]interface{}{
			":":   lib.Dict,
			"set": lib.Set,
		},
	})
	exp, _ := djson.Decode([]byte(`
		{"a": [null, null, "ok"]}
	`))
	assert.Equal(t, exp, out)
}
