package utils

import (
	"fmt"
	"net/mail"

	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/models"
)

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsValidRole(role string) bool {
	return role == "普通助理" || role == "资深助理" || role == "黑心"
}

func ValidateShifts(shifts []*models.ScheduleTemplateShift) error {
	for i, shift := range shifts {
		if shift.DayOfWeek < 1 || shift.DayOfWeek > 7 {
			return fmt.Errorf("第 %d 个班次的星期数必须在 1 到 7 之间，当前值为 %d", i+1, shift.DayOfWeek)
		}

		if shift.StartTime.IsZero() || shift.EndTime.IsZero() {
			return fmt.Errorf("第 %d 个班次的开始时间和结束时间不能为空", i+1)
		}

		if !shift.EndTime.After(shift.StartTime) {
			return fmt.Errorf("第 %d 个班次的结束时间必须晚于开始时间", i+1)
		}
	}

	// Check for overlapping shifts on the same day
	for i := 0; i < len(shifts); i++ {
		for j := i + 1; j < len(shifts); j++ {
			// Only compare shifts on the same day
			if shifts[i].DayOfWeek == shifts[j].DayOfWeek {
				// Check if shifts overlap
				if !(shifts[i].EndTime.Before(shifts[j].StartTime) || shifts[i].EndTime.Equal(shifts[j].StartTime) || shifts[j].EndTime.Before(shifts[i].StartTime) || shifts[j].EndTime.Equal(shifts[i].StartTime)) {
					return fmt.Errorf("第 %d 个班次和第 %d 个班次在同一天且时间重叠", i+1, j+1)
				}
			}
		}
	}

	return nil
}
