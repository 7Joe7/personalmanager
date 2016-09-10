package db

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/7joe7/personalmanager/resources"
)

func getMapEntitiesInner(bucketName []byte, entity interface{}, mapFunc func ()) func (*bolt.Tx) error {
	return func (tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		return b.ForEach(func (k, v []byte) error {
			return modifyEntityInner(b, k, v, entity, mapFunc)
		})
	}
}

func getModifyEntityInner(bucketName []byte, entityId []byte, entity interface{}, modify func ()) func (*bolt.Tx) error {
	return func (tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		return modifyEntityInner(b, entityId, b.Get(entityId), entity, modify)
	}
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

func getFilterEntitiesInner(bucketName []byte, entity interface{}, filterFunc func (string)) func (*bolt.Tx) error {
	return func (tx *bolt.Tx) error {
		return tx.Bucket(bucketName).ForEach(func (k, v []byte) error {
			key := string(k)
			if key == string(resources.DB_LAST_ID_KEY) {
				return nil
			}
			if err := json.Unmarshal(v, entity); err != nil {
				return err
			}
			filterFunc(key)
			return nil
		})
	}
}

func getRetrieveEntitiesInner(getObject func (string) interface{}, bucketName []byte) func (*bolt.Tx) error {
	return func (tx *bolt.Tx) error {
		return tx.Bucket(bucketName).ForEach(func (k, v []byte) error {
			key := string(k)
			if key == string(resources.DB_LAST_ID_KEY) {
				return nil
			}
			return json.Unmarshal(v, getObject(key))
		})
	}
}

func getAddEntityInner(entity interface{}, bucketName []byte, incrementedId *string) func (*bolt.Tx) error {
	return func (tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		id := getIncrementedId(bucket)
		incrementedId = &id
		value, err := json.Marshal(entity)
		if err != nil {
			return err
		}
		incrementedIdBytes := []byte(id)
		err = bucket.Put(incrementedIdBytes, value)
		if err != nil {
			return err
		}
		return bucket.Put(resources.DB_LAST_ID_KEY, incrementedIdBytes)
	}
}
