package self

import "reflect"

func CallMethod(obj interface{}, fn string, args []any) []reflect.Value {
	method := reflect.ValueOf(obj).MethodByName(fn)
	var inputs []reflect.Value
	for _, v := range args {
		inputs = append(inputs, reflect.ValueOf(v))
	}
	return method.Call(inputs)
}
