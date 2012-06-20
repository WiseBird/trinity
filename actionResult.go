package trinity

import (
	"errors"
	"net/http"
)

// ActionResultInterface generates the http-response using a given mvc infrastructure,
// controller, action, and request and writes it to the given response writer
type ActionResultInterface interface {
	Response(mvcI *MvcInfrastructure, c Controller, a Action, response http.ResponseWriter, request *http.Request)
}

// Error action result
type ErrorActionResult struct {
	err  error
	view bool // true, if the error occured during the error view generation.
	// an internal state variable used to avoid infinite loops
}

// Generates an error result representing a specified error
func ErrorResult(err error) ActionResultInterface {
	logger.Trace("")
	logger.Debug("err: %s", err.Error())

	return &ErrorActionResult{err, false}
}
func viewErrorResult(err error) ActionResultInterface {
	logger.Trace("")
	logger.Debug("err: %s", err.Error())

	return &ErrorActionResult{err, true}
}
func (result *ErrorActionResult) Response(mvcI *MvcInfrastructure, c Controller, a Action, response http.ResponseWriter, request *http.Request) {
	logger.Trace("")

	response.WriteHeader(500)

	if mvcI.internalErrorView == nil {
		defaultInternalError(response, request, result.err)
		return
	}

	if result.view && (c == mvcI.internalErrorView.C) && (a == mvcI.internalErrorView.A) {
		defaultInternalError(response, request, result.err)
		return
	}

	ShowView(mvcI.internalErrorView.C, mvcI.internalErrorView.A, result.err).
		Response(mvcI, mvcI.internalErrorView.C, mvcI.internalErrorView.A, response, request)
}

// Resource-not-found action result
type NotFoundActionResult struct {
}

// Generates a resource-not-found action result 
func NotFoundResult() ActionResultInterface {
	logger.Trace("")

	return &NotFoundActionResult{}
}
func (result *NotFoundActionResult) Response(mvcI *MvcInfrastructure, c Controller, a Action, response http.ResponseWriter, request *http.Request) {
	logger.Trace("")

	response.WriteHeader(404)

	if mvcI.notFoundView == nil {
		defaultNotFound(response, request)
		return
	}

	ShowView(mvcI.notFoundView.C, mvcI.notFoundView.A, request.URL.String()).
		Response(mvcI, mvcI.notFoundView.C, mvcI.notFoundView.A, response, request)
}

// Action result that performs a redirect to another controller/action
type RedirectToActionResult struct {
	c      Controller
	a      Action
	params map[string]string
}

// Creates a redirect action result
func RedirectToAction(c Controller, a Action, params map[string]string) ActionResultInterface {
	logger.Trace("")

	//c = toLowerC(c)
	//a = toLowerA(a)

	logger.Debug("c: %v, a: %v, p: %v", c, a, params)

	return &RedirectToActionResult{c, a, params}
}
func (result *RedirectToActionResult) Response(mvcI *MvcInfrastructure, c Controller, a Action, response http.ResponseWriter, request *http.Request) {
	logger.Trace("")

	if len(result.c) > 0 {
		c = result.c
	}
	if len(result.a) > 0 {
		a = result.a
	}

	response.Header().Set("Location", createURL(c, a, result.params))
	response.WriteHeader(302)
}

// ShowViewResult generates http-response using the view associated with the
// specified controller/action. Passes the specified vm as the view data.
type ShowViewResult struct {
	c  Controller
	a  Action
	vm interface{}
}

// Creates a ShowViewResult using the specified controller, action, vm. See the description of
// ShowViewResult.
func ShowView(c Controller, a Action, vm interface{}) ActionResultInterface {
	logger.Trace("")

	//c = toLowerC(c)
	//a = toLowerA(a)

	logger.Debug("c: %v, a: %v", c, a)

	return &ShowViewResult{c, a, vm}
}
func (result *ShowViewResult) Response(mvcI *MvcInfrastructure, c Controller, a Action, response http.ResponseWriter, request *http.Request) {
	logger.Trace("View: res.c = %v, res.a = %v, c = %v, a = %v", result.c, result.a, c, a)
	if len(result.c) > 0 {
		c = result.c
	}
	if len(result.a) > 0 {
		a = result.a
	}

	logger.Trace("get views")
	actions, exists := mvcI.views[c]
	if !exists {
		logger.Error("Controller not found: %v", c)
		viewErrorResult(errors.New("Controller not found")).Response(mvcI, c, a, response, request)
		return
	}

	logger.Trace("get actions")
	template, exists := actions[a]
	if !exists {
		logger.Error("Action not found: %v", a)
		viewErrorResult(errors.New("Action not found")).Response(mvcI, c, a, response, request)
		return
	}

	logger.Trace("render")
	html, err := renderPage(result.vm, template)
	if err != nil {
		logger.Error("%v", err)
		viewErrorResult(err).Response(mvcI, c, a, response, request)
		return
	}
	response.Write(html)
}
