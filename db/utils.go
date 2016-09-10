package db

import (
	"strconv"
	"encoding/json"
	"log"

	"github.com/7joe7/personalmanager/resources"
	"github.com/boltdb/bolt"
)

var (
	db *bolt.DB
)

func addEntity(entity interface{}, bucketName []byte) string {
	var incrementedId *string
	if err := db.Update(getAddEntityInner(entity, bucketName, incrementedId)); err != nil {
		log.Fatalf("Unable to add entity to bucket '%s'. %v", bucketName, err)
	}
	return *incrementedId
}

func deleteEntity(entityId []byte, bucketName []byte) {
	err := db.Update(func (tx *bolt.Tx) error {
		return tx.Bucket(bucketName).Delete(entityId)
	})
	if err != nil {
		log.Fatalf("Unable to delete entity id '%s' from bucket '%s'. %v", string(entityId), bucketName, err)
	}
}

func retrieveEntity(bucketName, entityId []byte, entity interface{}) {
	err := db.View(func (tx *bolt.Tx) error {
		return json.Unmarshal(tx.Bucket(bucketName).Get(entityId), entity)
	})
	if err != nil {
		log.Fatalf("Unable to retrieve entity id '%s' from bucket '%s'. %v", string(entityId), bucketName, err)
	}
}

func modifyEntity(bucketName []byte, entityId []byte, entity interface{}, modify func ()) {
	if err := db.Update(getModifyEntityInner(bucketName, entityId, entity, modify)); err != nil {
		log.Fatalf("Unable to modify entity id '%s' in bucket '%s'. %v", string(entityId), bucketName, err)
	}
}

func retrieveEntities(bucketName []byte, getObject func (string) interface{}) {
	if err := db.View(getRetrieveEntitiesInner(getObject, bucketName)); err != nil {
		log.Fatalf("Unable to retrieve entities from bucket '%s'. %v", bucketName, err)
	}
}

func mapEntities(mapFunc func (), entity interface{}, bucketName []byte) {
	if err := db.Update(getMapEntitiesInner(bucketName, entity, mapFunc)); err != nil {
		log.Fatalf("Unable to map entities from bucket '%s'. %v", bucketName, err)
	}
}

func filterEntities(bucketName []byte, entity interface{}, filterFunc func (string)) {
	if err := db.View(getFilterEntitiesInner(bucketName, entity, filterFunc)); err != nil {
		log.Fatalf("Unable to retrieve entities from bucket '%s'. %v", bucketName, err)
	}
}

func getIncrementedId(bucket *bolt.Bucket) string {
	lastId, err := strconv.Atoi(string(bucket.Get(resources.DB_LAST_ID_KEY)))
	if err != nil {
		return "0"
	}
	return strconv.Itoa(lastId + 1)
}

func open(path string) {
	var err error
	if db, err = bolt.Open(path, 0644, nil); err != nil {
		log.Fatalf("Unable to open database '%s'. %v", path, err)
	}
}

func executeMultiple(updates ...func (*bolt.Tx) error) {
	err := db.Update(func (tx *bolt.Tx) error {
		for i := 0; i < len(updates); i++ {
			if err := updates[i](tx); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Unable to modify database data. %v", err)
	}
}