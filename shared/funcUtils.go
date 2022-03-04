package shared

import (
	"errors"
	"reflect"
)

type stubMapping map[string]interface{}

var StubStorage = stubMapping{}

func CallServiceFunc(stub map[string]interface{}, funcName string, params ...interface{}) (result interface{}, err error) {
	StubStorage = stub
	f := reflect.ValueOf(StubStorage[funcName])
	if len(params) != f.Type().NumIn() {
		err = errors.New("The number of params is out of index.")
		return
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	var res []reflect.Value
	res = f.Call(in)
	result = res[0].Interface()
	return
}
