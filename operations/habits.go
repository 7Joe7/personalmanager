package operations

import (
	"time"

	"github.com/7joe7/personalmanager/resources"
	"github.com/7joe7/personalmanager/utils"
	"github.com/7joe7/personalmanager/db"
)

func getModifyHabitFunc(h *resources.Habit, name, repetition, deadline string, toggleActive, toggleDone, toggleDonePrevious bool, basePoints int, status *resources.Status) func () {
	return func () {
		if name != "" {
			h.Name = name
		}
		if toggleActive {
			if h.Active {
				deactivateHabit(h)
			} else {
				activateHabit(h, repetition)
			}
		}
		if h.Active {
			if basePoints != -1 {
				h.BasePoints = basePoints
			}
			if toggleDone {
				if h.Done {
					h.Done = false
					failHabit(h)
				} else {
					h.Done = true
					succeedHabit(h, h.LastStreakEnd)
				}
				change := h.ActualStreak * h.BasePoints
				status.Score += change
				status.Today += change
			}
			if deadline != "" {
				h.Deadline = utils.ParseTime(resources.DATE_FORMAT, deadline)
			}
			if toggleDonePrevious {
				succeedHabit(h, addPeriod(h.Repetition, h.LastStreakEnd))
				if h.Done {
					status.Score += h.ActualStreak * h.BasePoints + (h.ActualStreak - 1) * h.BasePoints
					status.Today += (h.ActualStreak - 1) * h.BasePoints

				} else {
					status.Score += (h.ActualStreak + 1) * h.BasePoints
				}
			}
		}
	}
}

func succeedHabit(h *resources.Habit, lastStreakEnd *time.Time) {
	if h.ActualStreak < 0 {
		h.ActualStreak = 0
	}
	if lastStreakEnd != nil && h.Deadline.Equal(*lastStreakEnd) {
		h.ActualStreak += h.LastStreak
	}
	h.Successes += 1
	h.ActualStreak += 1
}

func addPeriod(repetition string, deadline *time.Time) *time.Time {
	if deadline == nil {
		return nil
	}
	switch repetition {
	case resources.HBT_REPETITION_DAILY:
		return utils.GetTimePointer(deadline.Add(24 * time.Hour))
	case resources.HBT_REPETITION_WEEKLY:
		return utils.GetTimePointer(deadline.Add(7 * 24 * time.Hour))
	case resources.HBT_REPETITION_MONTHLY:
		return utils.GetTimePointer(deadline.AddDate(0, 1, 0))
	}
	return nil
}

func getSyncHabitFunc(h *resources.Habit, changeStatus *resources.Status) func () {
	return func () {
		if !h.Active {
			return
		}

		if h.Deadline.Before(time.Now()) {
			if !h.Done {
				failHabit(h)
				changeStatus.Score = h.ActualStreak * h.BasePoints
			}
			h.Deadline = addPeriod(h.Repetition, h.Deadline)
			h.Done = false
			h.Tries += 1
		}
	}
}

func failHabit(h *resources.Habit) {
	h.LastStreakEnd = utils.GetTimePointer(*h.Deadline)
	h.LastStreak = h.ActualStreak
	if h.ActualStreak > 0 {
		h.ActualStreak = 0
	}
	h.ActualStreak -= 1
}

func activateHabit(h *resources.Habit, repetition string) {
	h.Active = true
	if repetition == "" {
		repetition = resources.HBT_REPETITION_DAILY
	}
	h.Repetition = repetition
	h.Tries += 1
	h.Deadline = utils.GetTimePointer(time.Now().Add(24 * time.Hour).Truncate(24 * time.Hour))
}

func deactivateHabit(h *resources.Habit) {
	h.Active = false
	h.Deadline = nil
	h.Done = false
	h.ActualStreak = 0
	h.LastStreakEnd = nil
	h.LastStreak = 0
	h.Repetition = ""
	h.BasePoints = 0
}

func addHabit(name, repetition, deadline string, activeFlag bool, basePoints int) {
	h := resources.NewHabit(name)
	if activeFlag {
		activateHabit(h, repetition)
		h.Deadline = utils.ParseTime(resources.DATE_FORMAT, deadline)
		if basePoints != -1 {
			h.BasePoints = basePoints
		}
	}
	tr := db.NewTransaction()
	tr.Add(func () error {
		return tr.AddEntity(resources.DB_DEFAULT_HABITS_BUCKET_NAME, h)
	})
	tr.Execute()
}

func deleteHabit(habitId string) {
	db.DeleteEntity([]byte(habitId), resources.DB_DEFAULT_HABITS_BUCKET_NAME)
}

func modifyHabit(habitId, name, repetition, deadline string, toggleActive, toggleDone, toggleDonePrevious bool, basePoints int) {
	habit := &resources.Habit{}
	habitStatus := &resources.Status{}
	modifyHabit := getModifyHabitFunc(habit, name, repetition, deadline, toggleActive, toggleDone, toggleDonePrevious, basePoints, habitStatus)
	status := &resources.Status{}
	t := db.NewTransaction()
	t.Add(func () error {
		if err := t.ModifyEntity(resources.DB_DEFAULT_HABITS_BUCKET_NAME, []byte(habitId), habit, modifyHabit); err != nil {
			return err
		}
		if err := t.ModifyEntity(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_ACTUAL_STATUS_KEY, status, getAddScoreFunc(status, habitStatus)); err != nil {
			return err
		}
		return nil
	})
	t.Execute()
}

func getHabit(habitId string) *resources.Habit {
	habit := &resources.Habit{}
	tr := db.NewTransaction()
	tr.Add(func () error {
		return tr.RetrieveEntity(resources.DB_DEFAULT_HABITS_BUCKET_NAME, []byte(habitId), habit)
	})
	tr.Execute()
	return habit
}

func getHabits() map[string]*resources.Habit {
	habits := map[string]*resources.Habit{}
	db.RetrieveEntities(resources.DB_DEFAULT_HABITS_BUCKET_NAME, func (id string) resources.Entity {
		habits[id] = &resources.Habit{}
		return habits[id]
	})
	return habits
}

func filterHabits(filter func(*resources.Habit) bool) map[string]*resources.Habit {
	habits := map[string]*resources.Habit{}
	h := &resources.Habit{}
	copyFunc := func () {
		c := &resources.Habit{}
		*c = *h
		habits[h.Id] = c
	}
	db.FilterEntities(resources.DB_DEFAULT_HABITS_BUCKET_NAME, h, func () bool { return filter(h) }, copyFunc)
	return habits
}

func getActiveHabits() map[string]*resources.Habit {
	return FilterHabits(func (h *resources.Habit) bool { return h.Active })
}

func getNonActiveHabits() map[string]*resources.Habit {
	return FilterHabits(func (h *resources.Habit) bool { return !h.Active })
}
