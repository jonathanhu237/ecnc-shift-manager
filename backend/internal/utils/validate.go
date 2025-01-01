package utils

import (
	"fmt"
	"net/mail"
	"time"

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
		startTime, err := time.Parse("15:04", st.Shifts[i].StartTime)
		if err != nil {
			return fmt.Errorf("班次 %d 的开始时间格式无效", i)
		}
		endTime, err := time.Parse("15:04", st.Shifts[i].EndTime)
		if err != nil {
			return fmt.Errorf("班次 %d 的结束时间格式无效", i)
		}
		if startTime.After(endTime) {
			return fmt.Errorf("班次 %d 的开始时间晚于结束时间", i)
		}
	}

	for i := 0; i < len(st.Shifts); i++ {
		for j := i + 1; j < len(st.Shifts); j++ {
			iEndTime, _ := time.Parse("15:04", st.Shifts[i].EndTime)
			jStartTime, _ := time.Parse("15:04", st.Shifts[j].StartTime)
			iStartTime, _ := time.Parse("15:04", st.Shifts[i].StartTime)
			jEndTime, _ := time.Parse("15:04", st.Shifts[j].EndTime)

			if !(iEndTime.Before(jStartTime) ||
				iEndTime.Equal(jStartTime) ||
				iStartTime.After(jEndTime) ||
				iStartTime.Equal(jEndTime)) {
				return fmt.Errorf("班次 %d 与班次 %d 有时间冲突", i, j)
			}
		}
	}

	return nil
}
