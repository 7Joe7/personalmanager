package configuration

import (
	"fmt"

	"github.com/7joe7/personalmanager/db"
	"github.com/7joe7/personalmanager/resources"
)

func SetWeeksLeft(weeksLeft int) {
	tr := db.NewTransaction()
	tr.Add(func() error {
		return tr.SetValue(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_DEFAULT_EMAIL, []byte(fmt.Sprint(weeksLeft)))
	})
	tr.Execute()
}
