package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/models"
	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/utils"
)

func (h *Handlers) CreateScheduleTemplate(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Shifts      []struct {
			StartTime          string  `json:"startTime"`
			EndTime            string  `json:"endTime"`
			RequiredAssistants int32   `json:"requiredAssistants"`
			ApplicableDays     []int32 `json:"applicableDays"`
		} `json:"shifts"`
	}

	if err := h.readJSON(r, &payload); err != nil {
		h.errorResponse(w, r, err)
		return
	}

	// validate the input
	if payload.Name == "" {
		h.errorResponse(w, r, errors.New("班表模板名字为空"))
		return
	}

	for id, shift := range payload.Shifts {
		if shift.StartTime == "" {
			h.errorResponse(w, r, fmt.Errorf("班次 %d 的开始时间为空", id))
			return
		}
		if shift.EndTime == "" {
			h.errorResponse(w, r, fmt.Errorf("班次 %d 的结束时间为空", id))
			return
		}
		if shift.RequiredAssistants <= 0 {
			h.errorResponse(w, r, fmt.Errorf("班次 %d 的所需助理数必须大于 0", id))
			return
		}
		if len(shift.ApplicableDays) == 0 {
			h.errorResponse(w, r, fmt.Errorf("班次 %d 的适用日期为空", id))
			return
		}
		for _, day := range shift.ApplicableDays {
			if day < 1 || day > 7 {
				h.errorResponse(w, r, fmt.Errorf("班次 %d 的适用日期 %d 不在 1-7 之间", id, day))
				return
			}
		}
	}

	// create the schedule template instance
	st := &models.ScheduleTemplate{
		Name:        payload.Name,
		Description: payload.Description,
		Shifts:      make([]*models.ScheduleTemplateShift, 0, len(payload.Shifts)),
	}

	for _, shift := range payload.Shifts {
		sts := &models.ScheduleTemplateShift{
			StartTime:          shift.StartTime,
			EndTime:            shift.EndTime,
			RequiredAssistants: shift.RequiredAssistants,
			ApplicableDays:     shift.ApplicableDays,
		}

		sts.RequiredAssistants = shift.RequiredAssistants
		sts.ApplicableDays = shift.ApplicableDays

		st.Shifts = append(st.Shifts, sts)
	}

	// check the validity of the schedule template
	if err := utils.ValidateScheduleTemplate(st); err != nil {
		h.errorResponse(w, r, err)
		return
	}

	// insert the schedule template into the database
	if err := h.models.InsertScheduleTemplate(st); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.ConstraintName == "schedule_templates_name_key" {
			h.errorResponse(w, r, errors.New("班表模板名字重复"))
			return
		} else {
			h.internalServerError(w, r, err)
			return
		}
	}

	h.successResponse(w, r, "班表模板创建成功", st)
}

func (h *Handlers) GetScheduleTemplates(w http.ResponseWriter, r *http.Request) {
	scheduleTemplateIDAsString := chi.URLParam(r, "scheduleTemplateID")
	scheduleTemplateID, err := strconv.ParseInt(scheduleTemplateIDAsString, 10, 64)
	if err != nil {
		h.errorResponse(w, r, errors.New("班表模板 ID 非法"))
		return
	}

	sts, err := h.models.SelectScheduleTemplate(scheduleTemplateID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.errorResponse(w, r, errors.New("班表模板不存在"))
			return
		} else {
			h.internalServerError(w, r, err)
			return
		}
	}

	h.successResponse(w, r, "班表模板获取成功", sts)
}

func (h *Handlers) GetAllScheduleTemplateMeta(w http.ResponseWriter, r *http.Request) {
	sts, err := h.models.SelectAllScheduleTemplateMeta()
	if err != nil {
		h.internalServerError(w, r, err)
		return
	}

	h.successResponse(w, r, "班表模板元数据获取成功", sts)
}

func (h *Handlers) DeleteScheduleTemplate(w http.ResponseWriter, r *http.Request) {
	scheduleTemplateIDAsString := chi.URLParam(r, "scheduleTemplateID")
	scheduleTemplateID, err := strconv.ParseInt(scheduleTemplateIDAsString, 10, 64)
	if err != nil {
		h.errorResponse(w, r, errors.New("班表模板 ID 非法"))
		return
	}

	if err := h.models.DeleteScheduleTemplate(scheduleTemplateID); err != nil {
		h.internalServerError(w, r, err)
		return
	}

	h.successResponse(w, r, "班表模板删除成功", nil)
}
