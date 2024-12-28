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

func ValidateScheduleTemplate(st *models.ScheduleTemplate) error {
	for i := 0; i < len(st.Shifts); i++ {
		if st.Shifts[i].StartTime.After(st.Shifts[i].EndTime) {
			return fmt.Errorf("班次 %d 的开始时间晚于结束时间", i)
		}
	}

	for i := 0; i < len(st.Shifts); i++ {
		for j := i + 1; j < len(st.Shifts); j++ {
			if !(st.Shifts[i].EndTime.Before(st.Shifts[j].StartTime) ||
				st.Shifts[i].EndTime.Equal(st.Shifts[j].StartTime) ||
				st.Shifts[i].StartTime.After(st.Shifts[j].EndTime) ||
				st.Shifts[i].StartTime.Equal(st.Shifts[j].EndTime)) {
				return fmt.Errorf("班次 %d 与班次 %d 有时间冲突", i, j)
			}
		}
	}

	return nil
}
