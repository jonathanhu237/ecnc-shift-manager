package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/models"
	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/utils"
	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/workers"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/wneessen/go-mail"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		FullName string `json:"fullName"`
		Role     string `json:"role"`
	}

	if err := h.readJSON(r, &payload); err != nil {
		h.errorResponse(w, r, err)
		return
	}

	// check the payload
	switch {
	case (payload.Username == ""):
		h.errorResponse(w, r, errors.New("用户名为空"))
		return
	case (payload.Email == ""):
		h.errorResponse(w, r, errors.New("邮箱为空"))
		return
	case (!utils.IsValidEmail(payload.Email)):
		h.errorResponse(w, r, errors.New("邮箱非法"))
		return
	case (payload.FullName == ""):
		h.errorResponse(w, r, errors.New("姓名为空"))
		return
	case (payload.Role == ""):
		h.errorResponse(w, r, errors.New("角色为空"))
		return
	case (!utils.IsValidRole(payload.Role)):
		h.errorResponse(w, r, errors.New("角色非法"))
		return
	}

	// generate random password and hash
	random_password := utils.GenerateRandomPassword(12)
	passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte(random_password), bcrypt.DefaultCost)
	if err != nil {
		h.internalServerError(w, r, err)
		return
	}

	// insert the user to the database
	user := &models.User{
		Username:     payload.Username,
		Email:        payload.Email,
		PasswordHash: string(passwordHashBytes),
		FullName:     payload.FullName,
		Role:         payload.Role,
	}
	if err := h.models.InsertUser(user); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch {
			case (pgErr.ConstraintName == "users_username_key"):
				h.errorResponse(w, r, errors.New("用户名已存在"))
				return
			case (pgErr.ConstraintName == "users_email_key"):
				h.errorResponse(w, r, errors.New("邮箱已存在"))
				return
			default:
				h.internalServerError(w, r, err)
				return
			}
		} else {
			h.internalServerError(w, r, err)
			return
		}
	}

	// send the username and password to the e-mail
	message := mail.NewMsg()
	if err := message.From(h.config.MailClient.Sender); err != nil {
		h.internalServerError(w, r, err)
		return
	}
	if err := message.To(payload.Email); err != nil {
		h.internalServerError(w, r, err)
		return
	}

	mailPayload := workers.MailPayload{
		To:      payload.Email,
		Subject: "ECNC 假勤系统 - 您的账号信息",
		Body:    fmt.Sprintf("用户名: %s, 密码: %s", payload.Username, random_password),
	}
	jsonData, err := json.Marshal(mailPayload)
	if err != nil {
		h.internalServerError(w, r, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := h.emailChan.PublishWithContext(
		ctx,
		"",
		"mail_queue",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonData,
		},
	); err != nil {
		h.internalServerError(w, r, err)
		return
	}

	// return a success message
	h.successResponse(w, r, "创建用户成功", user)
}

func (h *Handlers) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.models.SelectAllUsers()
	if err != nil {
		h.internalServerError(w, r, err)
		return
	}

	h.successResponse(w, r, "获取所有用户信息成功", users)
}

func (h *Handlers) GetUser(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(userCtxKey).(*models.User)
	if !ok {
		h.internalServerError(w, r, errors.New("getUserHandler must be used after GetUserMiddleware"))
		return
	}

	h.successResponse(w, r, "获取用户信息成功", user)
}

func (h *Handlers) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(userCtxKey).(*models.User)
	if !ok {
		h.internalServerError(w, r, errors.New("UpdateUserRoleHandler must be used after GetUserMiddleware"))
		return
	}

	var payload struct {
		Role string `json:"role"`
	}
	if err := h.readJSON(r, &payload); err != nil {
		h.errorResponse(w, r, err)
		return
	}
	if !utils.IsValidRole(payload.Role) {
		h.errorResponse(w, r, errors.New("角色非法"))
		return
	}

	if user.Username == h.config.InitialAdmin.Username {
		h.errorResponse(w, r, errors.New("禁止修改初始管理员角色"))
		return
	}

	user.Role = payload.Role
	if err := h.models.UpdateUser(user); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			h.errorResponse(w, r, errors.New("用户已被修改或删除"))
			return
		default:
			h.internalServerError(w, r, err)
			return
		}
	}

	h.successResponse(w, r, "更新用户身份成功", user)
}

func (h *Handlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(userCtxKey).(*models.User)
	if !ok {
		h.internalServerError(w, r, errors.New("DeleteUser must be used after GetUserMiddleware"))
		return
	}

	if user.Username == h.config.InitialAdmin.Username {
		h.errorResponse(w, r, errors.New("禁止删除初始管理员"))
		return
	}

	if err := h.models.DeleteUser(user.ID); err != nil {
		h.internalServerError(w, r, err)
		return
	}

	h.successResponse(w, r, "删除用户成功", nil)
}
