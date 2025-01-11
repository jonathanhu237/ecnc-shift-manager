package handlers

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/models"
)

func (h *Handlers) CreateSchedulePlan(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Name                 string    `json:"name" validate:"required"`
		Description          string    `json:"description"`
		SubmissionStartTime  time.Time `json:"submissionStartTime" validate:"required"`
		SubmissionEndTime    time.Time `json:"submissionEndTime" validate:"required"`
		ActiveStartTime      time.Time `json:"activeStartTime" validate:"required"`
		ActiveEndTime        time.Time `json:"activeEndTime" validate:"required"`
		ScheduleTemplateName string    `json:"scheduleTemplateName" validate:"required"`
	}
	if err := h.readJSON(r, &payload); err != nil {
		h.errorResponse(w, r, err)
		return
	}
	if err := h.validate.Struct(payload); err != nil {
		h.errorResponse(w, r, err)
		return
	}

	sp := &models.SchedulePlan{
		Name:                 payload.Name,
		Description:          payload.Description,
		SubmissionStartTime:  payload.SubmissionStartTime,
		SubmissionEndTime:    payload.SubmissionEndTime,
		ActiveStartTime:      payload.ActiveStartTime,
		ActiveEndTime:        payload.ActiveEndTime,
		ScheduleTemplateName: payload.ScheduleTemplateName,
	}
	if err := h.models.InsertSchedulePlan(sp); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.ConstraintName == "schedule_plans_name_key" {
				h.errorResponse(w, r, errors.New("排班计划名已存在"))
				return
			} else if pgErr.ConstraintName == "schedule_plans_schedule_template_name_fkey" {
				h.errorResponse(w, r, errors.New("排班模板不存在"))
				return
			}
			h.internalServerError(w, r, err)
			return
		}
		h.internalServerError(w, r, err)
		return
	}

	h.successResponse(w, r, "创建排班计划成功", sp)
}

func (h *Handlers) GetSchedulePlanMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		schedulePlanIDParam := chi.URLParam(r, "schedulePlanID")
		schedulePlanID, err := uuid.Parse(schedulePlanIDParam)
		if err != nil {
			h.errorResponse(w, r, errors.New("无效的排班计划ID"))
			return
		}

		schedulePlan, err := h.models.SelectSchedulePlanByID(schedulePlanID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				h.errorResponse(w, r, errors.New("排班计划不存在"))
				return
			}
			h.internalServerError(w, r, err)
			return
		}

		ctx := context.WithValue(r.Context(), schedulePlanKey, schedulePlan)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handlers) GetSchedulePlan(w http.ResponseWriter, r *http.Request) {
	schedulePlan, ok := r.Context().Value(schedulePlanKey).(*models.SchedulePlan)
	if !ok {
		h.internalServerError(w, r, errors.New("GetSchedulePlan must be used after GetSchedulePlanMiddleware"))
		return
	}

	h.successResponse(w, r, "获取排班计划成功", schedulePlan)
}
