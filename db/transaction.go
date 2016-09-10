package db

import (
	"encoding/json"
	"log"

	"github.com/boltdb/bolt"
)

type Transaction struct {
	tx *bolt.Tx
	execs []func () error
}

func newTransaction() *Transaction {
	return &Transaction{execs:[]func () error {}}
}

func (t *Transaction) GetValue(bucketName, key []byte) []byte {
	return t.tx.Bucket(bucketName).Get(key)
}

func (t *Transaction) SetValue(bucketName, key, value []byte) error {
	return t.tx.Bucket(bucketName).Put(key, value)
}

func (t *Transaction) EnsureEntity(bucketName, key []byte, entity interface{}) error {
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

func (t *Transaction) AddEntity(bucketName []byte, entity interface{}) (*string, error) {
	var incrementedId *string
	if err := getAddEntityInner(entity, bucketName, incrementedId)(t.tx); err != nil {
		return nil, err
	}
	return incrementedId, nil
}

func (t *Transaction) RetrieveEntity(bucketName, id []byte, entity interface{}) error {
	return json.Unmarshal(t.tx.Bucket(bucketName).Get(id), entity)
}

func (t *Transaction) ModifyEntity(bucketName, key []byte, entity interface{}, modifyFunc func ()) error {
	b := t.tx.Bucket(bucketName)
	return modifyEntityInner(b, key, b.Get(key), entity, modifyFunc)
}

func (t *Transaction) MapEntities(bucketName []byte, entity interface{}, mapFunc func ()) error {
	b := t.tx.Bucket(bucketName)
	return b.ForEach(func (k, v []byte) error {
		return modifyEntityInner(b, k, v, entity, mapFunc)
	})
}

func (t *Transaction) InitializeBucket(bucketName []byte) error {
	_, err := t.tx.CreateBucketIfNotExists(bucketName)
	return err
}

func (t *Transaction) Execute() {
	err := db.Update(func (tx *bolt.Tx) error {
		for i := 0; i < len(t.execs); i++ {
			t.tx = tx
			if err := t.execs[1](); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Unable to execute transaction. %v", err)
	}
}

func (t *Transaction) Add(exec func () error) {
	t.execs = append(t.execs, exec)
}
