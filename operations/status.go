package operations

import (
	"github.com/7joe7/personalmanager/resources"
	"github.com/7joe7/personalmanager/db"
)

func getSyncStatusFunc(s *resources.Status, scoreChange int) func () {
	return func () {
		s.Score += scoreChange
		s.Today = 0
	}
}

func getAddScoreFunc(s *resources.Status, scoreChange int) func () {
	return func () {
		s.Score += scoreChange
		s.Today += scoreChange
	}
}

func GetStatus() *resources.Status {
	status := &resources.Status{}
	db.RetrieveEntity(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_ACTUAL_STATUS_KEY, status)
	return status
}
