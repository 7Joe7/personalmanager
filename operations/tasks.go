package operations

import (
	"time"

	"github.com/7joe7/personalmanager/anybar"
	"github.com/7joe7/personalmanager/db"
	"github.com/7joe7/personalmanager/resources"
	"github.com/7joe7/personalmanager/utils"
)

func getModifyTaskFunc(t *resources.Task, name, projectId, goalId, deadline, estimate, scheduled, taskType, note string, basePoints int, activeFlag, doneFlag bool, status *resources.Status) func() {
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
		if goalId != "" {
			t.Goal = &resources.Goal{Id: goalId}
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
	if t.TimeEstimate != nil {
		change += int(t.TimeEstimate.Minutes()) * t.BasePoints
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

func createTask(name, projectId, goalId, deadline, estimate, scheduled, taskType, note string, active bool, basePoints int, t resources.Transaction) (*resources.Task, error) {
	task := resources.NewTask(name)
	if projectId != "" {
		task.Project = &resources.Project{Id: projectId}
	}
	if goalId != "" {
		task.Goal = &resources.Goal{Id: goalId}
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
	if note != "" {
		task.Note = note
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

func addTask(name, projectId, goalId, deadline, estimate, scheduled, taskType, note string, active bool, basePoints int) string {
	var id string
	t := db.NewTransaction()
	t.Add(func() error {
		task, err := createTask(name, projectId, goalId, deadline, estimate, scheduled, taskType, note, active, basePoints, t)
		if err != nil {
			return err
		}
		err = t.AddEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, task)
		if err != nil {
			return err
		}
		if task.Goal != nil {
			return t.ModifyEntity(resources.DB_DEFAULT_GOALS_BUCKET_NAME, []byte(goalId), true, task.Goal, func () {
				task.Goal.Tasks = append(task.Goal.Tasks, task)
			})
		}
		return nil
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
			t.ModifyEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, value, true, task, func() {
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
		err := t.RetrieveEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, []byte(taskId), task, true)
		if err != nil {
			return err
		}
		if task.Goal != nil {
			goal := &resources.Goal{}
			err = t.ModifyEntity(resources.DB_DEFAULT_GOALS_BUCKET_NAME, []byte(task.Goal.Id), true, goal, func() {
				for i := 0; i < len(goal.Tasks); i++ {
					if goal.Tasks[i].Id == task.Goal.Id {
						goal.Tasks = append(goal.Tasks[:i], goal.Tasks[i+1:]...)
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
					if project.Tasks[i].Id == task.Goal.Id {
						project.Tasks = append(project.Tasks[:i], project.Tasks[i+1:]...)
					}
				}
			})
			if err != nil {
				return err
			}
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

func modifyTask(taskId, name, projectId, goalId, deadline, estimate, scheduled, taskType, note string, basePoints int, activeFlag, doneFlag bool) {
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
		err := t.ModifyEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, []byte(taskId), false, task, getModifyTaskFunc(task, name, projectId, goalId, deadline, estimate, scheduled, taskType, note, basePoints, activeFlag, doneFlag, changeStatus))
		if err != nil {
			return err
		}
		if projectId != "" {
			return t.ModifyEntity(resources.DB_DEFAULT_PROJECTS_BUCKET_NAME, []byte(projectId), true, task.Project, func () {
				task.Project.Tasks = append(task.Project.Tasks, task)
			})
		}
		if goalId != "" {
			return t.ModifyEntity(resources.DB_DEFAULT_GOALS_BUCKET_NAME, []byte(goalId), true, task.Goal, func () {
				task.Goal.Tasks = append(task.Goal.Tasks, task)
			})
		}
		return t.ModifyEntity(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_ACTUAL_STATUS_KEY, true, status, getAddScoreFunc(status, changeStatus))
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
	tr.Add(func () error { return filterTasksModal(tr, shallow, tasks, filter) })
	tr.Execute()
	return tasks
}

func filterTasksModal(tr resources.Transaction, shallow bool, tasks map[string]*resources.Task, filter func (*resources.Task) bool) error {
	var task *resources.Task
	getNewEntity := func () resources.Entity {
		task = &resources.Task{}
		return task
	}
	addEntity := func () { tasks[task.Id] = task }
	return tr.FilterEntities(resources.DB_DEFAULT_TASKS_BUCKET_NAME, shallow, addEntity, getNewEntity, func() bool { return filter(task) })
}

func filterTasksSlice(shallow bool, filter func(*resources.Task) bool) []*resources.Task {
	tasks := []*resources.Task{}
	var task *resources.Task
	getNewEntity := func () resources.Entity {
		task = &resources.Task{}
		return task
	}
	addEntity := func () { tasks = append(tasks, task) }
	db.FilterEntities(resources.DB_DEFAULT_TASKS_BUCKET_NAME, shallow, addEntity, getNewEntity, func() bool { return filter(task) })
	return tasks
}

func getTasksByGoal(goalId string) []*resources.Task {
	return filterTasksSlice(true, func (t *resources.Task) bool { return t.Goal.Id == goalId })
}

func getNextTasks() map[string]*resources.Task {
	return FilterTasks(func (t *resources.Task) bool { return t.Scheduled == resources.TASK_SCHEDULED_NEXT && (t.Type == "" || t.Type == resources.TASK_TYPE_PERSONAL) })
}

func getPersonalTasks() map[string]*resources.Task {
	return FilterTasks(func (t *resources.Task) bool { return t.Type == "" || t.Type == resources.TASK_TYPE_PERSONAL })
}

func getUnscheduledTasks() map[string]*resources.Task {
	return FilterTasks(func (t *resources.Task) bool { return (t.Scheduled == "" || t.Scheduled == resources.TASK_NOT_SCHEDULED) && (t.Type == "" || t.Type == resources.TASK_TYPE_PERSONAL) })
}

func getShoppingTasks() map[string]*resources.Task {
	return FilterTasks(func (t *resources.Task) bool { return (t.Type == resources.TASK_TYPE_SHOPPING )})
}

func getWorkNextTasks() map[string]*resources.Task {
	return FilterTasks(func (t *resources.Task) bool { return t.Scheduled == resources.TASK_SCHEDULED_NEXT && t.Type == resources.TASK_TYPE_WORK })
}

func getWorkUnscheduledTasks() map[string]*resources.Task {
	return FilterTasks(func (t *resources.Task) bool { return (t.Scheduled == "" || t.Scheduled == resources.TASK_NOT_SCHEDULED) && t.Type == resources.TASK_TYPE_WORK })
}
