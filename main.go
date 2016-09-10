package main

import (
	"flag"
	"fmt"

	"github.com/7joe7/personalmanager/alfred"
	"github.com/7joe7/personalmanager/db"
	"github.com/7joe7/personalmanager/resources"
	"log"
	"os"
	"runtime/debug"
	"io"
	"github.com/7joe7/personalmanager/operations"
)

var (
	action, id, name, projectId, repetition, deadline *string
	noneAllowed, activeFlag, doneFlag                 *bool
	basePoints                                        *int
)

func init() {
	action = flag.String("action", "", fmt.Sprintf("Provide action to be taken from this list: %v.", resources.ACTIONS))
	id = flag.String("id", "", fmt.Sprintf("Provide id of the entity you want to make the action for. Valid for these actions: ."))
	projectId = flag.String("projectId", "", fmt.Sprintf("Provide project id for task assignment."))
	name = flag.String("name", "", "Provide name.")
	activeFlag = flag.Bool("active", false, "Toggle active/show active only.")
	doneFlag = flag.Bool("done", false, "Toggle done.")
	repetition = flag.String("repetition", "", "Select repetition period.")
	basePoints = flag.Int("basePoints", -1, "Set base points for success/failure.")
	deadline = flag.String("deadline", "", "Specify deadine in format 'dd.MM.YYYY HH:mm'.")
	noneAllowed = flag.Bool("noneAllowed", false, "Provide information whether list should be retrieved with none value allowed.")

	db.Open()
	t := db.NewTransaction()
	operations.InitializeBuckets(t)
	operations.EnsureValues(t)
	operations.Synchronize(t)
	t.Execute()

	f, err := os.OpenFile(resources.LOG_FILE_PATH, os.O_APPEND|os.O_CREATE, 777)
	if err != nil {
		log.Fatalf("Unable to open log file. %v", err)
	}
	log.SetOutput(io.MultiWriter(os.Stdout, f))
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("Panicked. %v %s", r, string(debug.Stack()))
		}
	}()
	flag.Parse()
	switch *action {
	case resources.ACT_CREATE_TASK:
		alfred.PrintResult(fmt.Sprintf(resources.MSG_CREATE_SUCCESS, "task",
			operations.AddTask(*name, *projectId)))
	case resources.ACT_CREATE_PROJECT:
		alfred.PrintResult(fmt.Sprintf(resources.MSG_CREATE_SUCCESS, "project",
			operations.AddProject(*name)))
	case resources.ACT_CREATE_TAG:
		alfred.PrintResult(fmt.Sprintf(resources.MSG_CREATE_SUCCESS, "tag",
			operations.AddTag(*name)))
	case resources.ACT_CREATE_GOAL:
		alfred.PrintResult(fmt.Sprintf(resources.MSG_CREATE_SUCCESS, "goal",
			operations.AddGoal(*name)))
	case resources.ACT_CREATE_HABIT:
		alfred.PrintResult(fmt.Sprintf(resources.MSG_CREATE_SUCCESS, "habit",
			operations.AddHabit(*name, *repetition, *activeFlag)))
	case resources.ACT_PRINT_TASKS:
		alfred.PrintEntities(resources.Tasks{operations.GetTasks(), *noneAllowed})
	case resources.ACT_PRINT_PROJECTS:
		alfred.PrintEntities(resources.Projects{operations.GetProjects(), *noneAllowed})
	case resources.ACT_PRINT_TAGS:
		alfred.PrintEntities(resources.Tags{operations.GetTags(), *noneAllowed})
	case resources.ACT_PRINT_GOALS:
		alfred.PrintEntities(resources.Goals{operations.GetGoals(), *noneAllowed})
	case resources.ACT_PRINT_HABITS:
		if *activeFlag {
			alfred.PrintEntities(resources.Habits{operations.GetActiveHabits(), *noneAllowed, operations.GetStatus()})
		} else {
			alfred.PrintEntities(resources.Habits{operations.GetNonActiveHabits(), *noneAllowed, operations.GetStatus()})
		}
	case resources.ACT_DELETE_TASK:
		operations.DeleteTask(*id)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_DELETE_SUCCESS, "task"))
	case resources.ACT_DELETE_PROJECT:
		operations.DeleteProject(*id)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_DELETE_SUCCESS, "project"))
	case resources.ACT_DELETE_TAG:
		operations.DeleteTag(*id)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_DELETE_SUCCESS, "tag"))
	case resources.ACT_DELETE_GOAL:
		operations.DeleteGoal(*id)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_DELETE_SUCCESS, "goal"))
	case resources.ACT_DELETE_HABIT:
		operations.DeleteHabit(*id)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_DELETE_SUCCESS, "habit"))
	case resources.ACT_MODIFY_TASK:
		operations.ModifyTask(*id, *name, *projectId)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_MODIFY_SUCCESS, "task"))
	case resources.ACT_MODIFY_PROJECT:
		operations.ModifyProject(*id, *name)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_MODIFY_SUCCESS, "project"))
	case resources.ACT_MODIFY_TAG:
		operations.ModifyTag(*id, *name)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_MODIFY_SUCCESS, "tag"))
	case resources.ACT_MODIFY_GOAL:
		operations.ModifyGoal(*id, *name)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_MODIFY_SUCCESS, "goal"))
	case resources.ACT_MODIFY_HABIT:
		operations.ModifyHabit(*id, *name, *repetition, *deadline, *activeFlag, *doneFlag, *basePoints)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_MODIFY_SUCCESS, "habit"))
	default:
		flag.Usage()
	}
}
