package operations

import (
	"time"

	"github.com/7joe7/personalmanager/resources"
)

func synchronize(t resources.Transaction) {
	t.Add(func () error {
		lastSync := string(t.GetValue(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_LAST_SYNC_KEY))
		if lastSync == "" || isTimeForSync(lastSync) {
			habit := &resources.Habit{}
			habitStatus := &resources.Status{}
			err := t.MapEntities(resources.DB_DEFAULT_HABITS_BUCKET_NAME, true, habit, getSyncHabitFunc(habit, habitStatus, t))
			if err != nil {
				return err
			}
			err = addBonusIfAllHabitsDone(t, resources.HBT_REPETITION_DAILY, habitStatus)
			if err != nil {
				return err
			}
			err = addBonusIfAllHabitsDone(t, resources.HBT_REPETITION_WEEKLY, habitStatus)
			if err != nil {
				return err
			}
			err = addBonusIfAllHabitsDone(t, resources.HBT_REPETITION_MONTHLY, habitStatus)
			if err != nil {
				return err
			}
			status := &resources.Status{}
			err = t.ModifyEntity(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_ACTUAL_STATUS_KEY, true, status, getSyncStatusFunc(status, habitStatus))
			if err != nil {
				return err
			}
			err = t.SetValue(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_LAST_SYNC_KEY, []byte(time.Now().Format("Mon Jan 2 15:04:05 -0700 MST 2006")))
			if err != nil {
				return err
			}
			// TODO enable after cleaning db and
			//activeTaskId := t.GetValue(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_ACTUAL_ACTIVE_TASK_KEY)
			//if !anybar.Ping(resources.ANY_PORT_ACTIVE_TASK) && activeTaskId != nil && string(activeTaskId) != "" {
			//	activeTask := &resources.Task{}
			//	err = t.RetrieveEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, activeTaskId, activeTask)
			//	if err != nil {
			//		return err
			//	}
			//	resources.WaitGroup.Add(1)
			//	go anybar.StartWithIcon(resources.ANY_PORT_ACTIVE_TASK, activeTask.Name, resources.ANY_CMD_BLUE)
			//}
			//activePorts := anybar.GetActivePorts(t)
			//resources.WaitGroup.Add(1)
			//go anybar.EnsureActivePorts(activePorts)
		}
		return nil
	})
}

func addBonusIfAllHabitsDone(t resources.Transaction, repetition string, changeStatus *resources.Status) error {
	habits := map[string]*resources.Habit{}
	err := filterHabitsModal(t, true, habits, func (h *resources.Habit) bool { return h.Active && h.Repetition == repetition })
	if err != nil {
		return err
	}
	var undoneFound bool
	var pointsTogether int
	for _, habit := range habits {
		if habit.Done {
			pointsTogether += habit.ActualStreak * habit.ActualStreak * habit.BasePoints
		} else {
			undoneFound = true
			break
		}
	}
	if !undoneFound {
		changeStatus.Today += pointsTogether
		changeStatus.Score += pointsTogether
	}
	return nil
}

func isTimeForSync(lastSync string) bool {
	t, err := time.Parse("Mon Jan 2 15:04:05 -0700 MST 2006", lastSync)
	if err != nil {
		panic(err)
	}
	return t.Truncate(24 * time.Hour).Before(time.Now().Truncate(24 * time.Hour))
}
