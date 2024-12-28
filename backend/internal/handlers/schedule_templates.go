package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/models"
	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/utils"
)

func (h *Handlers) GetAllScheduleTemplateMeta(w http.ResponseWriter, r *http.Request) {
	templates, err := h.models.SelectAllScheduleTemplateMeta()
	if err != nil {
		h.internalServerError(w, r, err)
		return
	}

	h.successResponse(w, r, "获取所有班表模板元数据成功", templates)
}

func (h *Handlers) CreateScheduleTemplate(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Shifts      []struct {
			DayOfWeek          int32     `json:"dayOfWeek"`
			StartTime          time.Time `json:"startTime"`
			EndTime            time.Time `json:"endTime"`
			AssistantsRequired int32     `json:"assistantsRequired"`
		} `json:"shifts"`
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
	template := &models.ScheduleTemplate{
		Name:        payload.Name,
		Description: payload.Description,
		Shifts:      []*models.ScheduleTemplateShift{},
	}

	for _, shift := range payload.Shifts {
		shift := &models.ScheduleTemplateShift{
			DayOfWeek:          shift.DayOfWeek,
			StartTime:          shift.StartTime,
			EndTime:            shift.EndTime,
			AssistantsRequired: shift.AssistantsRequired,
		}
		template.Shifts = append(template.Shifts, shift)
	}
	if err := utils.ValidateShifts(template.Shifts); err != nil {
		h.errorResponse(w, r, err)
		return
	}

	if err := h.models.InsertScheduleTemplate(template); err != nil {
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
