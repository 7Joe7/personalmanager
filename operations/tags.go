package operations

import (
	"github.com/7joe7/personalmanager/db"
	"github.com/7joe7/personalmanager/resources"
)

func getModifyTagFunc(t *resources.Tag, name string) func() {
	return func() {
		if name != "" {
			t.Name = name
		}
	}
}

func AddTag(name string) {
	tr := db.NewTransaction()
	tr.Add(func() error {
		return tr.AddEntity(resources.DB_DEFAULT_TAGS_BUCKET_NAME, resources.NewTag(name))
	})
	tr.Execute()
}

func DeleteTag(tagId string) {
	db.DeleteEntity([]byte(tagId), resources.DB_DEFAULT_TAGS_BUCKET_NAME)
}

func ModifyTag(tagId, name string) {
	tag := &resources.Tag{}
	db.ModifyEntity(resources.DB_DEFAULT_TAGS_BUCKET_NAME, []byte(tagId), false, tag, getModifyTagFunc(tag, name))
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
