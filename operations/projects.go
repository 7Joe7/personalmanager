package operations

import (
	"github.com/7joe7/personalmanager/resources"
	"github.com/7joe7/personalmanager/db"
)

func getModifyProjectFunc(p *resources.Project, name string) func () {
	return func () {
		if name != "" {
			p.Name = name
		}
	}
}

func AddProject(name string) {
	tr := db.NewTransaction()
	tr.Add(func () error {
		return tr.AddEntity(resources.DB_DEFAULT_PROJECTS_BUCKET_NAME, resources.NewProject(name))
	})
	tr.Execute()
}

func DeleteProject(projectId string) {
	db.DeleteEntity(resources.DB_DEFAULT_PROJECTS_BUCKET_NAME, []byte(projectId))
}

func ModifyProject(projectId, name string) {
	project := &resources.Project{}
	db.ModifyEntity(resources.DB_DEFAULT_PROJECTS_BUCKET_NAME, []byte(projectId), false, project, getModifyProjectFunc(project, name))
}

func GetProject(projectId string) *resources.Project {
	project := &resources.Project{}
	tr := db.NewTransaction()
	tr.Add(func () error {
		return tr.RetrieveEntity(resources.DB_DEFAULT_PROJECTS_BUCKET_NAME, []byte(projectId), project, false)
	})
	tr.Execute()
	return project
}

func GetProjects() map[string]*resources.Project {
	projects := map[string]*resources.Project{}
	db.RetrieveEntities(resources.DB_DEFAULT_PROJECTS_BUCKET_NAME, false, func (id string) resources.Entity {
		projects[id] = &resources.Project{}
		return projects[id]
	})
	return projects
}
