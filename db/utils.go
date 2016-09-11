package db

import (
	"strconv"
	"encoding/json"

	"github.com/7joe7/personalmanager/resources"
	"github.com/boltdb/bolt"
)

var (
	db *bolt.DB
)

func addEntity(entity resources.Entity, bucketName []byte) (string, error) {
	var id string
	err := db.Update(func (tx *bolt.Tx) error {
		var err error
		id, err = getAddEntityInner(entity, bucketName)(tx)
		return err
	})
	if err != nil {
		return id, err
	}
	return id, nil
}

func deleteEntity(entityId []byte, bucketName []byte) error {
	return db.Update(func (tx *bolt.Tx) error {
		return tx.Bucket(bucketName).Delete(entityId)
	})
}

func retrieveEntity(bucketName, entityId []byte, entity interface{}) error {
	return db.View(func (tx *bolt.Tx) error {
		return json.Unmarshal(tx.Bucket(bucketName).Get(entityId), entity)
	})
}

func modifyEntity(bucketName []byte, entityId []byte, entity interface{}, modify func ()) error {
	return db.Update(getModifyEntityInner(bucketName, entityId, entity, modify))
}

func retrieveEntities(bucketName []byte, getObject func (string) interface{}) error {
	return db.View(getRetrieveEntitiesInner(getObject, bucketName))
}

func mapEntities(entity interface{}, bucketName []byte, mapFunc func ()) error {
	return db.Update(getMapEntitiesInner(bucketName, entity, mapFunc))
}

func filterEntities(bucketName []byte, entity interface{}, filterFunc func () bool, copyFunc func ()) error {
	return db.View(getFilterEntitiesInner(bucketName, entity, filterFunc, copyFunc))
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