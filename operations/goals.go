package operations

import (
	"github.com/7joe7/personalmanager/resources"
	"github.com/7joe7/personalmanager/db"
)

func getModifyGoalFunc(g *resources.Goal, name string) func () {
	return func () {
		if name != "" {
			g.Name = name
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
	db.DeleteEntity([]byte(goalId), resources.DB_DEFAULT_GOALS_BUCKET_NAME)
}

func ModifyGoal(goalId, name string) {
	goal := &resources.Goal{}
	db.ModifyEntity(resources.DB_DEFAULT_GOALS_BUCKET_NAME, []byte(goalId), goal, getModifyGoalFunc(goal, name))
}

func GetGoal(goalId string) *resources.Goal {
	goal := &resources.Goal{}
	tr := db.NewTransaction()
	tr.Add(func () error {
		return tr.RetrieveEntity(resources.DB_DEFAULT_GOALS_BUCKET_NAME, []byte(goalId), goal)
	})
	tr.Execute()
	return goal
}

func GetGoals() map[string]*resources.Goal {
	goals := map[string]*resources.Goal{}
	db.RetrieveEntities(resources.DB_DEFAULT_GOALS_BUCKET_NAME, func (id string) resources.Entity {
		goals[id] = &resources.Goal{}
		return goals[id]
	})
	return goals
}
