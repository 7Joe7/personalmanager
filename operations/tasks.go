package operations

import (
	"time"

	"github.com/7joe7/personalmanager/resources"
	"github.com/7joe7/personalmanager/db"
	"github.com/7joe7/personalmanager/utils"
)

func getModifyTaskFunc(t *resources.Task, name, projectId, deadline, estimate string, basePoints int, activeFlag, doneFlag bool, status *resources.Status) func () {
	return func () {
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
		if activeFlag {
			if t.InProgress {
				stopProgress(t)
			} else {
				t.InProgress = true
				t.InProgressSince = utils.GetTimePointer(time.Now())
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
				stopProgress(t)
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
	d, err := time.ParseDuration(utils.MinutesToHMFormat(t.InProgressSince.Sub(time.Now()).Minutes() + t.TimeSpent.Minutes()))
	if err != nil {
		panic(err)
	}
	t.TimeSpent = &d
}

func addTask(name, projectId, deadline, estimate string, active bool) string {
	var id string
	t := db.NewTransaction()
	t.Add(func () error {
		task := resources.NewTask(name)
		if projectId != "" {
			task.Project = &resources.Project{Id:projectId}
		}
		if deadline != "" {
			task.Deadline = utils.ParseTime(resources.DATE_FORMAT, deadline)
		}
		if estimate != "" {
			dur, err := time.ParseDuration(estimate)
			if err != nil {
				return err
			}
			task.TimeEstimate = &dur
		}
		if active {
			task.InProgress = true
			task.InProgressSince = utils.GetTimePointer(time.Now())
		}
		if err := task.Load(t); err != nil {
			return err
		}
		return t.AddEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, task)
	})
	t.Execute()
	return id
}

func deleteTask(taskId string) {
	t := db.NewTransaction()
	t.Add(func () error {
		return t.DeleteEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, []byte(taskId))
	})
	t.Execute()
}

func modifyTask(taskId, name, projectId, deadline, estimate string, basePoints int, activeFlag, doneFlag bool) {
	task := &resources.Task{}
	changeStatus := &resources.Status{}
	status := &resources.Status{}
	t := db.NewTransaction()
	t.Add(func () error {
		err := t.ModifyEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, []byte(taskId), task, getModifyTaskFunc(task, name, projectId, deadline, estimate, basePoints, activeFlag, doneFlag, changeStatus))
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
	t.Add(func () error {
		return t.RetrieveEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, []byte(taskId), task)
	})
	t.Execute()
	return task
}

func getTasks() map[string]*resources.Task {
	tasks := map[string]*resources.Task{}
	t := db.NewTransaction()
	t.Add(func () error {
		return t.RetrieveEntities(resources.DB_DEFAULT_TASKS_BUCKET_NAME, func (id string) resources.Entity {
			tasks[id] = &resources.Task{}
			return tasks[id]
		})
	})
	t.Execute()
	return tasks
}