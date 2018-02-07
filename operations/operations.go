package operations

import (
	"github.com/7joe7/personalmanager/operations/goals"
	"github.com/7joe7/personalmanager/resources"
)

func GetPointStats() map[string]*resources.PointStat {
	return getPointStats()
}

func GetDayPlan() map[string]resources.PlannedItem {
	return getDayPlan()
}

func GetReview() *resources.Review {
	return getReview()
}

func ModifyReview(cmd *resources.Command) {
	modifyReview(cmd)
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

func AddTask(cmd *resources.Command) {
	addTask(cmd)
}

func DeleteTask(taskId string) {
	deleteTask(taskId)
}

func ModifyTask(cmd *resources.Command) {
	modifyTask(cmd)
}

func GetTask(taskId string) *resources.Task {
	return getTask(taskId)
}

func GetTasks() map[string]*resources.Task {
	return getTasks()
}

func GetPersonalTasks() map[string]*resources.Task {
	return getPersonalTasks()
}

func FilterTasks(filter func(*resources.Task) bool) map[string]*resources.Task {
	return filterTasks(false, filter)
}

func FilterTasksModal(tr resources.Transaction, shallow bool, tasks map[string]*resources.Task, filter func(*resources.Task) bool) error {
	return filterTasksModal(tr, shallow, tasks, filter)
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

func GetGoalTasks(id string) map[string]*resources.Task {
	return getGoalTasks(id)
}

func GetProjectGoals(id string) map[string]*resources.Goal {
	return goals.GetProjectGoals(id)
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

func AddHabit(cmd *resources.Command) {
	addHabit(cmd)
}

func DeleteHabit(habitId string) {
	deleteHabit(habitId)
}

func ModifyHabit(cmd *resources.Command) {
	modifyHabit(cmd)
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

func AddProject(cmd *resources.Command) {
	addProject(cmd)
}

func DeleteProject(projectId string) {
	deleteProject(projectId)
}

func ModifyProject(cmd *resources.Command) {
	modifyProject(cmd)
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
