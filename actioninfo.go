package trinity

import (
	"reflect"
)

// Contains action information that is used during controller registration, including
// supported request type (GET/POST/...). 
// Note that if action for specific http-method is not found then action for GET method will be used.
type ActionInfo struct {
	handler reflect.Value
	method  string // http-method
	action  Action
}

// NewActionInfo constructs a new ActionInfo using a given func handler. Handler's
// kind must be equal to 'reflect.Func'.
func NewActionInfo(handler reflect.Value) *ActionInfo {
	logger.Trace("")

	if handler.Kind() != reflect.Func {
		panic("ActionInfo received unsupported value")
	}

	actionInfo := new(ActionInfo)

	actionInfo.handler = handler

	return actionInfo
}

// Method sets the http-method. Returns self (for chaining).
func (info *ActionInfo) Method(method string) *ActionInfo {
	info.method = method
	return info
}

// Action sets the action name. Can be used when method name is not the same as the action name or
// if one action is processed by multiple methods (with different request types)
// Returns self (for chaining).
func (info *ActionInfo) Action(a Action) *ActionInfo {
	info.action = a
	return info
}
