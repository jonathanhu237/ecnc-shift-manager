package models

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type ScheduleTemplateShift struct {
	ID                 int64     `json:"id,omitempty"`
	DayOfWeek          int32     `json:"dayOfWeek"`
	StartTime          time.Time `json:"startTime"`
	EndTime            time.Time `json:"endTime"`
	AssistantsRequired int32     `json:"assistantsRequired"`
}

type ScheduleTemplate struct {
	ID          int64                    `json:"id"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Shifts      []*ScheduleTemplateShift `json:"shifts"`
	CreatedAt   time.Time                `json:"createdAt"`
	Version     int32                    `json:"version"`
}

func (m *Models) SelectAllScheduleTemplateMeta() ([]*ScheduleTemplate, error) {
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

func (m *Models) InsertScheduleTemplate(st *ScheduleTemplate) error {
	// begin a transaction
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// insert the schedule template meta (e.g. title and description)
	query := `
		INSERT INTO schedule_templates (name, description)
		VALUES ($1, $2)
		RETURNING id, created_at, version
	`
	if err := tx.QueryRowContext(ctx, query, st.Name, st.Description).Scan(&st.ID, &st.CreatedAt, &st.Version); err != nil {
		return err
	}

	// insert the schedule template shifts
	if len(st.Shifts) == 0 {
		return tx.Commit()
	}

	var placeholders []string
	var params []interface{}

	for id, shift := range st.Shifts {
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", id*5+1, id*5+2, id*5+3, id*5+4, id*5+5))
		params = append(params, st.ID, shift.DayOfWeek, shift.StartTime, shift.EndTime, shift.AssistantsRequired)
	}

	query = fmt.Sprintf(`
		INSERT INTO schedule_template_shifts (
			schedule_template_id,
			day_of_week,
			start_time,
			end_time,
			assistants_required
		) VALUES %s
	`, strings.Join(placeholders, ", "))
	_, err = tx.ExecContext(ctx, query, params...)
	if err != nil {
		return err
	}

	// commit
	return tx.Commit()
}

func (m *Models) SelectScheduleTemplate(id int64) (*ScheduleTemplate, error) {
	query := `
		SELECT st.name, st.description, st.created_at, st.version, sts.day_of_week, sts.start_time, sts.end_time, sts.assistants_required
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

		if stsDayOfWeek.Valid && stsStartTime.Valid && stsEndTime.Valid && stsAssistantsRequired.Valid {
			sts := &ScheduleTemplateShift{
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

func (m *Models) UpdateScheduleTemplateShifts(st *ScheduleTemplate) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// delete all shifts
	query := `DELETE FROM schedule_template_shifts WHERE schedule_template_id = $1`
	if _, err := tx.ExecContext(ctx, query, st.ID); err != nil {
		return err
	}

	// insert new shifts
	query = `
		INSERT INTO schedule_template_shifts (schedule_template_id, day_of_week, start_time, end_time, assistants_required)
		VALUES ($1, $2, $3, $4, $5) 
	`
	for _, shift := range st.Shifts {
		args := []any{st.ID, shift.DayOfWeek, shift.StartTime, shift.EndTime, shift.AssistantsRequired}

		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			return err
		}
	}

	// update version
	query = `
		UPDATE schedule_templates
		SET version = version + 1
		WHERE id = $1 AND version = $2
	`
	res, err := m.db.ExecContext(ctx, query, st.ID, st.Version)
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

	return tx.Commit()
}
