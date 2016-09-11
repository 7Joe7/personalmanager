package db

import (
	"strconv"

	"github.com/7joe7/personalmanager/resources"
	"github.com/boltdb/bolt"
	"encoding/json"
)

var (
	db *bolt.DB
)

func deleteEntity(bucketName, id []byte) error {
	tr := newTransaction()
	tr.Add(func () error { return tr.DeleteEntity(bucketName, id) })
	return tr.execute()
}

func retrieveEntity(bucketName, id []byte, entity resources.Entity) error {
	tr := newTransaction()
	tr.Add(func () error { return tr.RetrieveEntity(bucketName, id, entity)})
	return tr.execute()
}

func modifyEntity(bucketName []byte, id []byte, entity resources.Entity, modify func ()) error {
	tr := newTransaction()
	tr.Add(func () error { return tr.ModifyEntity(bucketName, id, entity, modify)})
	return tr.execute()
}

func retrieveEntities(bucketName []byte, getObject func (string) resources.Entity) error {
	tr := newTransaction()
	tr.Add(func () error { return tr.RetrieveEntities(bucketName, getObject) })
	return tr.execute()
}

func mapEntities(entity resources.Entity, bucketName []byte, mapFunc func ()) error {
	tr := newTransaction()
	tr.Add(func () error { return tr.MapEntities(bucketName, entity, mapFunc) })
	return tr.execute()
}

func filterEntities(bucketName []byte, entity resources.Entity, filterFunc func () bool, copyFunc func ()) error {
	tr := newTransaction()
	tr.Add(func () error { return tr.FilterEntities(bucketName, entity, filterFunc, copyFunc)})
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

func modifyEntityInner(bucket *bolt.Bucket, key, value []byte, entity interface{}, modify func ()) error {
	if err := json.Unmarshal(value, entity); err != nil {
		return err
	}
	modify()
	resultValue, err := json.Marshal(entity)
	if err != nil {
		return err
	}
	return bucket.Put(key, resultValue)
}