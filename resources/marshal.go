package resources

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/7joe7/personalmanager/utils"
)

// MarshalJSON for PlannedItems
func (pi PlannedItems) MarshalJSON() ([]byte, error) {
	items := alfredItems{}
	var zeroCount int
	for id, t := range pi.PlannedItems {
		items = append(items, t.GetAlfredItem(id))
	}
	sort.Sort(items)
	if pi.Status != nil {
		items = append(alfredItems{pi.Status.GetAlfredItem()}, items...)
		zeroCount = -1
	}
	if pi.Sum {
		var sum float64
		for _, i := range pi.PlannedItems {
			if i.GetTimeEstimate() != nil {
				sum += i.GetTimeEstimate().Minutes()
			}
		}
		items = append(alfredItems{&AlfredItem{
			Name:  fmt.Sprintf("Count: %d, estimate: %s", len(pi.PlannedItems), utils.MinutesToHMFormat(sum)),
			Valid: false,
			Icon:  NewAlfredIcon(ICO_BLACK),
			Mods:  getEmptyMods()}}, items...)
	}
	if zeroItem := getZeroItem(pi.NoneAllowed, len(items) == zeroCount, "planned item"); zeroItem != nil {
		items = append(items, zeroItem)
	}
	return marshalItems(items)
}

func (r *Review) GetItem() *AlfredItem {
	return &AlfredItem{
		Name:  fmt.Sprintf("Review repeated %s, next: %s.", r.Repetition, r.Deadline.Format(DATE_FORMAT)),
		Icon:  NewAlfredIcon(ICO_BLACK),
		Valid: true}
}

func (t *Tag) GetAlfredItem(id string) *AlfredItem {
	return &AlfredItem{
		Name:     t.Name,
		Arg:      id,
		Subtitle: "",
		Icon:     NewAlfredIcon(""),
		Valid:    true}
}

func (s *Status) GetAlfredItem() *AlfredItem {
	return &AlfredItem{
		Name:  fmt.Sprintf(NAME_FORMAT_STATUS, s.Score, s.Today, s.Yesterday),
		Valid: false,
		Icon:  NewAlfredIcon(ICO_HABIT),
		Mods:  getEmptyMods()}
}

func (ts Tasks) MarshalJSON() ([]byte, error) {
	items := alfredItems{}
	var zeroCount int
	for id, t := range ts.Tasks {
		items = append(items, t.GetAlfredItem(id))
	}
	sort.Sort(items)
	if ts.Status != nil {
		items = append(alfredItems{ts.Status.GetAlfredItem()}, items...)
		zeroCount = -1
	}
	if ts.Sum {
		var sum float64
		for _, t := range ts.Tasks {
			if !t.Done && t.TimeEstimate != nil {
				sum += t.TimeEstimate.Minutes()
			}
		}
		items = append(alfredItems{&AlfredItem{
			Name:  fmt.Sprintf("Count: %d, estimate: %s", len(ts.Tasks), utils.MinutesToHMFormat(sum)),
			Valid: false,
			Icon:  NewAlfredIcon(ICO_BLACK),
			Mods:  getEmptyMods()}}, items...)
	}
	if zeroItem := getZeroItem(ts.NoneAllowed, len(items) == zeroCount, "task"); zeroItem != nil {
		items = append(items, zeroItem)
	}
	return marshalItems(items)
}

func (ps Projects) MarshalJSON() ([]byte, error) {
	items := alfredItems{}
	var zeroCount int
	for id, p := range ps.Projects {
		items = append(items, p.GetAlfredItem(id))
	}
	sort.Sort(items)
	if ps.Status != nil {
		items = append(alfredItems{ps.Status.GetAlfredItem()}, items...)
		zeroCount = 1
	}
	if zeroItem := getZeroItem(ps.NoneAllowed, len(items) == zeroCount, "project"); zeroItem != nil {
		items = append(items, zeroItem)
	}
	return marshalItems(items)
}

func (ts Tags) MarshalJSON() ([]byte, error) {
	items := alfredItems{}
	var zeroCount int
	for id, t := range ts.Tags {
		items = append(items, t.GetAlfredItem(id))
	}
	sort.Sort(items)
	if ts.Status != nil {
		items = append(alfredItems{ts.Status.GetAlfredItem()}, items...)
		zeroCount = 1
	}
	if zeroItem := getZeroItem(ts.NoneAllowed, len(items) == zeroCount, "tag"); zeroItem != nil {
		items = append(items, zeroItem)
	}
	return marshalItems(items)
}

func (gs Goals) MarshalJSON() ([]byte, error) {
	items := alfredItems{}
	var zeroCount int
	for id, g := range gs.Goals {
		items = append(items, g.GetAlfredItem(id))
	}
	sort.Sort(items)
	if gs.Status != nil {
		items = append(alfredItems{gs.Status.GetAlfredItem()}, items...)
		zeroCount = 1
	}
	if zeroItem := getZeroItem(gs.NoneAllowed, len(items) == zeroCount, "goal"); zeroItem != nil {
		items = append(items, zeroItem)
	}
	return marshalItems(items)
}

func (hs Habits) MarshalJSON() ([]byte, error) {
	items := alfredItems{}
	var zeroCount int
	for id, h := range hs.Habits {
		items = append(items, h.GetAlfredItem(id))
	}
	sort.Sort(items)
	if hs.Status != nil {
		items = append(alfredItems{hs.Status.GetAlfredItem()}, items...)
		zeroCount = 1
	}
	if hs.Overview {
		var dailyCount, dailyCountDone, weeklyCount, weeklyCountDone, monthlyCount, monthlyCountDone, yearlyCount, yearlyCountDone int
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
				case HBT_REPETITION_YEARLY:
					if h.Done {
						yearlyCountDone++
					}
					yearlyCount++
				}
			}
		}
		items = append(items, &AlfredItem{
			Name: fmt.Sprintf("D: %d/%d, W: %d/%d, M: %d/%d, Y: %d/%d",
				dailyCountDone, dailyCount,
				weeklyCountDone, weeklyCount,
				monthlyCountDone, monthlyCount,
				yearlyCountDone, yearlyCount),
			Valid: false,
			Icon:  NewAlfredIcon(ICO_BLACK),
			Mods:  getEmptyMods(),
		})
	}

	if zeroItem := getZeroItem(hs.NoneAllowed, len(items) == zeroCount, "habit"); zeroItem != nil {
		items = append(items, zeroItem)
	}
	return marshalItems(items)
}

func (ho alfredItems) Len() int           { return len(ho) }
func (ho alfredItems) Swap(i, j int)      { ho[i], ho[j] = ho[j], ho[i] }
func (ho alfredItems) Less(i, j int) bool { return ho[i] != nil && ho[i].entity.Less(ho[j].entity) }

func getZeroItem(noneAllowed, empty bool, elementType string) *AlfredItem {
	if noneAllowed {
		return &AlfredItem{
			Name:  "No " + elementType,
			Arg:   "-",
			Icon:  NewAlfredIcon(ICO_BLACK),
			Valid: true,
			Mods:  getEmptyMods()}
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
