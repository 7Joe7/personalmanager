package operations

import (
	"github.com/7joe7/personalmanager/resources"
	"github.com/7joe7/personalmanager/db"
	"github.com/7joe7/personalmanager/anybar"
)

func getModifyGoalFunc(g *resources.Goal, name, taskId, projectId, habitId string, activeFlag, doneFlag bool, habitRepetitionGoal, priority int, tr resources.Transaction) func () {
	return func () {
		if name != "" {
			g.Name = name
		}
		if taskId != "" {
			task := &resources.Task{}
			err := tr.ModifyEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, []byte(taskId), true, task, func() {
				task.Goal = g
				task.BasePoints = g.Priority
				if g.Active {
					task.Scheduled = resources.TASK_SCHEDULED_NEXT
				}
			})
			if err != nil {
				panic(err)
			}
			g.Tasks = append(g.Tasks, task)
		}
		if projectId != "" {
			project := &resources.Project{}
			err := tr.ModifyEntity(resources.DB_DEFAULT_PROJECTS_BUCKET_NAME, []byte(projectId), true, project, func () {
				project.Goals = append(project.Goals, g)
			})
			if err != nil {
				panic(err)
			}
			g.Project = project
		}
		if habitId != "" {
			habit := &resources.Habit{}
			err := tr.ModifyEntity(resources.DB_DEFAULT_HABITS_BUCKET_NAME, []byte(habitId), true, habit, func () {
				habit.Goal = g
				habit.BasePoints = priority
			})
			if err != nil {
				panic(err)
			}
			g.Habit = habit
		}
		if habitRepetitionGoal != -1 {
			g.HabitRepetitionGoal = habitRepetitionGoal
		}
		if priority != -1 {
			for i := 0; i < len(g.Tasks); i++ {
				if g.Tasks[i].BasePoints == g.Priority {
					task := &resources.Task{}
					err := tr.ModifyEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, []byte(g.Tasks[i].Id), true, task, func () {
						task.BasePoints = priority
					})
					if err != nil {
						panic(err)
					}
				}
			}
			g.Priority = priority
		}
		if activeFlag {
			if g.Active {
				toggleSubTasksScheduling(resources.TASK_SCHEDULED_NEXT, resources.TASK_NOT_SCHEDULED, g, tr)
				g.Active = false
				anybar.RemoveAndQuit(resources.DB_DEFAULT_GOALS_BUCKET_NAME, g.Id, tr)
			} else {
				toggleSubTasksScheduling(resources.TASK_NOT_SCHEDULED, resources.TASK_SCHEDULED_NEXT, g, tr)
				g.Active = true
				anybar.AddToActivePorts(g.Name, resources.ANY_CMD_CYAN, resources.DB_DEFAULT_GOALS_BUCKET_NAME, g.Id, tr)
			}
		}
		if doneFlag {
			var scoreChange int
			for i := 0; i < len(g.Tasks); i++ {
				if g.Tasks[i].Done {
					scoreChange += countScoreChange(g.Tasks[i])
				}
			}
			if g.Done {
				g.Done = false
				scoreChange = -scoreChange
				if g.Active {
					anybar.AddToActivePorts(g.Name, resources.ANY_CMD_CYAN, resources.DB_DEFAULT_GOALS_BUCKET_NAME, g.Id, tr)
				}
			} else {
				g.Done = true
				if g.Active {
					anybar.RemoveAndQuit(resources.DB_DEFAULT_GOALS_BUCKET_NAME, g.Id, tr)
				}
			}
			status := &resources.Status{}
			err := tr.ModifyEntity(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_ACTUAL_STATUS_KEY, true, status, func () {
				status.Score += scoreChange
				status.Today += scoreChange
			})
			if err != nil {
				panic(err)
			}
		}
	}
}

func toggleSubTasksScheduling(scheduledCriteria, scheduledSet string, g *resources.Goal, tr resources.Transaction) {
	for i := 0; i < len(g.Tasks); i++ {
		if g.Tasks[i].Scheduled == scheduledCriteria && !g.Tasks[i].Done {
			task := &resources.Task{}
			err := tr.ModifyEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, []byte(g.Tasks[i].Id), true, task, func () {
				task.Scheduled = scheduledSet
			})
			if err != nil {
				panic(err)
			}
		}
	}
}

func addGoal(name, projectId, habitId string, habitRepetitionGoal, priority int) string {
	goal := resources.NewGoal(name)
	if projectId != "" {
		goal.Project = &resources.Project{Id: projectId}
	}
	if habitId != "" {
		goal.Habit = &resources.Habit{Id: habitId}
	}
	if habitRepetitionGoal != -1 {
		goal.HabitRepetitionGoal = habitRepetitionGoal
	}
	if priority != -1 {
		goal.Priority = priority
	}
	tr := db.NewTransaction()
	tr.Add(func () error {
		return tr.AddEntity(resources.DB_DEFAULT_GOALS_BUCKET_NAME, goal)
	})
	if projectId != "" {
		tr.Add(func () error {
			project := &resources.Project{}
			err := tr.ModifyEntity(resources.DB_DEFAULT_PROJECTS_BUCKET_NAME, []byte(projectId), true, project, func () {
				project.Goals = append(project.Goals, goal)
			})
			return err
		})
	}
	if habitId != "" {
		tr.Add(func () error {
			habit := &resources.Habit{}
			err := tr.ModifyEntity(resources.DB_DEFAULT_HABITS_BUCKET_NAME, []byte(habitId), true, habit, func () {
				habit.Goal = goal
				habit.BasePoints = goal.Priority
			})
			return err
		})
	}
	tr.Execute()
	return goal.Id
}

func deleteGoal(goalId string) {
	tr := db.NewTransaction()
	tr.Add(func () error {
		goal := &resources.Goal{}
		err := tr.RetrieveEntity(resources.DB_DEFAULT_GOALS_BUCKET_NAME, []byte(goalId), goal, true)
		if err != nil {
			return err
		}
		for i := 0; i < len(goal.Tasks); i++ {
			task := &resources.Task{}
			err = tr.ModifyEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, []byte(goal.Tasks[i].Id), true, task, func () {
				task.Goal = nil
			})
			if err != nil {
				return err
			}
		}
		if goal.Project != nil {
			project := &resources.Project{}
			err = tr.ModifyEntity(resources.DB_DEFAULT_PROJECTS_BUCKET_NAME, []byte(goal.Project.Id), true, project, func () {
				for i := 0; i < len(project.Goals); i++ {
					if project.Goals[i].Id == goal.Id {
						project.Goals = append(project.Goals[:i], project.Goals[i+1:]...)
						break
					}
				}
			})
			if err != nil {
				return err
			}
		}
		if goal.Habit != nil {
			habit := &resources.Habit{}
			err := tr.ModifyEntity(resources.DB_DEFAULT_HABITS_BUCKET_NAME, []byte(goal.Habit.Id), true, habit, func () {
				habit.Goal = nil
			})
			if err != nil {
				return err
			}
		}
		if goal.Active {
			anybar.RemoveAndQuit(resources.DB_DEFAULT_GOALS_BUCKET_NAME, goalId, tr)
		}
		err = tr.DeleteEntity(resources.DB_DEFAULT_GOALS_BUCKET_NAME, []byte(goalId))
		if err != nil {
			return err
		}
		return nil
	})
	tr.Execute()
}

func modifyGoal(goalId, name, taskId, projectId, habitId string, activeFlag, doneFlag bool, habitRepetitionGoal, priority int) {
	goal := &resources.Goal{}
	tr := db.NewTransaction()
	tr.Add(func () error {
		return tr.ModifyEntity(resources.DB_DEFAULT_GOALS_BUCKET_NAME, []byte(goalId), false, goal, getModifyGoalFunc(goal, name, taskId, projectId, habitId, activeFlag, doneFlag, habitRepetitionGoal, priority, tr))
	})
	tr.Execute()
}

func getGoal(goalId string) *resources.Goal {
	goal := &resources.Goal{}
	tr := db.NewTransaction()
	tr.Add(func () error {
		return tr.RetrieveEntity(resources.DB_DEFAULT_GOALS_BUCKET_NAME, []byte(goalId), goal, false)
	})
	tr.Execute()
	return goal
}

func getGoals() map[string]*resources.Goal {
	goals := map[string]*resources.Goal{}
	db.RetrieveEntities(resources.DB_DEFAULT_GOALS_BUCKET_NAME, false, func (id string) resources.Entity {
		goals[id] = &resources.Goal{}
		return goals[id]
	})
	return goals
}

func getActiveGoals() map[string]*resources.Goal {
	return FilterGoals(func (g *resources.Goal) bool { return g.Active && !g.Done })
}

func getNonActiveGoals() map[string]*resources.Goal {
	return FilterGoals(func (g *resources.Goal) bool { return !g.Active })
}

func filterGoals(shallow bool, filter func(*resources.Goal) bool) map[string]*resources.Goal {
	goals := map[string]*resources.Goal{}
	var goal *resources.Goal
	getNewEntity := func () resources.Entity {
		goal = &resources.Goal{}
		return goal
	}
	addEntity := func () { goals[goal.Id] = goal }
	db.FilterEntities(resources.DB_DEFAULT_GOALS_BUCKET_NAME, shallow, addEntity, getNewEntity, func() bool { return filter(goal) })
	return goals
}
