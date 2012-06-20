package trinity

import (
	"net/http"
	"reflect"
	"strings"
)

var (
	Post   = method("POST")
	Get    = method("GET")
	Put    = method("PUT")
	Delete = method("DELETE")
	Head   = method("HEAD")

	emptyController = Controller("")
	emptyAction     = Action("")
	emptyValues = make([]reflect.Value, 0)
)

type httpHandler func(response http.ResponseWriter, request *http.Request) ActionResultInterface
type Controller string
type Action string
type method string

func C(s string) Controller {
	return Controller(s)
}
func A(s string) Action {
	return Action(s)
}

// Checks whether controller equals s. Case insensitive
func (c Controller) EqualsS(s string) bool {
	logger.Debug("c.EqualsS %v == %s", c, s)

	return strings.ToLower(s) == strings.ToLower(string(c))
}

// Checks whether action equals s. Case insensitive
func (a Action) EqualsS(s string) bool {
	logger.Debug("a.EqualsS %v == %s", a, s)

	return strings.ToLower(s) == strings.ToLower(string(a))
}

// Pair of controller and action
type ControllerAction struct {
	C Controller
	A Action
}

// Returns true if both controller and action are not empty
func (ca *ControllerAction) IsFull() bool {
	return ca.C != emptyController && ca.A != emptyAction
}

func toLowerA(a Action) Action {
	return Action(strings.ToLower(string(a)))
}

func toLowerC(c Controller) Controller {
	return Controller(strings.ToLower(string(c)))
}

func createURL(c Controller, a Action, params map[string]string) string {
	url := "/" + string(c) + "/" + string(a)

	query := ""
	if params != nil {
		for k, v := range params {
			if len(query) != 0 {
				query += "&"
			}
			query += k + "=" + v
		}
	}

	if len(query) != 0 {
		url += "?" + query
	}

	return url
}

func parseURL(url string) (c Controller, a Action) {
	logger.Trace("")

	if len(url) == 0 {
		return
	}

	pathParts := strings.Split(url[1:], "/")

	if len(pathParts) > 0 {
		c = Controller(pathParts[0])
	}
	if len(pathParts) > 1 {
		a = Action(pathParts[1])
	}

	return
}

func traceParamValue(paramValue reflect.Value) {
	logger.Trace("paramValue: Type.Name= %v Kind= %v", paramValue.Type().Name(), paramValue.Kind())
}

func traceParamType(paramType reflect.Type) {
	logger.Trace("paramType: Name= %v Kind= %v", paramType.Name(), paramType.Kind())
}

func stripPrefix(prefix string, h http.Handler, nf http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, prefix) {
			nf.ServeHTTP(w, r)
			return
		}
		r.URL.Path = r.URL.Path[len(prefix):]
		h.ServeHTTP(w, r)
	})
}