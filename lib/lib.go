package lib

import (
	"fmt"
	"math"

	"github.com/ysmood/gisp"
)

// MaxFnStackSize ...
const MaxFnStackSize = 17

// Raw ...
func Raw(ctx *gisp.Context) interface{} {
	return ctx.AST.([]interface{})[1]
}

// Throw ...
func Throw(ctx *gisp.Context) {
	ctx.Error(ctx.ArgStr(1))
}

// Get ...
func Get(ctx *gisp.Context) interface{} {
	obj := ctx.Arg(1)
	pathRaw := ctx.Arg(2)
	defaultVal := ctx.Arg(3)

	paths := toJSONPath(pathRaw)

	l := len(paths)

	if l == 0 {
		return defaultVal
	}

	for i := 0; i < l; i++ {
		p := paths[i]
		switch p.(type) {
		case string:
			var has bool
			dict, ok := obj.(map[string]interface{})

			if !ok {
				return defaultVal
			}

			obj, has = dict[p.(string)]

			if !has {
				return defaultVal
			}
		case uint64:
			arr, ok := obj.([]interface{})

			if !ok {
				return defaultVal
			}

			if int(p.(uint64)) >= len(arr) {
				return defaultVal
			}
			obj = arr[p.(uint64)]
		default:
			return defaultVal
		}
	}
	return obj
}

// Set ...
func Set(ctx *gisp.Context) interface{} {
	obj := ctx.Arg(1)
	pathRaw := ctx.Arg(2)

	// TODO: optimize circular detection
	// What a shame, go marshal doesn't support circular detection.
	// Here we use headless clone to break the links.
	val := clone(ctx.Arg(3))

	paths := toJSONPath(pathRaw)

	var p interface{}
	pathLen := len(paths)
	cur := obj
	last := pathLen - 1
	for i := 0; i < pathLen; i++ {
		p = paths[i]

		switch p.(type) {
		case string:
			index := p.(string)
			if i == last {
				cur.(map[string]interface{})[index] = val
			} else {
				next := cur.(map[string]interface{})[index]

				switch next.(type) {
				case map[string]interface{}:
				case []interface{}:
				default:
					if isUint64(paths[i+1]) {
						next = make([]interface{}, int(paths[i+1].(uint64)+1))
					} else {
						next = map[string]interface{}{}
					}

					cur.(map[string]interface{})[index] = next
				}
				cur = next
			}

		case uint64:
			index := p.(uint64)
			arr := cur.([]interface{})
			if i == last {
				arr[index] = val
			} else {
				l := uint64(len(arr))
				if index >= l {
					arr = append(arr, make([]interface{}, index-l+1)...)
					if i == 0 {
						obj = arr
					}
				}
				next := arr[index]

				switch next.(type) {
				case map[string]interface{}:
				case []interface{}:
				default:
					if isUint64(paths[i+1]) {
						next = make([]interface{}, paths[i+1].(uint64)+1)
					} else {
						next = map[string]interface{}{}
					}

					arr[index] = next
				}
				cur = next
			}

		default:
			ctx.Error(fmt.Sprintf("wrong path type: %T\nvalue: %v", p, p))
		}
	}
	return obj
}

// Str ...
func Str(ctx *gisp.Context) interface{} {
	return str(ctx.Arg(1))
}

// Includes ...
func Includes(ctx *gisp.Context) bool {
	list, isArr := ctx.Arg(1).([]interface{})

	if !isArr {
		return false
	}

	target := ctx.Arg(2)

	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}

// Arr ...
func Arr(ctx *gisp.Context) interface{} {
	l := ctx.Len()
	arr := make([]interface{}, l-1)

	for i := 1; i < l; i++ {
		arr[i-1] = ctx.Arg(i)
	}

	return arr
}

// Dict ...
func Dict(ctx *gisp.Context) interface{} {
	l := ctx.Len() - 1
	dict := make(map[string]interface{})
	for i := 1; i < l; i = i + 2 {
		dict[str(ctx.Arg(i))] = ctx.Arg(i + 1)
	}

	return dict
}

// Do ...
func Do(ctx *gisp.Context) interface{} {
	l := ctx.Len()
	var ret interface{}
	for i := 1; i < l; i++ {
		ret = ctx.Arg(i)
	}
	return ret
}

// Def define
func Def(ctx *gisp.Context) interface{} {
	val := ctx.Arg(2)
	ctx.Sandbox.Set(str(ctx.Arg(1)), val)
	return val
}

// Redef redefine
func Redef(ctx *gisp.Context) interface{} {
	val := ctx.Arg(2)
	ctx.Sandbox.Reset(str(ctx.Arg(1)), val)
	return val
}

// If ...
func If(ctx *gisp.Context) interface{} {
	if ctx.ArgBool(1) {
		return ctx.Arg(2)
	}
	return ctx.Arg(3)
}

// Add ...
func Add(ctx *gisp.Context) (ret interface{}) {
	ast := ctx.AST.([]interface{})
	l := len(ast)
	isStr := false

	if l == 1 {
		ret = float64(0)
	} else {
		arg := ctx.Arg(1)
		switch arg.(type) {
		case string:
			if l == 2 {
				ret = s2f(arg)
			} else {
				isStr = true
				ret = arg
			}
		default:
			ret = arg
		}
	}

	for i := 2; i < l; i++ {
		arg := ctx.Arg(i)

		switch arg.(type) {
		case string:
			if isStr {
				ret = ret.(string) + arg.(string)
			} else {
				ret = f2s(ret) + arg.(string)
			}
			isStr = true
		case float64:
			if isStr {
				ret = ret.(string) + f2s(arg)
			} else {
				ret = ret.(float64) + arg.(float64)
			}
		default:
			ret = fmt.Sprint(ret) + fmt.Sprint(arg)
		}
	}
	return
}

// Minus (- 10 2 3)
func Minus(ctx *gisp.Context) float64 {
	l := ctx.Len()
	o := ctx.ArgNum(1)
	for i := 2; i < l; i++ {
		o -= ctx.ArgNum(i)
	}
	return o
}

// Multiply (* 1 2 3)
func Multiply(ctx *gisp.Context) float64 {
	l := ctx.Len()
	o := ctx.ArgNum(1)
	for i := 2; i < l; i++ {
		o *= ctx.ArgNum(i)
	}
	return o
}

// Power ...
func Power(ctx *gisp.Context) interface{} {
	return math.Pow(ctx.ArgNum(1), ctx.ArgNum(2))
}

// Divide ...
func Divide(ctx *gisp.Context) interface{} {
	l := ctx.Len()
	o := ctx.ArgNum(1)
	for i := 2; i < l; i++ {
		o /= ctx.ArgNum(i)
	}
	return o
}

// Mod ...
func Mod(ctx *gisp.Context) interface{} {
	return math.Mod(ctx.ArgNum(1), ctx.ArgNum(2))
}

// Eq ...
func Eq(ctx *gisp.Context) interface{} {
	l := ctx.Len()
	last := ctx.Arg(1)
	for i := 2; i < l; i++ {
		if last != ctx.Arg(i) {
			return false
		}
	}

	return true
}

// Ne ...
func Ne(ctx *gisp.Context) interface{} {
	return ctx.Arg(1) != ctx.Arg(2)
}

// Lt ...
func Lt(ctx *gisp.Context) interface{} {
	l := ctx.Len()
	a := ctx.Arg(1)

	for i := 2; i < l; i++ {
		b := ctx.Arg(i)
		switch a.(type) {
		case string:
			if a.(string) >= b.(string) {
				return false
			}
		case float64:
			if a.(float64) >= b.(float64) {
				return false
			}
		default:
			if fmt.Sprint(a) >= fmt.Sprint(b) {
				return false
			}
		}
		a = b
	}

	return true
}

// Le ...
func Le(ctx *gisp.Context) interface{} {
	l := ctx.Len()
	a := ctx.Arg(1)

	for i := 2; i < l; i++ {
		b := ctx.Arg(i)
		switch a.(type) {
		case string:
			if a.(string) > b.(string) {
				return false
			}
		case float64:
			if a.(float64) > b.(float64) {
				return false
			}
		default:
			if fmt.Sprint(a) > fmt.Sprint(b) {
				return false
			}
		}
		a = b
	}

	return true
}

// Gt ...
func Gt(ctx *gisp.Context) interface{} {
	l := ctx.Len()
	a := ctx.Arg(1)

	for i := 2; i < l; i++ {
		b := ctx.Arg(i)
		switch a.(type) {
		case string:
			if a.(string) <= b.(string) {
				return false
			}
		case float64:
			if a.(float64) <= b.(float64) {
				return false
			}
		default:
			if fmt.Sprint(a) <= fmt.Sprint(b) {
				return false
			}
		}
		a = b
	}

	return true
}

// Ge ...
func Ge(ctx *gisp.Context) interface{} {
	l := ctx.Len()
	a := ctx.Arg(1)

	for i := 2; i < l; i++ {
		b := ctx.Arg(i)
		switch a.(type) {
		case string:
			if a.(string) < b.(string) {
				return false
			}
		case float64:
			if a.(float64) < b.(float64) {
				return false
			}
		default:
			if fmt.Sprint(a) < fmt.Sprint(b) {
				return false
			}
		}
		a = b
	}

	return true
}

// Not ...
func Not(ctx *gisp.Context) interface{} {
	return !ctx.ArgBool(1)
}

// And ...
func And(ctx *gisp.Context) interface{} {
	// for the laziness, we shouldn't use "ctx.Arg(1).(bool) && ctx.Arg(2).(bool)"
	l := ctx.Len()
	for i := 1; i < l; i++ {
		if !ctx.ArgBool(i) {
			return false
		}
	}

	return true
}

// Or ...
func Or(ctx *gisp.Context) interface{} {
	// for the laziness, we shouldn't use "ctx.Arg(1).(bool) || ctx.Arg(2).(bool)"
	l := ctx.Len()
	for i := 1; i < l; i++ {
		if ctx.ArgBool(i) {
			return true
		}
	}

	return false
}

// Switch ...
func Switch(ctx *gisp.Context) interface{} {

	if ctx.Len() == 1 {
		return nil
	}

	ast := ctx.AST.([]interface{})
	start := 1
	end := ctx.Len() - 1

	firstAst, firstOk := ast[start].([]interface{})
	var expr interface{}
	hasExpr := false

	if !firstOk {
		expr = ctx.Arg(start)
		hasExpr = true
		start++
	} else if len(firstAst) == 1 {
		if name, ok := firstAst[0].(string); !ok || name != "case" {
			expr = ctx.Arg(start)
			hasExpr = true
			start++
		}
	}

	lastAst, lastOk := ast[end].([]interface{})
	var defaultAst interface{}

	if lastOk && len(lastAst) == 2 {
		if name, ok := lastAst[0].(string); ok && name == "default" {
			defaultAst = lastAst[1]
			end--
		}
	}

	for i := start; i <= end; i++ {
		node, nodeOk := ast[i].([]interface{})
		if !nodeOk || len(node) != 3 {
			ctx.Error("switch unexpected identifier")
			return nil
		}
		name, nameOk := node[0].(string)
		if !nameOk || name != "case" {
			ctx.Error("switch unexpected identifier")
			return nil
		}
		itemValue := gisp.Run(&gisp.Context{
			AST:         node[1],
			Sandbox:     ctx.Sandbox,
			ENV:         ctx.ENV,
			Parent:      ctx.Parent,
			Index:       ctx.Index,
			IsLiftPanic: ctx.IsLiftPanic,
		})
		if hasExpr {
			if itemValue == expr {
				return gisp.Run(&gisp.Context{
					AST:         node[2],
					Sandbox:     ctx.Sandbox,
					ENV:         ctx.ENV,
					Parent:      ctx.Parent,
					Index:       ctx.Index,
					IsLiftPanic: ctx.IsLiftPanic,
				})
			}
		} else {
			if assert, ok := itemValue.(bool); ok && assert {
				return gisp.Run(&gisp.Context{
					AST:         node[2],
					Sandbox:     ctx.Sandbox,
					ENV:         ctx.ENV,
					Parent:      ctx.Parent,
					Index:       ctx.Index,
					IsLiftPanic: ctx.IsLiftPanic,
				})
			}
		}
	}

	return gisp.Run(&gisp.Context{
		AST:         defaultAst,
		Sandbox:     ctx.Sandbox,
		ENV:         ctx.ENV,
		Parent:      ctx.Parent,
		Index:       ctx.Index,
		IsLiftPanic: ctx.IsLiftPanic,
	})
}

// Fn Define a closure.
// (fn (a b ...) (exp))
func Fn(ctx *gisp.Context) interface{} {
	return func(this *gisp.Context) interface{} {
		if countStack(this) > MaxFnStackSize {
			this.Error("call stack overflow")
		}

		closure := ctx.Sandbox.Create()

		ast := ctx.AST.([]interface{})

		args := ast[1].([]interface{})

		for i, l := 0, len(args); i < l; i++ {
			closure.Set(args[i].(string), this.Arg(i+1))
		}

		return gisp.Run(&gisp.Context{
			AST:     ast[2],
			Sandbox: closure,
			ENV:     ctx.ENV,
			Parent:  this,
			Index:   ctx.Index,
		})
	}
}

func countStack(node *gisp.Context) int {
	count := 0

	for node != nil {
		count++
		node = node.Parent
	}

	return count
}

// For loop function that works like golang
// Example: (for i item (arr) (append (list) (item)))
func For(ctx *gisp.Context) {
	keyName := ctx.ArgStr(1)
	valName := ctx.ArgStr(2)
	arr := ctx.Arg(3)
	ast := ctx.AST.([]interface{})

	closure := ctx.Sandbox.Create()

	switch arr.(type) {
	case []interface{}:
		for i, item := range arr.([]interface{}) {
			closure.Set(keyName, i)
			closure.Set(valName, item)

			gisp.Run(&gisp.Context{
				AST:     ast[4],
				Sandbox: closure,
				ENV:     ctx.ENV,
				Parent:  ctx,
				Index:   ctx.Index,
			})
		}

	case map[string]interface{}:
		for i, item := range arr.(map[string]interface{}) {
			closure.Set(keyName, i)
			closure.Set(valName, item)

			gisp.Run(&gisp.Context{
				AST:     ast[4],
				Sandbox: closure,
				ENV:     ctx.ENV,
				Parent:  ctx,
				Index:   ctx.Index,
			})
		}

	default:
		ctx.Error("cannot iterate non-collection type")
	}
}

// Len get size of array, map or string.
// If type is not supported return -1.
func Len(ctx *gisp.Context) float64 {
	obj := ctx.Arg(1)

	switch obj.(type) {
	case []interface{}:
		return float64(len(obj.([]interface{})))
	case map[string]interface{}:
		return float64(len(obj.(map[string]interface{})))
	case string:
		return float64(len(obj.(string)))
	default:
		return -1
	}
}

// Concat ...
func Concat(ctx *gisp.Context) []interface{} {
	arr := []interface{}{}

	for i, l := 1, ctx.Len(); i < l; i++ {
		el := ctx.Arg(i)

		switch el.(type) {
		case []interface{}:
			arr = append(arr, el.([]interface{})...)

		default:
			arr = append(arr, el)
		}
	}

	return arr
}

// Append ...
func Append(ctx *gisp.Context) []interface{} {
	return append(ctx.ArgArr(1), ctx.Arg(2))
}
