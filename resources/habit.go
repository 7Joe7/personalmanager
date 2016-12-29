package resources

import (
	"fmt"
	"time"
	"encoding/json"
)

type Habit struct {
	Name          string     `json:",omitempty"`
	Active        bool       `json:",omitempty"`
	Done          bool       `json:",omitempty"`
	Negative      bool       `json:",omitempty"`
	Description   string     `json:",omitempty"`
	Deadline      *time.Time `json:",omitempty"`
	Tries         int        `json:",omitempty"`
	Successes     int        `json:",omitempty"`
	ActualStreak  int        `json:",omitempty"`
	LastStreak    int        `json:",omitempty"`
	LastStreakEnd *time.Time `json:",omitempty"`
	Repetition    string     `json:",omitempty"`
	BasePoints    int        `json:",omitempty"`
	Id            string     `json:",omitempty"`
	Goal          *Goal      `json:",omitempty"`
	Count         int        `json:",omitempty"`
	Limit         int        `json:",omitempty"`
	Average       float64    `json:",omitempty"`
}

func (h *Habit) SetId(id string) {
	h.Id = id
}

func (h *Habit) GetId() string {
	return h.Id
}

func (h *Habit) Load(tr Transaction) error {
	if h.Goal != nil {
		err := tr.RetrieveEntity(DB_DEFAULT_GOALS_BUCKET_NAME, []byte(h.Goal.Id), h.Goal, true)
		if err != nil {
			return err
		}
	}
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
			if h.Negative {
				if h.Count > h.Limit {
					ico = ICO_RED
					colour = "red"
				} else {
					ico = ICO_BLACK
					colour = "black"
				}
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

func (h *Habit) MarshalJSON() ([]byte, error) {
	type mHabit Habit
	if h.Goal != nil {
		h.Goal = &Goal{Id: h.Goal.Id}
	}
	return json.Marshal(mHabit(*h))
}

func (h *Habit) getItem(id string) *AlfredItem {
	var subtitle string
	switch {
	case h.Active:
		if h.Negative {
			subtitle = fmt.Sprintf(SUB_FORMAT_ACTIVE_BAD_HABIT, h.Limit - h.Count, h.Limit, h.Average, h.Successes,
				h.Tries, h.ActualStreak, h.Deadline.Format(DATE_FORMAT), h.BasePoints)
		} else {
			subtitle = fmt.Sprintf(SUB_FORMAT_ACTIVE_HABIT, h.Successes, h.Tries, h.ActualStreak,
				h.Deadline.Format(DATE_FORMAT), h.BasePoints)
		}
	default:
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
