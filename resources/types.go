package resources

import (
	"time"
)

type Items struct {
	Items []*AlfredItem `json:"items"`
}

type PlannedItems struct {
	PlannedItems map[string]PlannedItem
	NoneAllowed  bool
	Status       *Status
	Sum          bool
}

type Tasks struct {
	Tasks       map[string]*Task
	NoneAllowed bool
	Status      *Status
	Sum         bool
}

type Projects struct {
	Projects    map[string]*Project
	NoneAllowed bool
	Status      *Status
}

type Tags struct {
	Tags        map[string]*Tag
	NoneAllowed bool
	Status      *Status
}

type Goals struct {
	Goals       map[string]*Goal
	NoneAllowed bool
	Status      *Status
}

type Habits struct {
	Habits      map[string]*Habit
	NoneAllowed bool
	Status      *Status
	Overview    bool
}

type Mod struct {
	Valid    bool   `json:"valid"`
	Arg      string `json:"arg,omitempty"`
	Subtitle string `json:"subtitle"`
}

type Mods struct {
	Ctrl  *Mod `json:"ctrl,omitempty"`
	Alt   *Mod `json:"alt,omitempty"`
	Cmd   *Mod `json:"cmd,omitempty"`
	Fn    *Mod `json:"Fn,omitempty"`
	Shift *Mod `json:"Shift,omitempty"`
}

type AlfredItem struct {
	Name     string      `json:"title"`
	Arg      string      `json:"arg,omitempty"`
	Subtitle string      `json:"subtitle,omitempty"`
	Valid    bool        `json:"valid"`
	Icon     *AlfredIcon `json:"icon,omitempty"`
	Mods     *Mods       `json:"mods,omitempty"`
	entity   Entity
}

type alfredItems []*AlfredItem

type AlfredIcon struct {
	Type string `json:"type,omitempty"`
	Path string `json:"path"`
}

type Tag struct {
	Name string
	Id   string
}

func (t *Tag) SetId(id string) {
	t.Id = id
}

func (t *Tag) GetId() string {
	return t.Id
}

func (t *Tag) Load(tr Transaction) error {
	return nil
}

func (t *Tag) Less(entity Entity) bool {
	return true
}

type Status struct {
	Score         int
	Yesterday     int
	Today         int
	WorkDoneToday int `json:"-"`
}

func (s *Status) SetId(id string) {}

func (s *Status) GetId() string {
	return ""
}

func (s *Status) Load(tr Transaction) error {
	return nil
}

func (s *Status) Less(entity Entity) bool {
	return true
}

type ActivePorts []*ActivePort

func (ap ActivePorts) Len() int           { return len(ap) }
func (ap ActivePorts) Swap(i, j int)      { ap[i], ap[j] = ap[j], ap[i] }
func (ap ActivePorts) Less(i, j int) bool { return ap[i].Port < ap[j].Port }

type ActivePort struct {
	Port       int
	Name       string
	Colour     string
	BucketName []byte
	Id         string
}

type Review struct {
	Deadline   *time.Time
	Repetition string
}

func (r *Review) SetId(id string) {}

func (r *Review) GetId() string {
	return ""
}

func (r *Review) Load(tr Transaction) error {
	return nil
}

func (r *Review) Less(entity Entity) bool {
	return true
}
