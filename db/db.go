package db

import (
	"github.com/7joe7/personalmanager/resources"
)

func Open() {
	open(resources.DB_PATH)
}

func AddEntity(entity interface{}, bucketName []byte) string {
	return addEntity(entity, bucketName)
}

func DeleteEntity(entityId []byte, bucketName []byte) {
	deleteEntity(entityId, bucketName)
}

func RetrieveEntity(bucketName, entityId []byte, entity interface{}) {
	retrieveEntity(bucketName, entityId, entity)
}

func ModifyEntity(bucketName []byte, entityId []byte, entity interface{}, modify func ()) {
	modifyEntity(bucketName, entityId, entity, modify)
}

func RetrieveEntities(bucketName []byte, getObject func (string) interface{}) {
	retrieveEntities(bucketName, getObject)
}

func MapEntities(mapFunc func (), entity interface{}, bucketName []byte) {
	mapEntities(mapFunc, entity, bucketName)
}

func FilterEntities(bucketName []byte, entity interface{}, filterFunc func (string)) {
	filterEntities(bucketName, entity, filterFunc)
}

func NewTransaction() *Transaction {
	return newTransaction()
}