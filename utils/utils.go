package utils

import "time"

func GetTimePointer(t time.Time) *time.Time {
	return &t
}

func GetIntPointer(i int) *int {
	return &i
}