package trinity

import (
	"reflect"
)

// methodDescriptor contains reflect info for method extracted using reflect
type methodDescriptor struct {
	value   reflect.Value
	inTypes []reflect.Type
}

// newMethodDescriptorFromMethod is used to construct a method descriptor using a given method.
func newMethodDescriptorFromMethod(method interface{}) *methodDescriptor {
	return newMethodDescriptorFromValue(reflect.ValueOf(method))
}

// newMethodDescriptorFromValue is used to construct a method descriptor using a given reflect
// value (it should have Kind = reflect.Func).
func newMethodDescriptorFromValue(value reflect.Value) *methodDescriptor {
	descriptor := new(methodDescriptor)

	descriptor.value = value
	if descriptor.value.Kind() != reflect.Func {
		panic("Incorrect method kind. Expected - Func, goted - " + descriptor.value.Kind().String())
	}

	methodType := descriptor.value.Type()
	descriptor.inTypes = make([]reflect.Type, methodType.NumIn())
	for i := 0; i < methodType.NumIn(); i++ {
		descriptor.inTypes[i] = methodType.In(i)
	}

	return descriptor
}
