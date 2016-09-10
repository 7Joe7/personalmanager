package operations

import (
	"log"
	"time"

	"github.com/7joe7/personalmanager/resources"
	"github.com/7joe7/personalmanager/utils"
	"github.com/7joe7/personalmanager/db"
)

func getModifyHabitFunc(h *resources.Habit, name, repetition, deadline string, toggleActive, toggleDone bool, basePoints int, scoreChange *int) func () {
	return func () {
		if name != "" {
			h.Name = name
		}
		if toggleActive {
			if h.Active {
				h.Active = false
				h.Deadline = nil
				h.Done = false
				h.ActualStreak = 0
				h.LastStreakEnd = nil
				h.LastStreak = 0
				h.Repetition = ""
				h.BasePoints = 0
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
					if h.ActualStreak < 0 {
						h.ActualStreak = 0
					}
					if h.LastStreakEnd != nil && h.Deadline.Equal(*h.LastStreakEnd) {
						h.LastStreakEnd = nil
						h.ActualStreak = h.LastStreak
					}
					h.Successes += 1
					h.ActualStreak += 1
				}
				scoreChange = utils.GetIntPointer(h.ActualStreak * h.BasePoints)
			}
			if deadline != "" {
				t, err := time.Parse("2.1.2006 15:04", deadline)
				if err != nil {
					log.Fatalf("Unable to parse deadline. %v", err)
				}
				h.Deadline = &t
			}
		}
	}
}

func getSyncHabitFunc(h *resources.Habit, scoreChange *int) func () {
	return func () {
		if !h.Active {
			return
		}

		if h.Deadline.Before(time.Now()) {
			if !h.Done {
				failHabit(h)
				scoreChange = utils.GetIntPointer(h.ActualStreak * h.BasePoints)
			}
			switch h.Repetition {
			case resources.HBT_REPETITION_DAILY:
				h.Deadline = utils.GetTimePointer(h.Deadline.Add(24 * time.Hour))
			case resources.HBT_REPETITION_WEEKLY:
				h.Deadline = utils.GetTimePointer(h.Deadline.Add(7 * 24 * time.Hour))
			case resources.HBT_REPETITION_MONTHLY:
				h.Deadline = utils.GetTimePointer(h.Deadline.AddDate(0, 1, 0))
			}
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
	h.Successes -= 1
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

func addHabit(name, repetition string, activeFlag bool) string {
	h := resources.NewHabit(name)
	if activeFlag {
		activateHabit(h, repetition)
	}
	return db.AddEntity(h, resources.DB_DEFAULT_HABITS_BUCKET_NAME)
}

func deleteHabit(habitId string) {
	db.DeleteEntity([]byte(habitId), resources.DB_DEFAULT_HABITS_BUCKET_NAME)
}

func modifyHabit(habitId, name, repetition, deadline string, toggleActive, toggleDone bool, basePoints int) {
	var scoreChange *int
	habit := &resources.Habit{}
	modifyHabit := GetModifyHabitFunc(habit, name, repetition, deadline, toggleActive, toggleDone, basePoints, scoreChange)
	status := &resources.Status{}
	t := db.NewTransaction()
	t.Add(func () error {
		if err := t.ModifyEntity(resources.DB_DEFAULT_HABITS_BUCKET_NAME, []byte(habitId), habit, modifyHabit); err != nil {
			return err
		}
		if err := t.ModifyEntity(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_ACTUAL_STATUS_KEY, status, GetAddScoreFunc(status, *scoreChange)); err != nil {
			return err
		}
		return nil
	})
	t.Execute()
}

func getHabit(habitId string) *resources.Habit {
	habit := &resources.Habit{}
	db.RetrieveEntity(resources.DB_DEFAULT_HABITS_BUCKET_NAME, []byte(habitId), habit)
	return habit
}

func getHabits() map[string]*resources.Habit {
	habits := map[string]*resources.Habit{}
	db.RetrieveEntities(resources.DB_DEFAULT_HABITS_BUCKET_NAME, func (id string) interface{} {
		habits[id] = &resources.Habit{}
		return habits[id]
	})
	return habits
}

func filterHabits(filter func(*resources.Habit) bool) map[string]*resources.Habit {
	habits := map[string]*resources.Habit{}
	h := &resources.Habit{}
	db.FilterEntities(resources.DB_DEFAULT_HABITS_BUCKET_NAME, h, func (id string) {
		if filter(h) {
			hCopy := &resources.Habit{}
			*hCopy = *h
			habits[id] = hCopy
		}
	})
	return habits
}

func getActiveHabits() map[string]*resources.Habit {
	return FilterHabits(func (h *resources.Habit) bool {
		return h.Active
	})
}

func getNonActiveHabits() map[string]*resources.Habit {
	return FilterHabits(func (h *resources.Habit) bool {
		return !h.Active
	})
}
