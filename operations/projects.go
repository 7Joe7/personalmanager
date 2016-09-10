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

func AddProject(name string) string {
	return db.AddEntity(resources.NewProject(name), resources.DB_DEFAULT_PROJECTS_BUCKET_NAME)
}

func DeleteProject(projectId string) {
	db.DeleteEntity([]byte(projectId), resources.DB_DEFAULT_PROJECTS_BUCKET_NAME)
}

func ModifyProject(projectId, name string) {
	project := &resources.Project{}
	db.ModifyEntity(resources.DB_DEFAULT_PROJECTS_BUCKET_NAME, []byte(projectId), project, GetModifyProjectFunc(project, name))
}

func GetProject(projectId string) *resources.Project {
	project := &resources.Project{}
	db.RetrieveEntity(resources.DB_DEFAULT_PROJECTS_BUCKET_NAME, []byte(projectId), project)
	return project
}

func GetProjects() map[string]*resources.Project {
	projects := map[string]*resources.Project{}
	db.RetrieveEntities(resources.DB_DEFAULT_GOALS_BUCKET_NAME, func (id string) interface{} {
		projects[id] = &resources.Project{}
		return projects[id]
	})
	return projects
}
