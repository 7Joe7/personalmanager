package resources

import (
	"fmt"
	"time"
)

type PointStat struct {
	Value int
	Id    time.Time
}

func (ps *PointStat) GetAlfredItem(id string) *AlfredItem {
	return &AlfredItem{
		Name:     fmt.Sprintf("%s - %v", ps.Id.Format("2006-01-02"), ps.Value),
		Arg:      id,
		Subtitle: "",
		Icon:     NewAlfredIcon(ICO_GREEN),
		Valid:    true,
		entity:   ps}
}

func (ps *PointStat) SetId(id string) {
	var err error
	ps.Id, err = time.Parse("2006-01-02", id)
	if err != nil {
		panic(err)
	}
}

func (ps *PointStat) GetId() string {
	return ps.Id.Format("2006-01-02")
}

func (ps *PointStat) Load(tr Transaction) error {
	return nil
}

func (ps *PointStat) Less(entity Entity) bool {
	t, err := time.Parse("2006-01-02", entity.GetId())
	if err != nil {
		panic(err)
	}
	return ps.Id.After(t)
}
