package operations

import (
	"time"

	"fmt"
	"strconv"

	"github.com/7joe7/personalmanager/db"
	"github.com/7joe7/personalmanager/resources"
	"github.com/7joe7/personalmanager/utils"
)

func synchronize(t resources.Transaction, backup bool) {
	t.Add(func() error {
		lastSync := string(t.GetValue(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_LAST_SYNC_KEY))
		if lastSync == "" || isTimeForSync(lastSync) {
			if backup {
				db.BackupDatabase()
			}
			if utils.IsSunday() {
				weeksLeftText := string(t.GetValue(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_WEEKS_LEFT))
				if weeksLeftText != "" {
					weeksLeft, err := strconv.Atoi(weeksLeftText)
					if err != nil {
						panic(err)
					}
					weeksLeft -= 1
					err = t.SetValue(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_WEEKS_LEFT, []byte(fmt.Sprint(weeksLeft)))
					if err != nil {
						panic(err)
					}
					resources.WaitGroup.Add(2)
					go func() {
						resources.Abr.Quit(resources.ANY_PORT_WEEKS_LEFT)
						resources.Abr.StartWithIcon(resources.ANY_PORT_WEEKS_LEFT, weeksLeftText, resources.ANY_CMD_BROWN)
					}()
				}
			}
			habitStatus := &resources.Status{}
			err := t.MapEntities(resources.DB_DEFAULT_HABITS_BUCKET_NAME, true, getNewHabit, getSyncHabitFunc(habitStatus))
			if err != nil {
				return err
			}
			err = t.MapEntities(resources.DB_DEFAULT_TASKS_BUCKET_NAME, true, getNewTask, getSyncTaskFunc())
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
		}
		return nil
	})
}

func synchronizeAnybarPorts(t resources.Transaction) {
	t.Add(func() error {
		var err error
		activeTaskId := t.GetValue(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_ACTUAL_ACTIVE_TASK_KEY)
		if !resources.Abr.Ping(resources.ANY_PORT_ACTIVE_TASK) && activeTaskId != nil && string(activeTaskId) != "" {
			activeTask := &resources.Task{}
			err = t.RetrieveEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, activeTaskId, activeTask, true)
			if err != nil {
				return err
			}
			if activeTask.InProgress {
				resources.WaitGroup.Add(1)
				go resources.Abr.StartWithIcon(resources.ANY_PORT_ACTIVE_TASK, activeTask.Name, resources.ANY_CMD_BLUE)
			}
		}
		weeksLeft := string(t.GetValue(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_WEEKS_LEFT))
		if weeksLeft != "" && !resources.Abr.Ping(resources.ANY_PORT_WEEKS_LEFT) {
			resources.WaitGroup.Add(1)
			go resources.Abr.StartWithIcon(resources.ANY_PORT_WEEKS_LEFT, weeksLeft, resources.ANY_CMD_BROWN)
		}
		activePorts := resources.Abr.GetActivePorts(t)
		resources.WaitGroup.Add(1)
		go resources.Abr.EnsureActivePorts(activePorts)
		return nil
	})
}

func addBonusIfAllHabitsDone(t resources.Transaction, repetition string, changeStatus *resources.Status) error {
	habits := map[string]*resources.Habit{}
	err := filterHabitsModal(t, true, habits, func(h *resources.Habit) bool { return h.Active && h.Repetition == repetition })
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
