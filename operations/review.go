package operations

import (
	"time"

	"github.com/7joe7/personalmanager/resources"
	"github.com/7joe7/personalmanager/db"
)

func modifyReview(repetition, deadline string) {
	r := &resources.Review{}
	t := db.NewTransaction()
	t.Add(func () error {
		return t.ModifyEntity(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_REVIEW_SETTINGS_KEY, r, getModifyReviewFunc(r, repetition, deadline))
	})
	t.Execute()
}

func getModifyReviewFunc(r *resources.Review, repetition, deadline string) func () {
	return func () {
		if repetition != "" {
			r.Repetition = repetition
		}
		if deadline != "" {
			d, err := time.Parse(resources.DATE_FORMAT, deadline)
			if err != nil {
				panic(err)
			}
			r.Deadline = &d
		}
	}
}

func getReview() *resources.Review {
	r := &resources.Review{}
	t := db.NewTransaction()
	t.Add(func () error {
		return t.RetrieveEntity(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_REVIEW_SETTINGS_KEY, r)
	})
	t.Execute()
	return r
}
