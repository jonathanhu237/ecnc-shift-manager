package application

type appError struct {
	code    int
	message string
}

var (
	errUnauthorized     = appError{code: 401, message: "unauthorized"}
	errForbidden        = appError{code: 403, message: "forbidden"}
	errNotFound         = appError{code: 404, message: "not found"}
	errMethodNotAllowed = appError{code: 405, message: "method not allowed"}
	errInternalServer   = appError{code: 500, message: "internal server error"}

	errInvalidLogin = appError{code: 1001, message: "invalid username or password"}

	errTokenIsExpired = appError{code: 2001, message: "token is expired"}
	errInvalidToken   = appError{code: 2002, message: "invalid token"}

	errUsernameExistsInCreateUser = appError{code: 3001, message: "username already exists"}
	errEmailExistsInCreateUser    = appError{code: 3002, message: "email already exists"}

	errInvalidOldPasswordInResetPassword = appError{code: 4001, message: "invalid old password"}
)
