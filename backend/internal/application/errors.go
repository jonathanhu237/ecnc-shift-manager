package application

import "errors"

type appError struct {
	code    int
	message string
}

var (
	errUnauthorized                      = errors.New("用户未登录")
	errForbidden                         = errors.New("无权访问")
	errNotFound                          = errors.New("未找到资源")
	errMethodNotAllowed                  = errors.New("请求方式非法")
	errInternalServer                    = errors.New("服务器内部错误")
	errInvalidLogin                      = errors.New("用户名或密码错误")
	errTokenIsExpired                    = errors.New("访问令牌过期")
	errInvalidToken                      = errors.New("无效的访问令牌")
	errUsernameExistsInCreateUser        = errors.New("用户已存在")
	errEmailExistsInCreateUser           = errors.New("邮箱已存在")
	errInvalidOldPasswordInResetPassword = errors.New("旧密码错误")
	errUpdateInitialBlackCoreRole        = errors.New("禁止修改初始黑心的身份")
	errDeleteInitialBlackCore            = errors.New("禁止删除初始黑心")
	errScheduleTemplateExists            = errors.New("排班模板已存在")
)
