package gisp_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func IsFunc(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Func
}

func TestLab(t *testing.T) {
	// a := func() {
	// }
	// b := 10

	// fmt.Println(IsFunc(a), IsFunc(b))
	assert.Equal(t, 1, 1)
}
