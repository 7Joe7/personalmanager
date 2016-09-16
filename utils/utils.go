package utils

import (
	"time"
	"fmt"
)

func GetTimePointer(t time.Time) *time.Time {
	return &t
}

func GetDurationPointer(d time.Duration) *time.Duration {
	return &d
}

func GetFirstSaturday() *time.Time {
	now := time.Now().Truncate(24 * time.Hour)
	return GetTimePointer(now.Add(time.Duration(24 * (6 - int(now.Weekday()))) * time.Hour))
}

func ParseTime(format, deadline string) *time.Time {
	d, err := time.Parse(format, deadline)
	if err != nil {
		panic(err)
	}
	return &d
}

func DurationToHMFormat(d *time.Duration) string {
	return fmt.Sprintf("%dh%dm", int(d.Hours()), int(d.Minutes()) % 60)
}

func MinutesToHMFormat(minutes float64) string {
	t := int(minutes)
	ms := t % 60
	return fmt.Sprintf("%dh%dm", (t - ms) / 60, ms)
}