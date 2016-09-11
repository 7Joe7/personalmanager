package db

import (
	"fmt"

	"github.com/7joe7/personalmanager/resources"
)

func Open(path string) {
	if err := open(path); err != nil {
		panic(fmt.Errorf("Unable to open database '%s'. %v", resources.DB_PATH, err))
	}
}

func DeleteEntity(bucketName, id []byte) {
	if err := deleteEntity(bucketName, id); err != nil {
		panic(err)
	}
}

func ModifyEntity(bucketName []byte, entityId []byte, entity resources.Entity, modify func ()) {
	if err := modifyEntity(bucketName, entityId, entity, modify); err != nil {
		panic(fmt.Errorf("Unable to modify entity id '%s' in bucket '%s'. %v", string(entityId), bucketName, err))
	}
}

func RetrieveEntities(bucketName []byte, getObject func (string) resources.Entity) {
	if err := retrieveEntities(bucketName, getObject); err != nil {
		panic(fmt.Errorf("Unable to retrieve entities from bucket '%s'. %v", bucketName, err))
	}
}

func MapEntities(entity resources.Entity, bucketName []byte, mapFunc func ()) {
	if err := mapEntities(entity, bucketName, mapFunc); err != nil {
		panic(fmt.Errorf("Unable to map entities from bucket '%s'. %v", bucketName, err))
	}
}

func FilterEntities(bucketName []byte, entity resources.Entity, filterFunc func () bool, copyFunc func ()) {
	if err := filterEntities(bucketName, entity, filterFunc, copyFunc); err != nil {
		panic(fmt.Errorf("Unable to retrieve entities from bucket '%s'. %v", bucketName, err))
	}
}

func NewTransaction() resources.Transaction {
	return newTransaction()
}