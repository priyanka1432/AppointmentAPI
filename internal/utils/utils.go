package utils

import (
	"fmt"
	"strings"
	"time"
)

var DayNameToInt = map[string]int{
	"sunday":    0,
	"monday":    1,
	"tuesday":   2,
	"wednesday": 3,
	"thursday":  4,
	"friday":    5,
	"saturday":  6,
}

func DayStringToInt(day string) (int, error) {
	lowerday := strings.ToLower(strings.TrimSpace(day))
	if value, ok := DayNameToInt[lowerday]; ok {
		return value, nil
	}
	return -1, fmt.Errorf("invalid day: %s", day)
}
func AvailabilityAppliesToDate(storedDay string, date time.Time) (bool, error) {
	s := strings.TrimSpace(storedDay)
	if s == "" {
		return false, fmt.Errorf("empty day")
	}

	dayInt, err := DayStringToInt(s)
	if err != nil {
		return false, err
	}

	return int(date.Weekday()) == dayInt, nil
}

func ParseTimeOfDay(s string) (int, int, error) {
	time, err := time.Parse("15:04", s)
	if err != nil {

		return 0, 0, err
	}
	return time.Hour(), time.Minute(), nil
}

func Combine(date time.Time, hour, minute int) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), hour, minute, 0, 0, time.UTC)
}

func Is30(time time.Time) bool {
	min := time.Minute()
	return min == 0 || min == 30
}
