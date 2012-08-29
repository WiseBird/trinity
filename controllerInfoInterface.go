package trinity

import (
	"reflect"
	"strings"
)

// Provides invormation, that would be used during binding
type ControllerInfoInterface interface {
	GetController() Controller              // Controller name
	GetActionInfos() map[string]*ActionInfo // Actions information
}

var (
	controllerSuffix          = "controller"
	actionResultInterfaceTmp  ActionResultInterface
	actionResultInterfaceType = reflect.TypeOf(&actionResultInterfaceTmp).Elem()
)

// Base implementation of the ControllerInfoInterface. Basis for other extracter objects.
// Uses reflect package to extract actions.
// If controller name ends with "controller" then name will be cut off.
// Only methods that return mbc.ActionResultInterface can be used as actions.
type BaseControllerInfoExtracter struct {
	value interface{}

	Controller  Controller
	ActionInfos map[string]*ActionInfo
}

// BaseController constructor. BaseController inheritor is passed as the 'value'
// argument to the NewBaseController
func NewBaseControllerInfoExtracter(value interface{}) *BaseControllerInfoExtracter {
	logger.Trace("")

	baseController := new(BaseControllerInfoExtracter)
	baseController.value = value
	baseController.ActionInfos = make(map[string]*ActionInfo, 0)

	baseController.reflectValue()

	return baseController
}

func (baseController *BaseControllerInfoExtracter) reflectValue() {
	ptrType := reflect.TypeOf(baseController.value)
	if ptrType.Kind() != reflect.Ptr {
		panic("NewBaseControllerInfoExtracter received unsupported type")
	}

	structType := ptrType.Elem()

	structName := structType.Name()
	logger.Debugf("sN = %s", structName)

	if strings.HasSuffix(strings.ToLower(structName), controllerSuffix) {
		baseController.Controller = C(structName[:len(structName)-len(controllerSuffix)])
	} else {
		baseController.Controller = C(structName)
	}

	for i := 0; i < ptrType.NumMethod(); i++ {
		m := ptrType.Method(i)
		logger.Debugf("method= %s", m.Name)

		if !baseController.checkMethod(m) {
			logger.Trace("wrong signature")
			continue
		}

		baseController.addAction(m)
	}
}

func (baseController *BaseControllerInfoExtracter) addAction(method reflect.Method) *ActionInfo {
	logger.Trace("")

	actionInfo := NewActionInfo(method.Func).Action(A(method.Name))
	baseController.ActionInfos[method.Name] = actionInfo
	return actionInfo
}

func (baseController *BaseControllerInfoExtracter) checkMethod(method reflect.Method) bool {
	logger.Trace("")

	if method.Type.NumOut() != 1 {
		logger.Trace("Wrong out number")
		return false
	}

	retType := method.Type.Out(0)

	return reflect.DeepEqual(retType, actionResultInterfaceType)
}

// Adds specific method to action collection.
// Use when you want to add a method with arbitrary signature.
func (baseController *BaseControllerInfoExtracter) AddAction(methodName string) *ActionInfo {
	logger.Trace("")

	method := findMethod(baseController.value, methodName)
	if method == nil {
		return nil
	}

	return baseController.addAction(*method)
}

func findMethod(val interface{}, name string) *reflect.Method {
	ptrType := reflect.TypeOf(val)
	if ptrType.Kind() != reflect.Ptr {
		panic("FindMethod received unsupported type")
	}

	for i := 0; i < ptrType.NumMethod(); i++ {
		m := ptrType.Method(i)

		if m.Name == name {
			return &m
		}
	}

	return nil
}

// GetController returns the name of the controller
func (baseController *BaseControllerInfoExtracter) GetController() Controller {
	return baseController.Controller
}

// GetActionInfos returns information about actions
func (baseController *BaseControllerInfoExtracter) GetActionInfos() map[string]*ActionInfo {
	return baseController.ActionInfos
}
