package trinity

import (
	"net/http"
)

func (mvcI *MvcInfrastructure) bindAction(c Controller, a Action, m method, handler *methodDescriptor) {
	logger.Trace("")

	//c = toLowerC(c)
	//a = toLowerA(a)

	logger.Debugf("c: %v, a: %v, m: %v", c, a, m)

	mvcI.checkViewOnActionBind(c, a)

	actions, exists := mvcI.handlers[c]
	if !exists {
		logger.Trace("new Controller")
		actions = make(map[Action]map[method]*methodDescriptor, 0)
		mvcI.handlers[c] = actions
	}

	methods, exists := actions[a]
	if !exists {
		logger.Trace("new Action")
		methods = make(map[method]*methodDescriptor, 0)
		actions[a] = methods
	}

	methods[m] = handler

	url := createURL(c, a, nil)
	logger.Debugf("URL = %s", url)
	mvcI.Router.HandleFunc(url, mvcI.wrapHandler(c, a)).Methods(string(m))
}

func (mvcI *MvcInfrastructure) checkViewOnActionBind(c Controller, a Action) {
	actions, exists := mvcI.views[c]
	if !exists {
		logger.Warnf("Added handler for missing view: controller - %v, action - %v", c, a)
		return
	}

	_, exists = actions[a]
	if !exists {
		logger.Warnf("Added handler for missing view: controller - %v, action - %v", c, a)

		return
	}
}

// BindUrl bind a url to the controller/action pair. Works on top of the
// Gorilla HandleFunc method.
func (mvcI *MvcInfrastructure) BindUrl(c Controller, a Action, url string) {
	logger.Trace("")
	logger.Debugf("c: %v, a: %v, url: %v", c, a, url)

	mvcI.Router.HandleFunc(url, mvcI.wrapHandler(c, a))
}

func (mvcI *MvcInfrastructure) checkHandlerOnUrlBind(c Controller, a Action) {
	actions, exists := mvcI.handlers[c]
	if !exists {
		logger.Warnf("Added url for missing handler: controller - %v, action - %v", c, a)
		return
	}

	_, exists = actions[a]
	if !exists {
		logger.Warnf("Added url for missing handler: controller - %v, action - %v", c, a)
		return
	}
}

func (mvcI *MvcInfrastructure) callAction(c Controller, a Action, response http.ResponseWriter, request *http.Request) ActionResultInterface {
	logger.Tracef("c: %v, a: %v", c, a)

	//c = toLowerC(c)
	//a = toLowerA(a)

	logger.Trace("Search actions")
	actions, exists := mvcI.handlers[c]
	if !exists {
		logger.Errorf("controller not found: %v", c)
		return NotFoundResult()
	}

	logger.Trace("Search methods")
	methods, exists := actions[a]
	if !exists {
		logger.Errorf("action not found: %v", a)
		return NotFoundResult()
	}

	logger.Trace("Search handler")
	handler, exists := methods[method(request.Method)]
	if exists {
		return mvcI.callHandler(handler, c, a, response, request)
	}

	logger.Trace("Search handler for get")
	handler, exists = methods[Get]
	if exists {
		return mvcI.callHandler(handler, c, a, response, request)
	}

	return NotFoundResult()
}

func (mvcI *MvcInfrastructure) callHandler(handler *methodDescriptor, c Controller, a Action, response http.ResponseWriter, request *http.Request) ActionResultInterface {
	logger.Trace("")

	invoker := newHandlerInvoker(handler).
		AddParam(response).
		AddParam(request).
		AddParam(c).
		AddParam(a).
		AddValues(request.URL.Query())

	if request.Method == "POST" {
		err := request.ParseForm()
		if err != nil {
			logger.Error(err.Error())
			return ErrorResult(err)
		}

		invoker.AddValues(request.Form)
	}

	contr := mvcI.controllerConstructors[c].Call(emptyValues)[0].Interface().(ControllerInterface)
	contr.SetController(c)
	contr.SetAction(a)
	contr.SetRequest(request)
	contr.SetResponse(response)
	invoker.AddParam(contr)

	invoker.AddValue("Controller", string(c)).
		AddValue("Action", string(a))

	return invoker.Invoke()
}
