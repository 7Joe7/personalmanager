package main

import (
	"flag"
	"fmt"

	"github.com/7joe7/personalmanager/alfred"
	"github.com/7joe7/personalmanager/db"
	"github.com/7joe7/personalmanager/resources"
	"log"
	"os"
)

var (
	actions = []string{
		"create-task", "print-tasks", "delete-task", "modify-task",
		"create-project", "print-projects", "delete-project", "modify-project",
		"create-tag", "print-tags", "delete-tag", "modify-tag",
		"create-goal", "print-goals", "delete-goal", "modify-goal",
		"create-habit", "print-habits", "delete-habit", "modify-habit"}
	action, id, name, projectId, repetition *string
	noneAllowed, activeFlag, doneFlag       *bool
)

func init() {
	action = flag.String("action", "", fmt.Sprintf("Provide action to be taken from this list: %v.", actions))
	id = flag.String("id", "", fmt.Sprintf("Provide id of the entity you want to make the action for. Valid for these actions: ."))
	projectId = flag.String("projectId", "", fmt.Sprintf("Provide project id for task assignment."))
	name = flag.String("name", "", "Provide name.")
	activeFlag = flag.Bool("active", false, "Toggle active/show active only.")
	doneFlag = flag.Bool("done", false, "Toggle done.")
	repetition = flag.String("repetition", "", "Select repetition period.")
	noneAllowed = flag.Bool("noneAllowed", false, "Provide information whether list should be retrieved with none value allowed.")

	db.Open()
	db.InitializeBuckets()
	db.Synchronize()

	f, err := os.OpenFile(resources.LOG_FILE_PATH, os.O_APPEND|os.O_CREATE, 777)
	if err != nil {
		log.Fatalf("Unable to open log file. %v", err)
	}
	log.SetOutput(f)
}

func main() {
	flag.Parse()
	switch *action {
	case actions[0]:
		alfred.PrintResult(fmt.Sprintf(resources.MSG_CREATE_SUCCESS, "task",
			db.AddTask(resources.NewTask(*name, db.GetProject(*projectId)))))
	case actions[4]:
		alfred.PrintResult(fmt.Sprintf(resources.MSG_CREATE_SUCCESS, "project",
			db.AddProject(resources.NewProject(*name))))
	case actions[8]:
		alfred.PrintResult(fmt.Sprintf(resources.MSG_CREATE_SUCCESS, "tag",
			db.AddTag(resources.NewTag(*name))))
	case actions[12]:
		alfred.PrintResult(fmt.Sprintf(resources.MSG_CREATE_SUCCESS, "goal",
			db.AddGoal(resources.NewGoal(*name))))
	case actions[16]:
		alfred.PrintResult(fmt.Sprintf(resources.MSG_CREATE_SUCCESS, "habit",
			db.AddHabit(resources.NewHabit(*name))))
	case actions[1]:
		alfred.PrintEntities(resources.Tasks{db.GetTasks(), *noneAllowed})
	case actions[5]:
		alfred.PrintEntities(resources.Projects{db.GetProjects(), *noneAllowed})
	case actions[9]:
		alfred.PrintEntities(resources.Tags{db.GetTags(), *noneAllowed})
	case actions[13]:
		alfred.PrintEntities(resources.Goals{db.GetGoals(), *noneAllowed})
	case actions[17]:
		if *activeFlag {
			alfred.PrintEntities(resources.Habits{db.GetActiveHabits(), *noneAllowed})
		} else {
			alfred.PrintEntities(resources.Habits{db.GetHabits(), *noneAllowed})
		}
	case actions[2]:
		db.DeleteTask(*id)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_DELETE_SUCCESS, "task"))
	case actions[6]:
		db.DeleteProject(*id)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_DELETE_SUCCESS, "project"))
	case actions[10]:
		db.DeleteTag(*id)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_DELETE_SUCCESS, "tag"))
	case actions[14]:
		db.DeleteGoal(*id)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_DELETE_SUCCESS, "goal"))
	case actions[18]:
		db.DeleteHabit(*id)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_DELETE_SUCCESS, "habit"))
	case actions[3]:
		db.ModifyTask(*id, *name, *projectId)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_MODIFY_SUCCESS, "task"))
	case actions[7]:
		db.ModifyProject(*id, *name)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_MODIFY_SUCCESS, "project"))
	case actions[11]:
		db.ModifyTag(*id, *name)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_MODIFY_SUCCESS, "tag"))
	case actions[15]:
		db.ModifyGoal(*id, *name)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_MODIFY_SUCCESS, "goal"))
	case actions[19]:
		db.ModifyHabit(*id, *name, *repetition, *activeFlag, *doneFlag)
		alfred.PrintResult(fmt.Sprintf(resources.MSG_MODIFY_SUCCESS, "habit"))
	default:
		flag.Usage()
	}
}
