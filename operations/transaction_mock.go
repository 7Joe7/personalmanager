package operations

import (
	"fmt"
	"github.com/7joe7/personalmanager/resources"
)

const (
	GET_VALUE_CALLED_FORMAT         = "GetValue%s%s"
	SET_VALUE_CALLED_FORMAT         = "SetValue%s%s%s"
	ENSURE_VALUE_CALLED_FORMAT      = "EnsureValue%s%s%v"
	MODIFY_VALUE_CALLED_FORMAT      = "ModifyValue%s%s%v"
	ENSURE_ENTITY_CALLED_FORMAT     = "EnsureEntity%s%s%v"
	ADD_ENTITY_CALLED_FORMAT        = "AddEntity%s%v"
	DELETE_ENTITY_CALLED_FORMAT     = "DeleteEntity%s%s"
	RETRIEVE_ENTITY_CALLED_FORMAT   = "RetrieveEntity%s%s%v%v"
	RETRIEVE_ENTITIES_CALLED_FORMAT = "RetrieveEntities%s%v%v"
	MODIFY_ENTITY_CALLED_FORMAT     = "ModifyEntity%s%v%s%v%v"
	MAP_ENTITIES_CALLED_FORMAT      = "MapEntities%s%v%v%v"
	INITIALIZE_BUCKET_CALLED_FORMAT = "InitializeBucket%s"
	FILTER_ENTITIES_CALLED_FORMAT   = "FilterEntities%s%v%v%v%v"
	EXECUTE_CALLED_FORMAT           = "Execute"
	VIEW_CALLED_FORMAT              = "View"
	ADD_CALLED_FORMAT               = "Add"
)

type transactionMock struct {
	functionsCalled []string
	execs           []func() error
}

func (tm *transactionMock) GetValue(bucketName, key []byte) []byte {
	tm.functionsCalled = append(tm.functionsCalled, fmt.Sprintf(GET_VALUE_CALLED_FORMAT, string(bucketName), string(key)))
	return nil
}

func (tm *transactionMock) SetValue(bucketName, key, value []byte) error {
	tm.functionsCalled = append(tm.functionsCalled, fmt.Sprintf(SET_VALUE_CALLED_FORMAT, string(bucketName), string(key), string(value)))
	return nil
}

func (tm *transactionMock) EnsureValue(bucketName, key, defaultValue []byte) error {
	tm.functionsCalled = append(tm.functionsCalled, fmt.Sprintf(ENSURE_VALUE_CALLED_FORMAT, string(bucketName), string(key), string(defaultValue)))
	return nil
}

func (tm *transactionMock) ModifyValue(bucketName, key []byte, modify func ([]byte) []byte) error {
	tm.functionsCalled = append(tm.functionsCalled, fmt.Sprintf(MODIFY_VALUE_CALLED_FORMAT, string(bucketName), string(key), modify))
	return nil
}

func (tm *transactionMock) EnsureEntity(bucketName, key []byte, entity resources.Entity) error {
	tm.functionsCalled = append(tm.functionsCalled, fmt.Sprintf(ENSURE_ENTITY_CALLED_FORMAT, string(bucketName), string(key), entity))
	return nil
}

func (tm *transactionMock) AddEntity(bucketName []byte, entity resources.Entity) error {
	tm.functionsCalled = append(tm.functionsCalled, fmt.Sprintf(ADD_ENTITY_CALLED_FORMAT, string(bucketName), entity))
	return nil
}

func (tm *transactionMock) DeleteEntity(bucketName, id []byte) error {
	tm.functionsCalled = append(tm.functionsCalled, fmt.Sprintf(DELETE_ENTITY_CALLED_FORMAT, string(bucketName), string(id)))
	return nil
}

func (tm *transactionMock) RetrieveEntity(bucketName, id []byte, entity resources.Entity, shallow bool) error {
	tm.functionsCalled = append(tm.functionsCalled, fmt.Sprintf(RETRIEVE_ENTITY_CALLED_FORMAT, string(bucketName), string(id), entity, shallow))
	return nil
}

func (tm *transactionMock) ModifyEntity(bucketName, key []byte, shallow bool, entity resources.Entity, modifyFunc func()) error {
	tm.functionsCalled = append(tm.functionsCalled, fmt.Sprintf(MODIFY_ENTITY_CALLED_FORMAT, string(bucketName), shallow, string(key), entity, modifyFunc))
	return nil
}

func (tm *transactionMock) MapEntities(bucketName []byte, shallow bool, getNewEntity func () resources.Entity, mapFunc func(resources.Entity) func ()) error {
	tm.functionsCalled = append(tm.functionsCalled, fmt.Sprintf(MAP_ENTITIES_CALLED_FORMAT, string(bucketName), shallow, getNewEntity, mapFunc))
	return nil
}

func (tm *transactionMock) InitializeBucket(bucketName []byte) error {
	tm.functionsCalled = append(tm.functionsCalled, fmt.Sprintf(INITIALIZE_BUCKET_CALLED_FORMAT, string(bucketName)))
	return nil
}

func (tm *transactionMock) RetrieveEntities(bucketName []byte, shallow bool, getObject func(string) resources.Entity) error {
	tm.functionsCalled = append(tm.functionsCalled, fmt.Sprintf(RETRIEVE_ENTITIES_CALLED_FORMAT, string(bucketName), shallow, getObject))
	return nil
}

func (tm *transactionMock) FilterEntities(bucketName []byte, shallow bool, addEntity func (), getNewEntity func () resources.Entity, filterFunc func () bool) error {
	tm.functionsCalled = append(tm.functionsCalled, fmt.Sprintf(FILTER_ENTITIES_CALLED_FORMAT, string(bucketName), shallow, addEntity, getNewEntity, filterFunc))
	return nil
}

func (tm *transactionMock) Execute() {
	tm.functionsCalled = append(tm.functionsCalled, EXECUTE_CALLED_FORMAT)
	for i := 0; i < len(tm.execs); i++ { tm.execs[i]() }
}

func (tm *transactionMock) View() {
	tm.functionsCalled = append(tm.functionsCalled, VIEW_CALLED_FORMAT)
	for i := 0; i < len(tm.execs); i++ { tm.execs[i]() }
}

func (tm *transactionMock) Add(exec func() error) {
	tm.functionsCalled = append(tm.functionsCalled, ADD_CALLED_FORMAT)
	tm.execs = append(tm.execs, exec)
}