package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/debug"

	"github.com/7joe7/personalmanager/alfred"
	"github.com/7joe7/personalmanager/db"
	"github.com/7joe7/personalmanager/operations"
	"github.com/7joe7/personalmanager/resources"
	"github.com/7joe7/personalmanager/anybar"
)

var (
	action, id, name, projectId, repetition, deadline, estimate, scheduled, taskType, note *string
	noneAllowed, activeFlag, doneFlag, donePrevious                                        *bool
	basePoints                                                                             *int
)

func init() {
	action = flag.String("action", "", fmt.Sprintf("Provide action to be taken from this list: %v.", resources.ACTIONS))
	id = flag.String("id", "", fmt.Sprintf("Provide id of the entity you want to make the action for. Valid for these actions: ."))
	projectId = flag.String("projectId", "", fmt.Sprintf("Provide project id for task assignment."))
	name = flag.String("name", "", "Provide name.")
	activeFlag = flag.Bool("active", false, "Toggle active/show active only.")
	doneFlag = flag.Bool("done", false, "Toggle done.")
	donePrevious = flag.Bool("donePrevious", false, "Set done for previous period.")
	repetition = flag.String("repetition", "", "Select repetition period.")
	basePoints = flag.Int("basePoints", -1, "Set base points for success/failure.")
	deadline = flag.String("deadline", "", "Specify deadine in format 'dd.MM.YYYY HH:mm'.")
	estimate = flag.String("estimate", "", "Specify time estimate in format '2h45m'.")
	noneAllowed = flag.Bool("noneAllowed", false, "Provide information whether list should be retrieved with none value allowed.")
	scheduled = flag.String("scheduled", "", "Provide schedule period. (NEXT|NONE)")
	taskType = flag.String("taskType", "", "Provide task type. (PERSONAL|WORK)")
	note = flag.String("note", "", "Provide note.")

	anybar.Start(anybar.NewAnybarManager())

	db.Open(resources.DB_PATH)
	t := db.NewTransaction()
	operations.InitializeBuckets(t)
	operations.EnsureValues(t)
	operations.Synchronize(t)
	t.Execute()

	f, err := os.OpenFile(resources.LOG_FILE_PATH, os.O_APPEND|os.O_CREATE, 777)
	if err != nil {
		panic(err)
	}
	log.SetOutput(io.MultiWriter(os.Stdout, f))
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("Panicked. %v %s", r, string(debug.Stack()))
			os.Exit(3)
		}
	}()
	flag.Parse()
	switch *action {
	case resources.ACT_CREATE_TASK:
		operations.AddTask(*name, *projectId, *deadline, *estimate, *scheduled, *taskType, *activeFlag, *basePoints)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_CREATE_SUCCESS, "task"))
	case resources.ACT_CREATE_PROJECT:
		operations.AddProject(*name)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_CREATE_SUCCESS, "project"))
	case resources.ACT_CREATE_TAG:
		operations.AddTag(*name)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_CREATE_SUCCESS, "tag"))
	case resources.ACT_CREATE_GOAL:
		operations.AddGoal(*name)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_CREATE_SUCCESS, "goal"))
	case resources.ACT_CREATE_HABIT:
		operations.AddHabit(*name, *repetition, *deadline, *activeFlag, *basePoints)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_CREATE_SUCCESS, "habit"))
	case resources.ACT_PRINT_TASKS:
		alfred.PrintEntities(resources.Tasks{operations.GetTasks(), *noneAllowed, operations.GetStatus()})
	case resources.ACT_PRINT_PERSONAL_NEXT_TASKS:
		alfred.PrintEntities(resources.Tasks{operations.GetNextTasks(), *noneAllowed, operations.GetStatus()})
	case resources.ACT_PRINT_PERSONAL_UNSCHEDULED_TASKS:
		alfred.PrintEntities(resources.Tasks{operations.GetUnscheduledTasks(), *noneAllowed, operations.GetStatus()})
	case resources.ACT_PRINT_WORK_NEXT_TASKS:
		alfred.PrintEntities(resources.Tasks{operations.GetWorkNextTasks(), *noneAllowed, operations.GetStatus()})
	case resources.ACT_PRINT_WORK_UNSCHEDULED_TASKS:
		alfred.PrintEntities(resources.Tasks{operations.GetWorkUnscheduledTasks(), *noneAllowed, operations.GetStatus()})
	case resources.ACT_PRINT_TASK_NOTE:
		alfred.PrintResult(operations.GetTask(*id).Note)
	case resources.ACT_PRINT_PROJECTS:
		alfred.PrintEntities(resources.Projects{operations.GetProjects(), *noneAllowed, operations.GetStatus()})
	case resources.ACT_PRINT_TAGS:
		alfred.PrintEntities(resources.Tags{operations.GetTags(), *noneAllowed, operations.GetStatus()})
	case resources.ACT_PRINT_GOALS:
		alfred.PrintEntities(resources.Goals{operations.GetGoals(), *noneAllowed, operations.GetStatus()})
	case resources.ACT_PRINT_HABITS:
		if *activeFlag {
		alfred.PrintEntities(resources.Habits{operations.GetActiveHabits(), *noneAllowed, operations.GetStatus()})
	} else {
		alfred.PrintEntities(resources.Habits{operations.GetNonActiveHabits(), *noneAllowed, operations.GetStatus()})
	}
	case resources.ACT_PRINT_REVIEW:
		alfred.PrintEntities(resources.Items{[]*resources.AlfredItem{operations.GetReview().GetItem()}})
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
		operations.ModifyTask(*id, *name, *projectId, *deadline, *estimate, *scheduled, *taskType, *note, *basePoints, *activeFlag, *doneFlag)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_MODIFY_SUCCESS, "task"))
	case resources.ACT_MODIFY_PROJECT:
		operations.ModifyProject(*id, *name)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_MODIFY_SUCCESS, "project"))
	case resources.ACT_MODIFY_TAG:
		operations.ModifyTag(*id, *name)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_MODIFY_SUCCESS, "tag"))
	case resources.ACT_MODIFY_GOAL:
		operations.ModifyGoal(*id, *name, *deadline)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_MODIFY_SUCCESS, "goal"))
	case resources.ACT_MODIFY_HABIT:
		operations.ModifyHabit(*id, *name, *repetition, *deadline, *activeFlag, *doneFlag, *donePrevious, *basePoints)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_MODIFY_SUCCESS, "habit"))
	case resources.ACT_MODIFY_REVIEW:
		operations.ModifyReview(*repetition, *deadline)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_MODIFY_SUCCESS, "review"))
	case resources.ACT_DEBUG_DATABASE:
		db.PrintoutDbContents(*id)
	case resources.ACT_CUSTOM:

	default:
		flag.Usage()
	}
	resources.WaitGroup.Wait()
}
