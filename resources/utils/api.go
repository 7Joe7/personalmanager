package utils

import "github.com/7joe7/personalmanager/resources"

func RemoveTaskFromTasks(tasks []*resources.Task, taskToRemove *resources.Task) []*resources.Task {
	return removeTaskFromTasks(tasks, taskToRemove)
}

func RemoveHabitFromHabits(habits []*resources.Habit, habitToRemove *resources.Habit) []*resources.Habit {
	return removeHabitFromHabits(habits, habitToRemove)
}

func RemoveProjectFromProjects(projects []*resources.Project, projectToRemove *resources.Project) []*resources.Project {
	return removeProjectFromProjects(projects, projectToRemove)
}

func RemoveGoalFromGoals(goals []*resources.Goal, goalToRemove *resources.Goal) []*resources.Goal {
	return removeGoalFromGoals(goals, goalToRemove)
}

func RemoveTagFromTags(tags []*resources.Tag, tagToRemove *resources.Tag) []*resources.Tag {
	return removeTagFromTags(tags, tagToRemove)
}
