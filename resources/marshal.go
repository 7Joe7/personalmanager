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
		subtitle = fmt.Sprintf("%s, %d/%d, actual %d, %v, base points %d", h.Repetition, h.Tries, h.Successes, h.ActualStreak, h.Deadline.Format("2.1.2006 15:04"), h.BasePoints)
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

func (s *Status) getItem() *AlfredItem {
	return &AlfredItem{
		Name: fmt.Sprintf("Today %d, total %d.", s.Today, s.Score)}
}

func (ts Tasks) MarshalJSON() ([]byte, error) {
	items := []*AlfredItem{}
	for id, t := range ts.Tasks {
		items = append(items, t.getItem(id))
	}
	if zeroItem := getZeroItem(ts.NoneAllowed, len(items) == 0, "tasks"); zeroItem != nil {
		items = append(items, zeroItem)
	}
	return marshalItems(items)
}

func (ps Projects) MarshalJSON() ([]byte, error) {
	items := []*AlfredItem{}
	for id, p := range ps.Projects {
		items = append(items, p.getItem(id))
	}
	if zeroItem := getZeroItem(ps.NoneAllowed, len(items) == 0, "projects"); zeroItem != nil {
		items = append(items, zeroItem)
	}
	return marshalItems(items)
}

func (ts Tags) MarshalJSON() ([]byte, error) {
	items := []*AlfredItem{}
	for id, t := range ts.Tags {
		items = append(items, t.getItem(id))
	}
	if zeroItem := getZeroItem(ts.NoneAllowed, len(items) == 0, "tags"); zeroItem != nil {
		items = append(items, zeroItem)
	}
	return marshalItems(items)
}

func (gs Goals) MarshalJSON() ([]byte, error) {
	items := []*AlfredItem{}
	for id, g := range gs.Goals {
		items = append(items, g.getItem(id))
	}
	if zeroItem := getZeroItem(gs.NoneAllowed, len(items) == 0, "goals"); zeroItem != nil {
		items = append(items, zeroItem)
	}
	return marshalItems(items)
}

func (hs Habits) MarshalJSON() ([]byte, error) {
	items := []*AlfredItem{}
	var zeroCount int
	if hs.Status != nil {
		items = append(items, hs.Status.getItem())
		zeroCount = 1
	}
	for id, h := range hs.Habits {
		items = append(items, h.getItem(id))
	}
	if zeroItem := getZeroItem(hs.NoneAllowed, len(items) == zeroCount, "habits"); zeroItem != nil {
		items = append(items, zeroItem)
	}
	return marshalItems(items)
}

func getZeroItem(noneAllowed, empty bool, elementType string, ) *AlfredItem {
	if noneAllowed {
		return &AlfredItem{
			Name:  "None.",
			Arg:   "-1",
			Valid: true}
	} else if empty {
		return &AlfredItem{
			Name:  fmt.Sprintf("There are no %s.", elementType),
			Valid: false,
			Mods:  getEmptyMods()}
	}
	return nil
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