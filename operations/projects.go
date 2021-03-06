package operations

import (
	"github.com/7joe7/personalmanager/db"
	"github.com/7joe7/personalmanager/resources"
)

func getModifyProjectFunc(p *resources.Project, cmd *resources.Command, tr resources.Transaction) func() {
	return func() {
		if cmd.Name != "" {
			p.Name = cmd.Name
		}
		if cmd.TaskID != "" && cmd.TaskID != "-" {
			task := &resources.Task{}
			err := tr.ModifyEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, []byte(cmd.TaskID), true, task, func() {
				task.Project = p
			})
			if err != nil {
				panic(err)
			}
			p.Tasks = append(p.Tasks, task)
		}
		if cmd.GoalID != "" && cmd.GoalID != "-" {
			goal := &resources.Goal{}
			err := tr.ModifyEntity(resources.DB_DEFAULT_GOALS_BUCKET_NAME, []byte(cmd.GoalID), true, goal, func() {
				goal.Project = p
			})
			if err != nil {
				panic(err)
			}
			p.Goals = append(p.Goals, goal)
		}
		if cmd.ActiveFlag {
			p.Active = !p.Active
		}
		if cmd.DoneFlag {
			p.Done = !p.Done
		}
	}
}

func addProject(cmd *resources.Command) {
	tr := db.NewTransaction()
	tr.Add(func() error {
		return tr.AddEntity(resources.DB_DEFAULT_PROJECTS_BUCKET_NAME, resources.NewProject(cmd.Name))
	})
	tr.Execute()
}

func deleteProject(projectId string) {
	tr := db.NewTransaction()
	tr.Add(func() error {
		project := &resources.Project{}
		err := tr.RetrieveEntity(resources.DB_DEFAULT_PROJECTS_BUCKET_NAME, []byte(projectId), project, true)
		if err != nil {
			return err
		}
		for i := 0; i < len(project.Tasks); i++ {
			task := &resources.Task{}
			err = tr.ModifyEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, []byte(project.Tasks[i].Id), true, task, func() {
				task.Project = nil
			})
			if err != nil {
				return err
			}
		}
		for i := 0; i < len(project.Goals); i++ {
			goal := &resources.Goal{}
			err = tr.ModifyEntity(resources.DB_DEFAULT_GOALS_BUCKET_NAME, []byte(project.Goals[i].Id), true, goal, func() {
				goal.Project = nil
			})
			if err != nil {
				return err
			}
		}
		err = tr.DeleteEntity(resources.DB_DEFAULT_PROJECTS_BUCKET_NAME, []byte(projectId))
		if err != nil {
			return err
		}
		return nil
	})
	tr.Execute()
}

func modifyProject(cmd *resources.Command) {
	project := &resources.Project{}
	tr := db.NewTransaction()
	tr.Add(func() error {
		return tr.ModifyEntity(resources.DB_DEFAULT_PROJECTS_BUCKET_NAME, []byte(cmd.ID), false, project, getModifyProjectFunc(project, cmd, tr))
	})
	tr.Execute()
}

func getProject(projectId string) *resources.Project {
	project := &resources.Project{}
	tr := db.NewTransaction()
	tr.Add(func() error {
		return tr.RetrieveEntity(resources.DB_DEFAULT_PROJECTS_BUCKET_NAME, []byte(projectId), project, false)
	})
	tr.Execute()
	return project
}

func getProjects() map[string]*resources.Project {
	projects := map[string]*resources.Project{}
	db.RetrieveEntities(resources.DB_DEFAULT_PROJECTS_BUCKET_NAME, false, func(id string) resources.Entity {
		projects[id] = &resources.Project{}
		return projects[id]
	})
	return projects
}

func getActiveProjects() map[string]*resources.Project {
	return FilterProjects(func(p *resources.Project) bool { return p.Active && !p.Done })
}

func getInactiveProjects() map[string]*resources.Project {
	return FilterProjects(func(p *resources.Project) bool { return !p.Active && !p.Done })
}

func filterProjects(shallow bool, filter func(*resources.Project) bool) map[string]*resources.Project {
	projects := map[string]*resources.Project{}
	var project *resources.Project
	getNewEntity := func() resources.Entity {
		project = &resources.Project{}
		return project
	}
	addEntity := func() { projects[project.Id] = project }
	db.FilterEntities(resources.DB_DEFAULT_PROJECTS_BUCKET_NAME, shallow, addEntity, getNewEntity, func() bool { return filter(project) })
	return projects
}
