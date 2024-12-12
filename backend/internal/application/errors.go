package application

import "net/http"

type appError struct {
	code    int
	message string
}

func (app *Application) badRequest(message string) appError {
	return appError{code: http.StatusBadRequest, message: message}
}

var (
	errMethodNotAllowed  = appError{code: 405, message: "method not allowed"}
	errInternalServer    = appError{code: 500, message: "internal server error"}
	errInvalidLogin      = appError{code: 1001, message: "invalid username or password"}
	errAuthHeaderNotSet  = appError{code: 2001, message: "authorization header is not set"}
	errInvalidAuthHeader = appError{code: 2002, message: "authorization header is invalid"}
	errTokenIsExpired    = appError{code: 2003, message: "token is expired"}
	errInvalidToken      = appError{code: 2004, message: "invalid token"}
)
