package resources

import (
	"fmt"
	"encoding/json"
)

type Goal struct {
	Name     string
	Id       string
	Active   bool
	Tasks    []*Task
	Done     bool
}

func (g *Goal) SetId(id string) {
	g.Id = id
}

func (g *Goal) GetId() string {
	return g.Id
}

func (g *Goal) Load(tr Transaction) error {
	tasks := []*Task{}
	var task *Task
	getNewEntity := func () Entity {
		task = &Task{}
		return task
	}
	addEntity := func () { tasks = append(tasks, task) }
	err := tr.FilterEntities(DB_DEFAULT_TASKS_BUCKET_NAME, true, addEntity, getNewEntity, func () bool { return task.Goal != nil && task.Goal.Id == g.Id })
	if err != nil {
		return err
	}
	g.Tasks = tasks
	return nil
}

func (g *Goal) MarshalJSON() ([]byte, error) {
	type mGoal Goal
	if g.Tasks != nil && len(g.Tasks) != 0 {
		for i := 0; i < len(g.Tasks); i++ {
			g.Tasks[i] = &Task{Id: g.Tasks[i].Id}
		}
	}
	return json.Marshal(mGoal(*g))
}

func (g *Goal) getItem(id string) *AlfredItem {
	var iconPath string
	var order int
	switch {
	case g.Done:
		order = 7500
		iconPath = ICO_GREEN
	case g.Active:
		order = 100
		iconPath = ICO_CYAN
	default:
		order = 5000
		iconPath = ICO_BLACK
	}
	var doneTasksNumber int
	for i := 0; i < len(g.Tasks); i++ {
		if g.Tasks[i].Done {
			doneTasksNumber++
		}
	}
	subtitle := fmt.Sprintf("%d/%d", doneTasksNumber, len(g.Tasks))
	return &AlfredItem{
		Name:     g.Name,
		Arg:      id,
		Subtitle: subtitle,
		Icon:     NewAlfredIcon(iconPath),
		Valid:    true,
		order:    order}
}