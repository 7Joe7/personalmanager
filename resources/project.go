package resources

import (
	"fmt"
	"encoding/json"
)

type Project struct {
	Name string `json:",omitempty"`
	Note string `json:",omitempty"`
	Id   string `json:",omitempty"`
	Active bool `json:",omitempty"`
	Done   bool `json:",omitempty"`
	Tasks []*Task
}

func (p *Project) SetId(id string) {
	p.Id = id
}

func (p *Project) GetId() string {
	return p.Id
}

func (p *Project) Load(tr Transaction) error {
	return nil
}

func (p *Project) MarshalJSON() ([]byte, error) {
	type mProject Project
	if p.Tasks != nil && len(p.Tasks) != 0 {
		for i := 0; i < len(p.Tasks); i++ {
			p.Tasks[i] = &Task{Id: p.Tasks[i].Id}
		}
	}
	return json.Marshal(mProject(*p))
}

func (p *Project) getItem(id string) *AlfredItem {
	var iconPath string
	var order int
	switch {
	case p.Done:
		order = 7500
		iconPath = ICO_GREEN
	case p.Active:
		order = 100
		iconPath = ICO_CYAN
	default:
		order = 5000
		iconPath = ICO_BLACK
	}
	var doneTasksNumber int
	for i := 0; i < len(p.Tasks); i++ {
		if p.Tasks[i].Done {
			doneTasksNumber++
		}
	}
	subtitle := fmt.Sprintf("%d/%d", doneTasksNumber, len(p.Tasks))
	return &AlfredItem{
		Name:     p.Name,
		Arg:      id,
		Subtitle: subtitle,
		Icon:     NewAlfredIcon(iconPath),
		Valid:    true,
		order:    order}
}
