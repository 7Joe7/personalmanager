package db

import (
	"github.com/7joe7/personalmanager/resources"
	"encoding/json"
	"time"
)

func AddTask(task *resources.Task) string {
	return addEntity(task, resources.DB_DEFAULT_TASKS_BUCKET_NAME)
}

func AddTag(tag *resources.Tag) string {
	return addEntity(tag, resources.DB_DEFAULT_TAGS_BUCKET_NAME)
}

func AddProject(project *resources.Project) string {
	return addEntity(project, resources.DB_DEFAULT_PROJECTS_BUCKET_NAME)
}

func AddGoal(goal *resources.Goal) string {
	return addEntity(goal, resources.DB_DEFAULT_GOALS_BUCKET_NAME)
}

func AddHabit(habit *resources.Habit) string {
	return addEntity(habit, resources.DB_DEFAULT_HABITS_BUCKET_NAME)
}

func DeleteTag(tagId string) {
	deleteEntity(tagId, resources.DB_DEFAULT_TAGS_BUCKET_NAME)
}

func DeleteTask(taskId string) {
	deleteEntity(taskId, resources.DB_DEFAULT_TASKS_BUCKET_NAME)
}

func DeleteProject(projectId string) {
	deleteEntity(projectId, resources.DB_DEFAULT_PROJECTS_BUCKET_NAME)
}

func DeleteGoal(goalId string) {
	deleteEntity(goalId, resources.DB_DEFAULT_GOALS_BUCKET_NAME)
}

func DeleteHabit(habitId string) {
	deleteEntity(habitId, resources.DB_DEFAULT_HABITS_BUCKET_NAME)
}

func ModifyTask(taskId, name, projectId string) {
	task := &resources.Task{}
	modifyEntity(taskId, task, func() {
		if name != "" {
			task.Name = name
		}
		if projectId != "" {
			task.Project = GetProject(projectId)
		}
	}, resources.DB_DEFAULT_TASKS_BUCKET_NAME)
}

func ModifyProject(projectId, name string) {
	project := &resources.Project{}
	modifyEntity(projectId, project, func () {
		if name != "" {
			project.Name = name
		}
	}, resources.DB_DEFAULT_PROJECTS_BUCKET_NAME)
}

func ModifyTag(tagId, name string) {
	tag := &resources.Tag{}
	modifyEntity(tagId, tag, func () {
		if name != "" {
			tag.Name = name
		}
	}, resources.DB_DEFAULT_TAGS_BUCKET_NAME)
}

func ModifyGoal(goalId, name string) {
	goal := &resources.Goal{}
	modifyEntity(goalId, goal, func () {
		if name != "" {
			goal.Name = name
		}
	}, resources.DB_DEFAULT_GOALS_BUCKET_NAME)
}

func ModifyHabit(habitId, name, repetition string, toggleActive, toggleDone bool) {
	habit := &resources.Habit{}
	modifyEntity(habitId, habit, func () {
		if name != "" {
			habit.Name = name
		}
		if toggleActive {
			if habit.Active {
				habit.Active = false
				habit.Deadline = nil
				habit.Done = false
				habit.ActualStreak = 0
				habit.LastStreakEnd = nil
				habit.LastStreak = 0
				habit.Repetition = ""
			} else {
				habit.Active = true
				if repetition == "" {
					repetition = "Daily"
				}
				habit.Repetition = repetition
				habit.Deadline = getTimePointer(time.Now().Add(24 * time.Hour).Truncate(24 * time.Hour))
			}
		}
		if habit.Active && toggleDone {
			if habit.Done {
				habit.Done = false
				habit.LastStreakEnd = getTimePointer(*habit.Deadline)
				habit.LastStreak = habit.ActualStreak
				habit.ActualStreak -= 1
			} else {
				habit.Done = true
				if habit.Deadline.Equal(*habit.LastStreakEnd) {
					habit.LastStreakEnd = nil
					habit.ActualStreak = habit.LastStreak
				}
				habit.ActualStreak += 1
			}
		}
	}, resources.DB_DEFAULT_HABITS_BUCKET_NAME)
}

func GetProject(projectId string) *resources.Project {
	project := &resources.Project{}
	retrieveEntity(projectId, project, resources.DB_DEFAULT_PROJECTS_BUCKET_NAME)
	return project
}

func GetTask(taskId string) *resources.Task {
	task := &resources.Task{}
	retrieveEntity(taskId, task, resources.DB_DEFAULT_TASKS_BUCKET_NAME)
	return task
}

func GetTag(tagId string) *resources.Tag {
	tag := &resources.Tag{}
	retrieveEntity(tagId, tag, resources.DB_DEFAULT_TAGS_BUCKET_NAME)
	return tag
}

func GetGoal(goalId string) *resources.Goal {
	goal := &resources.Goal{}
	retrieveEntity(goalId, goal, resources.DB_DEFAULT_GOALS_BUCKET_NAME)
	return goal
}

func GetHabit(habitId string) *resources.Habit {
	habit := &resources.Habit{}
	retrieveEntity(habitId, habit, resources.DB_DEFAULT_HABITS_BUCKET_NAME)
	return habit
}

func GetTasks() map[string]*resources.Task {
	tasks := map[string]*resources.Task{}
	retrieveEntities(func (id string) interface{} {
		tasks[id] = &resources.Task{}
		return tasks[id]
	}, resources.DB_DEFAULT_TASKS_BUCKET_NAME)
	return tasks
}

func GetTags() map[string]*resources.Tag {
	tags := map[string]*resources.Tag{}
	retrieveEntities(func (id string) interface{} {
		tags[id] = &resources.Tag{}
		return tags[id]
	}, resources.DB_DEFAULT_TAGS_BUCKET_NAME)
	return tags
}

func GetProjects() map[string]*resources.Project {
	projects := map[string]*resources.Project{}
	retrieveEntities(func (id string) interface{} {
		projects[id] = &resources.Project{}
		return projects[id]
	}, resources.DB_DEFAULT_PROJECTS_BUCKET_NAME)
	return projects
}

func GetGoals() map[string]*resources.Goal {
	goals := map[string]*resources.Goal{}
	retrieveEntities(func (id string) interface{} {
		goals[id] = &resources.Goal{}
		return goals[id]
	}, resources.DB_DEFAULT_GOALS_BUCKET_NAME)
	return goals
}

func GetHabits() map[string]*resources.Habit {
	habits := map[string]*resources.Habit{}
	retrieveEntities(func (id string) interface{} {
		habits[id] = &resources.Habit{}
		return habits[id]
	}, resources.DB_DEFAULT_HABITS_BUCKET_NAME)
	return habits
}

func FilterHabits(filter func(*resources.Habit) bool) map[string]*resources.Habit {
	filteredHabits := map[string]*resources.Habit{}
	filterEntities(func (id string, value []byte) error {
		habit := &resources.Habit{}
		if err := json.Unmarshal(value, habit); err != nil {
			return err
		}
		if filter(habit) {
			filteredHabits[id] = habit
		}
		return nil
	}, resources.DB_DEFAULT_HABITS_BUCKET_NAME)
	return filteredHabits
}

func Open() {
	open(resources.DB_PATH)
}

func InitializeBuckets() {
	initializeBuckets(resources.BUCKETS_TO_INTIALIZE)
}

func Synchronize() {
	synchronize()
}
