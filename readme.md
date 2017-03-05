See js implementation for more info: https://github.com/ysmood/nisp

## Example

```go
import (
	"github.com/ysmood/gisp"
)

func main() {
    code := `["+", 1, 2]`

	out, _ := gisp.RunJSON([]byte(code), &gisp.Context{
		Sandbox: gisp.Sandbox{
            "+": func(ctx *gisp.Context) interface{} {
                a := ctx.Arg(1).(float64)
                b := ctx.Arg(2).(float64)
                return a + b
            },
        },
	})

	fmt.Println(out) // 3
}
```