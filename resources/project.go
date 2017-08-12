package resources

import (
	"encoding/json"
	"fmt"
)

type Project struct {
	Name   string `json:",omitempty"`
	Note   string `json:",omitempty"`
	Id     string `json:",omitempty"`
	Active bool   `json:",omitempty"`
	Done   bool   `json:",omitempty"`
	Tasks  []*Task
	Goals  []*Goal
}

func (p *Project) SetId(id string) {
	p.Id = id
}

func (p *Project) GetId() string {
	return p.Id
}

func (p *Project) Load(tr Transaction) error {
	tasks := []*Task{}
	var task *Task
	getNewEntity := func() Entity {
		task = &Task{}
		return task
	}
	addEntity := func() { tasks = append(tasks, task) }
	err := tr.FilterEntities(DB_DEFAULT_TASKS_BUCKET_NAME, true, addEntity, getNewEntity, func() bool { return task.Project != nil && task.Project.Id == p.Id })
	if err != nil {
		return err
	}
	p.Tasks = tasks
	goals := []*Goal{}
	var goal *Goal
	getNewEntity = func() Entity {
		goal = &Goal{}
		return goal
	}
	addEntity = func() { goals = append(goals, goal) }
	err = tr.FilterEntities(DB_DEFAULT_GOALS_BUCKET_NAME, true, addEntity, getNewEntity, func() bool { return goal.Project != nil && goal.Project.Id == p.Id })
	if err != nil {
		return err
	}
	p.Goals = goals
	return nil
}

func (p *Project) Less(entity Entity) bool {
	otherProject := entity.(*Project)
	if p.Done != otherProject.Done {
		return !p.Done
	}
	if p.Active != otherProject.Active {
		return p.Active
	}
	return true
}

func (p *Project) MarshalJSON() ([]byte, error) {
	type mProject Project
	if p.Tasks != nil && len(p.Tasks) != 0 {
		for i := 0; i < len(p.Tasks); i++ {
			p.Tasks[i] = &Task{Id: p.Tasks[i].Id}
		}
	}
	if p.Goals != nil && len(p.Goals) != 0 {
		for i := 0; i < len(p.Goals); i++ {
			p.Goals[i] = &Goal{Id: p.Goals[i].Id}
		}
	}
	return json.Marshal(mProject(*p))
}

func (p *Project) getItem(id string) *AlfredItem {
	var iconPath string
	switch {
	case p.Done:
		iconPath = ICO_GREEN
	case p.Active:
		iconPath = ICO_CYAN
	default:
		iconPath = ICO_BLACK
	}
	var doneTasksNumber int
	for i := 0; i < len(p.Tasks); i++ {
		if p.Tasks[i].Done {
			doneTasksNumber++
		}
	}
	var doneGoalsNumber int
	for i := 0; i < len(p.Goals); i++ {
		if p.Goals[i].Done {
			doneGoalsNumber++
		}
	}
	subtitle := fmt.Sprintf("%d/%d tasks, %d/%d goals", doneTasksNumber, len(p.Tasks), doneGoalsNumber, len(p.Goals))
	return &AlfredItem{
		Name:     p.Name,
		Arg:      id,
		Subtitle: subtitle,
		Icon:     NewAlfredIcon(iconPath),
		Valid:    true,
		entity:   p}
}
