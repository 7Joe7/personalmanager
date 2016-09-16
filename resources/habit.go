package resources

import (
	"time"
	"fmt"
)

type Habit struct {
	Name          string
	Active        bool
	Done          bool
	Deadline      *time.Time
	Tries         int
	Successes     int
	ActualStreak  int
	LastStreak    int
	LastStreakEnd *time.Time
	Repetition    string
	BasePoints    int
	Id            string
}

func (h *Habit) SetId(id string) {
	h.Id = id
}

func (h *Habit) GetId() string {
	return h.Id
}

func (h *Habit) Load(tr Transaction) error {
	return nil
}

func (h *Habit) getItem(id string) *AlfredItem {
	var subtitle string
	var icon *AlfredIcon
	var order int
	if h.Active {
		if h.Done {
			icon = NewAlfredIcon(ICO_GREEN)
			order = HBT_DONE_BASE_ORDER
		} else {
			switch h.Repetition {
			case HBT_REPETITION_DAILY:
				icon = NewAlfredIcon(ICO_RED)
				order = HBT_BASE_ORDER_DAILY
			case HBT_REPETITION_WEEKLY:
				icon = NewAlfredIcon(ICO_ORANGE)
				order = HBT_BASE_ORDER_WEEKLY
			case HBT_REPETITION_MONTHLY:
				icon = NewAlfredIcon(ICO_YELLOW)
				order = HBT_BASE_ORDER_MONTHLY
			}
		}
		subtitle = fmt.Sprintf(SUB_FORMAT_ACTIVE_HABIT, h.Successes, h.Tries, h.ActualStreak,
			h.Deadline.Format(DATE_FORMAT), h.BasePoints)
	} else {
		icon = NewAlfredIcon(ICO_BLACK)
		order = HBT_BASE_ORDER_DAILY
		subtitle = fmt.Sprintf(SUB_FORMAT_NON_ACTIVE_HABIT, h.Successes, h.Tries)
	}
	order -= h.BasePoints
	return &AlfredItem{
		Name:     h.Name,
		Arg:      id,
		Subtitle: subtitle,
		Icon:     icon,
		Valid:    true,
		order:    order}
}
