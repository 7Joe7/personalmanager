package db

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/7joe7/personalmanager/resources"
)

type Transaction interface {
	GetValue(bucketName, key []byte) []byte
	SetValue(bucketName, key, value []byte) error
	EnsureEntity(bucketName, key []byte, entity interface{}) error
	AddEntity(bucketName []byte, entity resources.Entity) (string, error)
	DeleteEntity(bucketName, id []byte) error
	RetrieveEntity(bucketName, id []byte, entity interface{}) error
	RetrieveEntities(bucketName []byte, getObject func (string) interface{}) error
	ModifyEntity(bucketName, key []byte, entity interface{}, modifyFunc func ()) error
	MapEntities(bucketName []byte, entity interface{}, mapFunc func ()) error
	InitializeBucket(bucketName []byte) error
	Execute()
	Add(exec func () error)
}

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

func (t *transaction) EnsureEntity(bucketName, key []byte, entity interface{}) error {
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

func (t *transaction) AddEntity(bucketName []byte, entity resources.Entity) (string, error) {
	return getAddEntityInner(entity, bucketName)(t.tx)
}

func (t *transaction) DeleteEntity(bucketName, id []byte) error {
	return t.tx.Bucket(bucketName).Delete(id)
}

func (t *transaction) RetrieveEntity(bucketName, id []byte, entity interface{}) error {
	return json.Unmarshal(t.tx.Bucket(bucketName).Get(id), entity)
}

func (t *transaction) RetrieveEntities(bucketName []byte, getObject func (string) interface{}) error {
	return getRetrieveEntitiesInner(getObject, bucketName)(t.tx)
}

func (t *transaction) ModifyEntity(bucketName, key []byte, entity interface{}, modifyFunc func ()) error {
	b := t.tx.Bucket(bucketName)
	return modifyEntityInner(b, key, b.Get(key), entity, modifyFunc)
}

func (t *transaction) MapEntities(bucketName []byte, entity interface{}, mapFunc func ()) error {
	b := t.tx.Bucket(bucketName)
	return b.ForEach(func (k, v []byte) error {
		if string(k) != string(resources.DB_LAST_ID_KEY) {
			return modifyEntityInner(b, k, v, entity, mapFunc)
		}
		return nil
	})
}

func (t *transaction) InitializeBucket(bucketName []byte) error {
	_, err := t.tx.CreateBucketIfNotExists(bucketName)
	return err
}

func (t *transaction) execute() error {
	return db.Update(func (tx *bolt.Tx) error {
		for i := 0; i < len(t.execs); i++ {
			t.tx = tx
			if err := t.execs[i](); err != nil {
				return err
			}
		}
		return nil
	})
}

func (t *transaction) Execute() {
	if err := t.execute(); err != nil {
		panic(err)
	}
}

func (t *transaction) Add(exec func () error) {
	t.execs = append(t.execs, exec)
}
