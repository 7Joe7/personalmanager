package resources

import (
	"time"
)

type Items struct {
	Items []*AlfredItem `json:"items"`
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

type Project struct {
	Name string `json:",omitempty"`
	Note string `json:",omitempty"`
	Id   string `json:",omitempty"`
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

type Goal struct {
	Name     string
	Deadline *time.Time
	Id       string
}

func (g *Goal) SetId(id string) {
	g.Id = id
}

func (g *Goal) GetId() string {
	return g.Id
}

func (g *Goal) Load(tr Transaction) error {
	return nil
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

func (h *Habit) SetId(id string) {
	h.Id = id
}

func (h *Habit) GetId() string {
	return h.Id
}

func (h *Habit) Load(tr Transaction) error {
	return nil
}

type Status struct {
	Score int
	Today int
}

func (s *Status) SetId(id string) {}

func (s *Status) GetId() string {
	return ""
}

func (s *Status) Load(tr Transaction) error {
	return nil
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
