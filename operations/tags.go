package operations

import (
	"github.com/7joe7/personalmanager/resources"
	"github.com/7joe7/personalmanager/db"
)

func getModifyTagFunc(t *resources.Tag, name string) func () {
	return func () {
		if name != "" {
			t.Name = name
		}
	}
}

func AddTag(name string) string {
	return db.AddEntity(resources.NewTag(name), resources.DB_DEFAULT_TAGS_BUCKET_NAME)
}

func DeleteTag(tagId string) {
	db.DeleteEntity([]byte(tagId), resources.DB_DEFAULT_TAGS_BUCKET_NAME)
}

func ModifyTag(tagId, name string) {
	tag := &resources.Tag{}
	db.ModifyEntity(resources.DB_DEFAULT_TAGS_BUCKET_NAME, []byte(tagId), tag, GetModifyTagFunc(tag, name))
}

func GetTag(tagId string) *resources.Tag {
	tag := &resources.Tag{}
	db.RetrieveEntity(resources.DB_DEFAULT_TAGS_BUCKET_NAME, []byte(tagId), tag)
	return tag
}

func GetTags() map[string]*resources.Tag {
	tags := map[string]*resources.Tag{}
	db.RetrieveEntities(resources.DB_DEFAULT_TAGS_BUCKET_NAME, func (id string) interface{} {
		tags[id] = &resources.Tag{}
		return tags[id]
	})
	return tags
}
