package utils

import "time"

func GetTimePointer(t time.Time) *time.Time {
	return &t
}

func GetFirstSaturday() *time.Time {
	now := time.Now().Truncate(24 * time.Hour)
	return GetTimePointer(now.Add(time.Duration(24 * (6 - int(now.Weekday()))) * time.Hour))
}