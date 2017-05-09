[![Build Status](https://travis-ci.org/ysmood/gisp.svg)](https://travis-ci.org/ysmood/gisp)

Gisp is a core scripting tool serving billions of
dynamic responses (such as AB test response) every day on our servers.

Read js implementation for detailed info: https://github.com/ysmood/nisp

## Example

```go
import (
    "github.com/ysmood/gisp"
)

func main() {
	//  arithmetic expression: 1 + (2 * 2)
	//    function expression: add(1, multiply(2, 2))
	// gisp (json) expression: ["+", 1, ["*", 2, 2]]
	code := `["+", 1, ["*", 2, 2]]`

	out, _ := gisp.RunJSON(code, &gisp.Context{
		Sandbox: gisp.Sandbox{
			"+": lib.Add,
			"*": lib.Multiply,
		},
	})

	fmt.Println(out) // 5
}
```

## Compare to Lua

Compare to normal gopher-lua, gisp is about 90 times faster with a much smaller memory footprint.
Though the test situation is very limited, but it reflects gisp is suitable for
simple (such as you don't want your users to take time to learn the complex grammar of Lua), performance demanded and embed scripting situation.

`go test -bench . -benchmem`

```
BenchmarkLua-8                	  100000	     23060 ns/op	   85464 B/op	      73 allocs/op
BenchmarkGisp-8               	 5000000	       248 ns/op	     264 B/op	       5 allocs/op
```