package trinity

import (
	"reflect"
)

// BindController registers a specified controller. 
// Expects parameterless constructor of controller.
// Iterates through the controller methods and registers its actions using GetActionInfos.
func (mvcI *MvcInfrastructure) BindController(constructor interface{}) {
	logger.Trace("")

	value, controllerInterface := mvcI.reflectAndCheck(constructor)
	controllerInfo := controllerInterface.GetInfo()

	controller := controllerInfo.GetController()
	mvcI.controllerConstructors[controller] = value
	
	for _, actionInfo := range controllerInfo.GetActionInfos() {
		httpMethod := "GET"
		if actionInfo.method != "" {
			httpMethod = actionInfo.method
		}

		mvcI.bindAction(controller, actionInfo.action, method(httpMethod), newMethodDescriptorFromValue(actionInfo.handler))
	}
}

func (mvcI *MvcInfrastructure) reflectAndCheck(constructor interface{}) (reflect.Value, ControllerInterface) {
	value := reflect.ValueOf(constructor)
	if value.Kind() != reflect.Func {
		panic("Incorrect constructor kind. Expected - Func, goted - " + value.Kind().String())
	}
	
	valueType := value.Type()
	if valueType.NumIn() > 0 {
		panic("Constructor must be parameterless.")
	}
	
	result := value.Call(emptyValues)
	if len(result) != 1 {
		panic("Constructor must return one value")
	}
	
	controllerInterface, ok := result[0].Interface().(ControllerInterface)
	if !ok {
		panic("Constructor return value must implement ControllerInterface")
	}
	
	return value, controllerInterface
}