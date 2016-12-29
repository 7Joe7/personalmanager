package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/7joe7/personalmanager/alfred"
	"github.com/7joe7/personalmanager/anybar"
	"github.com/7joe7/personalmanager/db"
	"github.com/7joe7/personalmanager/operations"
	"github.com/7joe7/personalmanager/resources"
	"github.com/7joe7/personalmanager/utils"
)

var (
	// parameters
	action, id, name, projectId, goalId, taskId, habitId, repetition, deadline, estimate, scheduled, taskType, note *string
	noneAllowed, activeFlag, doneFlag, donePrevious, undonePrevious, negativeFlag                                   *bool
	basePoints, habitRepetitionGoal                                                                                 *int
)

func init() {
	action = flag.String("action", "", fmt.Sprintf("Provide action to be taken from this list: %v.", resources.ACTIONS))
	id = flag.String("id", "", fmt.Sprintf("Provide id of the entity you want to make the action for. Valid for these actions: ."))
	name = flag.String("name", "", "Provide name.")
	projectId = flag.String("projectId", "", fmt.Sprintf("Provide project id for project assignment."))
	goalId = flag.String("goalId", "", fmt.Sprintf("Provide goal id for goal assignment."))
	taskId = flag.String("taskId", "", fmt.Sprintf("Provide task id for task assignment."))
	habitId = flag.String("habitId", "", fmt.Sprintf("Provide habit id for habit assignment."))
	repetition = flag.String("repetition", "", "Select repetition period.")
	deadline = flag.String("deadline", "", "Specify deadline in format 'dd.MM.YYYY HH:mm'.")
	estimate = flag.String("estimate", "", "Specify time estimate in format '2h45m'.")
	scheduled = flag.String("scheduled", "", "Provide schedule period. (NEXT|NONE)")
	taskType = flag.String("taskType", "", "Provide task type. (PERSONAL|WORK)")
	note = flag.String("note", "", "Provide note.")
	noneAllowed = flag.Bool("noneAllowed", false, "Provide information whether list should be retrieved with none value allowed.")
	activeFlag = flag.Bool("active", false, "Toggle active/show active only.")
	doneFlag = flag.Bool("done", false, "Toggle done.")
	donePrevious = flag.Bool("donePrevious", false, "Set done for previous period.")
	undonePrevious = flag.Bool("undonePrevious", false, "Set undone for previous period.")
	negativeFlag = flag.Bool("negative", false, "Set negative flag for habits.")
	basePoints = flag.Int("basePoints", -1, "Set base points for success/failure.")
	habitRepetitionGoal = flag.Int("habitRepetitionGoal", -1, "Set habit goal repetition number.")
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panicked. %v %s\n", r, string(debug.Stack()))
			log.Fatalf("Panicked. %v %s", r, string(debug.Stack()))
			os.Exit(3)
		}
	}()

	flag.Parse()

	f, err := os.OpenFile(fmt.Sprintf("%s/%s", utils.GetRunningBinaryPath(), resources.LOG_FILE_NAME), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	log.SetOutput(f)

	logBinaryCall()

	anybar.Start(anybar.NewAnybarManager())
	db.Open()
	t := db.NewTransaction()
	operations.InitializeBuckets(t)
	operations.EnsureValues(t)
	operations.Synchronize(t)
	t.Execute()

	switch *action {
	case resources.ACT_CREATE_TASK:
		operations.AddTask(*name, *projectId, *goalId, *deadline, *estimate, *scheduled, *taskType, *note, *activeFlag, *basePoints)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_CREATE_SUCCESS, "task"))
	case resources.ACT_CREATE_PROJECT:
		operations.AddProject(*name)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_CREATE_SUCCESS, "project"))
	case resources.ACT_CREATE_TAG:
		operations.AddTag(*name)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_CREATE_SUCCESS, "tag"))
	case resources.ACT_CREATE_GOAL:
		operations.AddGoal(*name, *projectId, *habitId, *habitRepetitionGoal)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_CREATE_SUCCESS, "goal"))
	case resources.ACT_CREATE_HABIT:
		operations.AddHabit(*name, *repetition, *note, *deadline, *goalId, *activeFlag, *negativeFlag, *basePoints, *habitRepetitionGoal)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_CREATE_SUCCESS, "habit"))
	case resources.ACT_PRINT_TASKS:
		alfred.PrintEntities(resources.Tasks{Tasks: operations.GetTasks(), NoneAllowed: *noneAllowed, Status: operations.GetStatus()})
	case resources.ACT_PRINT_PERSONAL_TASKS:
		alfred.PrintEntities(resources.Tasks{Tasks: operations.GetPersonalTasks(), NoneAllowed: *noneAllowed, Status: operations.GetStatus()})
	case resources.ACT_PRINT_PERSONAL_NEXT_TASKS:
		alfred.PrintEntities(resources.Tasks{Tasks: operations.GetNextTasks(), NoneAllowed: *noneAllowed, Status: operations.GetStatus(), Sum: true})
	case resources.ACT_PRINT_PERSONAL_UNSCHEDULED_TASKS:
		alfred.PrintEntities(resources.Tasks{Tasks: operations.GetUnscheduledTasks(), NoneAllowed: *noneAllowed, Status: operations.GetStatus()})
	case resources.ACT_PRINT_SHOPPING_TASKS:
		alfred.PrintEntities(resources.Tasks{Tasks: operations.GetShoppingTasks(), NoneAllowed: *noneAllowed, Status: operations.GetStatus()})
	case resources.ACT_PRINT_WORK_NEXT_TASKS:
		alfred.PrintEntities(resources.Tasks{Tasks: operations.GetWorkNextTasks(), NoneAllowed: *noneAllowed, Status: operations.GetStatus()})
	case resources.ACT_PRINT_WORK_UNSCHEDULED_TASKS:
		alfred.PrintEntities(resources.Tasks{Tasks: operations.GetWorkUnscheduledTasks(), NoneAllowed: *noneAllowed, Status: operations.GetStatus()})
	case resources.ACT_PRINT_TASK_NOTE:
		alfred.PrintResult(operations.GetTask(*id).Note)
	case resources.ACT_PRINT_PROJECTS:
		alfred.PrintEntities(resources.Projects{operations.GetProjects(), *noneAllowed, operations.GetStatus()})
	case resources.ACT_PRINT_ACTIVE_PROJECTS:
		alfred.PrintEntities(resources.Projects{operations.GetActiveProjects(), *noneAllowed, operations.GetStatus()})
	case resources.ACT_PRINT_INACTIVE_PROJECTS:
		alfred.PrintEntities(resources.Projects{operations.GetInactiveProjects(), *noneAllowed, operations.GetStatus()})
	case resources.ACT_PRINT_TAGS:
		alfred.PrintEntities(resources.Tags{operations.GetTags(), *noneAllowed, operations.GetStatus()})
	case resources.ACT_PRINT_GOALS:
		alfred.PrintEntities(resources.Goals{operations.GetGoals(), *noneAllowed, operations.GetStatus()})
	case resources.ACT_PRINT_ACTIVE_GOALS:
		alfred.PrintEntities(resources.Goals{operations.GetActiveGoals(), *noneAllowed, operations.GetStatus()})
	case resources.ACT_PRINT_NON_ACTIVE_GOALS:
		alfred.PrintEntities(resources.Goals{operations.GetNonActiveGoals(), *noneAllowed, operations.GetStatus()})
	case resources.ACT_PRINT_HABITS:
		if *activeFlag {
			alfred.PrintEntities(resources.Habits{operations.GetActiveHabits(), *noneAllowed, operations.GetStatus(), true})
		} else {
			alfred.PrintEntities(resources.Habits{operations.GetNonActiveHabits(), *noneAllowed, operations.GetStatus(), false})
		}
	case resources.ACT_PRINT_HABIT_DESCRIPTION:
		alfred.PrintResult(operations.GetHabit(*id).Description)
	case resources.ACT_PRINT_REVIEW:
		alfred.PrintEntities(resources.Items{[]*resources.AlfredItem{operations.GetReview().GetItem()}})
	case resources.ACT_EXPORT_SHOPPING_TASKS:
		operations.ExportShoppingTasks()
		alfred.PrintResult(fmt.Sprintf(resources.MSG_EXPORT_SUCCESS, "shopping tasks"))
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
		operations.ModifyTask(*id, *name, *projectId, *goalId, *deadline, *estimate, *scheduled, *taskType, *note, *basePoints, *activeFlag, *doneFlag)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_MODIFY_SUCCESS, "task"))
	case resources.ACT_MODIFY_PROJECT:
		operations.ModifyProject(*id, *name, *taskId, *goalId, *activeFlag, *doneFlag)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_MODIFY_SUCCESS, "project"))
	case resources.ACT_MODIFY_TAG:
		operations.ModifyTag(*id, *name)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_MODIFY_SUCCESS, "tag"))
	case resources.ACT_MODIFY_GOAL:
		operations.ModifyGoal(*id, *name, *taskId, *projectId, *habitId, *activeFlag, *doneFlag, *habitRepetitionGoal)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_MODIFY_SUCCESS, "goal"))
	case resources.ACT_MODIFY_HABIT:
		operations.ModifyHabit(*id, *name, *repetition, *note, *deadline, *goalId, *activeFlag, *doneFlag, *donePrevious, *undonePrevious, *negativeFlag, *basePoints, *habitRepetitionGoal)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_MODIFY_SUCCESS, "habit"))
	case resources.ACT_MODIFY_REVIEW:
		operations.ModifyReview(*repetition, *deadline)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_MODIFY_SUCCESS, "review"))
	case resources.ACT_SYNC_ANYBAR_PORTS:
		t := db.NewTransaction()
		operations.SynchronizeAnybarPorts(t)
		t.Execute()
	case resources.ACT_DEBUG_DATABASE:
		db.PrintoutDbContents(*id)
	case resources.ACT_SYNC_WITH_JIRA:
		operations.SyncWithJira()
	case resources.ACT_BACKUP_DATABASE:
		db.BackupDatabase()
	case resources.ACT_SET_EMAIL:
		operations.SetEmail(*name)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_SET_SUCCESS, "e-mail", *name))
	case resources.ACT_CUSTOM:
		t := db.NewTransaction()
		t.Add(func() error {
			habit := &resources.Habit{}
			err = t.ModifyEntity(resources.DB_DEFAULT_HABITS_BUCKET_NAME, []byte("137"), true, habit, func () {
				habit.Tries = 1
				habit.Count = 0
				habit.LastStreak = 0
			})
			if err != nil {
				return err
			}
			err = t.ModifyEntity(resources.DB_DEFAULT_HABITS_BUCKET_NAME, []byte("138"), true, habit, func () {
				habit.Tries = 1
				habit.Count = 0
				habit.LastStreak = 0
			})
			//getNewHabit := func() resources.Entity {
			//	return &resources.Habit{}
			//}
			//err := t.MapEntities(resources.DB_DEFAULT_HABITS_BUCKET_NAME, true, getNewHabit, func(e resources.Entity) func() {
			//	return func() {
			//		h := e.(*resources.Habit)
			//		if h.Goal != nil && h.Goal.Id == "52" {
			//			fmt.Printf("h: %v\n", h)
			//			h.Goal = nil
			//		}
			//	}
			//})
			if err != nil {
				return err
			}
			return nil
		})
		t.Execute()
	default:
		flag.Usage()
	}
	resources.WaitGroup.Wait()
}

func logBinaryCall() {
	log.Printf(`Called with string parameters:
		action: %s,
		id: %s,
		name: %s,
		projectId: %s,
		goalId: %s,
		taskId: %s,
		habitId: %s,
		repetition: %s,
		deadline: %s,
		estimate: %s,
		scheduled: %s,
		taskType: %s,
		note: %s,
		and with bool parameters:
		noneAllowed: %v,
		activeFlag: %v,
		doneFlag: %v,
		donePrevious: %v,
		undonePrevious: %v,
		and int parameters:
		basePoints: %v,
		repetitionGoal: %v.`, *action, *id, *name, *projectId, *goalId, *taskId,
		*habitId, *repetition, *deadline, *estimate, *scheduled, *taskType, *note, *noneAllowed, *activeFlag,
		*doneFlag, *donePrevious, *undonePrevious, *basePoints, *habitRepetitionGoal)
}
