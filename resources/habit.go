package resources

import (
	"time"
	"fmt"
)

type Habit struct {
	Name          string
	Active        bool
	Done          bool
	Description   string
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

func (h *Habit) GetIconColourAndOrder() (string, string, int) {
	if h.Active {
		if h.Done || h.ActualStreak > 49 {
			return ICO_GREEN, "green", HBT_DONE_BASE_ORDER
		} else {
			var order int
			var ico, colour string
			switch h.Repetition {
			case HBT_REPETITION_DAILY:
				ico, colour, order = ICO_RED, "red", HBT_BASE_ORDER_DAILY
			case HBT_REPETITION_WEEKLY:
				ico, colour, order = ICO_ORANGE, "orange", HBT_BASE_ORDER_WEEKLY
			case HBT_REPETITION_MONTHLY:
				ico, colour, order = ICO_YELLOW, "yellow", HBT_BASE_ORDER_MONTHLY
			}
			if h.ActualStreak > 21 {
				ico = ICO_BLUE
				colour = "blue"
				order += 1000
			}
			return ico, colour, order
		}
	} else {
		return ICO_BLACK, "black", HBT_BASE_ORDER_DAILY
	}
	return "", "", 0
}

func (h *Habit) getItem(id string) *AlfredItem {
	var subtitle string
	if h.Active {
		subtitle = fmt.Sprintf(SUB_FORMAT_ACTIVE_HABIT, h.Successes, h.Tries, h.ActualStreak,
			h.Deadline.Format(DATE_FORMAT), h.BasePoints)
	} else {
		subtitle = fmt.Sprintf(SUB_FORMAT_NON_ACTIVE_HABIT, h.Successes, h.Tries)
	}
	iconPath, _, order := h.GetIconColourAndOrder()
	icon := NewAlfredIcon(iconPath)
	order -= h.BasePoints
	return &AlfredItem{
		Name:     h.Name,
		Arg:      id,
		Subtitle: subtitle,
		Icon:     icon,
		Valid:    true,
		order:    order}
}
