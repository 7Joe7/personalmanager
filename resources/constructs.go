package resources

import (
	"github.com/7joe7/personalmanager/resources/validation"
)

func NewTask(name string) *Task {
	validation.VerifyTask(name)
	return &Task{Name: name, BasePoints: 1}
}

func NewProject(name string) *Project {
	validation.VerifyProject(name)
	return &Project{Name: name}
}

func NewTag(name string) *Tag {
	validation.VerifyTag(name)
	return &Tag{Name: name}
}

func NewGoal(name string) *Goal {
	validation.VerifyGoal(name)
	return &Goal{Name: name}
}

func NewHabit(name string) *Habit {
	validation.VerifyHabit(name)
	return &Habit{Name: name}
}

func NewAlfredIcon(path string) *AlfredIcon {
	return &AlfredIcon{Path: path}
}
