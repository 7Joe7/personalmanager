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

func (t *Task) MarshalJSON() ([]byte, error) {
	type mTask Task
	if t.Project != nil {
		t.Project = &Project{Id:t.Project.Id}
	}
	return json.Marshal(mTask(*t))
}

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
		Icon:     NewAlfredIcon(""),
		Valid:    true}
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

func (h *Habit) getItem(id string) *AlfredItem {
	var subtitle string
	var icon *AlfredIcon
	var order int
	if h.Active {
		if h.Done {
			icon = NewAlfredIcon(ICO_GREEN)
			order = HBT_DONE_BASE_ORDER
		} else {
			icon = NewAlfredIcon(ICO_RED)
			order = HBT_BASE_ORDER
		}
		subtitle = fmt.Sprintf(SUB_FORMAT_ACTIVE_HABIT, h.Repetition, h.Successes, h.Tries, h.ActualStreak,
			h.Deadline.Format(DEADLINE_FORMAT), h.BasePoints)
	} else {
		icon = NewAlfredIcon(ICO_BLACK)
		order = HBT_BASE_ORDER
		subtitle = fmt.Sprintf(SUB_FORMAT_NON_ACTIVE_HABIT, h.Successes, h.Tries)
	}
	order -= h.BasePoints
	return &AlfredItem{
		Name:     h.Name,
		Arg:      id,
		Subtitle: subtitle,
		Icon:     icon,
		Valid:    true,
		order:    order}
}

func (s *Status) getItem() *AlfredItem {
	return &AlfredItem{Name: fmt.Sprintf(NAME_FORMAT_STATUS, s.Score, s.Today), Icon:NewAlfredIcon(ICO_BLACK)}
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
	if zeroItem := getZeroItem(hs.NoneAllowed, len(items) == zeroCount, "habits"); zeroItem != nil {
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
			Name:  "None",
			Arg:   "-1",
			Icon:  NewAlfredIcon(ICO_BLACK),
			Valid: true}
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