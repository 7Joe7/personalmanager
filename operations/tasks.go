package operations

import (
	"github.com/7joe7/personalmanager/db"
	"github.com/7joe7/personalmanager/resources"
)

func getModifyTaskFunc(t *resources.Task, name, projectId string, tr resources.Transaction) func () {
	return func () {
		if name != "" {
			t.Name = name
		}
		if projectId != "" {
			t.Project = &resources.Project{}
			tr.RetrieveEntity(resources.DB_DEFAULT_PROJECTS_BUCKET_NAME, []byte(projectId), t.Project)
		}
	}
}

func addTask(name, projectId string) string {
	var id string
	t := db.NewTransaction()
	t.Add(func () error {
		var project *resources.Project
		if projectId != "" {
			project = &resources.Project{}
			if err := t.RetrieveEntity(resources.DB_DEFAULT_PROJECTS_BUCKET_NAME, []byte(projectId), project); err != nil {
				return err
			}
		}
		task := resources.NewTask(name, project)
		return t.AddEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, task)
	})
	t.Execute()
	return id
}

func deleteTask(taskId string) {
	t := db.NewTransaction()
	t.Add(func () error {
		return t.DeleteEntity([]byte(taskId), resources.DB_DEFAULT_TASKS_BUCKET_NAME)
	})
	t.Execute()
}

func modifyTask(taskId, name, projectId string) {
	task := &resources.Task{}
	t := db.NewTransaction()
	t.Add(func () error {
		return t.ModifyEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, []byte(taskId), task, getModifyTaskFunc(task, name, projectId, t))
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
	return tasks
}