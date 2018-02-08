package resources

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/7joe7/personalmanager/utils"
)

type Task struct {
	Id              string         `json:",omitempty"`
	Name            string         `json:",omitempty"`
	Note            string         `json:",omitempty"`
	BasePoints      int            `json:",omitempty"`
	InProgress      bool           `json:",omitempty"`
	Done            bool           `json:",omitempty"`
	Tags            []*Tag         `json:",omitempty"`
	Scheduled       string         `json:",omitempty"`
	Type            string         `json:",omitempty"`
	InProgressSince *time.Time     `json:",omitempty"`
	DoneTime        *time.Time     `json:",omitempty"`
	Deadline        *time.Time     `json:",omitempty"`
	TimeEstimate    *time.Duration `json:",omitempty"`
	TimeSpent       *time.Duration `json:",omitempty"`
	Project         *Project       `json:",omitempty"`
	Goal            *Goal          `json:",omitempty"`
}

func (t *Task) SetId(id string) {
	t.Id = id
}

func (t *Task) GetId() string {
	return t.Id
}

func (t *Task) GetTimeEstimate() *time.Duration {
	return t.TimeEstimate
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

func (t *Task) Less(entity Entity) bool {
	switch entity.(type) {
	case *Task:
		otherTask := entity.(*Task)
		if t.InProgress != otherTask.InProgress {
			return t.InProgress
		}
		if t.Done != otherTask.Done {
			return !t.Done
		}
		if (t.Goal != nil && t.Goal.Active) != (otherTask.Goal != nil && otherTask.Goal.Active) {
			return t.Goal != nil && t.Goal.Active
		}
		if t.BasePoints != otherTask.BasePoints {
			return t.BasePoints > otherTask.BasePoints
		}
		if (t.Deadline != nil) != (otherTask.Deadline != nil) {
			return t.Deadline != nil
		}
		if t.Deadline != nil && t.Deadline.Day() != otherTask.Deadline.Day() {
			return t.Deadline.Before(*otherTask.Deadline)
		}
		if (t.TimeEstimate == nil) != (otherTask.TimeEstimate == nil) {
			return t.TimeEstimate != nil
		}
		if t.TimeEstimate != nil && t.TimeEstimate.Minutes() != otherTask.TimeEstimate.Minutes() {
			return t.TimeEstimate.Minutes() < otherTask.TimeEstimate.Minutes()
		}
		return t.Name < otherTask.Name
	case *Habit:
		otherHabit := entity.(*Habit)
		if t.InProgress {
			return true
		}
		if otherHabit.Repetition == HBT_REPETITION_DAILY {
			return false
		}
		if t.BasePoints != otherHabit.BasePoints {
			return t.BasePoints > otherHabit.BasePoints
		}
		if (t.TimeEstimate == nil) != (otherHabit.TimeEstimate == nil) {
			return t.TimeEstimate != nil
		}
		if t.TimeEstimate != nil && t.TimeEstimate.Minutes() != otherHabit.TimeEstimate.Minutes() {
			return t.TimeEstimate.Minutes() < otherHabit.TimeEstimate.Minutes()
		}
		return false
	}
	return false
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
	result += "\n"
	return result
}

func (t *Task) GetAlfredItem(id string) *AlfredItem {
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
			subtitle += utils.MinutesToHMFormat(t.TimeSpent.Minutes()+time.Now().Sub(*t.InProgressSince).Minutes()) + "/"
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
		icoPath = ICO_BLUE
	} else if t.Done {
		icoPath = ICO_GREEN
	} else if t.Deadline != nil && t.Deadline.Before(time.Now().Add(time.Hour*24)) {
		icoPath = ICO_RED
	} else if todayTagPresent {
		icoPath = ICO_ORANGE
	} else if t.Goal != nil && t.Goal.Active {
		icoPath = ICO_CYAN
	} else {
		icoPath = ICO_YELLOW
	}

	return &AlfredItem{
		Name:     t.Name,
		Arg:      id,
		Subtitle: subtitle,
		Icon:     NewAlfredIcon(icoPath),
		Valid:    true,
		entity:   t}
}

func (t *Task) CountScoreChange(status *Status) int {
	change := t.BasePoints * 10
	if t.TimeEstimate != nil {
		change += int(t.TimeEstimate.Minutes()) * t.BasePoints
	}
	koef := 1
	for status.WorkDoneToday > koef {
		koef *= 2
	}
	return change * koef
}
