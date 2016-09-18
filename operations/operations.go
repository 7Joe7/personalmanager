package operations

import (
	"github.com/7joe7/personalmanager/resources"
)

func GetReview() *resources.Review {
	return getReview()
}

func ModifyReview(repetition, deadline string) {
	modifyReview(repetition, deadline)
}

func Synchronize(t resources.Transaction) {
	synchronize(t)
}

func InitializeBuckets(t resources.Transaction) {
	initializeBuckets(t, resources.BUCKETS_TO_INTIALIZE)
}

func EnsureValues(t resources.Transaction) {
	ensureValues(t)
}

func AddTask(name, projectId, deadline, estimate string, active bool, basePoints int) {
	addTask(name, projectId, deadline, estimate, active, basePoints)
}

func DeleteTask(taskId string) {
	deleteTask(taskId)
}

func ModifyTask(taskId, name, projectId, deadline, estimate, scheduled string, basePoints int, activeFlag, doneFlag bool) {
	modifyTask(taskId, name, projectId, deadline, estimate, scheduled, basePoints, activeFlag, doneFlag)
}

func GetTask(taskId string) *resources.Task {
	return getTask(taskId)
}

func GetTasks() map[string]*resources.Task {
	return getTasks()
}

func FilterTasks(filter func(*resources.Task) bool) map[string]*resources.Task {
	return filterTasks(filter)
}

func GetNextTasks() map[string]*resources.Task {
	return getNextTasks()
}

func GetUnscheduledTasks() map[string]*resources.Task {
	return getUnscheduledTasks()
}

func AddHabit(name, repetition, deadline string, activeFlag bool, basePoints int) {
	addHabit(name, repetition, deadline, activeFlag, basePoints)
}

func DeleteHabit(habitId string) {
	deleteHabit(habitId)
}

func ModifyHabit(habitId, name, repetition, deadline string, toggleActive, toggleDone, toggleDonePrevious bool, basePoints int) {
	modifyHabit(habitId, name, repetition, deadline, toggleActive, toggleDone, toggleDonePrevious, basePoints)
}

func GetHabit(habitId string) *resources.Habit {
	return getHabit(habitId)
}

func GetHabits() map[string]*resources.Habit {
	return getHabits()
}

func FilterHabits(filter func(*resources.Habit) bool) map[string]*resources.Habit {
	return filterHabits(filter)
}

func GetActiveHabits() map[string]*resources.Habit {
	return getActiveHabits()
}

func GetNonActiveHabits() map[string]*resources.Habit {
	return getNonActiveHabits()
}
