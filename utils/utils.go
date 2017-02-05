package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

func getTimePointer(t time.Time) *time.Time {
	return &t
}

func getDurationPointer(d time.Duration) *time.Duration {
	return &d
}

func getFirstSaturday() *time.Time {
	now := time.Now().Truncate(24 * time.Hour)
	return GetTimePointer(now.Add(time.Duration(24*(6-int(now.Weekday()))) * time.Hour))
}

func getDurationForRepetitionPeriod(repetition string) int {
	switch repetition {
	case "Daily":
		return int(time.Duration(int64(86400000000000)).Hours())
	case "Weekly":
		return int(time.Duration(int64(604800000000000)).Hours())
	case "Monthly": // approximation 1814400000000000
		return int(time.Duration(int64(2592000000000000)).Hours())
	}
	return 0
}

func parseTime(format, deadline string) *time.Time {
	d, err := time.Parse(format, deadline)
	if err != nil {
		panic(err)
	}
	return &d
}

func durationToHMFormat(d *time.Duration) string {
	return fmt.Sprintf("%dh%dm", int(d.Hours()), int(d.Minutes())%60)
}

func minutesToHMFormat(minutes float64) string {
	t := int(minutes)
	ms := t % 60
	return fmt.Sprintf("%dh%dm", (t-ms)/60, ms)
}

func getRunningBinaryPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}
