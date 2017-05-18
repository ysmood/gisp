package gisp_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ysmood/gisp"
)

func TestNewClosure(t *testing.T) {
	sandbox := gisp.New(gisp.Box{
		"foo": 1,
		"bar": 3,
	})

	newSandbox := sandbox.Create()

	newSandbox.Set("foo", 2)

	val1, _ := sandbox.Get("foo")
	val2, _ := newSandbox.Get("foo")
	val3, _ := newSandbox.Get("bar")

	assert.Equal(t, 1, val1)
	assert.Equal(t, 2, val2)
	assert.Equal(t, 3, val3)
}

func TestDeepClosure(t *testing.T) {
	c1 := gisp.New(gisp.Box{
		"foo": 1,
	})

	c2 := c1.Create()
	c3 := c2.Create()
	c4 := c3.Create()

	val, _ := c4.Get("foo")
	assert.Equal(t, 1, val)

	c4.Set("foo", 2)
	val2, _ := c1.Get("foo")
	assert.Equal(t, 1, val2)

	c3.Reset("foo", 2)
	val3, _ := c1.Get("foo")
	assert.Equal(t, 2, val3)

}

func TestClosureNames(t *testing.T) {
	c1 := gisp.New(gisp.Box{
		"a": 1,
	})

	c2 := c1.Create()

	c2.Set("b", 2)
	c2.Set("c", 3)

	c3 := c2.Create()

	c3.Set("d", 4)

	c4 := c3.Create()

	c4.Set("e", 5)

	assert.Equal(t, []string{"e", "d", "b", "c", "a"}, c4.Names())

}
