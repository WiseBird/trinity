package trinity

import (
	"code.google.com/p/gorilla/schema"
	"net/url"
	"reflect"
)

/*

Example of controller -> methodDescriptor -> handlerInvoker flow:

Example controller:

	type MyController struct {
		*mvc.ToLowerController
	}

	func NewMyController() *MyController {
		myController := new(MyController)

		myController.ToLowerController = mvc.NewToLowerController(accountController)
		myController.AddAction("Hello")

		return myController
	}

	func (myController *MyController) Action1(response http.ResponseWriter, request *http.Request) mvc.ActionResultInterface {
		...
	}

	func (myController *MyController) Action2(input *Action2Input) mvc.ActionResultInterface {
		...
	}

	type Action2Input struct {
		Name string
	}


For all Action* methods (there are two in this example) method descriptors will be created:

	methodDescriptor {
		value: Action1 reflected value,
		inTypes: [Typeof(*MyController), Typeof(http.ResponseWriter), Typeof(*http.Request)]
	}

	methodDescriptor {
		value: Action2 reflected value,
		inTypes: [Typeof(*MyController), Typeof(*Action2Input)]
	}

When handler invokers for these method descriptors are called parameters and values will be filled:

1. Parameters: *MyController, http.ResponseWriter, *http.Request, Controller, Action.

2. Values: ( url.Values ), extracted from URL and form data.

For the first method, only parameters are used.
For the second one, there is an argument (*Action2Input) which is not listed in parameters, so an empty
Action2Input object would be created and filled using form data and URL. Here the term 'Filled' 
means that 'Name' field of the created Action2Input struct would be set to a value with the same name
(URL value or form value "Name").

*/

var (
	decoder = schema.NewDecoder()
)

func decode(input interface{}, values url.Values) {
	removeQuotesFromStringValues(values)
	decoder.Decode(input, values)
}

// removeQuotesFromStringValues is used to remove quotes from strings that are added
// when javascript strings parameters are passed.
func removeQuotesFromStringValues(values url.Values) {
	for _, vals := range values {
		for i, val := range vals {
			if len(val) > 1 && val[0] == '"' && val[len(val)-1] == '"' {
				vals[i] = val[1 : len(val)-1]
			}
		}
	}
}

// invokerParam stores info needed to call a method with reflect.
type invokerParam struct {
	rValue reflect.Value // Непосредственно значение параметра
	rType  reflect.Type  // Тип значения, если исходное значение - указатель, то это тип значения под указателем
	isPtr  bool          // Является ли параметр указателем
}

// newInvokerParam creates a new invokerParam
func newInvokerParam(value interface{}) *invokerParam {
	param := new(invokerParam)

	param.rValue = reflect.ValueOf(value)
	param.rType = param.rValue.Type()
	traceParamType(param.rType)
	if param.rType.Kind() == reflect.Ptr {
		param.isPtr = true
		traceParamType(param.rType.Elem())
	}

	return param
}

// handlerInvoker is used to perform an action call using reflection.
// To get a parameters list two sources are used:
// 1. Predefined invokerParam objects
// 2. New objects that are extracted from the http request
type handlerInvoker struct {
	values  url.Values
	params  []*invokerParam
	handler *methodDescriptor
}

// newHandlerInvoker constructs a handlerInvoker that will be used to invoke 
// a method from the specified method descriptor.
func newHandlerInvoker(handler *methodDescriptor) *handlerInvoker {
	invoker := new(handlerInvoker)

	invoker.handler = handler
	invoker.values = make(url.Values, 0)
	invoker.params = make([]*invokerParam, 0)

	return invoker
}

// AddValue adds a string/value pair to the specified handler invoker. The invoker 
// is returned to provide convenient method call chaining.
func (invoker *handlerInvoker) AddValue(name string, val string) *handlerInvoker {
	logger.Trace("")
	logger.Debugf("[%s] %s", name, val)

	_, exists := invoker.values[name]
	if !exists {
		invoker.values[name] = make([]string, 0)
	}
	invoker.values[name] = append(invoker.values[name], val)

	return invoker
}

// AddValues calls AddValue on the given name/val pairs array. The invoker 
// is returned to provide convenient method call chaining.
func (invoker *handlerInvoker) AddValues(values url.Values) *handlerInvoker {
	logger.Trace("")

	for k, vs := range values {
		for _, v := range vs {
			invoker.AddValue(k, v)
		}
	}

	return invoker
}

// AddParam adds a parameter to the invoker. The invoker 
// is returned to provide convenient method call chaining.
func (invoker *handlerInvoker) AddParam(param interface{}) *handlerInvoker {
	logger.Trace("")

	invoker.params = append(invoker.params, newInvokerParam(param))

	return invoker
}

// Invoke calls the handler action and returns the invocation result.
func (invoker *handlerInvoker) Invoke() ActionResultInterface {
	logger.Trace("")

	args := make([]reflect.Value, len(invoker.handler.inTypes))
	for i, paramType := range invoker.handler.inTypes {
		paramValue := invoker.getParamValue(paramType)
		traceParamValue(paramValue)

		args[i] = paramValue
	}

	rets := invoker.handler.value.Call(args)
	if len(rets) > 0 {
		if res, ok := rets[0].Interface().(ActionResultInterface); ok {
			return res
		}
	}
	return nil
}

// getParamValue returns the specified invoker parameter value. If the value is not present in params, it will be
// created and filled from values
func (invoker *handlerInvoker) getParamValue(paramTypeRaw reflect.Type) reflect.Value {
	logger.Trace("")
	traceParamType(paramTypeRaw)

	isPtr := paramTypeRaw.Kind() == reflect.Ptr
	paramType := paramTypeRaw
	if isPtr {
		paramType = paramType.Elem()
		traceParamType(paramType)
	}

	for _, param := range invoker.params {
		if reflect.DeepEqual(paramType, param.rType) {
			logger.Trace("exactly from params")
			return param.rValue
		}

		if isPtr && param.isPtr {
			p := param.rType.Elem()
			if reflect.DeepEqual(paramType, p) {
				logger.Trace("exactly from params")
				return param.rValue
			}
		}

		if paramType.Kind() == reflect.Interface {
			logger.Trace("interface")
			if param.rType.Implements(paramType) {
				logger.Trace("implemented in params")
				return param.rValue
			}
		}

		/*if paramType.AssignableTo(param.rType) {
			logger.Trace("assignable from params")
			return param.rValue
		}*/
	}

	logger.Trace("decode from values")
	paramValue := getParamNewValue(paramType, isPtr)
	logger.Debugf("%v", paramValue)
	decode(paramValue.Interface(), invoker.values)
	logger.Debugf("%v", paramValue)

	return paramValue
}

func getParamNewValue(paramType reflect.Type, isPtr bool) reflect.Value {
	if isPtr {
		return reflect.New(paramType)
	} else {
		return reflect.New(paramType).Elem()
	}

	panic("impossible")
}
