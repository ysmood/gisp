package lib

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/ysmood/gisp"
)

// MaxFnStackSize ...
var MaxFnStackSize = int(17)

// MaxStringLen ...
var MaxStringLen = int(1e6)

// Raw ...
func Raw(ctx *gisp.Context) interface{} {
	return ctx.AST.([]interface{})[1]
}

// Throw ...
func Throw(ctx *gisp.Context) interface{} {
	ctx.Error(ctx.ArgStr(1))
	return nil
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
			switch obj.(type) {
			case []interface{}:
				arr := obj.([]interface{})
				if int(p.(uint64)) >= len(arr) {
					return defaultVal
				}
				obj = arr[p.(uint64)]
			case map[string]interface{}:
				var has bool
				dict := obj.(map[string]interface{})
				index := strconv.FormatUint(p.(uint64), 10)
				obj, has = dict[index]

				if !has {
					return defaultVal
				}
			default:
				return defaultVal
			}
		default:
			return defaultVal
		}
	}
	return obj
}

// Set ...
func Set(ctx *gisp.Context) interface{} {
	obj := ctx.Arg(1)
	if obj == nil {
		return nil
	}
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
			self := cur.(map[string]interface{})
			if i == last {
				self[index] = val
			} else {
				next := self[index]

				switch next.(type) {
				case map[string]interface{}:
				case []interface{}:
				default:
					if isUint64(paths[i+1]) {
						next = make([]interface{}, int(paths[i+1].(uint64)+1))
					} else {
						next = map[string]interface{}{}
					}

					self[index] = next
				}
				cur = next
			}

		case uint64:
			index := p.(uint64)
			switch cur.(type) {
			case []interface{}:
				arr := cur.([]interface{})
				l := uint64(len(arr))
				if index >= l {
					arr = append(arr, make([]interface{}, index-l+1)...)
					if i == 0 {
						obj = arr
					}
				}
				if i == last {
					arr[index] = val
				} else {
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
			case map[string]interface{}:
				item := strconv.FormatUint(index, 10)
				self := cur.(map[string]interface{})
				if i == last {
					self[item] = val
				} else {
					next := self[item]

					switch next.(type) {
					case map[string]interface{}:
					case []interface{}:
					default:
						if isUint64(paths[i+1]) {
							next = make([]interface{}, int(paths[i+1].(uint64)+1))
						} else {
							next = map[string]interface{}{}
						}

						self[item] = next
					}
					cur = next
				}
			}

		default:
			ctx.Error(fmt.Sprintf("wrong path type: %T\nvalue: %v", p, p))
		}
	}
	return obj
}

// Del ...
func Del(ctx *gisp.Context) interface{} {
	obj := ctx.Arg(1)
	pathRaw := ctx.Arg(2)

	paths := toJSONPath(pathRaw)

	l := len(paths)
	last := l - 1
	cur := obj
	var prev interface{}
	var prevKey interface{}

	for i := 0; i < l; i++ {
		p := paths[i]
		switch p.(type) {
		case string:
			var has bool
			dict, ok := cur.(map[string]interface{})

			if !ok {
				return obj
			}

			key := p.(string)

			if i == last {
				delete(dict, key)
				return obj
			}

			prev, prevKey = cur, key
			cur, has = dict[key]

			if !has {
				return obj
			}

		case uint64:
			switch cur.(type) {
			case []interface{}:
				arr := cur.([]interface{})
				arrLen := len(arr)
				index := int(p.(uint64))
				if index >= arrLen {
					return obj
				}
				if i == last {
					var newArr []interface{}

					if index == arrLen-1 {
						newArr = arr[:index]
					} else if index == 0 {
						newArr = arr[index+1:]
					} else {
						newArr = append(arr[:index], arr[index+1:]...)
					}

					if prev == nil {
						return newArr
					}

					switch prevKey.(type) {
					case string:
						prevDict := prev.(map[string]interface{})
						prevDict[prevKey.(string)] = newArr
					case uint64:
						prevArr := prev.([]interface{})
						prevArr[prevKey.(uint64)] = newArr
					}

					return obj
				}
				prev, prevKey = cur, index
				cur = arr[p.(uint64)]
			case map[string]interface{}:
				var has bool
				dict := cur.(map[string]interface{})
				key := strconv.FormatUint(p.(uint64), 10)

				if i == last {
					delete(dict, key)
					return obj
				}
				prev, prevKey = cur, key
				cur, has = dict[key]

				if !has {
					return obj
				}
			default:
				return obj
			}
		default:
			return obj
		}
	}
	return obj
}

// Str ...
func Str(ctx *gisp.Context) interface{} {
	return str(ctx.Arg(1))
}

// Includes ...
func Includes(ctx *gisp.Context) interface{} {
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

		switch ret.(type) {
		case string:
			if len(ret.(string)) > MaxStringLen {
				ctx.Error(fmt.Sprintf("max string length exceeded %v", MaxStringLen))
			}
		}
	}
	return
}

// Minus (- 10 2 3)
func Minus(ctx *gisp.Context) interface{} {
	l := ctx.Len()
	o := ctx.ArgNum(1)
	for i := 2; i < l; i++ {
		o -= ctx.ArgNum(i)
	}
	return o
}

// Multiply (* 1 2 3)
func Multiply(ctx *gisp.Context) interface{} {
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
			PreRun:      ctx.PreRun,
			PostRun:     ctx.PostRun,
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
					PreRun:      ctx.PreRun,
					PostRun:     ctx.PostRun,
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
					PreRun:      ctx.PreRun,
					PostRun:     ctx.PostRun,
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
		PreRun:      ctx.PreRun,
		PostRun:     ctx.PostRun,
	})
}

// Fn Define a closure.
// (fn (a b ...) (exp))
func Fn(ctx *gisp.Context) interface{} {
	return func(this *gisp.Context) interface{} {
		// count stack
		node := this
		count := 0
		for node != nil {
			count++
			node = node.Parent
		}

		if count > MaxFnStackSize {
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
			PreRun:  ctx.PreRun,
			PostRun: ctx.PostRun,
		})
	}
}

// For loop function that works like golang
// Example: (for i item (arr) (append (list) (item)))
func For(ctx *gisp.Context) interface{} {
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
				PreRun:  ctx.PreRun,
				PostRun: ctx.PostRun,
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
				PreRun:  ctx.PreRun,
				PostRun: ctx.PostRun,
			})
		}

	default:
		ctx.Error("cannot iterate non-collection type")
	}

	return nil
}

// Len get size of array, map or string.
// If type is not supported return -1.
func Len(ctx *gisp.Context) interface{} {
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
func Concat(ctx *gisp.Context) interface{} {
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
func Append(ctx *gisp.Context) interface{} {
	return append(ctx.ArgArr(1), ctx.Arg(2))
}

// Split ...
func Split(ctx *gisp.Context) interface{} {
	arr := strings.Split(ctx.ArgStr(1), ctx.ArgStr(2))

	ret := make([]interface{}, len(arr))

	for i, el := range arr {
		ret[i] = el
	}

	return ret
}

// Slice ...
func Slice(ctx *gisp.Context) interface{} {
	arr := ctx.Arg(1)

	switch arr.(type) {
	case string:
		return arr.(string)[int(ctx.ArgNum(2)):int(ctx.ArgNum(3))]
	case []interface{}:
		return arr.([]interface{})[int(ctx.ArgNum(2)):int(ctx.ArgNum(3))]
	default:
		return nil
	}
}

// IndexOf ...
func IndexOf(ctx *gisp.Context) interface{} {
	arr := ctx.Arg(1)

	switch arr.(type) {
	case string:
		return float64(strings.Index(arr.(string), ctx.ArgStr(2)))
	case []interface{}:
		list := arr.([]interface{})
		target := ctx.Arg(2)
		for i, el := range list {
			if el == target {
				return float64(i)
			}
		}
		return -1
	default:
		return -1
	}
}
