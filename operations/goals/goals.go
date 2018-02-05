package goals

import (
	"github.com/7joe7/personalmanager/db"
	"github.com/7joe7/personalmanager/resources"
	rutils "github.com/7joe7/personalmanager/resources/utils"
)

func getModifyGoalFunc(g *resources.Goal, cmd *resources.Command, tr resources.Transaction) func() {
	return func() {
		var err error
		if cmd.Name != "" {
			g.Name = cmd.Name
		}
		if cmd.TaskID != "" && cmd.TaskID != "-" {
			task := &resources.Task{}
			err = tr.ModifyEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, []byte(cmd.TaskID), true, task, func() {
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
		switch cmd.ProjectID {
		case "-":
			if g.Project != nil {
				project := &resources.Project{}
				err = tr.ModifyEntity(resources.DB_DEFAULT_PROJECTS_BUCKET_NAME, []byte(g.Project.Id), true, project, func() {
					project.Goals = rutils.RemoveGoalFromGoals(project.Goals, g)
				})
				if err != nil {
					panic(err)
				}
				g.Project = nil
			}
		case "":
		default:
			project := &resources.Project{}
			err = tr.ModifyEntity(resources.DB_DEFAULT_PROJECTS_BUCKET_NAME, []byte(cmd.ProjectID), true, project, func() {
				project.Goals = append(project.Goals, g)
			})
			if err != nil {
				panic(err)
			}
			g.Project = project
		}
		switch cmd.HabitID {
		case "-":
			if g.Habit != nil {
				habit := &resources.Habit{}
				err = tr.ModifyEntity(resources.DB_DEFAULT_HABITS_BUCKET_NAME, []byte(g.Habit.Id), true, habit, func() {
					habit.Goal = nil
				})
				if err != nil {
					panic(err)
				}
				g.Habit = nil
			}
		case "":
		default:
			habit := &resources.Habit{}
			err = tr.ModifyEntity(resources.DB_DEFAULT_HABITS_BUCKET_NAME, []byte(cmd.HabitID), true, habit, func() {
				habit.Goal = g
				habit.BasePoints = cmd.BasePoints
			})
			if err != nil {
				panic(err)
			}
			g.Habit = habit
		}
		if cmd.HabitRepetitionGoal != -1 {
			g.HabitRepetitionGoal = cmd.HabitRepetitionGoal
		}
		if cmd.BasePoints != -1 {
			for i := 0; i < len(g.Tasks); i++ {
				if g.Tasks[i].BasePoints == g.Priority {
					task := &resources.Task{}
					err := tr.ModifyEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, []byte(g.Tasks[i].Id), true, task, func() {
						task.BasePoints = cmd.BasePoints
					})
					if err != nil {
						panic(err)
					}
				}
			}
			g.Priority = cmd.BasePoints
		}
		if cmd.ActiveFlag {
			if g.Active {
				toggleSubTasksScheduling(resources.TASK_SCHEDULED_NEXT, resources.TASK_NOT_SCHEDULED, g, tr)
				g.Active = false
				resources.Abr.RemoveAndQuit(resources.DB_DEFAULT_GOALS_BUCKET_NAME, g.Id, tr)
			} else {
				toggleSubTasksScheduling(resources.TASK_NOT_SCHEDULED, resources.TASK_SCHEDULED_NEXT, g, tr)
				g.Active = true
				resources.Abr.AddToActivePorts(g.Name, resources.ANY_CMD_CYAN, resources.DB_DEFAULT_GOALS_BUCKET_NAME, g.Id, tr)
			}
		}
		status := &resources.Status{}
		if cmd.DoneFlag {
			var scoreChange int
			for i := 0; i < len(g.Tasks); i++ {
				if g.Tasks[i].Done {
					scoreChange += g.Tasks[i].CountScoreChange(status)
				}
			}
			if g.Done {
				g.Done = false
				scoreChange = -scoreChange
				if g.Active {
					resources.Abr.AddToActivePorts(g.Name, resources.ANY_CMD_CYAN, resources.DB_DEFAULT_GOALS_BUCKET_NAME, g.Id, tr)
				}
			} else {
				g.Done = true
				if g.Active {
					resources.Abr.RemoveAndQuit(resources.DB_DEFAULT_GOALS_BUCKET_NAME, g.Id, tr)
				}
			}
			err = tr.ModifyEntity(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_ACTUAL_STATUS_KEY, true, status, func() {
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
			err := tr.ModifyEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, []byte(g.Tasks[i].Id), true, task, func() {
				task.Scheduled = scheduledSet
			})
			if err != nil {
				panic(err)
			}
		}
	}
}

func AddGoal(cmd *resources.Command) string {
	goal := resources.NewGoal(cmd.Name)
	if cmd.ProjectID != "" && cmd.ProjectID != "-" {
		goal.Project = &resources.Project{Id: cmd.ProjectID}
	}
	if cmd.HabitID != "" && cmd.HabitID != "-" {
		goal.Habit = &resources.Habit{Id: cmd.HabitID}
	}
	if cmd.HabitRepetitionGoal != -1 {
		goal.HabitRepetitionGoal = cmd.HabitRepetitionGoal
	}
	if cmd.BasePoints != -1 {
		goal.Priority = cmd.BasePoints
	}
	tr := db.NewTransaction()
	tr.Add(func() error {
		return tr.AddEntity(resources.DB_DEFAULT_GOALS_BUCKET_NAME, goal)
	})
	if cmd.ProjectID != "" && cmd.ProjectID != "-" {
		tr.Add(func() error {
			project := &resources.Project{}
			err := tr.ModifyEntity(resources.DB_DEFAULT_PROJECTS_BUCKET_NAME, []byte(cmd.ProjectID), true, project, func() {
				project.Goals = append(project.Goals, goal)
			})
			return err
		})
	}
	if cmd.HabitID != "" && cmd.HabitID != "-" {
		tr.Add(func() error {
			habit := &resources.Habit{}
			err := tr.ModifyEntity(resources.DB_DEFAULT_HABITS_BUCKET_NAME, []byte(cmd.HabitID), true, habit, func() {
				habit.Goal = goal
				habit.BasePoints = goal.Priority
			})
			return err
		})
	}
	tr.Execute()
	return goal.Id
}

func DeleteGoal(goalId string) {
	tr := db.NewTransaction()
	tr.Add(func() error {
		goal := &resources.Goal{}
		err := tr.RetrieveEntity(resources.DB_DEFAULT_GOALS_BUCKET_NAME, []byte(goalId), goal, true)
		if err != nil {
			return err
		}
		for i := 0; i < len(goal.Tasks); i++ {
			task := &resources.Task{}
			err = tr.ModifyEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, []byte(goal.Tasks[i].Id), true, task, func() {
				if task.Goal.Id == goal.Id {
					task.Goal = nil
				}
			})
			if err != nil {
				return err
			}
		}
		if goal.Project != nil {
			project := &resources.Project{}
			err = tr.ModifyEntity(resources.DB_DEFAULT_PROJECTS_BUCKET_NAME, []byte(goal.Project.Id), true, project, func() {
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
			err := tr.ModifyEntity(resources.DB_DEFAULT_HABITS_BUCKET_NAME, []byte(goal.Habit.Id), true, habit, func() {
				habit.Goal = nil
			})
			if err != nil {
				return err
			}
		}
		if goal.Active {
			resources.Abr.RemoveAndQuit(resources.DB_DEFAULT_GOALS_BUCKET_NAME, goalId, tr)
		}
		err = tr.DeleteEntity(resources.DB_DEFAULT_GOALS_BUCKET_NAME, []byte(goalId))
		if err != nil {
			return err
		}
		return nil
	})
	tr.Execute()
}

func ModifyGoal(cmd *resources.Command) {
	goal := &resources.Goal{}
	tr := db.NewTransaction()
	tr.Add(func() error {
		return tr.ModifyEntity(resources.DB_DEFAULT_GOALS_BUCKET_NAME, []byte(cmd.ID), false, goal, getModifyGoalFunc(goal, cmd, tr))
	})
	tr.Execute()
}

func GetGoal(goalId string) *resources.Goal {
	goal := &resources.Goal{}
	tr := db.NewTransaction()
	tr.Add(func() error {
		return tr.RetrieveEntity(resources.DB_DEFAULT_GOALS_BUCKET_NAME, []byte(goalId), goal, false)
	})
	tr.Execute()
	return goal
}

func GetGoals() map[string]*resources.Goal {
	goals := map[string]*resources.Goal{}
	db.RetrieveEntities(resources.DB_DEFAULT_GOALS_BUCKET_NAME, false, func(id string) resources.Entity {
		goals[id] = &resources.Goal{}
		return goals[id]
	})
	return goals
}

func GetActiveGoals() map[string]*resources.Goal {
	return FilterGoals(false, func(g *resources.Goal) bool { return g.Active && !g.Done })
}

func GetNonActiveGoals() map[string]*resources.Goal {
	return FilterGoals(false, func(g *resources.Goal) bool { return !g.Active && !g.Done })
}

func GetIncompleteGoals() map[string]*resources.Goal {
	return FilterGoals(false, func(g *resources.Goal) bool { return !g.Done })
}

func GetDoneGoals() map[string]*resources.Goal {
	return FilterGoals(false, func(g *resources.Goal) bool { return g.Done })
}

func GetProjectGoals(id string) map[string]*resources.Goal {
	return FilterGoals(false, func(g *resources.Goal) bool { return !g.Done && g.Project != nil && g.Project.Id == id })
}

func FilterGoals(shallow bool, filter func(*resources.Goal) bool) map[string]*resources.Goal {
	goals := map[string]*resources.Goal{}
	var goal *resources.Goal
	getNewEntity := func() resources.Entity {
		goal = &resources.Goal{}
		return goal
	}
	addEntity := func() { goals[goal.Id] = goal }
	db.FilterEntities(resources.DB_DEFAULT_GOALS_BUCKET_NAME, shallow, addEntity, getNewEntity, func() bool { return filter(goal) })
	return goals
}
