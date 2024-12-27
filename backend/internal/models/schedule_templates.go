package models

import (
	"context"
	"database/sql"
	"time"
)

type ScheduleTemplateShift struct {
	ID                 int64     `json:"id"`
	DayOfWeek          int32     `json:"dayOfWeek"`
	StartTime          time.Time `json:"startTime"`
	EndTime            time.Time `json:"endTime"`
	AssistantsRequired int32     `json:"assistantsRequired"`
}

type ScheduleTemplate struct {
	ID          int64                    `json:"id"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Shifts      []*ScheduleTemplateShift `json:"shifts,omitempty"`
	CreatedAt   time.Time                `json:"createdAt"`
	Version     int32                    `json:"-"`
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

func (m *Models) SelectScheduleTemplate(id int64) (*ScheduleTemplate, error) {
	query := `
		SELECT st.name, st.description, st.created_at, st.version, sts.id, sts.day_of_week, sts.start_time, sts.end_time, sts.assistants_required
		FROM schedule_templates st
		LEFT JOIN schedule_template_shifts sts ON st.id = sts.schedule_template_id
		WHERE st.id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := m.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	scheduleTemplate := &ScheduleTemplate{
		ID:     id,
		Shifts: make([]*ScheduleTemplateShift, 0),
	}
	for rows.Next() {
		var (
			stName                string
			stDescription         string
			stCreatedAt           time.Time
			stVersion             int32
			stsID                 sql.NullInt64
			stsDayOfWeek          sql.NullInt32
			stsStartTime          sql.NullTime
			stsEndTime            sql.NullTime
			stsAssistantsRequired sql.NullInt32
		)
		if err := rows.Scan(
			&stName,
			&stDescription,
			&stCreatedAt,
			&stVersion,
			&stsID,
			&stsDayOfWeek,
			&stsStartTime,
			&stsEndTime,
			&stsAssistantsRequired,
		); err != nil {
			return nil, err
		}

		if scheduleTemplate.Name == "" {
			scheduleTemplate.Name = stName
			scheduleTemplate.Description = stDescription
			scheduleTemplate.CreatedAt = stCreatedAt
			scheduleTemplate.Version = stVersion
		}

		if stsID.Valid && stsDayOfWeek.Valid && stsStartTime.Valid && stsEndTime.Valid && stsAssistantsRequired.Valid {
			sts := &ScheduleTemplateShift{
				ID:                 stsID.Int64,
				DayOfWeek:          stsDayOfWeek.Int32,
				StartTime:          stsStartTime.Time,
				EndTime:            stsEndTime.Time,
				AssistantsRequired: stsAssistantsRequired.Int32,
			}
			scheduleTemplate.Shifts = append(scheduleTemplate.Shifts, sts)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if scheduleTemplate.Name == "" {
		return nil, sql.ErrNoRows
	}

	return scheduleTemplate, nil
}

func (m *Models) UpdateScheduleTemplateMeta(st *ScheduleTemplate) error {
	query := `
		UPDATE schedule_templates
		SET description = $1, version = version + 1
		WHERE id = $2 AND version = $3
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := m.db.ExecContext(ctx, query, st.Description, st.ID, st.Version)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
