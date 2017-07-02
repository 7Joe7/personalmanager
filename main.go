package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"time"

	"github.com/7joe7/personalmanager/operations/alfred"
	"github.com/7joe7/personalmanager/operations/anybar"
	"github.com/7joe7/personalmanager/db"
	"github.com/7joe7/personalmanager/operations"
	"github.com/7joe7/personalmanager/operations/goals"
	"github.com/7joe7/personalmanager/operations/configuration"
	"github.com/7joe7/personalmanager/operations/exporter"
	"github.com/7joe7/personalmanager/resources"
	"github.com/7joe7/personalmanager/utils"
	rutils "github.com/7joe7/personalmanager/resources/utils"
)

var (
	// parameters
	action, id, name, projectId, goalId, taskId, habitId, repetition, deadline, estimate, scheduled, taskType, note *string
	noneAllowed, activeFlag, doneFlag, donePrevious, undonePrevious, negativeFlag, learnedFlag                      *bool
	basePoints, habitRepetitionGoal                                                                                 *int
)

func init() {
	action = flag.String("action", "", fmt.Sprintf("Provide action to be taken from this list: %v.", resources.ACTIONS))
	id = flag.String("id", "", "Provide id of the entity you want to make the action for. Valid for these actions: .")
	name = flag.String("name", "", "Provide name.")
	projectId = flag.String("projectId", "", "Provide project id for project assignment.")
	goalId = flag.String("goalId", "", "Provide goal id for goal assignment.")
	taskId = flag.String("taskId", "", "Provide task id for task assignment.")
	habitId = flag.String("habitId", "", "Provide habit id for habit assignment.")
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
	learnedFlag = flag.Bool("learned", false, "Set learned flag for habits.")
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

	err := ensureAppSupportFolder(rutils.GetAppSupportFolderPath())
	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile(fmt.Sprintf("%s/%s", rutils.GetAppSupportFolderPath(), resources.LOG_FILE_NAME), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	log.SetOutput(f)
	resources.Alf = alfred.NewAlfred(os.Stdout)
	resources.Abr = anybar.NewAnybarManager(utils.GetRunningBinaryPath())

	logBinaryCall()

	db.Open()
	t := db.NewTransaction()
	operations.InitializeBuckets(t)
	operations.EnsureValues(t)
	operations.Synchronize(t)
	t.Execute()

	switch *action {
	case resources.ACT_CREATE_TASK:
		operations.AddTask(*name, *projectId, *goalId, *deadline, *estimate, *scheduled, *taskType, *note, *activeFlag, *basePoints)
		resources.Alf.PrintResult(fmt.Sprintf(resources.MSG_CREATE_SUCCESS, "task"))
	case resources.ACT_CREATE_PROJECT:
		operations.AddProject(*name)
		resources.Alf.PrintResult(fmt.Sprintf(resources.MSG_CREATE_SUCCESS, "project"))
	case resources.ACT_CREATE_TAG:
		operations.AddTag(*name)
		resources.Alf.PrintResult(fmt.Sprintf(resources.MSG_CREATE_SUCCESS, "tag"))
	case resources.ACT_CREATE_GOAL:
		goals.AddGoal(*name, *projectId, *habitId, *habitRepetitionGoal, *basePoints)
		resources.Alf.PrintResult(fmt.Sprintf(resources.MSG_CREATE_SUCCESS, "goal"))
	case resources.ACT_CREATE_HABIT:
		operations.AddHabit(*name, *repetition, *note, *deadline, *goalId, *activeFlag, *negativeFlag, *basePoints, *habitRepetitionGoal)
		resources.Alf.PrintResult(fmt.Sprintf(resources.MSG_CREATE_SUCCESS, "habit"))
	case resources.ACT_PRINT_TASKS:
		resources.Alf.PrintEntities(resources.Tasks{Tasks: operations.GetTasks(), NoneAllowed: *noneAllowed, Status: operations.GetStatus()})
	case resources.ACT_PRINT_PERSONAL_TASKS:
		resources.Alf.PrintEntities(resources.Tasks{Tasks: operations.GetPersonalTasks(), NoneAllowed: *noneAllowed, Status: operations.GetStatus()})
	case resources.ACT_PRINT_PERSONAL_NEXT_TASKS:
		resources.Alf.PrintEntities(resources.Tasks{Tasks: operations.GetNextTasks(), NoneAllowed: *noneAllowed, Status: operations.GetStatus(), Sum: true})
	case resources.ACT_PRINT_PERSONAL_UNSCHEDULED_TASKS:
		resources.Alf.PrintEntities(resources.Tasks{Tasks: operations.GetUnscheduledTasks(), NoneAllowed: *noneAllowed, Status: operations.GetStatus()})
	case resources.ACT_PRINT_SHOPPING_TASKS:
		resources.Alf.PrintEntities(resources.Tasks{Tasks: operations.GetShoppingTasks(), NoneAllowed: *noneAllowed, Status: operations.GetStatus()})
	case resources.ACT_PRINT_WORK_NEXT_TASKS:
		resources.Alf.PrintEntities(resources.Tasks{Tasks: operations.GetWorkNextTasks(), NoneAllowed: *noneAllowed, Status: operations.GetStatus()})
	case resources.ACT_PRINT_WORK_UNSCHEDULED_TASKS:
		resources.Alf.PrintEntities(resources.Tasks{Tasks: operations.GetWorkUnscheduledTasks(), NoneAllowed: *noneAllowed, Status: operations.GetStatus()})
	case resources.ACT_PRINT_TASK_NOTE:
		resources.Alf.PrintResult(operations.GetTask(*id).Note)
	case resources.ACT_PRINT_PROJECTS:
		resources.Alf.PrintEntities(resources.Projects{operations.GetProjects(), *noneAllowed, operations.GetStatus()})
	case resources.ACT_PRINT_ACTIVE_PROJECTS:
		resources.Alf.PrintEntities(resources.Projects{operations.GetActiveProjects(), *noneAllowed, operations.GetStatus()})
	case resources.ACT_PRINT_INACTIVE_PROJECTS:
		resources.Alf.PrintEntities(resources.Projects{operations.GetInactiveProjects(), *noneAllowed, operations.GetStatus()})
	case resources.ACT_PRINT_TAGS:
		resources.Alf.PrintEntities(resources.Tags{operations.GetTags(), *noneAllowed, operations.GetStatus()})
	case resources.ACT_PRINT_GOALS:
		resources.Alf.PrintEntities(resources.Goals{goals.GetGoals(), *noneAllowed, operations.GetStatus()})
	case resources.ACT_PRINT_ACTIVE_GOALS:
		resources.Alf.PrintEntities(resources.Goals{goals.GetActiveGoals(), *noneAllowed, operations.GetStatus()})
	case resources.ACT_PRINT_NON_ACTIVE_GOALS:
		resources.Alf.PrintEntities(resources.Goals{goals.GetNonActiveGoals(), *noneAllowed, operations.GetStatus()})
	case resources.ACT_PRINT_INCOMPLETE_GOALS:
		resources.Alf.PrintEntities(resources.Goals{goals.GetIncompleteGoals(), *noneAllowed, operations.GetStatus()})
	case resources.ACT_PRINT_HABITS:
		if *activeFlag {
			resources.Alf.PrintEntities(resources.Habits{operations.GetActiveHabits(), *noneAllowed, operations.GetStatus(), true})
		} else {
			resources.Alf.PrintEntities(resources.Habits{operations.GetNonActiveHabits(), *noneAllowed, operations.GetStatus(), false})
		}
	case resources.ACT_PRINT_HABIT_DESCRIPTION:
		resources.Alf.PrintResult(operations.GetHabit(*id).Description)
	case resources.ACT_PRINT_REVIEW:
		resources.Alf.PrintEntities(resources.Items{[]*resources.AlfredItem{operations.GetReview().GetItem()}})
	case resources.ACT_EXPORT_SHOPPING_TASKS:
		exporter.ExportShoppingTasks(resources.CFG_EXPORT_CONFIG_PATH)
		resources.Alf.PrintResult(fmt.Sprintf(resources.MSG_EXPORT_SUCCESS, "shopping tasks"))
	case resources.ACT_DELETE_TASK:
		operations.DeleteTask(*id)
		resources.Alf.PrintResult(fmt.Sprintf(resources.MSG_DELETE_SUCCESS, "task"))
	case resources.ACT_DELETE_PROJECT:
		operations.DeleteProject(*id)
		resources.Alf.PrintResult(fmt.Sprintf(resources.MSG_DELETE_SUCCESS, "project"))
	case resources.ACT_DELETE_TAG:
		operations.DeleteTag(*id)
		resources.Alf.PrintResult(fmt.Sprintf(resources.MSG_DELETE_SUCCESS, "tag"))
	case resources.ACT_DELETE_GOAL:
		goals.DeleteGoal(*id)
		resources.Alf.PrintResult(fmt.Sprintf(resources.MSG_DELETE_SUCCESS, "goal"))
	case resources.ACT_DELETE_HABIT:
		operations.DeleteHabit(*id)
		resources.Alf.PrintResult(fmt.Sprintf(resources.MSG_DELETE_SUCCESS, "habit"))
	case resources.ACT_MODIFY_TASK:
		operations.ModifyTask(*id, *name, *projectId, *goalId, *deadline, *estimate, *scheduled, *taskType, *note, *basePoints, *activeFlag, *doneFlag)
		resources.Alf.PrintResult(fmt.Sprintf(resources.MSG_MODIFY_SUCCESS, "task"))
	case resources.ACT_MODIFY_PROJECT:
		operations.ModifyProject(*id, *name, *taskId, *goalId, *activeFlag, *doneFlag)
		resources.Alf.PrintResult(fmt.Sprintf(resources.MSG_MODIFY_SUCCESS, "project"))
	case resources.ACT_MODIFY_TAG:
		operations.ModifyTag(*id, *name)
		resources.Alf.PrintResult(fmt.Sprintf(resources.MSG_MODIFY_SUCCESS, "tag"))
	case resources.ACT_MODIFY_GOAL:
		goals.ModifyGoal(*id, *name, *taskId, *projectId, *habitId, *activeFlag, *doneFlag, *habitRepetitionGoal, *basePoints)
		resources.Alf.PrintResult(fmt.Sprintf(resources.MSG_MODIFY_SUCCESS, "goal"))
	case resources.ACT_MODIFY_HABIT:
		operations.ModifyHabit(*id, *name, *repetition, *note, *deadline, *goalId, *activeFlag, *doneFlag, *donePrevious, *undonePrevious, *negativeFlag, *learnedFlag, *basePoints, *habitRepetitionGoal)
		resources.Alf.PrintResult(fmt.Sprintf(resources.MSG_MODIFY_SUCCESS, "habit"))
	case resources.ACT_MODIFY_REVIEW:
		operations.ModifyReview(*repetition, *deadline)
		resources.Alf.PrintResult(fmt.Sprintf(resources.MSG_MODIFY_SUCCESS, "review"))
	case resources.ACT_SYNC_ANYBAR_PORTS:
		t := db.NewTransaction()
		operations.SynchronizeAnybarPorts(t)
		t.Execute()
		resources.Alf.PrintResult(fmt.Sprintf(resources.MSG_SYNC_SUCCESS, "AnyBar ports"))
	case resources.ACT_DEBUG_DATABASE:
		db.PrintoutDbContents(*id)
	case resources.ACT_SYNC_WITH_JIRA:
		operations.SyncWithJira()
	case resources.ACT_BACKUP_DATABASE:
		db.BackupDatabase()
	case resources.ACT_SET_CONFIG_VALUE:
		switch *id {
		case string(resources.DB_DEFAULT_EMAIL):
			exporter.SetEmail(*name)
			resources.Alf.PrintResult(fmt.Sprintf(resources.MSG_SET_SUCCESS, "e-mail", *name))
		case string(resources.DB_WEEKS_LEFT):
			configuration.SetWeeksLeft(*basePoints)
			weeksLeft := fmt.Sprint(*basePoints)
			if !resources.Abr.Ping(resources.ANY_PORT_WEEKS_LEFT) {
				resources.WaitGroup.Add(1)
				go resources.Abr.StartWithIcon(resources.ANY_PORT_WEEKS_LEFT, weeksLeft, resources.ANY_CMD_BROWN)
			}
			resources.Alf.PrintResult(fmt.Sprintf(resources.MSG_SET_SUCCESS, "weeks left", weeksLeft))
		}
	case resources.ACT_CUSTOM:
		t := db.NewTransaction()
		t.Add(func() error {
			getNewHabit := func() resources.Entity {
				return &resources.Habit{}
			}
			err := t.MapEntities(resources.DB_DEFAULT_HABITS_BUCKET_NAME, true, getNewHabit, func(e resources.Entity) func() {
				return func() {
					h := e.(*resources.Habit)
					if !h.Active {
						return
					}
					if h.Repetition == resources.HBT_REPETITION_DAILY {
						for h.Deadline.Before(time.Now()) {
							h.Deadline = addPeriod(resources.HBT_REPETITION_DAILY, h.Deadline)
						}
					}
				}
			})
			//habit := &resources.Habit{}
			//err = t.ModifyEntity(resources.DB_DEFAULT_HABITS_BUCKET_NAME, []byte("137"), true, habit, func () {
			//	habit.Successes = 6
			//	habit.LastStreak = 4
			//	habit.ActualStreak = 2
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

func addPeriod(repetition string, deadline *time.Time) *time.Time {
	if deadline == nil {
		return nil
	}
	switch repetition {
	case resources.HBT_REPETITION_DAILY:
		return utils.GetTimePointer(deadline.Add(24 * time.Hour))
	case resources.HBT_REPETITION_WEEKLY:
		return utils.GetTimePointer(deadline.Add(7 * 24 * time.Hour))
	case resources.HBT_REPETITION_MONTHLY:
		return utils.GetTimePointer(deadline.AddDate(0, 1, 0))
	}
	return nil
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

// Creating application support folder if it doesn't exist
func ensureAppSupportFolder(appSupportFolderPath string) error {
	_, err := os.Stat(appSupportFolderPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(appSupportFolderPath, os.FileMode(0744))
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}
