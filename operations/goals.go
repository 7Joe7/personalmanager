package operations

import (
	"github.com/7joe7/personalmanager/resources"
	"github.com/7joe7/personalmanager/db"
	"github.com/7joe7/personalmanager/anybar"
)

func getModifyGoalFunc(g *resources.Goal, name, taskId string, activeFlag, doneFlag bool, tr resources.Transaction) func () {
	return func () {
		if name != "" {
			g.Name = name
		}
		if taskId != "" {
			task := &resources.Task{}
			err := tr.ModifyEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, []byte(taskId), true, task, func() {
				task.Goal = g
				if g.Active {
					task.Scheduled = resources.TASK_SCHEDULED_NEXT
				}
			})
			if err != nil {
				panic(err)
			}
			g.Tasks = append(g.Tasks, task)
		}
		if activeFlag {
			if g.Active {
				toggleSubTasksScheduling(resources.TASK_SCHEDULED_NEXT, resources.TASK_NOT_SCHEDULED, g, tr)
				g.Active = false
				anybar.RemoveAndQuit(resources.DB_DEFAULT_GOALS_BUCKET_NAME, g.Id, tr)
			} else {
				toggleSubTasksScheduling(resources.TASK_NOT_SCHEDULED, resources.TASK_SCHEDULED_NEXT, g, tr)
				g.Active = true
				anybar.AddToActivePorts(g.Name, resources.ANY_CMD_YELLOW, resources.DB_DEFAULT_GOALS_BUCKET_NAME, g.Id, tr)
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
					anybar.AddToActivePorts(g.Name, resources.ANY_CMD_YELLOW, resources.DB_DEFAULT_GOALS_BUCKET_NAME, g.Id, tr)
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

func AddGoal(name string) string {
	goal := resources.NewGoal(name)
	tr := db.NewTransaction()
	tr.Add(func () error {
		return tr.AddEntity(resources.DB_DEFAULT_GOALS_BUCKET_NAME, goal)
	})
	tr.Execute()
	return goal.Id
}

func DeleteGoal(goalId string) {
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

func ModifyGoal(goalId, name, taskId string, activeFlag, doneFlag bool) {
	goal := &resources.Goal{}
	tr := db.NewTransaction()
	tr.Add(func () error {
		return tr.ModifyEntity(resources.DB_DEFAULT_GOALS_BUCKET_NAME, []byte(goalId), false, goal, getModifyGoalFunc(goal, name, taskId, activeFlag, doneFlag, tr))
	})
	tr.Execute()
}

func GetGoal(goalId string) *resources.Goal {
	goal := &resources.Goal{}
	tr := db.NewTransaction()
	tr.Add(func () error {
		return tr.RetrieveEntity(resources.DB_DEFAULT_GOALS_BUCKET_NAME, []byte(goalId), goal, false)
	})
	tr.Execute()
	return goal
}

func GetGoals() map[string]*resources.Goal {
	goals := map[string]*resources.Goal{}
	db.RetrieveEntities(resources.DB_DEFAULT_GOALS_BUCKET_NAME, false, func (id string) resources.Entity {
		goals[id] = &resources.Goal{}
		return goals[id]
	})
	return goals
}

func getActiveGoals() map[string]*resources.Goal {
	return FilterGoals(func (g *resources.Goal) bool { return g.Active })
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
