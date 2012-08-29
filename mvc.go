package trinity

import (
	"code.google.com/p/gorilla/mux"
	"fmt"
	"net/http"
	"reflect"
)

// MvcInfrastructure stores the data needed to create and support the MVC environment.
// Used to register controllers, actions, views.
type MvcInfrastructure struct {
	accessChecker AccessCheckerInterface // used to check access to specific controller/action pairs
	viewsFolder   string                 // path to the views folder

	notFoundView      *ControllerAction // used to show the url-not-found error
	internalErrorView *ControllerAction // used to show internal server errors

	handlers    map[Controller]map[Action]map[method]*methodDescriptor // action handlers
	views       map[Controller]map[Action]*templateDescriptor          // views
	controllerConstructors map[Controller]reflect.Value // Controller ctors


	Router *mux.Router // The main routing object
}

// NewMvcInfrastructure creates a new MvcInfrastructure object
func NewMvcInfrastructure() *MvcInfrastructure {
	mvcI := new(MvcInfrastructure)

	mvcI.handlers = make(map[Controller]map[Action]map[method]*methodDescriptor, 0)
	mvcI.views = make(map[Controller]map[Action]*templateDescriptor, 0)
	//mvcI.controllers = make([]ControllerInterface, 0)
	mvcI.controllerConstructors = make(map[Controller]reflect.Value, 0)

	mvcI.Router = mux.NewRouter()
	mvcI.Router.NotFoundHandler = NewNotFoundHandler(mvcI)

	return mvcI
}

// SetAccessChecker sets the access checker object
func (mvcI *MvcInfrastructure) SetAccessChecker(accessChecker AccessCheckerInterface) {
	mvcI.accessChecker = accessChecker
}

// Если запрашиваемый URL начинается с указанного префикса, то выдает статические файлы находящиеся по указанному пути.
// Для всех остальных запросов отвечает mvcI.Router.NotFoundHandler
func (mvcI *MvcInfrastructure) ServeStatic(prefix string, path string) {
	h := stripPrefix(prefix, http.FileServer(http.Dir(path)), mvcI.Router.NotFoundHandler)
	mvcI.Router.NewRoute().Handler(h)
}

// ParseViewsFolder iterates through files in the current view folder and registers 
// views from it.
func (mvcI *MvcInfrastructure) ParseViewsFolder(viewsFolder string) error {
	logger.Trace("")

	mvcI.viewsFolder = viewsFolder

	err := newViewFolderParser(mvcI).parse()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

func (mvcI *MvcInfrastructure) bindView(c Controller, a Action, templatePath string) error {
	logger.Trace("")

	//c = toLowerC(c)
	//a = toLowerA(a)

	logger.Debugf("c = %v, a = %v", c, a)

	template, err := newTemplateDescriptor(mvcI.viewsFolder, templatePath)
	if err != nil {
		logger.Errorf("%v", err)
		return err
	}

	actions, exists := mvcI.views[c]
	if !exists {
		actions = make(map[Action]*templateDescriptor, 0)
		mvcI.views[c] = actions
	}

	actions[a] = template
	return nil
}

// SetNotFoundView sets the url-not-found view.
func (mvcI *MvcInfrastructure) SetNotFoundView(notFoundView *ControllerAction) {
	if notFoundView != nil && !notFoundView.IsFull() {
		panic("NotFoundView must contains controller and action")
	}

	mvcI.notFoundView = notFoundView
}

// SetInternalErrorView sets the view shown on internal server errors.
func (mvcI *MvcInfrastructure) SetInternalErrorView(internalErrorView *ControllerAction) {
	if internalErrorView != nil && !internalErrorView.IsFull() {
		panic("InternalErrorView must contains controller and action")
	}

	mvcI.internalErrorView = internalErrorView
}

func defaultNotFound(response http.ResponseWriter, request *http.Request) {
	logger.Trace("")
	response.Write([]byte("<html><body>Not found</body></html>"))
}

func defaultInternalError(response http.ResponseWriter, request *http.Request, err interface{}) {
	logger.Debugf("%v", err)
	response.Write([]byte(fmt.Sprintf("<html><body>Internal Error: %v</body></html>", err)))
}

// wrapHandler creates and internal func for a specified controller/action pair
func (mvcI *MvcInfrastructure) wrapHandler(c Controller, a Action) func(response http.ResponseWriter, request *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		logger.Debugf("c: %v, a: %v, m: %s", c, a, request.Method)

		mvcI.handleRequest(c, a, response, request)
	}
}

func (mvcI *MvcInfrastructure) handleRequest(c Controller, a Action, response http.ResponseWriter, request *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			logger.Trace("recovered from panic")
			ErrorResult(err).Response(mvcI, c, a, response, request)
		}
	}()
	
	var res ActionResultInterface
	
	if mvcI.accessChecker != nil {
		logger.Trace("check access")
		res = mvcI.accessChecker.IsAccessAllowed(c, a, response, request)
		if res != nil {
			logger.Trace("access denied")
		}
	}

	if res == nil {
		res = mvcI.callAction(c, a, response, request)
	}
	
	if res != nil {
		res.Response(mvcI, c, a, response, request)
	}

	logger.Trace("Done")
}
