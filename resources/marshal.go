package resources

import (
	"encoding/json"
	"fmt"
	"sort"
)

func (r *Review) GetItem() *AlfredItem {
	return &AlfredItem{
		Name: fmt.Sprintf("Review repeated %s, next: %s.", r.Repetition, r.Deadline.Format(DATE_FORMAT)),
		Icon: NewAlfredIcon(ICO_BLACK),
		Valid: true}
}

func (p *Project) getItem(id string) *AlfredItem {
	return &AlfredItem{
		Name:     p.Name,
		Arg:      id,
		Subtitle: p.Note,
		Icon:     NewAlfredIcon(""),
		Valid:    true}
}

func (t *Tag) getItem(id string) *AlfredItem {
	return &AlfredItem{
		Name:     t.Name,
		Arg:      id,
		Subtitle: "",
		Icon:     NewAlfredIcon(""),
		Valid:    true}
}

func (g *Goal) getItem(id string) *AlfredItem {
	return &AlfredItem{
		Name:     g.Name,
		Arg:      id,
		Subtitle: "",
		Icon:     NewAlfredIcon(""),
		Valid:    true}
}

func (s *Status) getItem() *AlfredItem {
	return &AlfredItem{
		Name: fmt.Sprintf(NAME_FORMAT_STATUS, s.Score, s.Today),
		Icon: NewAlfredIcon(ICO_HABIT),
		Mods:getEmptyMods()}
}

func (ts Tasks) MarshalJSON() ([]byte, error) {
	items := items{}
	var zeroCount int
	if ts.Status != nil {
		items = append(items, ts.Status.getItem())
		zeroCount = 1
	}
	for id, h := range ts.Tasks {
		items = append(items, h.getItem(id))
	}
	if zeroItem := getZeroItem(ts.NoneAllowed, len(items) == zeroCount, "task"); zeroItem != nil {
		items = append(items, zeroItem)
	}
	sort.Sort(items)
	return marshalItems(items)
}

func (ps Projects) MarshalJSON() ([]byte, error) {
	items := []*AlfredItem{}
	for id, p := range ps.Projects {
		items = append(items, p.getItem(id))
	}
	if zeroItem := getZeroItem(ps.NoneAllowed, len(items) == 0, "project"); zeroItem != nil {
		items = append(items, zeroItem)
	}
	return marshalItems(items)
}

func (ts Tags) MarshalJSON() ([]byte, error) {
	items := []*AlfredItem{}
	for id, t := range ts.Tags {
		items = append(items, t.getItem(id))
	}
	if zeroItem := getZeroItem(ts.NoneAllowed, len(items) == 0, "tag"); zeroItem != nil {
		items = append(items, zeroItem)
	}
	return marshalItems(items)
}

func (gs Goals) MarshalJSON() ([]byte, error) {
	items := []*AlfredItem{}
	for id, g := range gs.Goals {
		items = append(items, g.getItem(id))
	}
	if zeroItem := getZeroItem(gs.NoneAllowed, len(items) == 0, "goal"); zeroItem != nil {
		items = append(items, zeroItem)
	}
	return marshalItems(items)
}

func (hs Habits) MarshalJSON() ([]byte, error) {
	items := items{}
	var zeroCount int
	if hs.Status != nil {
		items = append(items, hs.Status.getItem())
		zeroCount = 1
	}
	for id, h := range hs.Habits {
		items = append(items, h.getItem(id))
	}
	sort.Sort(items)
	if zeroItem := getZeroItem(hs.NoneAllowed, len(items) == zeroCount, "habit"); zeroItem != nil {
		items = append(items, zeroItem)
	}
	return marshalItems(items)
}

func (ho items) Len() int { return len(ho) }
func (ho items) Swap(i, j int) { ho[i], ho[j] = ho[j], ho[i] }
func (ho items) Less(i, j int) bool {
	if ho[i].order == ho[j].order {
		return ho[i].Name < ho[j].Name
	}
	return ho[i].order < ho[j].order
}

func getZeroItem(noneAllowed, empty bool, elementType string) *AlfredItem {
	if noneAllowed {
		return &AlfredItem{
			Name:  "No " + elementType,
			Arg:   "",
			Icon:  NewAlfredIcon(ICO_BLACK),
			Valid: true,
			Mods:  getEmptyMods(),
			order: 10000}
	} else if empty {
		return &AlfredItem{
			Name:  fmt.Sprintf(NAME_FORMAT_EMPTY, elementType),
			Valid: false,
			Icon:  NewAlfredIcon(ICO_BLACK),
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