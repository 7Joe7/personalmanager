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

func AddTask(name, projectId string) {
	addTask(name, projectId)
}

func DeleteTask(taskId string) {
	deleteTask(taskId)
}

func ModifyTask(taskId, name, projectId string) {
	modifyTask(taskId, name, projectId)
}

func GetTask(taskId string) *resources.Task {
	return getTask(taskId)
}

func GetTasks() map[string]*resources.Task {
	return getTasks()
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
