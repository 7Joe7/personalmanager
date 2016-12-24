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
	synchronize(t, true)
}

func SynchronizeAnybarPorts(t resources.Transaction) {
	synchronizeAnybarPorts(t)
}

func InitializeBuckets(t resources.Transaction) {
	initializeBuckets(t, resources.BUCKETS_TO_INTIALIZE)
}

func EnsureValues(t resources.Transaction) {
	ensureValues(t)
}

func AddTask(name, projectId, goalId, deadline, estimate, scheduled, taskType, note string, active bool, basePoints int) {
	addTask(name, projectId, goalId, deadline, estimate, scheduled, taskType, note, active, basePoints)
}

func DeleteTask(taskId string) {
	deleteTask(taskId)
}

func ModifyTask(taskId, name, projectId, goalId, deadline, estimate, scheduled, taskType, note string, basePoints int, activeFlag, doneFlag bool) {
	modifyTask(taskId, name, projectId, goalId, deadline, estimate, scheduled, taskType, note, basePoints, activeFlag, doneFlag)
}

func GetTask(taskId string) *resources.Task {
	return getTask(taskId)
}

func GetTasks() map[string]*resources.Task {
	return getTasks()
}

func GetPersonalTasks() map[string] *resources.Task {
	return getPersonalTasks()
}

func FilterTasks(filter func(*resources.Task) bool) map[string]*resources.Task {
	return filterTasks(false, filter)
}

func FilterTasksSlice(filter func(*resources.Task) bool) []*resources.Task {
	return filterTasksSlice(false, filter)
}

func GetNextTasks() map[string]*resources.Task {
	return getNextTasks()
}

func GetUnscheduledTasks() map[string]*resources.Task {
	return getUnscheduledTasks()
}

func GetShoppingTasks() map[string]*resources.Task {
	return getShoppingTasks()
}

func GetWorkNextTasks() map[string]*resources.Task {
	return getWorkNextTasks()
}

func GetWorkUnscheduledTasks() map[string]*resources.Task {
	return getWorkUnscheduledTasks()
}

func AddHabit(name, repetition, description, deadline, goalId string, activeFlag, negativeFlag bool, basePoints int) {
	addHabit(name, repetition, description, deadline, goalId, activeFlag, negativeFlag, basePoints)
}

func DeleteHabit(habitId string) {
	deleteHabit(habitId)
}

func ModifyHabit(habitId, name, repetition, description, deadline, goalId string, toggleActive, toggleDone, toggleDonePrevious, toggleUndonePrevious, negativeFlag bool, basePoints int) {
	modifyHabit(habitId, name, repetition, description, deadline, goalId, toggleActive, toggleDone, toggleDonePrevious, toggleUndonePrevious, negativeFlag, basePoints)
}

func GetHabit(habitId string) *resources.Habit {
	return getHabit(habitId)
}

func GetHabits() map[string]*resources.Habit {
	return getHabits()
}

func FilterHabits(filter func(*resources.Habit) bool) map[string]*resources.Habit {
	return filterHabits(false, filter)
}

func GetActiveHabits() map[string]*resources.Habit {
	return getActiveHabits()
}

func GetNonActiveHabits() map[string]*resources.Habit {
	return getNonActiveHabits()
}

func AddGoal(name, projectId, habitId string, habitRepetitionGoal int) string {
	return addGoal(name, projectId, habitId, habitRepetitionGoal)
}

func DeleteGoal(goalId string) {
	deleteGoal(goalId)
}

func ModifyGoal(goalId, name, taskId, projectId, habitId string, activeFlag, doneFlag bool, habitRepetitionGoal int) {
	modifyGoal(goalId, name, taskId, projectId, habitId, activeFlag, doneFlag, habitRepetitionGoal)
}

func GetGoal(goalId string) *resources.Goal {
	return getGoal(goalId)
}

func GetGoals() map[string]*resources.Goal {
	return getGoals()
}

func FilterGoals(filter func(*resources.Goal) bool) map[string]*resources.Goal {
	return filterGoals(false, filter)
}

func GetActiveGoals() map[string]*resources.Goal {
	return getActiveGoals()
}

func GetNonActiveGoals() map[string]*resources.Goal {
	return getNonActiveGoals()
}

func AddProject(name string) {
	addProject(name)
}

func DeleteProject(projectId string) {
	deleteProject(projectId)
}

func ModifyProject(projectId, name, taskId, goalId string, activeFlag, doneFlag bool) {
	modifyProject(projectId, name, taskId, goalId, activeFlag, doneFlag)
}

func GetProject(projectId string) *resources.Project {
	return getProject(projectId)
}

func GetProjects() map[string]*resources.Project {
	return getProjects()
}

func FilterProjects(filter func(*resources.Project) bool) map[string]*resources.Project {
	return filterProjects(false, filter)
}

func GetActiveProjects() map[string]*resources.Project {
	return getActiveProjects()
}

func GetInactiveProjects() map[string]*resources.Project {
	return getInactiveProjects()
}

func SyncWithJira() {
	syncWithJira()
}

func ExportShoppingTasks() {
	exportTasks()
}

func SetEmail(email string) {
	setEmail(email)
}