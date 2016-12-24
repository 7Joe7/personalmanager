package resources

import (
	"encoding/json"
	"time"
	"strconv"

	"github.com/7joe7/personalmanager/utils"
	"fmt"
)

type Task struct {
	Id              string `json:",omitempty"`
	Name            string `json:",omitempty"`
	Note            string `json:",omitempty"`
	BasePoints      int    `json:",omitempty"`
	InProgress      bool   `json:",omitempty"`
	Done            bool   `json:",omitempty"`
	Tags            []*Tag `json:",omitempty"`
	Scheduled       string `json:",omitempty"`
	Type            string `json:",omitempty"`
	InProgressSince *time.Time
	DoneTime        *time.Time
	Deadline        *time.Time
	TimeEstimate    *time.Duration
	TimeSpent       *time.Duration
	Project         *Project
	Goal            *Goal
}

func (t *Task) SetId(id string) {
	t.Id = id
}

func (t *Task) GetId() string {
	return t.Id
}

func (t *Task) Load(tr Transaction) error {
	if t.Project != nil {
		err := tr.RetrieveEntity(DB_DEFAULT_PROJECTS_BUCKET_NAME, []byte(t.Project.Id), t.Project, true)
		if err != nil {
			return err
		}
	}
	if t.Goal != nil {
		err := tr.RetrieveEntity(DB_DEFAULT_GOALS_BUCKET_NAME, []byte(t.Goal.Id), t.Goal, true)
		if err != nil {
			return err
		}
	}
	for i := 0; i < len(t.Tags); i++ {
		if err := tr.RetrieveEntity(DB_DEFAULT_TAGS_BUCKET_NAME, []byte(t.Tags[i].Id), t.Tags[i], true); err != nil {
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
	if t.Goal != nil {
		t.Goal = &Goal{Id: t.Goal.Id}
	}
	return json.Marshal(mTask(*t))
}

func (t *Task) Export() string {
	if t.Done {
		return ""
	}
	result := t.Name
	if t.Note != "" {
		result += "\nNote: " + t.Note
	}
	if t.Deadline != nil {
		result += fmt.Sprintf("\nDeadline: %v", t.Deadline)
	}
	if t.Goal != nil {
		result += "\nGoal: " + t.Goal.Name
	}
	result += "\n"
	return result
}

func (t *Task) getItem(id string) *AlfredItem {
	var subtitle string
	var comma bool
	var order int

	if t.Deadline != nil {
		subtitle = t.Deadline.Format(DATE_FORMAT)
		comma = true
	}

	if t.Project != nil {
		if comma {
			subtitle += "; "
		}
		comma = true
		subtitle += t.Project.Name
	}

	if t.Goal != nil {
		if t.Project != nil {
			subtitle += ": "
		} else if comma {
			subtitle += "; "
		}
		comma = true
		subtitle += t.Goal.Name
	}

	if comma {
		subtitle += "; "
	}
	subtitle += strconv.Itoa(t.BasePoints)

	var todayTagPresent bool
	if len(t.Tags) > 0 {
		subtitle += "; Tags: "
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

	subtitle += "; Spent: "

	if t.TimeSpent == nil {
		if t.InProgress && t.InProgressSince != nil {
			subtitle += utils.DurationToHMFormat(utils.GetDurationPointer(time.Now().Sub(*t.InProgressSince))) + "/"
		} else {
			subtitle += "0h0m/"
		}
	} else {
		if t.InProgress {
			subtitle += utils.MinutesToHMFormat(t.TimeSpent.Minutes() + time.Now().Sub(*t.InProgressSince).Minutes()) + "/"
		} else {
			subtitle += utils.DurationToHMFormat(t.TimeSpent) + "/"
		}
	}
	if t.TimeEstimate == nil {
		subtitle += "?"
	} else {
		subtitle += utils.DurationToHMFormat(t.TimeEstimate)
	}

	var icoPath string
	if t.InProgress {
		order = 2
		icoPath = ICO_BLUE
	} else if t.Done {
		order = 2000 - t.BasePoints
		icoPath = ICO_GREEN
	} else if t.Deadline != nil && t.Deadline.Before(time.Now().Add(time.Hour * 24)) {
		order = 250 - t.BasePoints
		icoPath = ICO_RED
	} else if todayTagPresent {
		order = 500 - t.BasePoints
		icoPath = ICO_ORANGE
	} else if t.Goal != nil && t.Goal.Active {
		order = 750 - t.BasePoints
		icoPath = ICO_CYAN
	} else {
		order = 1000 - t.BasePoints
		icoPath = ICO_YELLOW
	}

	return &AlfredItem{
		Name:     t.Name,
		Arg:      id,
		Subtitle: subtitle,
		Icon:     NewAlfredIcon(icoPath),
		Valid:    true,
		order:    order}
}
