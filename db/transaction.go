package db

import (
	"encoding/json"
	"log"

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
	log.Printf(`Getting value:
	bucketName: %s,
	key: %s.`, string(bucketName), string(key))
	return t.tx.Bucket(bucketName).Get(key)
}

func (t *transaction) SetValue(bucketName, key, value []byte) error {
	log.Printf(`Setting value:
	bucketName: %s,
	key: %s,
	value: %s.`, string(bucketName), string(key), string(value))
	return t.tx.Bucket(bucketName).Put(key, value)
}

func (t *transaction) ModifyValue(bucketName, key []byte, modify func ([]byte) []byte) error {
	log.Printf(`Modifying value:
	bucketName: %s,
	key: %s.`, string(bucketName), string(key))
	b := t.tx.Bucket(bucketName)
	return b.Put(key, modify(b.Get(key)))
}

func (t *transaction) EnsureValue(bucketName, key, defaultValue []byte) error {
	log.Printf(`Ensuring value:
	bucketName: %s,
	key: %s,
	defaultValue: %s.`, string(bucketName), string(key), string(defaultValue))
	b := t.tx.Bucket(bucketName)
	if b.Get(key) == nil {
		return b.Put(key, defaultValue)
	}
	return nil
}

func (t *transaction) EnsureEntity(bucketName, key []byte, entity resources.Entity) error {
	log.Printf(`Ensuring entity:
	bucketName: %s,
	key: %s,
	entity: %v.`, string(bucketName), string(key), entity)
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
	log.Printf(`Adding entity:
	bucketName: %s,
	entity: %v.`, bucketName, entity)
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
	log.Printf(`Deleting entity:
	bucketName: %s,
	id: %s.`, string(bucketName), string(id))
	return t.tx.Bucket(bucketName).Delete(id)
}

func (t *transaction) RetrieveEntity(bucketName, id []byte, entity resources.Entity, shallow bool) error {
	log.Printf(`Retriving entity:
	bucketName: %s,
	id: %s,
	entity: %v,
	shallow: %v.`, string(bucketName), string(id), entity, shallow)
	if err := json.Unmarshal(t.tx.Bucket(bucketName).Get(id), entity); err != nil {
		return err
	}
	if shallow {
		return nil
	}
	return entity.Load(t)
}

func (t *transaction) RetrieveEntities(bucketName []byte, shallow bool, getObject func (string) resources.Entity) error {
	log.Printf(`Retrieving entities:
	bucketName: %s,
	shallow: %v`, string(bucketName), shallow)
	return t.tx.Bucket(bucketName).ForEach(func (k, v []byte) error {
		key := string(k)
		if key == string(resources.DB_LAST_ID_KEY) {
			return nil
		}
		entity := getObject(key)
		if err := json.Unmarshal(v, entity); err != nil {
			return err
		}
		if shallow {
			return nil
		}
		return entity.Load(t)
	})
}

func (t *transaction) ModifyEntity(bucketName, key []byte, shallow bool, entity resources.Entity, modifyFunc func ()) error {
	log.Printf(`Modifying entity:
	bucketName: %s,
	key: %s,
	shallow: %v,
	entity: %v.`, string(bucketName), string(key), shallow, entity)
	b := t.tx.Bucket(bucketName)
	return t.modifyEntityInner(b, key, b.Get(key), shallow, entity, modifyFunc)
}

func (t *transaction) MapEntities(bucketName []byte, shallow bool, getNewEntity func () resources.Entity, mapFunc func (resources.Entity) func ()) error {
	log.Printf(`Mapping entities:
	bucketName: %s,
	shallow: %v.`, string(bucketName), shallow)
	b := t.tx.Bucket(bucketName)
	return b.ForEach(func (k, v []byte) error {
		if string(k) != string(resources.DB_LAST_ID_KEY) {
			entity := getNewEntity()
			return t.modifyEntityInner(b, k, v, shallow, entity, mapFunc(entity))
		}
		return nil
	})
}

func (t *transaction) modifyEntityInner(bucket *bolt.Bucket, key, value []byte, shallow bool, entity resources.Entity, modify func ()) error {
	if err := json.Unmarshal(value, entity); err != nil {
		return err
	}
	if !shallow {
		entity.Load(t)
	}
	modify()
	resultValue, err := json.Marshal(entity)
	if err != nil {
		return err
	}
	return bucket.Put(key, resultValue)
}

func (t *transaction) FilterEntities(bucketName []byte, shallow bool, addEntity func (), getNewEntity func () resources.Entity, filterFunc func () bool) error {
	log.Printf(`Filtering entities:
	bucketName: %s,
	shallow: %v`, string(bucketName), shallow)
	return t.tx.Bucket(bucketName).ForEach(func (k, v []byte) error {
		key := string(k)
		if key == string(resources.DB_LAST_ID_KEY) {
			return nil
		}
		entity := getNewEntity()
		err := json.Unmarshal(v, entity)
		if err != nil {
			return err
		}
		if filterFunc() {
			addEntity()
			if shallow {
				return nil
			}
			err = entity.Load(t)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (t *transaction) InitializeBucket(bucketName []byte) error {
	log.Printf("Initializing bucket %s\n", string(bucketName))
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
