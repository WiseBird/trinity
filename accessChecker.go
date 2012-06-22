package trinity

import (
	"net/http"
)

// ActionResultInterface represents security objects that check whether this controller/action
// is allowed in context of response, request, app state.
//
// Returns nil if access is allowed
//
// Returns ActionReult if access is denied. In this case ActionReult is used to response
type AccessCheckerInterface interface {
	IsAccessAllowed(c Controller, a Action, response http.ResponseWriter, request *http.Request) ActionResultInterface
}
