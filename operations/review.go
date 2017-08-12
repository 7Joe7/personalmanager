package operations

import (
	"github.com/7joe7/personalmanager/db"
	"github.com/7joe7/personalmanager/resources"
	"github.com/7joe7/personalmanager/utils"
)

func modifyReview(repetition, deadline string) {
	r := &resources.Review{}
	t := db.NewTransaction()
	t.Add(func() error {
		return t.ModifyEntity(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_REVIEW_SETTINGS_KEY, false, r, getModifyReviewFunc(r, repetition, deadline))
	})
	t.Execute()
}

func getModifyReviewFunc(r *resources.Review, repetition, deadline string) func() {
	return func() {
		if repetition != "" {
			r.Repetition = repetition
		}
		if deadline == "-" {
			r.Deadline = nil
		} else if deadline != "" {
			r.Deadline = utils.ParseTime(resources.DATE_FORMAT, deadline)
		}
	}
}

func getReview() *resources.Review {
	r := &resources.Review{}
	t := db.NewTransaction()
	t.Add(func() error {
		return t.RetrieveEntity(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_REVIEW_SETTINGS_KEY, r, false)
	})
	t.Execute()
	return r
}
