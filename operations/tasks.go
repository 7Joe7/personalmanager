package operations

import (
	"github.com/7joe7/personalmanager/db"
	"github.com/7joe7/personalmanager/resources"
)

func getModifyTaskFunc(t *resources.Task, name, projectId string) func () {
	return func () {
		if name != "" {
			t.Name = name
		}
		if projectId != "" {
			t.Project = &resources.Project{}
			db.RetrieveEntity(resources.DB_DEFAULT_PROJECTS_BUCKET_NAME, []byte(projectId), t.Project)
		}
	}
}

func addTask(name, projectId string) string {
	var id *string
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
		var err error
		id, err = t.AddEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, task)
		return err
	})
	t.Execute()
	return *id
}

func deleteTask(taskId string) {
	db.DeleteEntity([]byte(taskId), resources.DB_DEFAULT_TASKS_BUCKET_NAME)
}

func modifyTask(taskId, name, projectId string) {
	task := &resources.Task{}
	db.ModifyEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, []byte(taskId), task, GetModifyTaskFunc(task, name, projectId))
}

func getTask(taskId string) *resources.Task {
	task := &resources.Task{}
	db.RetrieveEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, []byte(taskId), task)
	return task
}

func getTasks() map[string]*resources.Task {
	tasks := map[string]*resources.Task{}
	db.RetrieveEntities(resources.DB_DEFAULT_TASKS_BUCKET_NAME, func (id string) interface{} {
		tasks[id] = &resources.Task{}
		return tasks[id]
	})
	return tasks
}