package handlers

import (
	"errors"
	"net/http"

	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handlers) GetMyInfo(w http.ResponseWriter, r *http.Request) {
	myInfo, ok := r.Context().Value(requesterCtxKey).(*models.User)
	if !ok {
		panic("GetMyInfoHandler should be used after GetRequesterMiddleware")
	}

	h.successResponse(w, r, "获取个人信息成功", myInfo)
}

func (h *Handlers) UpdateMyPassword(w http.ResponseWriter, r *http.Request) {
	requester, ok := r.Context().Value(requesterCtxKey).(*models.User)
	if !ok {
		panic("UpdateMyPasswordHandler should be used after GetMyInfoHandler")
	}

	var payload struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}
	if err := h.readJSON(r, &payload); err != nil {
		h.errorResponse(w, r, err)
		return
	}
	switch {
	case payload.OldPassword == "":
		h.errorResponse(w, r, errors.New("旧密码为空"))
		return
	case payload.NewPassword == "":
		h.errorResponse(w, r, errors.New("新密码为空"))
		return
	}

	// verify the old password
	if err := bcrypt.CompareHashAndPassword([]byte(requester.PasswordHash), []byte(payload.OldPassword)); err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			h.errorResponse(w, r, errors.New("密码错误"))
			return
		default:
			h.internalServerError(w, r, err)
			return
		}
	}

	// update the password
	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(payload.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		h.internalServerError(w, r, err)
		return
	}

	requester.PasswordHash = string(newPasswordHash)
	if err := h.models.UpdateUser(requester); err != nil {
		h.internalServerError(w, r, err)
		return
	}

	// response
	h.successResponse(w, r, "修改密码成功", nil)
}
