package resources

import (
	"log"
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

func NewAlfredIcon(color string) *AlfredIcon {
	switch color {
	case "black_alt","black","blue","cyan","exclamation","green","orange","purple","question","question_alt","red","white","white_alt","yellow":
		return &AlfredIcon{Path: "./icons/" + color + "@2x.png"}
	case "special":
		return &AlfredIcon{Path: "./icons/" + color + ".png"}
	default:
		log.Fatalf("Icon '%s' not supported. Add icon to icons folder and add it to the supported icons switch.", color)
		return nil
	}
}
