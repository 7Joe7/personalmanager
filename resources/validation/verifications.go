package validation

import (
	"errors"
)

func verifyTask(name string) error {
	if name == "" {
		return errors.New("Task name is empty.")
	}
	return nil
}

func verifyProject(name string) error {
	if name == "" {
		return errors.New("Project name is empty.")
	}
	return nil
}

func verifyTag(name string) error {
	if name == "" {
		return errors.New("Tag name is empty.")
	}
	return nil
}

func verifyGoal(name string) error {
	if name == "" {
		return errors.New("Goal name is empty.")
	}
	return nil
}

func verifyHabit(name string) error {
	if name == "" {
		return errors.New("Habit name is empty.")
	}
	return nil
}
