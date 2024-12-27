package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/models"
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch {
			case pgErr.ConstraintName == "schedule_templates_name_key":
				h.errorResponse(w, r, errors.New("班表模板名已存在"))
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

	h.successResponse(w, r, "创建班表模板成功", template)
}

func (h *Handlers) GetScheduleTemplate(w http.ResponseWriter, r *http.Request) {
	scheduleTemplate, ok := r.Context().Value(scheduleTemplateKey).(*models.ScheduleTemplate)
	if !ok {
		h.internalServerError(w, r, errors.New("GetScheduleTemplate must be called after GetScheduleTemplateMiddleware"))
		return
	}

	h.successResponse(w, r, "获取班表模板成功", scheduleTemplate)
}

func (h *Handlers) UpdateScheduleTemplateDescription(w http.ResponseWriter, r *http.Request) {
	scheduleTemplate, ok := r.Context().Value(scheduleTemplateKey).(*models.ScheduleTemplate)
	if !ok {
		h.internalServerError(w, r, errors.New("UpdateScheduleTemplateDescription must be called after GetScheduleTemplateMiddleware"))
		return
	}

	var payload struct {
		Description string `json:"description"`
	}

	if err := h.readJSON(r, &payload); err != nil {
		h.errorResponse(w, r, err)
		return
	}
	scheduleTemplate.Description = payload.Description

	if err := h.models.UpdateScheduleTemplateMeta(scheduleTemplate); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			h.errorResponse(w, r, errors.New("班表模板已被修改或删除，请重试"))
			return
		default:
			h.internalServerError(w, r, err)
			return
		}
	}

	h.successResponse(w, r, "更新班表模板成功", scheduleTemplate)
}
