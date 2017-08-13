package resources

import (
	"encoding/json"
	"fmt"
)

type Goal struct {
	Name                string   `json:",omitempty"`
	Id                  string   `json:",omitempty"`
	Active              bool     `json:",omitempty"`
	Project             *Project `json:",omitempty"`
	Priority            int      `json:",omitempty"`
	Tasks               []*Task  `json:",omitempty"`
	Habit               *Habit   `json:",omitempty"`
	HabitRepetitionGoal int      `json:",omitempty"`
	Done                bool     `json:",omitempty"`
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
	getNewEntity := func() Entity {
		task = &Task{}
		return task
	}
	addEntity := func() { tasks = append(tasks, task) }
	err := tr.FilterEntities(DB_DEFAULT_TASKS_BUCKET_NAME, true, addEntity, getNewEntity, func() bool { return task.Goal != nil && task.Goal.Id == g.Id })
	if err != nil {
		return err
	}
	g.Tasks = tasks
	if g.Habit != nil {
		err = tr.RetrieveEntity(DB_DEFAULT_HABITS_BUCKET_NAME, []byte(g.Habit.Id), g.Habit, true)
		if err != nil {
			return err
		}
	}
	if g.Project != nil {
		err = tr.RetrieveEntity(DB_DEFAULT_PROJECTS_BUCKET_NAME, []byte(g.Project.Id), g.Project, true)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Goal) Less(entity Entity) bool {
	otherGoal := entity.(*Goal)
	if g.Done != otherGoal.Done {
		return otherGoal.Done
	}
	if g.Active != otherGoal.Active {
		return g.Active
	}
	if g.Priority != otherGoal.Priority {
		return g.Priority > otherGoal.Priority
	}
	return g.Name < otherGoal.Name
}

func (g *Goal) MarshalJSON() ([]byte, error) {
	type mGoal Goal
	if g.Tasks != nil && len(g.Tasks) != 0 {
		for i := 0; i < len(g.Tasks); i++ {
			g.Tasks[i] = &Task{Id: g.Tasks[i].Id}
		}
	}
	if g.Project != nil {
		g.Project = &Project{Id: g.Project.Id}
	}
	if g.Habit != nil {
		g.Habit = &Habit{Id: g.Habit.Id}
	}
	return json.Marshal(mGoal(*g))
}

func (g *Goal) getItem(id string) *AlfredItem {
	var iconPath string
	switch {
	case g.Done:
		iconPath = ICO_GREEN
	case g.Active:
		iconPath = ICO_CYAN
	default:
		iconPath = ICO_BLACK
	}
	var doneNumber int
	for i := 0; i < len(g.Tasks); i++ {
		if g.Tasks[i].Done {
			doneNumber++
		}
	}
	if g.Habit != nil {
		if g.Habit.ActualStreak > 0 {
			doneNumber += g.Habit.ActualStreak
		}
	}
	subtitle := fmt.Sprintf("Priority %d, %d/%d", g.Priority, doneNumber, len(g.Tasks)+g.HabitRepetitionGoal)
	if g.Habit != nil {
		subtitle = fmt.Sprintf("%s, %s", g.Habit.Name, subtitle)
	}
	if g.Project != nil {
		subtitle = fmt.Sprintf("%s, %s", g.Project.Name, subtitle)
	}
	return &AlfredItem{
		Name:     g.Name,
		Arg:      id,
		Subtitle: subtitle,
		Icon:     NewAlfredIcon(iconPath),
		Valid:    true,
		entity:   g}
}
