package utils

import "reflect"

func NonNil[T any](k T) T {
	v := reflect.ValueOf(k)

	//nolint:exhaustive
	switch v.Kind() {
	case reflect.Slice:
		if v.IsNil() {
			return reflect.MakeSlice(v.Type(), 0, 0).Interface().(T)
		}
	case reflect.Ptr, reflect.Map, reflect.Chan, reflect.Func, reflect.Interface:
		if v.IsNil() {
			return reflect.New(v.Type().Elem()).Interface().(T)
		}
	default:
		panic("unhandled default case: " + v.Kind().String())
	}

	return k
}
