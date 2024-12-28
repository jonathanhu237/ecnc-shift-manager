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
		sts := &ScheduleTemplateShift{
			StartTime:          shift.StartTime,
			EndTime:            shift.EndTime,
			RequiredAssistants: shift.RequiredAssistants,
		}
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
		if err := tx.QueryRowContext(ctx, query, st.ID, sts.StartTime, sts.EndTime, sts.RequiredAssistants).Scan(&sts.ID); err != nil {
			return err
		}

		// insert the applicable days
		for _, day := range shift.ApplicableDays {
			query := `
				INSERT INTO schedule_template_shifts_availability (schedule_template_shift_id, day_of_week)
				VALUES ($1, $2)
			`
			if _, err := tx.ExecContext(ctx, query, sts.ID, day); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}
