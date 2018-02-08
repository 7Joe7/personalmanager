package exporter

import (
	"io/ioutil"

	"encoding/json"

	"github.com/7joe7/personalmanager/db"
	"github.com/7joe7/personalmanager/operations"
	"github.com/7joe7/personalmanager/resources"
)

type exportConfig struct {
	SmtpAddress        string `json:"smtpAddress"`
	SmtpPort           int    `json:"smtpPort"`
	AdminEmailAddress  string `json:"adminEmailAddress"`
	AdminEmailPassword string `json:"adminEmailPassword"`
}

func ExportShoppingTasks() string {
	return exportTasks()
}

func SetEmail(email string) {
	setEmail(email)
}

// TODO change to exportEntities - add Export method to all entities
func exportTasks() string {
	tasks := map[string]*resources.Task{}
	var email string
	tr := db.NewTransaction()
	tr.Add(func() error {
		return operations.FilterTasksModal(tr, false, tasks, func(t *resources.Task) bool { return t.Type == resources.TASK_TYPE_SHOPPING })
	})
	tr.Add(func() error {
		email = string(tr.GetValue(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_DEFAULT_EMAIL))
		return nil
	})
	tr.Execute()
	var message string
	for _, task := range tasks {
		message += task.Export()
	}
	return message
}

func setEmail(email string) {
	tr := db.NewTransaction()
	tr.Add(func() error {
		return tr.SetValue(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_DEFAULT_EMAIL, []byte(email))
	})
	tr.Execute()
}

// temporary solution, should be done through a request to a server which will send the email,
// this way the end customer may still be easily able to sign in to the jot.company@gmail.com
func readExportConfig(address string) *exportConfig {
	exportConfigText, err := ioutil.ReadFile(address)
	if err != nil {
		panic(err)
	}
	exportConfig := &exportConfig{}
	err = json.Unmarshal(exportConfigText, exportConfig)
	if err != nil {
		panic(err)
	}
	return exportConfig
}
