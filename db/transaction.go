package db

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/7joe7/personalmanager/resources"
)

type transaction struct {
	tx *bolt.Tx
	execs []func () error
}

func newTransaction() *transaction {
	return &transaction{execs:[]func () error {}}
}

func (t *transaction) GetValue(bucketName, key []byte) []byte {
	return t.tx.Bucket(bucketName).Get(key)
}

func (t *transaction) SetValue(bucketName, key, value []byte) error {
	return t.tx.Bucket(bucketName).Put(key, value)
}

func (t *transaction) EnsureEntity(bucketName, key []byte, entity resources.Entity) error {
	b := t.tx.Bucket(bucketName)
	if b.Get(key) == nil {
		v, err := json.Marshal(entity)
		if err != nil {
			return err
		}
		return b.Put(key, v)
	}
	return nil
}

func (t *transaction) AddEntity(bucketName []byte, entity resources.Entity) error {
	bucket := t.tx.Bucket(bucketName)
	id := getIncrementedId(bucket)
	entity.SetId(id)
	value, err := json.Marshal(entity)
	if err != nil {
		return err
	}
	incrementedIdBytes := []byte(id)
	if err = bucket.Put(incrementedIdBytes, value); err != nil {
		return err
	}
	return bucket.Put(resources.DB_LAST_ID_KEY, incrementedIdBytes)
}

func (t *transaction) DeleteEntity(bucketName, id []byte) error {
	return t.tx.Bucket(bucketName).Delete(id)
}

func (t *transaction) RetrieveEntity(bucketName, id []byte, entity resources.Entity) error {
	if err := json.Unmarshal(t.tx.Bucket(bucketName).Get(id), entity); err != nil {
		return err
	}
	return entity.Load(t)
}

func (t *transaction) RetrieveEntities(bucketName []byte, getObject func (string) resources.Entity) error {
	return t.tx.Bucket(bucketName).ForEach(func (k, v []byte) error {
		key := string(k)
		if key == string(resources.DB_LAST_ID_KEY) {
			return nil
		}
		entity := getObject(key)
		if err := json.Unmarshal(v, entity); err != nil {
			return err
		}
		return entity.Load(t)
	})
}

func (t *transaction) ModifyEntity(bucketName, key []byte, entity resources.Entity, modifyFunc func ()) error {
	b := t.tx.Bucket(bucketName)
	return modifyEntityInner(b, key, b.Get(key), entity, modifyFunc)
}

func (t *transaction) MapEntities(bucketName []byte, entity resources.Entity, mapFunc func ()) error {
	b := t.tx.Bucket(bucketName)
	return b.ForEach(func (k, v []byte) error {
		if string(k) != string(resources.DB_LAST_ID_KEY) {
			return modifyEntityInner(b, k, v, entity, mapFunc)
		}
		return nil
	})
}

func (t *transaction) FilterEntities(bucketName []byte, entity resources.Entity, filterFunc func () bool, copyFunc func ()) error {
	return t.tx.Bucket(bucketName).ForEach(func (k, v []byte) error {
		key := string(k)
		if key == string(resources.DB_LAST_ID_KEY) {
			return nil
		}
		if err := json.Unmarshal(v, entity); err != nil {
			return err
		}
		if filterFunc() {
			copyFunc()
		}
		return nil
	})
}

func (t *transaction) InitializeBucket(bucketName []byte) error {
	_, err := t.tx.CreateBucketIfNotExists(bucketName)
	return err
}

func (t *transaction) execute() error {
	return db.Update(t.executeAll)
}

func (t *transaction) view() error {
	return db.View(t.executeAll)
}

func (t *transaction) executeAll(tx *bolt.Tx) error {
	for i := 0; i < len(t.execs); i++ {
		t.tx = tx
		if err := t.execs[i](); err != nil {
			return err
		}
	}
	return nil
}

func (t *transaction) Execute() {
	if err := t.execute(); err != nil {
		panic(err)
	}
}

func (t *transaction) View() {
	if err := t.view(); err != nil {
		panic(err)
	}
}

func (t *transaction) Add(exec func () error) {
	t.execs = append(t.execs, exec)
}
