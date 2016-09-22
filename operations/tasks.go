package operations

import (
	"time"

	"github.com/7joe7/personalmanager/anybar"
	"github.com/7joe7/personalmanager/db"
	"github.com/7joe7/personalmanager/resources"
	"github.com/7joe7/personalmanager/utils"
)

func getModifyTaskFunc(t *resources.Task, name, projectId, deadline, estimate, scheduled, taskType, note string, basePoints int, activeFlag, doneFlag bool, status *resources.Status) func() {
	return func() {
		if name != "" {
			t.Name = name
		}
		if basePoints != -1 {
			t.BasePoints = basePoints
		}
		if projectId != "" {
			t.Project = &resources.Project{Id: projectId}
		}
		if deadline != "" {
			t.Deadline = utils.ParseTime(resources.DATE_FORMAT, deadline)
		}
		if estimate != "" {
			dur, err := time.ParseDuration(estimate)
			if err != nil {
				panic(err)
			}
			t.TimeEstimate = &dur
		}
		if scheduled != "" {
			t.Scheduled = scheduled
		}
		if taskType != "" {
			t.Type = taskType
		}
		if note != "" {
			t.Note = note
		}
		if activeFlag {
			if t.InProgress {
				stopProgress(t)
			} else {
				startProgress(t)
			}
		}
		if doneFlag {
			change := countScoreChange(t)
			if t.Done {
				t.Done = false
				status.Score -= change
				status.Today -= change
			} else {
				t.Done = true
				if t.InProgress {
					stopProgress(t)
				}
				status.Score += change
				status.Today += change
			}
		}
	}
}

func countScoreChange(t *resources.Task) int {
	change := t.BasePoints * 10
	if t.TimeSpent != nil {
		change += int(t.TimeSpent.Minutes()) * t.BasePoints
	}
	return change
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
	go anybar.Quit(resources.ANY_PORT_ACTIVE_TASK)
}

func startProgress(t *resources.Task) {
	t.InProgress = true
	t.InProgressSince = utils.GetTimePointer(time.Now())
	resources.WaitGroup.Add(1)
	go anybar.StartWithIcon(resources.ANY_PORT_ACTIVE_TASK, t.Name, resources.ANY_CMD_BLUE)
}

func createTask(name, projectId, deadline, estimate, scheduled, taskType string, active bool, basePoints int, t resources.Transaction) (*resources.Task, error) {
	task := resources.NewTask(name)
	if projectId != "" {
		task.Project = &resources.Project{Id: projectId}
	}
	if deadline != "" {
		task.Deadline = utils.ParseTime(resources.DATE_FORMAT, deadline)
	}
	if estimate != "" {
		dur, err := time.ParseDuration(estimate)
		if err != nil {
			return nil, err
		}
		task.TimeEstimate = &dur
	}
	if scheduled != "" {
		task.Scheduled = scheduled
	}
	if taskType != "" {
		task.Type = taskType
	}
	if basePoints != -1 {
		task.BasePoints = basePoints
	}
	if active {
		task.InProgress = true
		task.InProgressSince = utils.GetTimePointer(time.Now())
	}
	if err := task.Load(t); err != nil {
		return nil, err
	}
	return task, nil
}

func addTask(name, projectId, deadline, estimate, scheduled, taskType string, active bool, basePoints int) string {
	var id string
	t := db.NewTransaction()
	t.Add(func() error {
		task, err := createTask(name, projectId, deadline, estimate, scheduled, taskType, active, basePoints, t)
		if err != nil {
			return err
		}
		return t.AddEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, task)
	})
	t.Execute()
	return id
}

func moveActiveTask(t resources.Transaction, toggledTaskId string) error {
	value := t.GetValue(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_ACTUAL_ACTIVE_TASK_KEY)
	var actualActiveTask []byte
	valueStr := string(value)
	if valueStr == toggledTaskId {
		actualActiveTask = []byte{}
	} else {
		actualActiveTask = []byte(toggledTaskId)
		if valueStr != "" {
			task := &resources.Task{}
			t.ModifyEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, value, task, func() {
				stopProgress(task)
			})
		}
	}
	return t.SetValue(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_ACTUAL_ACTIVE_TASK_KEY, actualActiveTask)
}

func deleteTask(taskId string) {
	t := db.NewTransaction()
	t.Add(func() error {
		task := &resources.Task{}
		if err := t.RetrieveEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, []byte(taskId), task); err != nil {
			return err
		}
		if task.InProgress {
			if err := t.SetValue(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_ACTUAL_ACTIVE_TASK_KEY, []byte{}); err != nil {
				return err
			}
		}
		return t.DeleteEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, []byte(taskId))
	})
	t.Execute()
}

func modifyTask(taskId, name, projectId, deadline, estimate, scheduled, taskType, note string, basePoints int, activeFlag, doneFlag bool) {
	task := &resources.Task{}
	changeStatus := &resources.Status{}
	status := &resources.Status{}
	t := db.NewTransaction()
	t.Add(func() error {
		if activeFlag {
			err := moveActiveTask(t, taskId)
			if err != nil {
				return err
			}
		}
		err := t.ModifyEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, []byte(taskId), task, getModifyTaskFunc(task, name, projectId, deadline, estimate, scheduled, taskType, note, basePoints, activeFlag, doneFlag, changeStatus))
		if err != nil {
			return err
		}
		return t.ModifyEntity(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_ACTUAL_STATUS_KEY, status, getAddScoreFunc(status, changeStatus))
	})
	t.Execute()
}

func getTask(taskId string) *resources.Task {
	task := &resources.Task{}
	t := db.NewTransaction()
	t.Add(func() error {
		return t.RetrieveEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, []byte(taskId), task)
	})
	t.Execute()
	return task
}

func getTasks() map[string]*resources.Task {
	tasks := map[string]*resources.Task{}
	t := db.NewTransaction()
	t.Add(func() error {
		return t.RetrieveEntities(resources.DB_DEFAULT_TASKS_BUCKET_NAME, func(id string) resources.Entity {
			tasks[id] = &resources.Task{}
			return tasks[id]
		})
	})
	t.Execute()
	return tasks
}

func filterTasks(filter func(*resources.Task) bool) map[string]*resources.Task {
	tasks := map[string]*resources.Task{}
	var task *resources.Task
	getNewEntity := func () resources.Entity {
		task = &resources.Task{}
		return task
	}
	addEntity := func () { tasks[task.Id] = task }
	db.FilterEntities(resources.DB_DEFAULT_TASKS_BUCKET_NAME, addEntity, getNewEntity, func() bool { return filter(task) })
	return tasks
}

func getNextTasks() map[string]*resources.Task {
	return FilterTasks(func (t *resources.Task) bool { return t.Scheduled == resources.TASK_SCHEDULED_NEXT && (t.Type == "" || t.Type == resources.TASK_TYPE_PERSONAL) })
}

func getUnscheduledTasks() map[string]*resources.Task {
	return FilterTasks(func (t *resources.Task) bool { return (t.Scheduled == "" || t.Scheduled == resources.TASK_NOT_SCHEDULED) && (t.Type == "" || t.Type == resources.TASK_TYPE_PERSONAL) })
}

func getWorkNextTasks() map[string]*resources.Task {
	return FilterTasks(func (t *resources.Task) bool { return t.Scheduled == resources.TASK_SCHEDULED_NEXT && t.Type == resources.TASK_TYPE_WORK })
}

func getWorkUnscheduledTasks() map[string]*resources.Task {
	return FilterTasks(func (t *resources.Task) bool { return (t.Scheduled == "" || t.Scheduled == resources.TASK_NOT_SCHEDULED) && t.Type == resources.TASK_TYPE_WORK })
}
