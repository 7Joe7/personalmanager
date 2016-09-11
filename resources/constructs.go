package resources

import (
	"github.com/7joe7/personalmanager/checks"
)

func NewTask(name string, project *Project) *Task {
	checks.VerifyTask(name)
	return &Task{Name:name, Project: project}
}

func NewProject(name string) *Project {
	checks.VerifyProject(name)
	return &Project{Name:name}
}

func NewTag(name string) *Tag {
	checks.VerifyTag(name)
	return &Tag{Name:name}
}

func NewGoal(name string) *Goal {
	checks.VerifyGoal(name)
	return &Goal{Name:name}
}

func NewHabit(name string) *Habit {
	checks.VerifyHabit(name)
	return &Habit{Name:name}
}

func NewAlfredIcon(path string) *AlfredIcon {
	return &AlfredIcon{Path:path}
}
