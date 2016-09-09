package checks

import (
	"log"
)

func VerifyTask(name string) {
	if name == "" {
		log.Fatal("Task name is empty.")
	}
}

func VerifyProject(name string) {
	if name == "" {
		log.Fatal("Project name is empty.")
	}
}

func VerifyTag(name string) {
	if name == "" {
		log.Fatal("Tag name is empty.")
	}
}

func VerifyGoal(name string) {
	if name == "" {
		log.Fatal("Goal name is empty.")
	}
}

func VerifyHabit(name string) {
	if name == "" {
		log.Fatal("Habit name is empty.")
	}
}
