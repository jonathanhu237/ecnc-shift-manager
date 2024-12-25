package application

type appError struct {
	code    int
	message string
}

var (
	errUnauthorized                      = appError{code: 401, message: "用户未登录"}
	errForbidden                         = appError{code: 403, message: "无权访问"}
	errNotFound                          = appError{code: 404, message: "未找到资源"}
	errMethodNotAllowed                  = appError{code: 405, message: "请求方式非法"}
	errInternalServer                    = appError{code: 500, message: "服务器内部错误"}
	errInvalidLogin                      = appError{code: 1001, message: "用户名或密码错误"}
	errTokenIsExpired                    = appError{code: 2001, message: "访问令牌过期"}
	errInvalidToken                      = appError{code: 2002, message: "无效的访问令牌"}
	errUsernameExistsInCreateUser        = appError{code: 3001, message: "用户已存在"}
	errEmailExistsInCreateUser           = appError{code: 3002, message: "邮箱已存在"}
	errInvalidOldPasswordInResetPassword = appError{code: 4001, message: "旧密码错误"}
	errUpdateInitialBlackCoreRole        = appError{code: 5001, message: "禁止修改初始黑心的身份"}
	errDeleteInitialBlackCore            = appError{code: 6001, message: "禁止删除初始黑心"}
	errScheduleTemplateExists            = appError{code: 7001, message: "排班模板已存在"}
)
