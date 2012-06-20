package trinity

import (
	"net/http"
)

// ControllerInterface represents mvc controller objects.
type ControllerInterface interface {
	SetController(Controller)
	SetAction(Action)
	SetRequest(*http.Request)
	SetResponse(http.ResponseWriter)
	
	GetInfo() ControllerInfoInterface
}

type BaseController struct {
	C Controller
	A Action
	Request *http.Request
	Response http.ResponseWriter
}

func NewBaseController() *BaseController {
	return new(BaseController)
}

func (baseController *BaseController) SetController(c Controller) {
	baseController.C = c
}
func (baseController *BaseController) SetAction(a Action) {
	baseController.A = a
}
func (baseController *BaseController) SetRequest(request *http.Request) {
	baseController.Request = request
}
func (baseController *BaseController) SetResponse(response http.ResponseWriter) {
	baseController.Response = response
}