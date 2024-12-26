package handlers

import (
	"errors"
	"net/http"
)

func (h *Handlers) GetAllScheduleTemplates(w http.ResponseWriter, r *http.Request) {
	templates, err := h.models.SelectScheduleTemplates()
	if err != nil {
		h.internalServerError(w, r, err)
		return
	}

	h.successResponse(w, r, "获取班表模板成功", templates)
}

func (h *Handlers) CreateScheduleTemplate(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := h.readJSON(r, &payload); err != nil {
		h.errorResponse(w, r, err)
		return
	}
	if payload.Name == "" {
		h.errorResponse(w, r, errors.New("班表模板名为空"))
		return
	}

	// insert the template
	template, err := h.models.InsertScheduleTemplate(payload.Name, payload.Description)
	if err != nil {
		h.internalServerError(w, r, err)
		return
	}

	h.successResponse(w, r, "创建班表模板成功", template)
}