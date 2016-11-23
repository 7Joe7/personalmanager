package db

import (
	"fmt"

	"github.com/7joe7/personalmanager/resources"
)

func Open(path string) {
	if err := open(path); err != nil {
		panic(fmt.Errorf("Unable to open database '%s'. %v", path, err))
	}
}

func DeleteEntity(bucketName, id []byte) {
	if err := deleteEntity(bucketName, id); err != nil {
		panic(err)
	}
}

func ModifyEntity(bucketName, id []byte, shallow bool, entity resources.Entity, modify func ()) {
	if err := modifyEntity(bucketName, id, shallow, entity, modify); err != nil {
		panic(fmt.Errorf("Unable to modify entity id '%s' in bucket '%s'. %v", string(id), bucketName, err))
	}
}

func RetrieveEntities(bucketName []byte, shallow bool, getObject func (string) resources.Entity) {
	if err := retrieveEntities(bucketName, shallow, getObject); err != nil {
		panic(fmt.Errorf("Unable to retrieve entities from bucket '%s'. %v", bucketName, err))
	}
}

func MapEntities(entity resources.Entity, bucketName []byte, shallow bool, mapFunc func ()) {
	if err := mapEntities(entity, bucketName, shallow, mapFunc); err != nil {
		panic(fmt.Errorf("Unable to map entities from bucket '%s'. %v", bucketName, err))
	}
}

func FilterEntities(bucketName []byte, shallow bool, addEntity func(), getNewEntity func () resources.Entity, filterFunc func () bool) {
	if err := filterEntities(bucketName, shallow, addEntity, getNewEntity, filterFunc); err != nil {
		panic(fmt.Errorf("Unable to retrieve entities from bucket '%s'. %v", bucketName, err))
	}
}

func NewTransaction() resources.Transaction {
	return newTransaction()
}