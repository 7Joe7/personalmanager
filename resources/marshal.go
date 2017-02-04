package resources

import (
	"encoding/json"
	"fmt"
	"sort"
	"github.com/7joe7/personalmanager/utils"
)

func (r *Review) GetItem() *AlfredItem {
	return &AlfredItem{
		Name:  fmt.Sprintf("Review repeated %s, next: %s.", r.Repetition, r.Deadline.Format(DATE_FORMAT)),
		Icon:  NewAlfredIcon(ICO_BLACK),
		Valid: true}
}

func (t *Tag) getItem(id string) *AlfredItem {
	return &AlfredItem{
		Name:     t.Name,
		Arg:      id,
		Subtitle: "",
		Icon:     NewAlfredIcon(""),
		Valid:    true}
}

func (s *Status) getItem() *AlfredItem {
	return &AlfredItem{
		Name:  fmt.Sprintf(NAME_FORMAT_STATUS, s.Score, s.Today, s.Yesterday),
		Valid: false,
		Icon:  NewAlfredIcon(ICO_HABIT),
		Mods:  getEmptyMods(),
		order: 0}
}

func (ts Tasks) MarshalJSON() ([]byte, error) {
	items := items{}
	var zeroCount int
	if ts.Status != nil {
		items = append(items, ts.Status.getItem())
		zeroCount = -1
	}
	if ts.Sum {
		var sum float64
		for _, t := range ts.Tasks {
			if !t.Done && t.TimeEstimate != nil {
				sum += t.TimeEstimate.Minutes()
			}
		}
		items = append(items, &AlfredItem{
			Name: fmt.Sprintf("Total estimate: %s", utils.MinutesToHMFormat(sum)),
			Valid: false,
			Icon: NewAlfredIcon(ICO_BLACK),
			Mods: getEmptyMods(),
			order: 1})
	}
	for id, t := range ts.Tasks {
		items = append(items, t.getItem(id))
	}
	if zeroItem := getZeroItem(ts.NoneAllowed, len(items) == zeroCount, "task"); zeroItem != nil {
		items = append(items, zeroItem)
	}
	sort.Sort(items)
	return marshalItems(items)
}

func (ps Projects) MarshalJSON() ([]byte, error) {
	items := items{}
	var zeroCount int
	if ps.Status != nil {
		items = append(items, ps.Status.getItem())
		zeroCount = 1
	}
	for id, p := range ps.Projects {
		items = append(items, p.getItem(id))
	}
	if zeroItem := getZeroItem(ps.NoneAllowed, len(items) == zeroCount, "project"); zeroItem != nil {
		items = append(items, zeroItem)
	}
	sort.Sort(items)
	return marshalItems(items)
}

func (ts Tags) MarshalJSON() ([]byte, error) {
	items := items{}
	var zeroCount int
	if ts.Status != nil {
		items = append(items, ts.Status.getItem())
		zeroCount = 1
	}
	for id, t := range ts.Tags {
		items = append(items, t.getItem(id))
	}
	if zeroItem := getZeroItem(ts.NoneAllowed, len(items) == zeroCount, "tag"); zeroItem != nil {
		items = append(items, zeroItem)
	}
	sort.Sort(items)
	return marshalItems(items)
}

func (gs Goals) MarshalJSON() ([]byte, error) {
	items := items{}
	var zeroCount int
	if gs.Status != nil {
		items = append(items, gs.Status.getItem())
		zeroCount = 1
	}
	for id, g := range gs.Goals {
		items = append(items, g.getItem(id))
	}
	if zeroItem := getZeroItem(gs.NoneAllowed, len(items) == zeroCount, "goal"); zeroItem != nil {
		items = append(items, zeroItem)
	}
	sort.Sort(items)
	return marshalItems(items)
}

func (hs Habits) MarshalJSON() ([]byte, error) {
	items := items{}
	var zeroCount int
	if hs.Status != nil {
		items = append(items, hs.Status.getItem())
		zeroCount = 1
	}
	if hs.Overview {
		var dailyCount, dailyCountDone, weeklyCount, weeklyCountDone, monthlyCount, monthlyCountDone int
		for _, h := range hs.Habits {
			if !h.Negative {
				switch h.Repetition {
				case HBT_REPETITION_DAILY:
					if h.Done {
						dailyCountDone++
					}
					dailyCount++
				case HBT_REPETITION_WEEKLY:
					if h.Done {
						weeklyCountDone++
					}
					weeklyCount++
				case HBT_REPETITION_MONTHLY:
					if h.Done {
						monthlyCountDone++
					}
					monthlyCount++
				}
			}
		}
		items = append(items, &AlfredItem{
			Name: fmt.Sprintf("D: %d/%d, W: %d/%d, M: %d/%d",
				dailyCountDone, dailyCount,
				weeklyCountDone, weeklyCount,
				monthlyCountDone, monthlyCount),
			Valid: false,
			Icon: NewAlfredIcon(ICO_BLACK),
			Mods: getEmptyMods(),
			order: 1,
		})
	}
	for id, h := range hs.Habits {
		items = append(items, h.getItem(id))
	}
	if zeroItem := getZeroItem(hs.NoneAllowed, len(items) == zeroCount, "habit"); zeroItem != nil {
		items = append(items, zeroItem)
	}
	sort.Sort(items)
	return marshalItems(items)
}

func (ho items) Len() int      { return len(ho) }
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
			Mods:  getEmptyMods(),
			order: 10000}
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
