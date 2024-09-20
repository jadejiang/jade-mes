package util

import (
	"fmt"
	"strings"
	"time"
)

// IsInPeriods 当前是否在时间区间，并可选地检查是否符合周几，
// 其中多区间，可以用”,“分割。weekday，有多个可以直接拼接，0为周末，
func IsInPeriods(startTime, endTime string, weekday string, from *time.Time, to *time.Time) bool {
	now := time.Now().Local()

	// weekday格式：0123456，其中0为周末
	if weekday != "" && !isWeekDayMatched(weekday, &now) {
		return false
	}

	if (from != nil && now.Before(*from)) || (to != nil && now.After(*to)) {
		return false
	}

	if startTime != "" && endTime != "" {
		startTimes := strings.Split(startTime, ",")
		endTimes := strings.Split(endTime, ",")

		currentTime := now.Format("15:04")

		for i, start := range startTimes {
			// in case of "index out of range"
			if len(endTimes) < i {
				break
			}
			end := endTimes[i]

			// HH:mm
			if len(start) != 5 || len(end) != 5 {
				continue
			}

			if currentTime >= start && currentTime <= end {
				return true
			}
		}

		return false
	}

	return true
}

// isWeekDayMatched weekday格式：0123456，其中0为周末
func isWeekDayMatched(weekday string, t *time.Time) bool {
	if t == nil {
		now := time.Now()
		t = &now
	}

	if weekday == "" {
		return true
	}

	weekDay := fmt.Sprintf("%d", t.Weekday())
	return strings.Contains(weekday, weekDay)
}
