package models

import (
	"context"
	"time"
)

type ScheduleTemplateShift struct {
	ID                 int64     `json:"id"`
	StartTime          time.Time `json:"startTime"`
	EndTime            time.Time `json:"endTime"`
	RequiredAssistants int32     `json:"requiredAssistants"`
	ApplicableDays     []int32   `json:"applicableDays"`
}

type ScheduleTemplate struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Shifts      []*ScheduleTemplateShift
	CreatedAt   time.Time `json:"createdAt"`
	Version     int32     `json:"version"`
}

func (m *Models) InsertScheduleTemplate(st *ScheduleTemplate) error {
	// begin transaction
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// insert the meta
	query := `
		INSERT INTO schedule_templates (name, description)
		VALUES ($1, $2)
		RETURNING id, created_at, version
	`
	if err := tx.QueryRowContext(ctx, query, st.Name, st.Description).Scan(&st.ID, &st.CreatedAt, &st.Version); err != nil {
		return err
	}

	// insert the shifts
	for _, shift := range st.Shifts {
		query := `
			INSERT INTO 
				schedule_template_shifts (
					schedule_template_id, 
					start_time,
					end_time,
					required_assistants
				)
			VALUES ($1, $2, $3, $4)
			RETURNING id
		`
		if err := tx.QueryRowContext(ctx, query, st.ID, shift.StartTime, shift.EndTime, shift.RequiredAssistants).Scan(&shift.ID); err != nil {
			return err
		}

		// insert the applicable days
		for _, day := range shift.ApplicableDays {
			query := `
				INSERT INTO schedule_template_shifts_availability (schedule_template_shift_id, day_of_week)
				VALUES ($1, $2)
			`
			if _, err := tx.ExecContext(ctx, query, shift.ID, day); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (m *Models) SelectScheduleTemplate(id int64) (*ScheduleTemplate, error) {
	st := &ScheduleTemplate{
		ID:     id,
		Shifts: make([]*ScheduleTemplateShift, 0),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// query the meta
	query := `
		SELECT name, description, created_at, version
		FROM schedule_templates
		WHERE id = $1
	`
	if err := m.db.QueryRowContext(ctx, query, id).Scan(&st.Name, &st.Description, &st.CreatedAt, &st.Version); err != nil {
		return nil, err
	}

	// query the shifts
	query = `
		SELECT id, start_time, end_time, required_assistants
		FROM schedule_template_shifts
		WHERE schedule_template_id = $1
	`
	rows, err := m.db.QueryContext(ctx, query, st.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		sts := &ScheduleTemplateShift{
			ApplicableDays: make([]int32, 0),
		}
		if err := rows.Scan(&sts.ID, &sts.StartTime, &sts.EndTime, &sts.RequiredAssistants); err != nil {
			return nil, err
		}

		// query the applicable days
		query := `
			SELECT day_of_week
			FROM schedule_template_shifts_availability
			WHERE schedule_template_shift_id = $1
		`

		rows, err := m.db.QueryContext(ctx, query, sts.ID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var day_of_week int32
			if err := rows.Scan(&day_of_week); err != nil {
				return nil, err
			}
			sts.ApplicableDays = append(sts.ApplicableDays, day_of_week)
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return st, nil
}
