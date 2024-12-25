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
	ID          int64                    `json:"id"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Shifts      []*ScheduleTemplateShift `json:"shifts"`
	CreatedAt   time.Time                `json:"createdAt"`
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

func (m *ScheduleTemplateModel) InsertScheduleTemplate(name, description string) (*ScheduleTemplate, error) {
	query := `
		INSERT INTO schedule_templates (name, description)
		VALUES ($1, $2)
		RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	st := &ScheduleTemplate{
		Name:        name,
		Description: description,
	}
	if err := m.DB.QueryRowContext(ctx, query, st.Name, st.Description).Scan(&st.ID, &st.CreatedAt); err != nil {
		return nil, err
	}

	return st, nil
}

func (m *ScheduleTemplateModel) CheckTemplateNameValid(name string) (bool, error) {
	query := `
		SELECT EXISTS (SELECT 1 FROM schedule_templates WHERE name = $1)
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var exists bool
	if err := m.DB.QueryRowContext(ctx, query, name).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}
