package operations

import (
	"net/smtp"

	"github.com/7joe7/personalmanager/resources"
	"github.com/7joe7/personalmanager/db"
)

// TODO change to exportEntities - add Export method to all entities
// TODO take emails and password from DB
// TODO allow to send somebody else
func exportTasks() {
	tasks := map[string]*resources.Task{}
	var email string
	tr := db.NewTransaction()
	tr.Add(func () error { return filterTasksModal(tr, false, tasks, func (t *resources.Task) bool { return t.Type == resources.TASK_TYPE_SHOPPING }) })
	tr.Add(func () error {
		email = string(tr.GetValue(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_DEFAULT_EMAIL))
		return nil
	})
	tr.Execute()
	var message string
	for _, task := range tasks {
		message += task.Export()
	}
	err := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", "jot.company@gmail.com", "moderator7", "smtp.gmail.com"), "jot.company@gmail.com", []string{email}, []byte(message))
	if err != nil {
		panic(err)
	}
}

func setEmail(email string) {
	tr := db.NewTransaction()
	tr.Add(func () error {
		return tr.SetValue(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_DEFAULT_EMAIL, []byte(email))
	})
	tr.Execute()
}
