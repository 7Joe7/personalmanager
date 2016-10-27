package operations

import (
	"net/smtp"

	"github.com/7joe7/personalmanager/resources"
)

// TODO change to exportEntities - add Export method to all entities
// TODO take emails and password from DB
// TODO allow to send somebody else
func exportTasks(shoppingTasks map[string]*resources.Task) {
	var message string
	for _, task := range shoppingTasks {
		message += task.Export()
	}
	err := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", "jot.company@gmail.com", "moderator7", "smtp.gmail.com"), "jot.company@gmail.com", []string{"josef.erneker@gmail.com"}, []byte(message))
	if err != nil {
		panic(err)
	}
}

func setEmail(email string) {

}
