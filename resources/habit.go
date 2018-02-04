package resources

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/7joe7/personalmanager/utils"
)

type Habit struct {
	Name          string         `json:",omitempty"`
	Active        bool           `json:",omitempty"`
	Done          bool           `json:",omitempty"`
	Negative      bool           `json:",omitempty"`
	Learned       bool           `json:",omitempty"`
	Description   string         `json:",omitempty"`
	Deadline      *time.Time     `json:",omitempty"`
	Tries         int            `json:",omitempty"`
	Successes     int            `json:",omitempty"`
	ActualStreak  int            `json:",omitempty"`
	LastStreak    int            `json:",omitempty"`
	LastStreakEnd *time.Time     `json:",omitempty"`
	Repetition    string         `json:",omitempty"`
	BasePoints    int            `json:",omitempty"`
	Id            string         `json:",omitempty"`
	Goal          *Goal          `json:",omitempty"`
	Count         int            `json:",omitempty"`
	Limit         int            `json:",omitempty"`
	Average       float64        `json:",omitempty"`
	AlarmTime     *time.Time     `json:",omitempty"`
	TimeEstimate  *time.Duration `json:",omitempty"`
}

func (h *Habit) SetId(id string) {
	h.Id = id
}

func (h *Habit) GetId() string {
	return h.Id
}

func (h *Habit) GetTimeEstimate() *time.Duration {
	return h.TimeEstimate
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

func (h *Habit) Less(entity Entity) bool {
	switch entity.(type) {
	case *Habit:
		otherHabit := entity.(*Habit)
		if h.Done != otherHabit.Done {
			return otherHabit.Done
		}
		if h.Active != otherHabit.Active {
			return h.Active
		}
		if h.Repetition != otherHabit.Repetition {
			return h.Repetition == HBT_REPETITION_DAILY ||
				(otherHabit.Repetition != HBT_REPETITION_DAILY &&
					h.Repetition == HBT_REPETITION_WEEKLY)
		}
		if (h.AlarmTime != nil) != (otherHabit.AlarmTime != nil) {
			return h.AlarmTime != nil
		}
		if h.AlarmTime != nil && (h.AlarmTime.Day() != otherHabit.AlarmTime.Day() || h.AlarmTime.Hour() != otherHabit.AlarmTime.Hour() || h.AlarmTime.Minute() != otherHabit.AlarmTime.Minute()) {
			return h.AlarmTime.Before(*otherHabit.AlarmTime)
		}
		if h.BasePoints != otherHabit.BasePoints {
			return h.BasePoints > otherHabit.BasePoints
		}
		if (h.TimeEstimate == nil) != (otherHabit.TimeEstimate == nil) {
			return h.TimeEstimate != nil
		}
		if h.TimeEstimate != nil && h.TimeEstimate.Minutes() != otherHabit.TimeEstimate.Minutes() {
			return h.TimeEstimate.Minutes() < otherHabit.TimeEstimate.Minutes()
		}
		return h.Name < otherHabit.Name
	case *Task:
		otherTask := entity.(*Task)
		if otherTask.InProgress {
			return false
		}
		if h.Repetition == HBT_REPETITION_DAILY {
			return true
		}
		if h.BasePoints != otherTask.BasePoints {
			return h.BasePoints > otherTask.BasePoints
		}
		if (h.TimeEstimate == nil) != (otherTask.TimeEstimate == nil) {
			return h.TimeEstimate != nil
		}
		if h.TimeEstimate != nil && h.TimeEstimate.Minutes() != otherTask.TimeEstimate.Minutes() {
			return h.TimeEstimate.Minutes() < otherTask.TimeEstimate.Minutes()
		}
		return true
	}
	return false
}

func (h *Habit) GetIconColourAndOrder() (string, string) {
	if h.Active {
		if h.Done || h.Learned {
			return ICO_GREEN, "green"
		} else {
			var ico, colour string
			switch h.Repetition {
			case HBT_REPETITION_DAILY:
				ico, colour = ICO_RED, "red"
			case HBT_REPETITION_WEEKLY:
				ico, colour = ICO_ORANGE, "orange"
			case HBT_REPETITION_MONTHLY:
				ico, colour = ICO_YELLOW, "yellow"
			case HBT_REPETITION_YEARLY:
				ico, colour = ICO_PURPLE, "purple"
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
			}
			return ico, colour
		}
	} else {
		return ICO_BLACK, "black"
	}
}

func (h *Habit) MarshalJSON() ([]byte, error) {
	type mHabit Habit
	if h.Goal != nil {
		h.Goal = &Goal{Id: h.Goal.Id}
	}
	return json.Marshal(mHabit(*h))
}

func (h *Habit) GetAlfredItem(id string) *AlfredItem {
	var subtitle string
	switch {
	case h.Active:
		if h.Negative {
			if h.Repetition == HBT_REPETITION_DAILY {
				subtitle = fmt.Sprintf(SUB_FORMAT_ACTIVE_BAD_HABIT, h.Successes, h.Tries, h.ActualStreak, h.Count,
					h.Limit, h.Average, h.BasePoints)
			} else {
				subtitle = fmt.Sprintf(SUB_FORMAT_ACTIVE_BAD_HABIT_NOT_DAILY, h.Successes, h.Tries, h.ActualStreak,
					h.Count, h.Limit, h.Average, h.Deadline.Format(DATE_FORMAT), h.BasePoints)
			}
		} else if h.Repetition == HBT_REPETITION_DAILY {
			if h.AlarmTime == nil {
				subtitle = fmt.Sprintf(SUB_FORMAT_ACTIVE_DAILY_HABIT, h.Successes, h.Tries, h.ActualStreak, h.BasePoints)
			} else {
				subtitle = fmt.Sprintf(SUB_FORMAT_ACTIVE_DAILY_HABIT_WITH_ALARM, h.Successes, h.Tries, h.ActualStreak,
					h.AlarmTime.Format(HOUR_MINUTE_FORMAT), h.BasePoints)
			}
		} else {
			if h.AlarmTime == nil {
				subtitle = fmt.Sprintf(SUB_FORMAT_ACTIVE_NOT_DAILY, h.Successes, h.Tries, h.ActualStreak,
					h.Deadline.Format(DATE_FORMAT), h.BasePoints)
			} else {
				subtitle = fmt.Sprintf(SUB_FORMAT_ACTIVE_NOT_DAILY_WITH_ALARM, h.Successes, h.Tries, h.ActualStreak,
					h.AlarmTime.Format(DATE_HOUR_MINUTE_FORMAT), h.Deadline.Format(DATE_FORMAT), h.BasePoints)
			}
		}
		if h.Learned {
			subtitle = "Learned, " + subtitle
		}
		if h.TimeEstimate != nil {
			subtitle += fmt.Sprintf(", %s", utils.DurationToHMFormat(h.TimeEstimate))
		}
	default:
		subtitle = fmt.Sprintf(SUB_FORMAT_NON_ACTIVE_HABIT, h.Successes, h.Tries)
	}
	iconPath, _ := h.GetIconColourAndOrder()
	icon := NewAlfredIcon(iconPath)
	return &AlfredItem{
		Name:     h.Name,
		Arg:      id,
		Subtitle: subtitle,
		Icon:     icon,
		Valid:    true,
		entity:   h}
}
