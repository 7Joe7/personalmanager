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

func AddEntity(entity resources.Entity, bucketName []byte) string {
	id, err := addEntity(entity, bucketName)
	if err != nil {
		panic(fmt.Errorf("Unable to add entity to bucket '%s'. %v", bucketName, err))
	}
	return id
}

func DeleteEntity(entityId []byte, bucketName []byte) {
	if err := deleteEntity(entityId, bucketName); err != nil {
		panic(fmt.Errorf("Unable to delete entity id '%s' from bucket '%s'. %v", string(entityId), bucketName, err))
	}
}

func RetrieveEntity(bucketName, entityId []byte, entity interface{}) {
	if err := retrieveEntity(bucketName, entityId, entity); err != nil {
		panic(fmt.Errorf("Unable to retrieve entity id '%s' from bucket '%s'. %v", string(entityId), bucketName, err))
	}
}

func ModifyEntity(bucketName []byte, entityId []byte, entity interface{}, modify func ()) {
	if err := modifyEntity(bucketName, entityId, entity, modify); err != nil {
		panic(fmt.Errorf("Unable to modify entity id '%s' in bucket '%s'. %v", string(entityId), bucketName, err))
	}
}

func RetrieveEntities(bucketName []byte, getObject func (string) interface{}) {
	if err := retrieveEntities(bucketName, getObject); err != nil {
		panic(fmt.Errorf("Unable to retrieve entities from bucket '%s'. %v", bucketName, err))
	}
}

func MapEntities(entity interface{}, bucketName []byte, mapFunc func ()) {
	if err := mapEntities(entity, bucketName, mapFunc); err != nil {
		panic(fmt.Errorf("Unable to map entities from bucket '%s'. %v", bucketName, err))
	}
}

func FilterEntities(bucketName []byte, entity interface{}, filterFunc func () bool, copyFunc func ()) {
	if err := filterEntities(bucketName, entity, filterFunc, copyFunc); err != nil {
		panic(fmt.Errorf("Unable to retrieve entities from bucket '%s'. %v", bucketName, err))
	}
}

func NewTransaction() Transaction {
	return newTransaction()
}