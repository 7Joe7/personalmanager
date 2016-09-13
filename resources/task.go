package resources

import (
	"encoding/json"
	"time"

	"github.com/7joe7/personalmanager/utils"
)

type Task struct {
	Id              string `json:",omitempty"`
	Name            string `json:",omitempty"`
	Note            string `json:",omitempty"`
	BasePoints      int    `json:",omitempty"`
	InProgress      bool   `json:",omitempty"`
	Done            bool   `json:",omitempty"`
	Tags            []*Tag `json:",omitempty"`
	InProgressSince *time.Time
	DoneTime        *time.Time
	Deadline        *time.Time
	TimeEstimate    *time.Duration
	TimeSpent       *time.Duration
	Project         *Project
}

func (t *Task) SetId(id string) {
	t.Id = id
}

func (t *Task) GetId() string {
	return t.Id
}

func (t *Task) Load(tr Transaction) error {
	if t.Project != nil {
		return tr.RetrieveEntity(DB_DEFAULT_PROJECTS_BUCKET_NAME, []byte(t.Project.Id), t.Project)
	}
	for i := 0; i < len(t.Tags); i++ {
		if err := tr.RetrieveEntity(DB_DEFAULT_TAGS_BUCKET_NAME, []byte(t.Tags[i].Id), t.Tags[i]); err != nil {
			return err
		}
	}
	return nil
}

func (t *Task) MarshalJSON() ([]byte, error) {
	type mTask Task
	if t.Project != nil {
		t.Project = &Project{Id: t.Project.Id}
	}
	return json.Marshal(mTask(*t))
}

func (t *Task) getItem(id string) *AlfredItem {
	var subtitle string
	var comma bool
	if t.Deadline != nil {
		subtitle = t.Deadline.Format(DATE_FORMAT)
		comma = true
	}

	if t.Project != nil {
		if comma {
			subtitle += "; "
		}
		comma = true
		subtitle = t.Project.Name
	}

	var todayTagPresent bool
	if len(t.Tags) > 0 {
		if comma {
			subtitle += "; "
		}
		comma = true
		subtitle += "Tags: "
	}
	for i := 0; i < len(t.Tags); i++ {
		if i > 0 {
			subtitle += ", "
		}
		subtitle += t.Tags[i].Name
		if t.Tags[i].Name == "TODAY" {
			todayTagPresent = true
		}
	}

	if comma {
		subtitle += "; "
	}
	subtitle += "Spent: "

	if t.TimeSpent == nil {
		subtitle += "?/"
	} else {
		subtitle += utils.DurationToHMFormat(t.TimeSpent) + "/"
	}
	if t.TimeEstimate == nil {
		subtitle += "?"
	} else {
		subtitle += utils.DurationToHMFormat(t.TimeEstimate)
	}

	var icoPath string
	if t.InProgress {
		icoPath = ICO_BLUE
	} else if t.Done {
		icoPath = ICO_GREEN
	} else if t.Deadline != nil && t.Deadline.Before(time.Now()) {
		icoPath = ICO_RED
	} else if todayTagPresent {
		icoPath = ICO_ORANGE
	} else {
		icoPath = ICO_BLACK
	}

	return &AlfredItem{
		Name:     t.Name,
		Arg:      id,
		Subtitle: subtitle,
		Icon:     NewAlfredIcon(icoPath),
		Valid:    true}
}
