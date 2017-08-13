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

func NewCommand(
	action, id, name, projectId, goalId, taskId,
	repetition, deadline, alarm, estimate, scheduled,
	taskType, note string, noneAllowed, activeFlag,
	doneFlag, donePrevious, undonePrevious, negativeFlag,
	learnedFlag bool, basePoints, habitRepetitionGoal int) *Command {
	return &Command{
		Action: action, ID: id, Name: name, ProjectID: projectId, GoalID: goalId,
		TaskID: taskId, Repetition: repetition, Deadline: deadline, Alarm: alarm, Estimate: estimate,
		Scheduled: scheduled, TaskType: taskType, Note: note, NoneAllowed: noneAllowed,
		ActiveFlag: activeFlag, DoneFlag: doneFlag, DonePrevious: donePrevious,
		UndonePrevious: undonePrevious, NegativeFlag: negativeFlag, LearnedFlag: learnedFlag,
		BasePoints: basePoints, HabitRepetitionGoal: habitRepetitionGoal,
	}
}
