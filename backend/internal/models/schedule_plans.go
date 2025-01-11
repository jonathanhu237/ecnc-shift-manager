package models

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type SchedulePlan struct {
	ID                   uuid.UUID `json:"id"`
	Name                 string    `json:"name"`
	Description          string    `json:"description"`
	SubmissionStartTime  time.Time `json:"submissionStartTime"`
	SubmissionEndTime    time.Time `json:"submissionEndTime"`
	ActiveStartTime      time.Time `json:"activeStartTime"`
	ActiveEndTime        time.Time `json:"activeEndTime"`
	ScheduleTemplateName string    `json:"scheduleTemplateName"`
	CreatedAt            time.Time `json:"created_at"`
	Version              int32     `json:"version"`
}

func (m *Models) InsertSchedulePlan(sp *SchedulePlan) error {
	query := `
		INSERT INTO schedule_plans (
			name,
			description,
			submission_start_time,
			submission_end_time,
			active_start_time,
			active_end_time,
			schedule_template_name
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, version
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	args := []any{sp.Name, sp.Description, sp.SubmissionStartTime, sp.SubmissionEndTime, sp.ActiveStartTime, sp.ActiveEndTime, sp.ScheduleTemplateName}
	if err := m.db.QueryRowContext(ctx, query, args...).Scan(&sp.ID, &sp.CreatedAt, &sp.Version); err != nil {
		return err
	}

	return nil
}

func (m *Models) SelectSchedulePlanByID(id uuid.UUID) (*SchedulePlan, error) {
	query := `
		SELECT 
			name, 
			description, 
			submission_start_time, 
			submission_end_time,
			active_start_time,
			active_end_time,
			schedule_template_name,
			created_at,
			version
		FROM schedule_plans
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sp := &SchedulePlan{
		ID: id,
	}
	dest := []any{
		&sp.Name,
		&sp.Description,
		&sp.SubmissionStartTime,
		&sp.SubmissionEndTime,
		&sp.ActiveStartTime,
		&sp.ActiveEndTime,
		&sp.ScheduleTemplateName,
		&sp.CreatedAt,
		&sp.Version,
	}

	if err := m.db.QueryRowContext(ctx, query, id).Scan(dest...); err != nil {
		return nil, err
	}

	return sp, nil
}
