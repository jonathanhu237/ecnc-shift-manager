package models

import (
	"context"
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
	Shifts      []*ScheduleTemplateShift `json:"shifts,omitempty"`
	CreatedAt   time.Time                `json:"createdAt"`
}

func (m *Models) SelectScheduleTemplates() ([]*ScheduleTemplate, error) {
	query := `
		SELECT id, name, description, created_at
		FROM schedule_templates
		ORDER BY created_at DESC
	`

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	scheduleTemplates := make([]*ScheduleTemplate, 0)
	for rows.Next() {
		scheduleTemplate := &ScheduleTemplate{}

		if err := rows.Scan(
			&scheduleTemplate.ID,
			&scheduleTemplate.Name,
			&scheduleTemplate.Description,
			&scheduleTemplate.CreatedAt,
		); err != nil {
			return nil, err
		}

		scheduleTemplates = append(scheduleTemplates, scheduleTemplate)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return scheduleTemplates, nil
}

func (m *Models) InsertScheduleTemplate(name, description string) (*ScheduleTemplate, error) {
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
	if err := m.db.QueryRowContext(ctx, query, st.Name, st.Description).Scan(&st.ID, &st.CreatedAt); err != nil {
		return nil, err
	}

	return st, nil
}
