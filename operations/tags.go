package operations

import (
	"github.com/7joe7/personalmanager/db"
	"github.com/7joe7/personalmanager/resources"
)

func getModifyTagFunc(t *resources.Tag, cmd *resources.Command) func() {
	return func() {
		if cmd.Name != "" {
			t.Name = cmd.Name
		}
	}
}

func AddTag(cmd *resources.Command) {
	tr := db.NewTransaction()
	tr.Add(func() error {
		return tr.AddEntity(resources.DB_DEFAULT_TAGS_BUCKET_NAME, resources.NewTag(cmd.Name))
	})
	tr.Execute()
}

func DeleteTag(tagId string) {
	db.DeleteEntity([]byte(tagId), resources.DB_DEFAULT_TAGS_BUCKET_NAME)
}

func ModifyTag(cmd *resources.Command) {
	tag := &resources.Tag{}
	db.ModifyEntity(resources.DB_DEFAULT_TAGS_BUCKET_NAME, []byte(cmd.ID), false, tag, getModifyTagFunc(tag, cmd))
}

func GetTag(tagId string) *resources.Tag {
	tag := &resources.Tag{}
	tr := db.NewTransaction()
	tr.Add(func() error {
		return tr.RetrieveEntity(resources.DB_DEFAULT_TAGS_BUCKET_NAME, []byte(tagId), tag, false)
	})
	tr.Execute()
	return tag
}

func GetTags() map[string]*resources.Tag {
	tags := map[string]*resources.Tag{}
	db.RetrieveEntities(resources.DB_DEFAULT_TAGS_BUCKET_NAME, false, func(id string) resources.Entity {
		tags[id] = &resources.Tag{}
		return tags[id]
	})
	return tags
}
