package lib

import (
	"fmt"
	"strconv"
	"strings"
)

func f2s(v interface{}) string {
	return strconv.FormatFloat(v.(float64), 'f', -1, 64)
}

func s2f(v interface{}) (f float64) {
	f, _ = strconv.ParseFloat(v.(string), 64)
	return
}

func isUint64(v interface{}) (is bool) {
	_, is = v.(uint64)
	return
}

func str(val interface{}) (str string) {
	switch val.(type) {
	case string:
		str = val.(string)
	case float64:
		str = f2s(val)
	case []byte:
		str = string(val.([]byte))
	default:
		str = fmt.Sprint(val)
	}
	return
}

func clone(obj interface{}) interface{} {
	switch obj.(type) {
	case map[string]interface{}:
		new := map[string]interface{}{}
		for k, v := range obj.(map[string]interface{}) {
			new[k] = clone(v)
		}
		return new

	case []interface{}:
		new := make([]interface{}, len(obj.([]interface{})))
		for k, v := range obj.([]interface{}) {
			new[k] = clone(v)
		}
		return new

	default:
		return obj
	}
}

func toJSONPath(pathRaw interface{}) (paths []interface{}) {
	switch pathRaw.(type) {
	case string:
		strArr := strings.Split(pathRaw.(string), ".")
		paths = make([]interface{}, len(strArr))
		for i, p := range strArr {
			f, err := strconv.ParseUint(p, 10, 32)
			if err == nil {
				paths[i] = f
			} else {
				paths[i] = p
			}
		}
	case []interface{}:
		paths = pathRaw.([]interface{})
	case float64:
		paths = []interface{}{
			uint64(pathRaw.(float64)),
		}
	default:
		paths = []interface{}{}
	}
	return
}
