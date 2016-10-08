package db

import (
	"strconv"

	"github.com/7joe7/personalmanager/resources"
	"github.com/boltdb/bolt"
)

var (
	db *bolt.DB
)

func deleteEntity(bucketName, id []byte) error {
	tr := newTransaction()
	tr.Add(func () error { return tr.DeleteEntity(bucketName, id) })
	return tr.execute()
}

func retrieveEntity(bucketName, id []byte, entity resources.Entity, shallow bool) error {
	tr := newTransaction()
	tr.Add(func () error { return tr.RetrieveEntity(bucketName, id, entity, shallow)})
	return tr.execute()
}

func modifyEntity(bucketName, id []byte, shallow bool, entity resources.Entity, modify func ()) error {
	tr := newTransaction()
	tr.Add(func () error { return tr.ModifyEntity(bucketName, id, shallow, entity, modify)})
	return tr.execute()
}

func retrieveEntities(bucketName []byte, shallow bool, getObject func (string) resources.Entity) error {
	tr := newTransaction()
	tr.Add(func () error { return tr.RetrieveEntities(bucketName, shallow, getObject) })
	return tr.execute()
}

func mapEntities(entity resources.Entity, bucketName []byte, shallow bool, mapFunc func ()) error {
	tr := newTransaction()
	tr.Add(func () error { return tr.MapEntities(bucketName, shallow, entity, mapFunc) })
	return tr.execute()
}

func filterEntities(bucketName []byte, shallow bool, addEntity func (), getNewEntity func () resources.Entity, filterFunc func () bool) error {
	tr := newTransaction()
	tr.Add(func () error { return tr.FilterEntities(bucketName, shallow, addEntity, getNewEntity, filterFunc)})
	return tr.view()
}

func getIncrementedId(bucket *bolt.Bucket) string {
	lastId, err := strconv.Atoi(string(bucket.Get(resources.DB_LAST_ID_KEY)))
	if err != nil {
		return "0"
	}
	return strconv.Itoa(lastId + 1)
}

func open(path string) error {
	var err error
	if db, err = bolt.Open(path, 0644, nil); err != nil {
		return err
	}
	return nil
}