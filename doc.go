//

/*
A helper package which extends net/http and Gorilla web toolkit to
simplify MVC web apps creation. 

Routing

1. Automatic: during controller/action binding each action for each controller is mapped to "controller/action" URL.

2. Manual: BindUrl method could be used to bind any specific URL to a specified controller/action.

Controllers

Controller is a type which implements ControllerInterface. It is
registered using a controller bind call with constructor:

mvcI.BindController(NewMyController)

BindController expects parameterless constructor of ControllerInterface which provides the necessary information.

To simplify controller declaration helpers can be used: BaseController and ToLowerControllerInfoExtracter.
Example:

	type MyController struct {
		*mvc.BaseController
	}

	func NewMyController() *MyController {
		myController := new(MyController)
		myController.BaseController = mvc.NewBaseController()
		return myController
	}
	
	func (myController *MyController) GetInfo() mvc.ControllerInfoInterface {
		return mvc.NewToLowerControllerInfoExtracter(myController)
	}

	func (myController *MyController) SomeAction() mvc.ActionResultInterface {
		...
	}

	func (myController *MyController) SomeMethod() {
		...
	}

Type name is used for controller name and if it ends with "controller" then suffix will be discarded.
Actions must return mvc.ActionResultInterface to be detected in the BaseController constructor.

For example from above:
	GetController() returns "my"
	GetActionInfos() returns { "SomeAction" : { _handler_, "GET", "someaction" } }

So if we call mvcI.BindController(NewMyController) then action "my/someaction" will be binded.

If needed, action with any signature can be added manually using AddAction method:

	type MyController struct {
		*mvc.BaseController
	}

	func NewMyController() *MyController {
		myController := new(MyController)
		myController.BaseController = mvc.NewBaseController()
		return myController
	}
	
	func (myController *MyController) GetInfo() mvc.ControllerInfoInterface {
		myControllerInfo := mvc.NewToLowerControllerInfoExtracter(myController)
		myControllerInfo.AddAction("Hello")
		return myControllerInfo
	}

	func (myController *MyController) Hello(response http.ResponseWriter) {
		response.Write([]byte("<html><body>Hello!</body></html>"))
	}

Views 

Views are implemented using html/template. Views are registered by ParseViewsFolder,
which takes views folder path as an argument.

Any file with its extension listed in ViewsSuffix can be used as a view.

View folder structure:
	Controller1
		->  Action1.ghtml
		->  Action2.ghtml
	Controller2
		->  Action1.ghtml

View page should consist of two parts: a define for options and a define for contents.
Options are used to set master page, additional templates, etc.

Options define name must start ViewOptionsTemplate constant, but should be different
for different views to avoid "redefinition of template" if multiple templates are
included.

File paths in Options section are set relative to the path passed to the ParseViewsFolder.

Example:

	{{define "ViewOptions_MyView"}}
		MasterPage=master.ghtml
		AdditionalTemplate=shared/additionalTemplate.pghtml
	{{end}}

*/
package trinity
