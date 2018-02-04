package operations

import (
	"fmt"

	"time"

	"github.com/7joe7/personalmanager/db"
	"github.com/7joe7/personalmanager/resources"
)

func getDayPlan() map[string]resources.PlannedItem {
	tasks := map[string]*resources.Task{}
	habits := map[string]*resources.Habit{}
	tr := db.NewTransaction()
	tr.Add(
		func() error {
			return filterTasksModal(tr, false, tasks, func(task *resources.Task) bool {
				return task.Scheduled == resources.TASK_SCHEDULED_NEXT && !task.Done && task.Type == resources.TASK_TYPE_PERSONAL
			})
		})
	tr.Add(
		func() error {
			return filterHabitsModal(tr, false, habits, func(habit *resources.Habit) bool {
				return !habit.Done && habit.Active && (habit.AlarmTime == nil || habit.AlarmTime.Before(time.Now()))
			})
		})
	tr.Execute()
	plannedItems := map[string]resources.PlannedItem{}
	for id, habit := range habits {
		plannedItems[fmt.Sprintf("%sH", id)] = habit
	}
	for id, task := range tasks {
		plannedItems[fmt.Sprintf("%sT", id)] = task
	}
	return plannedItems
}
