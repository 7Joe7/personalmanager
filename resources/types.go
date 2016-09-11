package resources

import "time"

type Entity interface {
	SetId(string)
}

type Tasks struct {
	Tasks       map[string]*Task
	NoneAllowed bool
}

type Projects struct {
	Projects    map[string]*Project
	NoneAllowed bool
}

type Tags struct {
	Tags        map[string]*Tag
	NoneAllowed bool
}

type Goals struct {
	Goals       map[string]*Goal
	NoneAllowed bool
}

type Habits struct {
	Habits      map[string]*Habit
	NoneAllowed bool
	Status      *Status
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
	order    int
}

type items []*AlfredItem

type AlfredIcon struct {
	Type string `json:"type,omitempty"`
	Path string `json:"path"`
}

type Task struct {
	Name    string `json:"title"`
	Note    string
	Tags    []*Tag
	Project *Project
	Id      string
}

func (t *Task) SetId(id string) {
	t.Id = id
}

type Tag struct {
	Name string
	Id   string
}

func (t *Tag) SetId(id string) {
	t.Id = id
}

type Project struct {
	Name string
	Note string
	Id   string
}

func (t *Project) SetId(id string) {
	t.Id = id
}

type Goal struct {
	Name string
	Id   string
}

func (t *Goal) SetId(id string) {
	t.Id = id
}

type Habit struct {
	Name          string
	Active        bool
	Done          bool
	Deadline      *time.Time
	Tries         int
	Successes     int
	ActualStreak  int
	LastStreak    int
	LastStreakEnd *time.Time
	Repetition    string
	BasePoints    int
	Id            string
}

func (t *Habit) SetId(id string) {
	t.Id = id
}

type Status struct {
	Score int
	Today int
}
