package operations

import (
	"github.com/7joe7/personalmanager/resources"
	"github.com/7joe7/personalmanager/db"
)

func getSyncStatusFunc(s, changeStatus *resources.Status) func () {
	return func () {
		s.Score += changeStatus.Score
		s.Today = 0
	}
}

func getAddScoreFunc(s, changeStatus *resources.Status) func () {
	return func () {
		s.Score += changeStatus.Score
		s.Today += changeStatus.Today
	}
}

func GetStatus() *resources.Status {
	status := &resources.Status{}
	tr := db.NewTransaction()
	tr.Add(func () error {
		return tr.RetrieveEntity(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_ACTUAL_STATUS_KEY, status)
	})
	tr.Execute()
	return status
}
