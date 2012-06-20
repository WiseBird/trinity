package trinity

import (
	"net/http"
)

// NotFoundHandler implements the http.Handler interface and provides a NotFoundActionResult response.
type NotFoundHandler struct {
	mvcI *MvcInfrastructure
}

func NewNotFoundHandler(mvcI *MvcInfrastructure) *NotFoundHandler {
	handler := new(NotFoundHandler)

	handler.mvcI = mvcI

	return handler
}

func (handler *NotFoundHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	NotFoundResult().Response(handler.mvcI, emptyController, emptyAction, response, request)
}
