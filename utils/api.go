package utils

import (
	"time"
)

func GetTimePointer(t time.Time) *time.Time {
	return getTimePointer(t)
}

func GetDurationPointer(d time.Duration) *time.Duration {
	return getDurationPointer(d)
}

func GetFirstSaturday() *time.Time {
	return getFirstSaturday()
}

func GetDurationForRepetitionPeriod(repetition string) int {
	return getDurationForRepetitionPeriod(repetition)
}

func ParseTime(format, deadline string) *time.Time {
	return parseTime(format, deadline)
}

func DurationToHMFormat(d *time.Duration) string {
	return durationToHMFormat(d)
}

func MinutesToHMFormat(minutes float64) string {
	return minutesToHMFormat(minutes)
}

func GetRunningBinaryPath() string {
	return getRunningBinaryPath()
}
