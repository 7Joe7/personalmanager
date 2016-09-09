package db

import (
	"strconv"
	"encoding/json"

	"github.com/7joe7/personalmanager/resources"

	"github.com/boltdb/bolt"
	"log"
)

var (
	db *bolt.DB
)

func addEntity(entity interface{}, bucketName []byte) string {
	var incrementedId string
	err := db.Update(func (tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		incrementedId = getIncrementedId(bucket)
		value, err := json.Marshal(entity)
		if err != nil {
			return err
		}
		incrementedIdBytes := []byte(incrementedId)
		err = bucket.Put(incrementedIdBytes, value)
		if err != nil {
			return err
		}
		return bucket.Put(resources.DB_LAST_ID_KEY, incrementedIdBytes)
	})
	if err != nil {
		log.Fatalf("Unable to add entity to bucket '%s'. %v", bucketName, err)
	}
	return incrementedId
}

func deleteEntity(entityId string, bucketName []byte) {
	err := db.Update(func (tx *bolt.Tx) error {
		return tx.Bucket(bucketName).Delete([]byte(entityId))
	})
	if err != nil {
		log.Fatalf("Unable to delete entity id '%s' from bucket '%s'. %v", entityId, bucketName, err)
	}
}

func retrieveEntity(entityId string, entity interface{}, bucketName []byte) {
	err := db.View(func (tx *bolt.Tx) error {
		return json.Unmarshal(tx.Bucket(bucketName).Get([]byte(entityId)), entity)
	})
	if err != nil {
		log.Fatalf("Unable to retrieve entity id '%s' from bucket '%s'. %v", entityId, bucketName, err)
	}
}

func modifyEntity(entityId string, entity interface{}, modify func (), bucketName []byte) {
	err := db.Update(func (tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		key := []byte(entityId)
		err := json.Unmarshal(bucket.Get(key), entity)
		if err != nil {
			return err
		}
		modify()
		resultBalue, err := json.Marshal(entity)
		if err != nil {
			return err
		}
		return bucket.Put(key, resultBalue)
	})
	if err != nil {
		log.Fatalf("Unable to modify entity id '%s' in bucket '%s'. %v", entityId, bucketName, err)
	}
}

func retrieveEntities(getObject func (string) interface{}, bucketName []byte) {
	err := db.View(func (tx *bolt.Tx) error {
		return tx.Bucket(bucketName).ForEach(func (k, v []byte) error {
			key := string(k)
			if key == string(resources.DB_LAST_ID_KEY) {
				return nil
			}
			return json.Unmarshal(v, getObject(key))
		})
	})
	if err != nil {
		log.Fatalf("Unable to retrieve entities from bucket '%s'. %v", bucketName, err)
	}
}

func mapEntities(mapFunc func (string, []byte) error, bucketName []byte) {
	err := db.Update(func (tx *bolt.Tx) error {
		return tx.Bucket(bucketName).ForEach(func (k, v []byte) error {
			key := string(k)
			if key == string(resources.DB_LAST_ID_KEY) {
				return nil
			}
			return mapFunc(key, v)
		})
	})
	if err != nil {
		log.Fatalf("Unable to map entities from bucket '%s'. %v", bucketName, err)
	}
}

func filterEntities(filter func (string, []byte) error, bucketName []byte) {
	err := db.View(func (tx *bolt.Tx) error {
		return tx.Bucket(bucketName).ForEach(func (k, v []byte) error {
			key := string(k)
			if key == string(resources.DB_LAST_ID_KEY) {
				return nil
			}
			return filter(key, v)
		})
	})
	if err != nil {
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
	db, err = bolt.Open(path, 0644, nil)
	if err != nil {
		log.Fatalf("Unable to open database '%s'. %v", path, err)
	}
}

func initializeBuckets(bucketsToInitialize [][]byte) {
	err := db.Update(func (tx *bolt.Tx) error {
		for i := 0; i < len(bucketsToInitialize); i++ {
			_, err := tx.CreateBucketIfNotExists(bucketsToInitialize[i])
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Unable to initialize buckets %v. %v", bucketsToInitialize, err)
	}
}