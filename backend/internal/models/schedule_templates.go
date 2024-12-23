package models

import (
	"context"
	"database/sql"
	"time"
)

type ScheduleTemplateShift struct {
	ID                 int64     `json:"id"`
	DayOfWeek          int       `json:"dayOfWeek"`
	StartTime          time.Time `json:"startTime"`
	EndTime            time.Time `json:"endTime"`
	AssistantsRequired int       `json:"assistantsRequired"`
}

type ScheduleTemplate struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Shifts      []*ScheduleTemplateShift
	CreatedAt   time.Time
}

type ScheduleTemplateModel struct {
	DB *sql.DB
}

func (m *ScheduleTemplateModel) SelectAll() ([]*ScheduleTemplate, error) {
	query := `
		SELECT
			st.id,
			st.name,
			st.description,
			sts.id,
			sts.day_of_week,
			sts.start_time,
			sts.end_time,
			st.created_at
		FROM schedule_templates AS st
		LEFT JOIN schedule_template_shifts AS sts
			ON st.id = sts.schedule_template_id
		ORDER BY st.created_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	scheduleTemplatesMap := make(map[int64]*ScheduleTemplate)
	for rows.Next() {
		var (
			templateID          int64
			templateName        string
			templateDescription string
			shiftID             int64
			shiftDayOfWeek      int
			ShiftStartTime      time.Time
			ShiftEndTime        time.Time
			templateCreatedAt   time.Time
		)

		if err := rows.Scan(
			&templateID,
			&templateName,
			&templateDescription,
			&shiftID,
			&shiftDayOfWeek,
			&ShiftStartTime,
			&ShiftEndTime,
			&templateCreatedAt,
		); err != nil {
			return nil, err
		}

		if _, exists := scheduleTemplatesMap[templateID]; !exists {
			scheduleTemplatesMap[templateID] = &ScheduleTemplate{
				ID:          templateID,
				Name:        templateName,
				Description: templateDescription,
				Shifts:      make([]*ScheduleTemplateShift, 0),
				CreatedAt:   templateCreatedAt,
			}
		}

		scheduleTemplatesMap[templateID].Shifts = append(scheduleTemplatesMap[templateID].Shifts, &ScheduleTemplateShift{
			ID:        shiftID,
			DayOfWeek: shiftDayOfWeek,
			StartTime: ShiftStartTime,
			EndTime:   ShiftEndTime,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	scheduleTemplates := make([]*ScheduleTemplate, 0, len(scheduleTemplatesMap))
	for _, template := range scheduleTemplatesMap {
		scheduleTemplates = append(scheduleTemplates, template)
	}

	return scheduleTemplates, nil
}
