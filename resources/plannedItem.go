package resources

import "time"

type PlannedItem interface {
	SetId(id string)
	GetId() string
	GetAlfredItem(id string) *AlfredItem
	GetTimeEstimate() *time.Duration
}
