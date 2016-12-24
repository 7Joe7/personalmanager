package operations

import (
	"time"

	"github.com/7joe7/personalmanager/resources"
	"github.com/7joe7/personalmanager/utils"
	"github.com/7joe7/personalmanager/db"
	"github.com/7joe7/personalmanager/anybar"
)

func getModifyHabitFunc(h *resources.Habit, name, repetition, description, deadline, goalId string, toggleActive, toggleDone, toggleDonePrevious, toggleUndonePrevious, negativeFlag bool, basePoints, repetitionGoal int, status *resources.Status, tr resources.Transaction) func () {
	return func () {
		if name != "" {
			h.Name = name
		}
		if description != "" {
			h.Description = description
		}
		if goalId != "" {
			h.Goal = &resources.Goal{Id: goalId}
		}
		if toggleActive {
			if h.Active {
				if h.Negative {
					tr.Add(func () error {
						return tr.DeleteEntity(resources.DB_DEFAULT_HABITS_BUCKET_NAME, []byte(h.Id))
					})
				} else {
					deactivateHabit(h)
				}
				//anybar.RemoveAndQuit(resources.DB_DEFAULT_HABITS_BUCKET_NAME, h.Id, tr)
			} else {
				activateHabit(h, repetition)
				//_, colour, _ := h.GetIconColourAndOrder()
				//anybar.AddToActivePorts(h.Name, colour, resources.DB_DEFAULT_HABITS_BUCKET_NAME, h.Id, tr)
			}
		}
		if negativeFlag {
			h.Negative = !h.Negative
		}
		if h.Active {
			if basePoints != -1 {
				h.BasePoints = basePoints
			}
			if repetitionGoal != -1 {
				h.Limit = repetitionGoal
			}
			if toggleDone {
				if h.Negative {
					h.Count++
					h.Average = (h.Average * float32(h.Tries) + 1) / float32(h.Tries)
					if h.Count - 1 <= h.Limit && h.Count > h.Limit {
						failHabit(h)
						countPointChange(h, status, 1)
					} else if h.Count > h.Limit {
						countPointChange(h, status, h.Count - h.Limit + 1)
					}
				} else {
					if h.Done {
						h.Done = false
						failHabit(h)
						h.Successes -= 1
						//_, colour, _ := h.GetIconColourAndOrder()
						//anybar.AddToActivePorts(h.Name, colour, resources.DB_DEFAULT_HABITS_BUCKET_NAME, h.Id, tr)
					} else {
						h.Done = true
						succeedHabit(h, h.Deadline)
						//anybar.RemoveAndQuit(resources.DB_DEFAULT_HABITS_BUCKET_NAME, h.Id, tr)
					}
					countPointChange(h, status, 1)
				}
			}
			if deadline != "" {
				h.Deadline = utils.ParseTime(resources.DATE_FORMAT, deadline)
			}
			if toggleDonePrevious {
				previousActualStreak := h.ActualStreak
				succeedHabit(h, removePeriod(h.Repetition, h.Deadline))
				if h.Done {
					if previousActualStreak == 1 {
						status.Score += h.BasePoints
						status.Today += h.ActualStreak * h.ActualStreak * h.BasePoints - h.BasePoints
					}
					status.Score += h.ActualStreak * h.ActualStreak * h.BasePoints + (h.ActualStreak - 1) * (h.ActualStreak - 1) * h.BasePoints
				} else {
					if previousActualStreak < 0 {
						status.Score += previousActualStreak * previousActualStreak * h.BasePoints
					}
					status.Score += (h.ActualStreak + 1) * (h.ActualStreak + 1) * h.BasePoints
				}
			}
			if toggleUndonePrevious {
				h.LastStreakEnd = removePeriod(h.Repetition, h.Deadline)
				h.LastStreak = h.ActualStreak - 1
				h.ActualStreak = -1
				h.Successes -= 1
				if h.Done {
					status.Score -= h.ActualStreak * h.ActualStreak * h.BasePoints - h.BasePoints
					status.Today -= h.ActualStreak * h.ActualStreak * h.BasePoints - h.BasePoints
					h.ActualStreak = 1
					h.LastStreak = h.LastStreak - 1
				}
				status.Score -= (h.LastStreak + 1) * (h.LastStreak + 1) * h.BasePoints + h.BasePoints
			}
		}
	}
}

func countPointChange(h *resources.Habit, status *resources.Status, coefficient int) {
	change := h.ActualStreak * h.ActualStreak * h.BasePoints * coefficient
	switch h.Repetition {
	case resources.HBT_REPETITION_WEEKLY:
		change *= 2
	case resources.HBT_REPETITION_MONTHLY:
		change *= 3
	}
	if h.Done {
		status.Score += change
		status.Today += change
	} else {
		status.Score -= change
		status.Today -= change
	}
}

func succeedHabit(h *resources.Habit, deadline *time.Time) {
	if h.ActualStreak < 0 {
		h.ActualStreak = 0
	}
	if deadline != nil && h.LastStreakEnd != nil && deadline.Equal(*h.LastStreakEnd) && h.LastStreak > 0 {
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

func removePeriod(repetition string, deadline *time.Time) *time.Time {
	if deadline == nil {
		return nil
	}
	switch repetition {
	case resources.HBT_REPETITION_DAILY:
		return utils.GetTimePointer(deadline.Add(-24 * time.Hour))
	case resources.HBT_REPETITION_WEEKLY:
		return utils.GetTimePointer(deadline.Add(-7 * 24 * time.Hour))
	case resources.HBT_REPETITION_MONTHLY:
		return utils.GetTimePointer(deadline.AddDate(0, -1, 0))
	}
	return nil
}

func getNewHabit() resources.Entity {
	return &resources.Habit{}
}

func getSyncHabitFunc(changeStatus *resources.Status) func (resources.Entity) func () {
	return func (entity resources.Entity) func () {
		return func () {
			h := entity.(*resources.Habit)
			if !h.Active {
				return
			}

			if h.Deadline.Before(time.Now()) {
				if h.ActualStreak > 49 && *h.LastStreakEnd != *h.Deadline {
					h.Done = true
					h.Tries += 1
					succeedHabit(h, h.Deadline)
				} else if h.Negative && h.Limit > h.Count {
					h.Done = true
					h.Tries += 1
					succeedHabit(h, h.Deadline)
					changeStatus.Score += h.ActualStreak * h.ActualStreak * h.BasePoints
				} else {
					numberOfMissedDeadlines := getNumberOfMissedDeadlines(h)
					for i := 0; i < numberOfMissedDeadlines; i++ {
						// if the last period
						if i == numberOfMissedDeadlines - 1 {
							// not done or not already failed
							if !h.Done && (h.LastStreakEnd == nil || *h.LastStreakEnd != *h.Deadline) {
								failHabit(h)
								changeStatus.Score -= h.ActualStreak * h.ActualStreak * h.BasePoints
							}
						} else {
							failHabit(h)
							changeStatus.Score -= h.ActualStreak * h.ActualStreak * h.BasePoints
						}
						h.Deadline = addPeriod(h.Repetition, h.Deadline)
					}
					h.Done = false
					h.Tries += numberOfMissedDeadlines
				}
			}
		}
	}
}

func getNumberOfMissedDeadlines(h *resources.Habit) int {
	return (int(time.Now().Sub(*(h.Deadline)).Hours()) / utils.GetDurationForRepetitionPeriod(h.Repetition)) + 1
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

func addHabit(name, repetition, description, deadline, goalId string, activeFlag, negativeFlag bool, basePoints, repetitionGoal int) {
	h := resources.NewHabit(name)
	if description != "" {
		h.Description = description
	}
	if goalId != "" {
		h.Goal = &resources.Goal{Id: goalId}
	}
	if negativeFlag {
		h.Negative = true
	}
	if repetitionGoal != -1 {
		h.Limit = repetitionGoal
	}
	if activeFlag {
		activateHabit(h, repetition)
		if repetition != resources.HBT_REPETITION_DAILY {
			if deadline == "" {
				h.Deadline = utils.GetFirstSaturday()
			} else {
				h.Deadline = utils.ParseTime(resources.DATE_FORMAT, deadline)
			}
		}
		if basePoints != -1 {
			h.BasePoints = basePoints
		}
	}
	tr := db.NewTransaction()
	tr.Add(func () error {
		err := tr.AddEntity(resources.DB_DEFAULT_HABITS_BUCKET_NAME, h)
		if err != nil {
			return err
		}
		//if activeFlag {
		//	_, colour, _ := h.GetIconColourAndOrder()
		//	anybar.AddToActivePorts(h.Name, colour, resources.DB_DEFAULT_HABITS_BUCKET_NAME, h.Id, tr)
		//}
		return nil
	})
	tr.Execute()
}

func deleteHabit(habitId string) {
	t := db.NewTransaction()
	t.Add(func () error {
		h := &resources.Habit{}
		err := t.RetrieveEntity(resources.DB_DEFAULT_HABITS_BUCKET_NAME, []byte(habitId), h, true)
		if err != nil {
			return err
		}
		if h.Goal != nil {
			goal := &resources.Goal{}
			err = t.ModifyEntity(resources.DB_DEFAULT_GOALS_BUCKET_NAME, []byte(h.Goal.Id), true, goal, func () {
				goal.Habit = nil
			})
			if err != nil {
				return err
			}
		}
		if h.Active {
			anybar.RemoveAndQuit(resources.DB_DEFAULT_HABITS_BUCKET_NAME, habitId, t)
		}
		return t.DeleteEntity(resources.DB_DEFAULT_HABITS_BUCKET_NAME, []byte(habitId))
	})
	t.Execute()
}

func modifyHabit(habitId, name, repetition, description, deadline, goalId string, toggleActive, toggleDone, toggleDonePrevious, toggleUndonePrevious, negativeFlag bool, basePoints, repetitionGoal int) {
	habit := &resources.Habit{}
	habitStatus := &resources.Status{}
	status := &resources.Status{}
	t := db.NewTransaction()
	t.Add(func () error {
		modifyHabit := getModifyHabitFunc(habit, name, repetition, description, deadline, goalId, toggleActive, toggleDone, toggleDonePrevious, toggleUndonePrevious, negativeFlag, basePoints, repetitionGoal, habitStatus, t)
		if err := t.ModifyEntity(resources.DB_DEFAULT_HABITS_BUCKET_NAME, []byte(habitId), false, habit, modifyHabit); err != nil {
			return err
		}
		if err := t.ModifyEntity(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_ACTUAL_STATUS_KEY, true, status, getAddScoreFunc(status, habitStatus)); err != nil {
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
		return tr.RetrieveEntity(resources.DB_DEFAULT_HABITS_BUCKET_NAME, []byte(habitId), habit, false)
	})
	tr.Execute()
	return habit
}

func getHabits() map[string]*resources.Habit {
	habits := map[string]*resources.Habit{}
	db.RetrieveEntities(resources.DB_DEFAULT_HABITS_BUCKET_NAME, false, func (id string) resources.Entity {
		habits[id] = &resources.Habit{}
		return habits[id]
	})
	return habits
}

func filterHabits(shallow bool, filter func (*resources.Habit) bool) map[string]*resources.Habit {
	habits := map[string]*resources.Habit{}
	tr := db.NewTransaction()
	tr.Add(func () error { return filterHabitsModal(tr, shallow, habits, filter) })
	tr.Execute()
	return habits
}

func filterHabitsModal(t resources.Transaction, shallow bool, habits map[string]*resources.Habit, filter func (*resources.Habit) bool) error {
	var entity *resources.Habit
	getNewEntityFunc := func () resources.Entity {
		entity = &resources.Habit{}
		return entity
	}
	addEntityFunc := func () { habits[entity.Id] = entity }
	entityFilter := func () bool { return filter(entity) }
	return t.FilterEntities(resources.DB_DEFAULT_HABITS_BUCKET_NAME, shallow, addEntityFunc, getNewEntityFunc, entityFilter)
}

func getActiveHabits() map[string]*resources.Habit {
	return FilterHabits(func (h *resources.Habit) bool { return h.Active })
}

func getNonActiveHabits() map[string]*resources.Habit {
	return FilterHabits(func (h *resources.Habit) bool { return !h.Active })
}
