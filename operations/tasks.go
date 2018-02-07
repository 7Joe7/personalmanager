package operations

import (
	"time"

	"fmt"
	"strconv"

	"github.com/7joe7/personalmanager/db"
	"github.com/7joe7/personalmanager/resources"
	rutils "github.com/7joe7/personalmanager/resources/utils"
	"github.com/7joe7/personalmanager/utils"
)

func getModifyTaskFunc(t *resources.Task, cmd *resources.Command, status *resources.Status) func() {
	return func() {
		if cmd.Name != "" {
			t.Name = cmd.Name
		}
		if cmd.BasePoints != -1 {
			t.BasePoints = cmd.BasePoints
		}
		switch cmd.ProjectID {
		case "-":
			t.Project = nil
		case "":
		default:
			t.Project = &resources.Project{Id: cmd.ProjectID}
		}
		switch cmd.GoalID {
		case "-":
			t.Goal = nil
		case "":
		default:
			t.Goal = &resources.Goal{Id: cmd.GoalID}
		}
		if cmd.Deadline == "-" {
			t.Deadline = nil
		} else if cmd.Deadline != "" {
			t.Deadline = utils.ParseTime(resources.DATE_FORMAT, cmd.Deadline)
		}
		if cmd.Estimate == "-" {
			t.TimeEstimate = nil
		} else if cmd.Estimate != "" {
			dur, err := time.ParseDuration(cmd.Estimate)
			if err != nil {
				panic(err)
			}
			t.TimeEstimate = &dur
		}
		scheduleTask(cmd.Scheduled, t)
		if cmd.TaskType != "" {
			t.Type = cmd.TaskType
		}
		if cmd.Note != "" {
			t.Note = cmd.Note
		}
		if cmd.ActiveFlag {
			if t.InProgress {
				stopProgress(t)
			} else {
				startProgress(t)
			}
		}
		if cmd.DoneFlag {
			change := t.CountScoreChange(status)
			if t.Done {
				t.Done = false
				t.DoneTime = nil
				status.Score -= change
				status.Today -= change
			} else {
				t.Done = true
				t.DoneTime = utils.GetTimePointer(time.Now())
				if t.InProgress {
					stopProgress(t)
				}
				status.Score += change
				status.Today += change
			}
		}
	}
}

func scheduleTask(scheduled string, t *resources.Task) {
	if scheduled != "" {
		switch scheduled {
		case resources.TASK_SCHEDULED_NEXT:
			t.Scheduled = resources.TASK_SCHEDULED_NEXT
		case resources.TASK_NOT_SCHEDULED:
			t.Scheduled = resources.TASK_NOT_SCHEDULED
			if t.InProgress {
				stopProgress(t)
			}
		}
	}
}

func stopProgress(t *resources.Task) {
	t.InProgress = false
	var length float64
	if t.InProgressSince != nil {
		length = time.Now().Sub(*t.InProgressSince).Minutes()
	}
	if t.TimeSpent != nil {
		length += t.TimeSpent.Minutes()
	}
	d, err := time.ParseDuration(utils.MinutesToHMFormat(length))
	if err != nil {
		panic(err)
	}
	t.TimeSpent = &d
	resources.WaitGroup.Add(1)
	go resources.Abr.Quit(resources.ANY_PORT_ACTIVE_TASK)
}

func startProgress(t *resources.Task) {
	t.InProgress = true
	t.InProgressSince = utils.GetTimePointer(time.Now())
	resources.WaitGroup.Add(1)
	go resources.Abr.StartWithIcon(resources.ANY_PORT_ACTIVE_TASK, t.Name, resources.ANY_CMD_BLUE)
}

func getNewTask() resources.Entity {
	return &resources.Task{}
}

func getSyncTaskFunc() func(resources.Entity) func() {
	return func(entity resources.Entity) func() {
		return func() {
			t := entity.(*resources.Task)
			if t.Scheduled != resources.TASK_SCHEDULED_NEXT && t.Deadline != nil && t.Deadline.Before(time.Now().Add(time.Hour*24)) {
				scheduleTask(resources.TASK_SCHEDULED_NEXT, t)
			}
		}
	}
}

func createTask(cmd *resources.Command, goal *resources.Goal, t resources.Transaction) (*resources.Task, error) {
	task := resources.NewTask(cmd.Name)
	if cmd.ProjectID != "" && cmd.ProjectID != "-" {
		task.Project = &resources.Project{Id: cmd.ProjectID}
	}
	if goal != nil {
		task.Goal = goal
		task.BasePoints = goal.Priority
	}
	if cmd.Deadline != "" && cmd.Deadline != "-" {
		task.Deadline = utils.ParseTime(resources.DATE_FORMAT, cmd.Deadline)
	}
	if cmd.Estimate != "" {
		dur, err := time.ParseDuration(cmd.Estimate)
		if err != nil {
			return nil, err
		}
		task.TimeEstimate = &dur
	}
	if cmd.Scheduled != "" {
		task.Scheduled = cmd.Scheduled
	}
	if cmd.TaskType != "" {
		task.Type = cmd.TaskType
	}
	if cmd.Note != "" {
		task.Note = cmd.Note
	}
	if cmd.BasePoints != -1 {
		task.BasePoints = cmd.BasePoints
	}
	if cmd.ActiveFlag {
		task.InProgress = true
		task.InProgressSince = utils.GetTimePointer(time.Now())
	}
	if err := task.Load(t); err != nil {
		return nil, err
	}
	return task, nil
}

func addTask(cmd *resources.Command) string {
	var id string
	t := db.NewTransaction()
	t.Add(func() error {
		var goal *resources.Goal
		if cmd.GoalID != "-" && cmd.GoalID != "" {
			goal = &resources.Goal{}
			err := t.RetrieveEntity(resources.DB_DEFAULT_GOALS_BUCKET_NAME, []byte(cmd.GoalID), goal, false)
			if err != nil {
				return err
			}
		}
		task, err := createTask(cmd, goal, t)
		if err != nil {
			return err
		}
		err = t.AddEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, task)
		if err != nil {
			return err
		}
		if task.Project != nil {
			err = t.ModifyEntity(resources.DB_DEFAULT_PROJECTS_BUCKET_NAME, []byte(cmd.ProjectID), true, task.Project, func() {
				task.Project.Tasks = append(task.Project.Tasks, task)
			})
			if err != nil {
				return err
			}
		}
		if task.Goal != nil {
			err = t.ModifyEntity(resources.DB_DEFAULT_GOALS_BUCKET_NAME, []byte(cmd.GoalID), true, task.Goal, func() {
				task.Goal.Tasks = append(task.Goal.Tasks, task)
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	t.Execute()
	return id
}

func toggleActiveTask(t resources.Transaction, toggledTaskId string) error {
	value := t.GetValue(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_ACTUAL_ACTIVE_TASK_KEY)
	var actualActiveTask []byte
	valueStr := string(value)
	if valueStr == toggledTaskId {
		actualActiveTask = []byte{}
	} else { // deactivate old task and set new task to actualActiveTask
		actualActiveTask = []byte(toggledTaskId)
		if valueStr != "" {
			task := &resources.Task{}
			t.ModifyEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, value, true, task, func() {
				stopProgress(task)
			})
		}
	}
	return t.SetValue(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_ACTUAL_ACTIVE_TASK_KEY, actualActiveTask)
}

func setActiveTaskDone(t resources.Transaction, doneTaskId string) error {
	value := t.GetValue(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_ACTUAL_ACTIVE_TASK_KEY)
	if string(value) == doneTaskId {
		return t.SetValue(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_ACTUAL_ACTIVE_TASK_KEY, []byte{})
	}
	return nil
}

func deleteTask(taskId string) {
	t := db.NewTransaction()
	t.Add(func() error {
		task := &resources.Task{}
		err := t.RetrieveEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, []byte(taskId), task, true)
		if err != nil {
			return err
		}
		if task.Goal != nil {
			goal := &resources.Goal{}
			err = t.ModifyEntity(resources.DB_DEFAULT_GOALS_BUCKET_NAME, []byte(task.Goal.Id), true, goal, func() {
				for i := 0; i < len(goal.Tasks); i++ {
					if goal.Tasks[i].Id == task.Id {
						goal.Tasks = append(goal.Tasks[:i], goal.Tasks[i+1:]...)
						break
					}
				}
			})
			if err != nil {
				return err
			}
		}
		if task.Project != nil {
			project := &resources.Project{}
			err = t.ModifyEntity(resources.DB_DEFAULT_PROJECTS_BUCKET_NAME, []byte(task.Project.Id), true, project, func() {
				for i := 0; i < len(project.Tasks); i++ {
					if project.Tasks[i].Id == task.Id {
						project.Tasks = append(project.Tasks[:i], project.Tasks[i+1:]...)
						break
					}
				}
			})
			if err != nil {
				return err
			}
		}
		if task.InProgress {
			stopProgress(task)
			if err := t.SetValue(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_ACTUAL_ACTIVE_TASK_KEY, []byte{}); err != nil {
				return err
			}
		}
		return t.DeleteEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, []byte(taskId))
	})
	t.Execute()
}

func modifyTask(cmd *resources.Command) {
	task := &resources.Task{}
	changeStatus := &resources.Status{}
	status := &resources.Status{}
	t := db.NewTransaction()
	t.Add(func() error {
		var err error
		if cmd.ActiveFlag {
			err = toggleActiveTask(t, cmd.ID)
			if err != nil {
				return err
			}
		}
		if cmd.DoneFlag {
			err = setActiveTaskDone(t, cmd.ID)
			if err != nil {
				return err
			}
		}
		err = t.RetrieveEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, []byte(cmd.ID), task, true)
		if err != nil {
			return err
		}
		switch cmd.ProjectID {
		case "-":
			if task.Project != nil {
				err = t.ModifyEntity(resources.DB_DEFAULT_PROJECTS_BUCKET_NAME, []byte(task.Project.Id), true, task.Project, func() {
					task.Project.Tasks = rutils.RemoveTaskFromTasks(task.Project.Tasks, task)
				})
				if err != nil {
					return err
				}
			}
		case "":
		default:
			if task.Project != nil && task.Project.Id != cmd.ProjectID {
				err = t.ModifyEntity(resources.DB_DEFAULT_PROJECTS_BUCKET_NAME, []byte(task.Project.Id), true, task.Project, func() {
					task.Project.Tasks = rutils.RemoveTaskFromTasks(task.Project.Tasks, task)
				})
				if err != nil {
					return err
				}
			}
			task.Project = &resources.Project{}
			err = t.ModifyEntity(resources.DB_DEFAULT_PROJECTS_BUCKET_NAME, []byte(cmd.ProjectID), true, task.Project, func() {
				task.Project.Tasks = append(task.Project.Tasks, task)
			})
			if err != nil {
				return err
			}
		}
		switch cmd.GoalID {
		case "-":
			if task.Goal != nil {
				err = t.ModifyEntity(resources.DB_DEFAULT_GOALS_BUCKET_NAME, []byte(task.Goal.Id), true, task.Goal, func() {
					task.Goal.Tasks = rutils.RemoveTaskFromTasks(task.Goal.Tasks, task)
				})
				if err != nil {
					return err
				}
			}
		case "":
		default:
			if task.Goal != nil && task.Goal.Id != cmd.GoalID {
				err = t.ModifyEntity(resources.DB_DEFAULT_GOALS_BUCKET_NAME, []byte(task.Goal.Id), true, task.Goal, func() {
					task.Goal.Tasks = rutils.RemoveTaskFromTasks(task.Goal.Tasks, task)
				})
				if err != nil {
					return err
				}
			}
			task.Goal = &resources.Goal{}
			err = t.ModifyEntity(resources.DB_DEFAULT_GOALS_BUCKET_NAME, []byte(cmd.GoalID), true, task.Goal, func() {
				task.Goal.Tasks = append(task.Goal.Tasks, task)
			})
			if err != nil {
				return err
			}
			cmd.BasePoints = task.Goal.Priority
		}
		lastSync := string(t.GetValue(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_LAST_SYNC_KEY))
		tasksDoneToday := map[string]*resources.Task{}
		var ts *resources.Task
		getNewEntity := func() resources.Entity {
			ts = &resources.Task{}
			return ts
		}
		addEntity := func() { tasksDoneToday[task.Id] = ts }
		t.FilterEntities(resources.DB_DEFAULT_TASKS_BUCKET_NAME, true, addEntity, getNewEntity, func() bool {
			if !ts.Done || ts.DoneTime == nil {
				return false
			}
			if lastSync == "" {
				return true
			}
			lastSyncTime, err := time.Parse("Mon Jan 2 15:04:05 -0700 MST 2006", lastSync)
			if err != nil {
				panic(err)
			}
			if ts.DoneTime.Before(lastSyncTime) {
				return false
			}
			return ts.DoneTime.Before(lastSyncTime.AddDate(0, 0, 1))
		})
		for _, taskDoneToday := range tasksDoneToday {
			if taskDoneToday.TimeEstimate != nil {
				changeStatus.WorkDoneToday += int(taskDoneToday.TimeEstimate.Hours())
			}
		}
		err = t.ModifyEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, []byte(cmd.ID), false, task, getModifyTaskFunc(task, cmd, changeStatus))
		if err != nil {
			return err
		}
		err = t.ModifyEntity(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_ACTUAL_STATUS_KEY, true, status, getAddScoreFunc(status, changeStatus))
		if err != nil {
			return err
		}
		if err := t.ModifyValue(resources.DB_DEFAULT_POINTS_BUCKET_NAME, []byte(time.Now().Format("2006-01-02")), func(formerValue []byte) []byte {
			if len(formerValue) == 0 {
				return []byte(fmt.Sprint(changeStatus.Today))
			}
			before, err := strconv.Atoi(string(formerValue))
			if err != nil {
				panic(err)
			}
			return []byte(fmt.Sprint(before + changeStatus.Today))
		}); err != nil {
			return err
		}
		if err := t.ModifyValue(resources.DB_DEFAULT_POINTS_BUCKET_NAME, []byte(time.Now().AddDate(0, 0, -1).Format("2006-01-02")), func(formerValue []byte) []byte {
			if len(formerValue) == 0 {
				return []byte(fmt.Sprint(changeStatus.Yesterday))
			}
			before, err := strconv.Atoi(string(formerValue))
			if err != nil {
				panic(err)
			}
			return []byte(fmt.Sprint(before + changeStatus.Yesterday))
		}); err != nil {
			return err
		}
		return nil
	})
	t.Execute()
}

func getTask(taskId string) *resources.Task {
	task := &resources.Task{}
	t := db.NewTransaction()
	t.Add(func() error {
		return t.RetrieveEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, []byte(taskId), task, false)
	})
	t.Execute()
	return task
}

func getTasks() map[string]*resources.Task {
	tasks := map[string]*resources.Task{}
	t := db.NewTransaction()
	t.Add(func() error {
		return t.RetrieveEntities(resources.DB_DEFAULT_TASKS_BUCKET_NAME, false, func(id string) resources.Entity {
			tasks[id] = &resources.Task{}
			return tasks[id]
		})
	})
	t.Execute()
	return tasks
}

func filterTasks(shallow bool, filter func(*resources.Task) bool) map[string]*resources.Task {
	tasks := map[string]*resources.Task{}
	tr := db.NewTransaction()
	tr.Add(func() error { return filterTasksModal(tr, shallow, tasks, filter) })
	tr.Execute()
	return tasks
}

func filterTasksModal(tr resources.Transaction, shallow bool, tasks map[string]*resources.Task, filter func(*resources.Task) bool) error {
	var task *resources.Task
	getNewEntity := func() resources.Entity {
		task = &resources.Task{}
		return task
	}
	addEntity := func() { tasks[task.Id] = task }
	return tr.FilterEntities(resources.DB_DEFAULT_TASKS_BUCKET_NAME, shallow, addEntity, getNewEntity, func() bool { return filter(task) })
}

func filterTasksSlice(shallow bool, filter func(*resources.Task) bool) []*resources.Task {
	tasks := []*resources.Task{}
	var task *resources.Task
	getNewEntity := func() resources.Entity {
		task = &resources.Task{}
		return task
	}
	addEntity := func() { tasks = append(tasks, task) }
	db.FilterEntities(resources.DB_DEFAULT_TASKS_BUCKET_NAME, shallow, addEntity, getNewEntity, func() bool { return filter(task) })
	return tasks
}

func getTasksByGoal(goalId string) []*resources.Task {
	return filterTasksSlice(true, func(t *resources.Task) bool { return t.Goal.Id == goalId })
}

func getNextTasks() map[string]*resources.Task {
	return FilterTasks(func(t *resources.Task) bool {
		return !t.Done && t.Scheduled == resources.TASK_SCHEDULED_NEXT && (t.Type == "" || t.Type == resources.TASK_TYPE_PERSONAL)
	})
}

func getPersonalTasks() map[string]*resources.Task {
	return FilterTasks(func(t *resources.Task) bool { return t.Type == "" || t.Type == resources.TASK_TYPE_PERSONAL })
}

func getUnscheduledTasks() map[string]*resources.Task {
	return FilterTasks(func(t *resources.Task) bool {
		return !t.Done && (t.Scheduled == "" || t.Scheduled == resources.TASK_NOT_SCHEDULED) && (t.Type == "" || t.Type == resources.TASK_TYPE_PERSONAL)
	})
}

func getGoalTasks(id string) map[string]*resources.Task {
	return FilterTasks(func(t *resources.Task) bool {
		return !t.Done && (t.Goal != nil && t.Goal.Id == id)
	})
}

func getShoppingTasks() map[string]*resources.Task {
	return FilterTasks(func(t *resources.Task) bool { return !t.Done && (t.Type == resources.TASK_TYPE_SHOPPING) })
}

func getWorkNextTasks() map[string]*resources.Task {
	return FilterTasks(func(t *resources.Task) bool {
		return !t.Done && t.Scheduled == resources.TASK_SCHEDULED_NEXT && t.Type == resources.TASK_TYPE_WORK
	})
}

func getWorkUnscheduledTasks() map[string]*resources.Task {
	return FilterTasks(func(t *resources.Task) bool {
		return !t.Done && (t.Scheduled == "" || t.Scheduled == resources.TASK_NOT_SCHEDULED) && t.Type == resources.TASK_TYPE_WORK
	})
}
