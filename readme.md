[![Build Status](https://travis-ci.org/ysmood/gisp.svg)](https://travis-ci.org/ysmood/gisp)

Gisp is a core scripting tool serving billions of
dynamic responses (such as AB test response) every day on our servers.

Read js implementation for detailed info: https://github.com/ysmood/nisp

## Example

See the examples folder.

## Compare to Lua

Compare to normal gopher-lua, gisp is about 90 times faster with a much smaller memory footprint.
Though the test situation is very limited, but it reflects gisp is suitable for
simple (such as you don't want your users to take time to learn the complex grammar of Lua), performance demanded and embed scripting situation.

`go test -bench . -benchmem`

```
BenchmarkLua-8                	  100000	     23060 ns/op	   85464 B/op	      73 allocs/op
BenchmarkGisp-8               	 5000000	       248 ns/op	     264 B/op	       5 allocs/op
```