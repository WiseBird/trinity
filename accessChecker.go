package trinity

import (
	"net/http"
)

// ActionResultInterface represents security objects that check whether this controller/action
// is allowed in context of response, request, app state.
//
// Returns nil, nil, true if access is allowed
//
// Returns cTo, aTo, false if access is denied. In this case cTo, aTo are the redirected
// controller/action (usually authorization/authentication or error controller/action)
type AccessCheckerInterface interface {
	IsAccessAllowed(c Controller, a Action, response http.ResponseWriter, request *http.Request) (Controller, Action, bool)
}
