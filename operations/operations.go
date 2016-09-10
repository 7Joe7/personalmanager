package operations

import (
	"github.com/7joe7/personalmanager/resources"
	"github.com/7joe7/personalmanager/db"
)

func GetModifyHabitFunc(h *resources.Habit, name, repetition, deadline string, toggleActive, toggleDone bool, basePoints int, scoreChange *int) func () {
	return getModifyHabitFunc(h, name, repetition, deadline, toggleActive, toggleDone, basePoints, scoreChange)
}

func GetSyncHabitFunc(h *resources.Habit, scoreChange *int) func () {
	return getSyncHabitFunc(h, scoreChange)
}

func GetAddScoreFunc(s *resources.Status, scoreChange int) func () {
	return getAddScoreFunc(s, scoreChange)
}

func GetModifyTaskFunc(t *resources.Task, name, projectId string) func () {
	return getModifyTaskFunc(t, name, projectId)
}

func GetModifyProjectFunc(p *resources.Project, name string) func () {
	return getModifyProjectFunc(p, name)
}

func GetModifyTagFunc(t *resources.Tag, name string) func () {
	return getModifyTagFunc(t, name)
}

func GetModifyGoalFunc(g *resources.Goal, name string) func () {
	return getModifyGoalFunc(g, name)
}

func GetSyncStatusFunc(s *resources.Status, scoreChange int) func () {
	return getSyncStatusFunc(s, scoreChange)
}

func Synchronize(t *db.Transaction) {
	synchronize(t)
}

func InitializeBuckets(t *db.Transaction) {
	initializeBuckets(t, resources.BUCKETS_TO_INTIALIZE)
}

func EnsureValues(t *db.Transaction) {
	ensureValues(t)
}

func AddTask(name, projectId string) string {
	return addTask(name, projectId)
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

func AddHabit(name, repetition string, activeFlag bool) string {
	return addHabit(name, repetition, activeFlag)
}

func DeleteHabit(habitId string) {
	deleteHabit(habitId)
}

func ModifyHabit(habitId, name, repetition, deadline string, toggleActive, toggleDone bool, basePoints int) {
	modifyHabit(habitId, name, repetition, deadline, toggleActive, toggleDone, basePoints)
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
