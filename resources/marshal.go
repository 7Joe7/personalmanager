package resources

import (
	"encoding/json"
	"fmt"
)

func (t *Task) getItem(id string) *AlfredItem {
	var subtitle string
	if t.Project != nil {
		subtitle = t.Project.Name + " "
	}
	subtitle = subtitle + t.Note

	return &AlfredItem{
		Name:     t.Name,
		Arg:      id,
		Subtitle: subtitle,
		Valid:    true}
}

func (p *Project) getItem(id string) *AlfredItem {
	return &AlfredItem{
		Name:     p.Name,
		Arg:      id,
		Subtitle: p.Note,
		Valid:    true}
}

func (t *Tag) getItem(id string) *AlfredItem {
	return &AlfredItem{
		Name:     t.Name,
		Arg:      id,
		Subtitle: "",
		Valid:    true}
}

func (g *Goal) getItem(id string) *AlfredItem {
	return &AlfredItem{
		Name:     g.Name,
		Arg:      id,
		Subtitle: "",
		Valid:    true}
}

func (h *Habit) getItem(id string) *AlfredItem {
	var subtitle string
	var icon *AlfredIcon
	if h.Active {
		if h.Done {
			icon = NewAlfredIcon("green")
		} else {
			icon = NewAlfredIcon("red")
		}
		subtitle = fmt.Sprintf("%s, %d/%d, %v", h.Repetition, h.Tries, h.Successes, h.Deadline.Format("2.1.2006 15:04"))
	} else {
		icon = NewAlfredIcon("black")
		subtitle = fmt.Sprintf("%d/%d", h.Tries, h.Successes)
	}
	return &AlfredItem{
		Name:     h.Name,
		Arg:      id,
		Subtitle: subtitle,
		Icon:     icon,
		Valid:    true}
}

func (ts Tasks) MarshalJSON() ([]byte, error) {
	items := []*AlfredItem{}
	for id, t := range ts.Tasks {
		items = append(items, t.getItem(id))
	}
	return marshalItems(addZeroElement(items, ts.NoneAllowed, "tasks"))
}

func (ps Projects) MarshalJSON() ([]byte, error) {
	items := []*AlfredItem{}
	for id, p := range ps.Projects {
		items = append(items, p.getItem(id))
	}
	return marshalItems(addZeroElement(items, ps.NoneAllowed, "projects"))
}

func (ts Tags) MarshalJSON() ([]byte, error) {
	items := []*AlfredItem{}
	for id, t := range ts.Tags {
		items = append(items, t.getItem(id))
	}
	return marshalItems(addZeroElement(items, ts.NoneAllowed, "tags"))
}

func (gs Goals) MarshalJSON() ([]byte, error) {
	items := []*AlfredItem{}
	for id, g := range gs.Goals {
		items = append(items, g.getItem(id))
	}
	return marshalItems(addZeroElement(items, gs.NoneAllowed, "goals"))
}

func (hs Habits) MarshalJSON() ([]byte, error) {
	items := []*AlfredItem{}
	for id, h := range hs.Habits {
		items = append(items, h.getItem(id))
	}
	return marshalItems(addZeroElement(items, hs.NoneAllowed, "habits"))
}

func addZeroElement(items []*AlfredItem, noneAllowed bool, elementType string) []*AlfredItem {
	if noneAllowed {
		items = append(items, &AlfredItem{
			Name:  "None.",
			Arg:   "-1",
			Valid: true})
	} else if len(items) == 0 {
		items = append(items, &AlfredItem{
			Name:  fmt.Sprintf("There are no %s.", elementType),
			Valid: false,
			Mods:  getEmptyMods()})
	}
	return items
}

func marshalItems(items []*AlfredItem) ([]byte, error) {
	return json.Marshal(&struct {
		Items []*AlfredItem `json:"items"`
	}{
		Items: items,
	})
}

func getEmptyMods() *Mods {
	return &Mods{Ctrl: &Mod{}, Alt: &Mod{}, Cmd: &Mod{}, Fn: &Mod{}, Shift: &Mod{}}
}