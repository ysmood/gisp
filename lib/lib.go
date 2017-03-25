package lib

import (
	"fmt"
	"math"

	"github.com/ysmood/gisp"
)

// Raw ...
func Raw(ctx *gisp.Context) interface{} {
	return ctx.AST.([]interface{})[1]
}

// Get ...
func Get(ctx *gisp.Context) interface{} {
	obj := ctx.Arg(1)
	pathRaw := ctx.Arg(2)
	defaultVal := ctx.Arg(3)

	paths := toJSONPath(pathRaw)

	for i, l := 0, len(paths); i < l; i++ {
		p := paths[i]
		switch p.(type) {
		case string:
			var has bool
			obj, has = obj.(map[string]interface{})[p.(string)]

			if !has {
				return defaultVal
			}
		case uint64:
			if int(p.(uint64)) >= len(obj.([]interface{})) {
				return defaultVal
			}
			obj = obj.([]interface{})[p.(uint64)]
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
	val := ctx.Arg(3)

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
			if i == last {
				cur.([]interface{})[index] = val
			} else {
				next := cur.([]interface{})[index]

				switch next.(type) {
				case map[string]interface{}:
				case []interface{}:
				default:
					if isUint64(paths[i+1]) {
						next = make([]interface{}, paths[i+1].(uint64)+1)
					} else {
						next = map[string]interface{}{}
					}

					cur.([]interface{})[index] = next
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
func Includes(ctx *gisp.Context) interface{} {
	list := ctx.Arg(1).([]interface{})
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
	var arr []interface{}

	for i := 1; i < l; i++ {
		arr = append(arr, ctx.Arg(i))
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

// Def ...
func Def(ctx *gisp.Context) interface{} {
	val := ctx.Arg(2)
	ctx.Sandbox[str(ctx.Arg(1))] = func(ctx *gisp.Context) interface{} {
		return val
	}
	return val
}

// If ...
func If(ctx *gisp.Context) interface{} {
	if ctx.Arg(1).(bool) {
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

// Minus ...
func Minus(ctx *gisp.Context) interface{} {
	return ctx.Arg(1).(float64) - ctx.Arg(2).(float64)
}

// Multiply ...
func Multiply(ctx *gisp.Context) interface{} {
	return ctx.Arg(1).(float64) * ctx.Arg(2).(float64)
}

// Power ...
func Power(ctx *gisp.Context) interface{} {
	return math.Pow(ctx.Arg(1).(float64), ctx.Arg(2).(float64))
}

// Divide ...
func Divide(ctx *gisp.Context) interface{} {
	return ctx.Arg(1).(float64) / ctx.Arg(2).(float64)
}

// Mod ...
func Mod(ctx *gisp.Context) interface{} {
	return math.Mod(ctx.Arg(1).(float64), ctx.Arg(2).(float64))
}

// Eq ...
func Eq(ctx *gisp.Context) interface{} {
	return ctx.Arg(1) == ctx.Arg(2)
}

// Ne ...
func Ne(ctx *gisp.Context) interface{} {
	return ctx.Arg(1) != ctx.Arg(2)
}

// Lt ...
func Lt(ctx *gisp.Context) interface{} {
	a := ctx.Arg(1)
	b := ctx.Arg(2)
	switch a.(type) {
	case string:
		return a.(string) < b.(string)
	case float64:
		return a.(float64) < b.(float64)
	default:
		return fmt.Sprint(a) < fmt.Sprint(b)
	}
}

// Le ...
func Le(ctx *gisp.Context) interface{} {
	a := ctx.Arg(1)
	b := ctx.Arg(2)
	switch a.(type) {
	case string:
		return a.(string) <= b.(string)
	case float64:
		return a.(float64) <= b.(float64)
	default:
		return fmt.Sprint(a) <= fmt.Sprint(b)
	}
}

// Gt ...
func Gt(ctx *gisp.Context) interface{} {
	a := ctx.Arg(1)
	b := ctx.Arg(2)
	switch a.(type) {
	case string:
		return a.(string) > b.(string)
	case float64:
		return a.(float64) > b.(float64)
	default:
		return fmt.Sprint(a) > fmt.Sprint(b)
	}
}

// Ge ...
func Ge(ctx *gisp.Context) interface{} {
	a := ctx.Arg(1)
	b := ctx.Arg(2)
	switch a.(type) {
	case string:
		return a.(string) >= b.(string)
	case float64:
		return a.(float64) >= b.(float64)
	default:
		return fmt.Sprint(a) >= fmt.Sprint(b)
	}
}

// Not ...
func Not(ctx *gisp.Context) interface{} {
	return !ctx.Arg(1).(bool)
}

// And ...
func And(ctx *gisp.Context) interface{} {
	// for the laziness, we shouldn't use "ctx.Arg(1).(bool) && ctx.Arg(2).(bool)"
	if ctx.Arg(1).(bool) {
		return ctx.Arg(2).(bool)
	}
	return false
}

// Or ...
func Or(ctx *gisp.Context) interface{} {
	// for the laziness, we shouldn't use "ctx.Arg(1).(bool) || ctx.Arg(2).(bool)"
	if ctx.Arg(1).(bool) {
		return true
	}
	return ctx.Arg(2).(bool)
}
