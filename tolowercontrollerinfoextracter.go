package trinity

// ToLowerControllerInfoExtracter is a BaseControllerInfoExtracter extension used to convert action names
// to lower case.
type ToLowerControllerInfoExtracter struct {
	*BaseControllerInfoExtracter
}

// ToLowerControllerInfoExtracter constructor.
func NewToLowerControllerInfoExtracter(value interface{}) *ToLowerControllerInfoExtracter {
	logger.Trace("")

	toLowerController := new(ToLowerControllerInfoExtracter)
	toLowerController.BaseControllerInfoExtracter = NewBaseControllerInfoExtracter(value)

	toLowerController.Controller = toLowerC(toLowerController.Controller)
	for _, actionInfo := range toLowerController.ActionInfos {
		actionInfo.action = toLowerA(actionInfo.action)
	}

	return toLowerController
}

// Adds specific method to action collection.
// Use when you want to add a method with arbitrary signature.
func (toLowerController *ToLowerControllerInfoExtracter) AddAction(methodName string) *ActionInfo {
	logger.Trace("")

	actionInfo := toLowerController.BaseControllerInfoExtracter.AddAction(methodName)
	actionInfo.action = toLowerA(actionInfo.action)
	return actionInfo
}